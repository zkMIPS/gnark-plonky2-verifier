package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"github.com/succinctlabs/gnark-plonky2-verifier/certificate/data"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")

	inputDir   = flag.String("input_dir", "/efs/zkm/test/test_proof/proof/test_cache3/aggregate", "The request input dir")
	outputPath = flag.String("output_dir", "testdata/mips/groth16.proof", "The request output dir")
)

func reqSnarkProof(client pb.ProverServiceClient, req *pb.FinalProofRequest) {
	log.Printf("Requesting snark proof for stark proof")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.FinalProof(ctx, req)
	if err != nil {
		log.Fatalf("client.FinalProof failed: %v", err)
	}
	log.Printf("resp:%+v", resp)
}

func getSnarkProofJobTask(client pb.ProverServiceClient, req *pb.GetTaskResultRequest) {
	log.Printf("Query snark proof job result")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetTaskResult(ctx, req)
	if err != nil {
		log.Fatalf("client.GetTaskResult failed: %v", err)
	}
	log.Printf("resp:%+v", resp)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewProverServiceClient(conn)

	namespace := uuid.New()
	name := []byte("zkm")
	u := uuid.NewSHA1(namespace, name)

	commonCircuitData, err := os.ReadFile(*inputDir + "/common_circuit_data.json")
	if err != nil {
		log.Fatalf("fail to read common_circuit_data: %v", err)
	}
	verifierOnlyCircuitData, err := os.ReadFile(*inputDir + "/verifier_only_circuit_data.json")
	if err != nil {
		log.Fatalf("fail to read verifier_only_circuit_data: %v", err)
	}
	proofWithPublicInputs, err := os.ReadFile(*inputDir + "/proof_with_public_inputs.json")
	if err != nil {
		log.Fatalf("fail to read proof_with_public_inputs: %v", err)
	}
	// Request snark proof
	reqSnarkProof(client, &pb.FinalProofRequest{
		ChainId:                 11155111,
		Timestamp:               uint64(time.Now().Unix()),
		ProofId:                 u.String(),
		ComputedRequestId:       u.String(),
		CommonCircuitData:       commonCircuitData,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
		ProofWithPublicInputs:   proofWithPublicInputs,
	})

	time.Sleep(2*time.Minute)

	// Query snark proof
	getSnarkProofJobTask(client, &pb.GetTaskResultRequest{
		ProofId:           u.String(),
		ComputedRequestId: u.String(),
	})
}
