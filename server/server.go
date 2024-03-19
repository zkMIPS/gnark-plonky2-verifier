package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/succinctlabs/gnark-plonky2-verifier/certificate/data"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	tls             = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile        = flag.String("cert_file", "", "The TLS cert file")
	keyFile         = flag.String("key_file", "", "The TLS key file")
	port            = flag.Int("port", 50051, "The server port")
	workerName      = flag.String("prover_worker_name", "groth16_prover", "The prover worker name")
	proverCycleTime = flag.Uint64("prover_cycle_time", 1000, "The prover cycle time")
	proverHeartBeat = flag.Uint64("prover_heart_beat", 2000, "The prover heart beat")
	proverTimeout   = flag.Uint64("prover_time_out", 60*60, "The prover time out")
	dbUser          = flag.String("db_user", "root", "The database username")
	dbPassword      = flag.String("db_password", "123456", "The database password")
	dbHost          = flag.String("db_host", "127.0.0.1", "The datbase host")
	dbPort          = flag.String("db_port", "3306", "The database port")
	dbName          = flag.String("db_name", "zkm", "The database name")

	profileCircuit = flag.Bool("profile", false, "profile the circuit")
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
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
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	connectDatabase()
	defer db.Close()

	go func(workerName string, interval uint64, timeout uint64, heartBeat uint64) {
		proverWorkCycle(workerName, interval, timeout, heartBeat)
	}(*workerName, *proverCycleTime, *proverTimeout, *proverHeartBeat)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterProverServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
