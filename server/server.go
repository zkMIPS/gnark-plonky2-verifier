package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/succinctlabs/gnark-plonky2-verifier/certificate/data"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"sync"
)

var (
	tls             = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile        = flag.String("cert_file", "", "The TLS cert file")
	keyFile         = flag.String("key_file", "", "The TLS key file")
	port            = flag.Int("port", 50051, "The server port")
	workerName      = flag.String("prover_worker_name", "groth16_prover", "The prover worker name")
	proverCycleTime = flag.Uint64("prover_cycle_time", 1000, "The prover cycle time")
	proverTimeout   = flag.Uint64("prover_time_out", 60*60, "The prover time out")
)

type proverService struct {
	pb.UnimplementedProverServiceServer

	mu sync.Mutex
}

func (s *proverService) GetStatus(ctx context.Context, in *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	// TODO
	return &pb.GetStatusResponse{}, nil
}

func (s *proverService) SplitElf(ctx context.Context, in *pb.SplitElfRequest) (*pb.SplitElfResponse, error) {
	// TODO
	return &pb.SplitElfResponse{}, nil
}

func (s *proverService) Prove(ctx context.Context, in *pb.ProveRequest) (*pb.ProveResponse, error) {
	// TODO
	return &pb.ProveResponse{}, nil
}

func (s *proverService) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	// TODO
	return &pb.AggregateResponse{}, nil
}

func (s *proverService) AggregateAll(ctx context.Context, in *pb.AggregateAllRequest) (*pb.AggregateAllResponse, error) {
	// TODO

	return &pb.AggregateAllResponse{}, nil
}

func (s *proverService) GetTaskResult(context.Context, *pb.GetTaskResultRequest) (*pb.GetTaskResultResponse, error) {
	// TODO

	return &pb.GetTaskResultResponse{}, nil
}

func (s *proverService) FinalProof(ctx context.Context, req *pb.FinalProofRequest) (*pb.FinalProofResponse, error) {
	addProverJobToQueue(ctx, req)

	return &pb.FinalProofResponse{ProofId: req.ProofId, ComputedRequestId: req.ComputedRequestId, Result: getSuccessResult(pb.ResultCode_RESULT_OK, "get task successfully")}, nil
}

func newServer() *proverService {
	s := &proverService{}
	return s
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

	go func(workerName string, interval uint64, timeout uint64) {
		proverWorkCycle(workerName, interval, timeout)
	}(*workerName, *proverCycleTime, *proverTimeout)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterProverServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
