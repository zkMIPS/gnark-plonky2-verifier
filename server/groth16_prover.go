package main

import (
	"context"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"log"
	"time"
)

type ProverInputResponse struct {
	JobId             uint64
	SnarkProofRequest *pb.FinalProofRequest
}

func proverWorkCycle(workerName string, interval uint64, proverTimeout uint64) {
	log.Printf("Running worker cycle")
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		log.Printf("Prover cycle started.")

		resp, err := getJob(workerName)
		if err != nil {
			log.Printf("failed to get prover job,err :%+v", err)
		}
		if err != nil || resp.SnarkProofRequest == nil {
			continue
		}

		ch := make(chan error)

		go computeProof(ch)

		select {
		case result := <-ch:
			if result != nil {
				log.Printf("failed to compute groth16 proof,err: %+v", result)
			} else {
				log.Printf("groth16 proof was computed.")
			}
		case <-time.After(time.Duration(proverTimeout) * time.Second):
			log.Printf("Prover timeout.")
		}
	}
}

func getJob(proverName string) (ProverInputResponse, error) {
	// TODO

	return ProverInputResponse{SnarkProofRequest: &pb.FinalProofRequest{}}, nil
}

func computeProof(ch chan error) {
	// TODO

	ch <- nil
}

func addProverJobToQueue(ctx context.Context, req *pb.FinalProofRequest) {

}

func getSuccessResult(code pb.ResultCode, msg string) *pb.Result {
	return &pb.Result{Code: code, Message: msg}
}
