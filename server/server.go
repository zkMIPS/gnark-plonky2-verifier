package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	log "github.com/sirupsen/logrus"
	"github.com/succinctlabs/gnark-plonky2-verifier/certificate/data"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/go-sql-driver/mysql"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tls             = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile        = flag.String("cert_file", "", "The TLS cert file")
	keyFile         = flag.String("key_file", "", "The TLS key file")
	port            = flag.Int("port", 50051, "The server port")
	metricsPort     = flag.Int("metrics_port", 50061, "The metrics port")
	workerName      = flag.String("prover_worker_name", "groth16_prover", "The prover worker name")
	proverCycleTime = flag.Uint64("prover_cycle_time", 1000, "The prover cycle time")
	proverHeartBeat = flag.Uint64("prover_heart_beat", 2000, "The prover heart beat")
	proverTimeout   = flag.Uint64("prover_time_out", 60*60, "The prover time out")
	dbUser          = flag.String("db_user", "root", "The database username")
	dbPassword      = flag.String("db_password", "123456", "The database password")
	dbHost          = flag.String("db_host", "127.0.0.1", "The datbase host")
	dbPort          = flag.String("db_port", "3306", "The database port")
	dbName          = flag.String("db_name", "zkm", "The database name")
	logLevel        = flag.Uint64("log_level", uint64(log.InfoLevel), "The log level")
	cacheDir        = flag.String("cache_dir", "/efs/zkm/test/test_proof/cache_proof", "The circuit and key cache dir")
	inputParentDir  = flag.String("input_parent_dir", "/efs/zkm/test/test_proof/proof", "The stark proof parent dir")

	profileCircuit = flag.Bool("profile", false, "profile the circuit")
)

var (
	r1csCircuit constraint.ConstraintSystem

	pk groth16.ProvingKey
	vk groth16.VerifyingKey
)

var db *sql.DB = nil

type proverService struct {
	pb.UnimplementedProverServiceServer

	mu sync.Mutex
}

func (s *proverService) GetStatus(ctx context.Context, in *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	// Unsupported
	return &pb.GetStatusResponse{}, nil
}

func (s *proverService) SplitElf(ctx context.Context, in *pb.SplitElfRequest) (*pb.SplitElfResponse, error) {
	// Unsupported
	return &pb.SplitElfResponse{}, nil
}

func (s *proverService) Prove(ctx context.Context, in *pb.ProveRequest) (*pb.ProveResponse, error) {
	// Unsupported
	return &pb.ProveResponse{}, nil
}

func (s *proverService) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	// Unsupported
	return &pb.AggregateResponse{}, nil
}

func (s *proverService) AggregateAll(ctx context.Context, in *pb.AggregateAllRequest) (*pb.AggregateAllResponse, error) {
	// Unsupported

	return &pb.AggregateAllResponse{}, nil
}

func (s *proverService) GetTaskResult(ctx context.Context, req *pb.GetTaskResultRequest) (*pb.GetTaskResultResponse, error) {
	result := queryProverJobStatus(req)

	return &pb.GetTaskResultResponse{ProofId: req.ProofId, ComputedRequestId: req.ComputedRequestId, Result: result}, nil
}

func (s *proverService) FinalProof(ctx context.Context, req *pb.FinalProofRequest) (*pb.FinalProofResponse, error) {
	result := addProverJobToQueue(ctx, req)

	return &pb.FinalProofResponse{ProofId: req.ProofId, ComputedRequestId: req.ComputedRequestId, Result: result}, nil
}

func newServer() *proverService {
	s := &proverService{}
	return s
}

func connectDatabase() {
	var err error = nil
	db, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, *dbName))
	if err != nil {
		log.Fatalf("Failed to ")
	}
}

func initCircuitKeys() {
	if r1csCircuit != nil {
		return
	}

	circuitPath := *cacheDir + "/circuit"
	pkPath := *cacheDir + "/proving.key"
	vkPath := *cacheDir + "/verifying.key"

	_, err := os.Stat(circuitPath)

	// generate circuit if it is not exist by precomputed proof
	if os.IsNotExist(err) {
		commonCircuitData, err := types.ReadCommonCircuitData(*cacheDir + "/common_circuit_data.json")
		if err != nil {
			logger().Errorln(err)
			return
		}

		proofWithPisData, err := types.ReadProofWithPublicInputs(*cacheDir + "/proof_with_public_inputs.json")
		if err != nil {
			logger().Errorln(err)
			return
		}
		proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

		verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData(*cacheDir + "/verifier_only_circuit_data.json")
		if err != nil {
			logger().Errorln(err)
			return
		}
		verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)

		circuit := verifier.ExampleVerifierCircuit{
			Proof:                   proofWithPis.Proof,
			PublicInputs:            proofWithPis.PublicInputs,
			VerifierOnlyCircuitData: verifierOnlyCircuitData,
			CommonCircuitData:       commonCircuitData,
		}

		var builder frontend.NewBuilder = r1cs.NewBuilder
		r1csCircuit, err = frontend.Compile(ecc.BN254.ScalarField(), builder, &circuit)
		fR1CS, _ := os.Create(circuitPath)
		r1csCircuit.WriteTo(fR1CS)
		fR1CS.Close()
	} else {
		fCircuit, err := os.Open(circuitPath)
		if err != nil {
			logger().Errorln(err)
			return
		}

		r1csCircuit = groth16.NewCS(ecc.BN254)
		r1csCircuit.ReadFrom(fCircuit)
		fCircuit.Close()
	}

	_, err = os.Stat(pkPath)
	if os.IsNotExist(err) {
		pk, vk, err = groth16.Setup(r1csCircuit)
		if err != nil {
			logger().Errorln(err)
			return
		}

		fPK, _ := os.Create(pkPath)
		pk.WriteTo(fPK)
		fPK.Close()

		if vk != nil {
			fVK, _ := os.Create(vkPath)
			vk.WriteTo(fVK)
			fVK.Close()
		}
	} else {
		fPk, err := os.Open(pkPath)
		if err != nil {
			logger().Errorln(err)
			return
		}
		pk = groth16.NewProvingKey(ecc.BN254)
		pk.ReadFrom(fPk)

		fVk, err := os.Open(vkPath)
		if err != nil {
			logger().Errorln(err)
			return
		}
		vk = groth16.NewVerifyingKey(ecc.BN254)
		vk.ReadFrom(fVk)
		defer fVk.Close()
	}
}

func init() {
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(uint32(*logLevel)))

	initCircuitKeys()
}

func logger() *log.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return log.WithField("file", filename).WithField("function", fn)
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{
			grpc.Creds(creds),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		}
	}

	connectDatabase()
	defer db.Close()

	go func(workerName string, interval uint64, timeout uint64, heartBeat uint64) {
		proverWorkCycle(workerName, interval, timeout, heartBeat)
	}(*workerName, *proverCycleTime, *proverTimeout, *proverHeartBeat)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterProverServiceServer(grpcServer, newServer())

	grpc_prometheus.Register(grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *metricsPort), nil))
	}()

	grpcServer.Serve(lis)
}
