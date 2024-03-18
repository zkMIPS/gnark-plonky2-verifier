package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"log"
	"time"
)

type ProverJobStatus int

const (
	Idle ProverJobStatus = iota
	InProgress
	Done
)

type ProverJobType int

const (
	SingleProof ProverJobType = iota
	AggregatedProof
)

const SingleProofJobPriority = 1
const AggregatedProofJobPriority = 0

type ProverInputResponse struct {
	JobId             int
	SnarkProofRequest *pb.FinalProofRequest
}

func proverWorkCycle(workerName string, interval uint64, proverTimeout uint64) {
	log.Printf("Running worker cycle")
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		log.Printf("Prover cycle started.")

		resp, err := getJob(workerName)
		if err != nil {
			log.Printf("Failed to get prover job,err :%+v", err)
		}
		if err != nil || resp.SnarkProofRequest == nil {
			continue
		}

		ch := make(chan error)

		go computeProof(resp, ch)

		select {
		case result := <-ch:
			if result != nil {
				log.Printf("Failed to compute groth16 proof,err: %+v", result)
			} else {
				log.Printf("groth16 proof was computed.")
			}
		case <-time.After(time.Duration(proverTimeout) * time.Second):
			log.Printf("Prover timeout.")
		}
	}
}

func getJob(proverName string) (ProverInputResponse, error) {
	log.Printf("Request stark proof to prove from worker: %s", proverName)
	if len(proverName) == 0 {
		return ProverInputResponse{}, fmt.Errorf("empty prover worker name")
	}

	return loadIdleProverJobFromQueue()
}

func loadIdleProverJobFromQueue() (ProverInputResponse, error) {
	var resp = ProverInputResponse{}
	tx, err := db.Begin()
	if err != nil {
		return resp, err
	}

	rows, err := db.Query(
		"SELECT id,job_data FROM prover_job_queue "+
			"WHERE job_status = ? ORDER BY job_priority, id, proof_id, computed_request_id LIMIT 1", Idle)
	if err != nil {
		tx.Rollback()
		return resp, err
	}
	defer rows.Close()

	if rows.Next() { // if and only if one result
		var id int
		var jobData string
		err := rows.Scan(&id, &jobData)
		if err != nil {
			tx.Rollback()
			return resp, err
		}

		_, err = db.Exec("UPDATE prover_job_queue SET (job_status, updated_at, updated_by) = "+
			"?, now(), 'server_give_job') WHERE id = ?", InProgress, id)
		if err != nil {
			tx.Rollback()
			return resp, err
		}

		proofReq := &pb.FinalProofRequest{}

		if err := jsonpb.UnmarshalString(jobData, proofReq); err != nil {
			tx.Rollback()
			return resp, err
		}

		resp = ProverInputResponse{JobId: id, SnarkProofRequest: proofReq}
	}

	if err := tx.Commit(); err != nil {
		return resp, fmt.Errorf("failed to commit transaction for loadIdleProverJobFromQueue")
	}

	return resp, nil
}

func computeProof(job ProverInputResponse, ch chan error) {
	// TODO

	ch <- nil
}

func addProverJobToQueue(ctx context.Context, req *pb.FinalProofRequest) *pb.Result {
	jobData, err := proto.Marshal(req)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		log.Printf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}
	insertQuery := "INSERT INTO prover_job_queue (job_status, job_priority, job_type, updated_by, proof_id, computed_request_id, job_data) VALUES(?,?,?,?,?,?,?)"
	_, err = db.Exec(insertQuery, Idle, SingleProofJobPriority, SingleProof, "server_add_job", req.ProofId, req.ComputedRequestId, string(jobData))
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		log.Printf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}

	return getSuccessResult("get task successfully.")
}

func queryProverJobStatus(req *pb.GetTaskResultRequest) *pb.Result {
	query := "SELECT * FROM proofs WHERE proof_id = ? and computed_request_id = ?"

	rows, err := db.Query(query, req.ProofId, req.ComputedRequestId)
	if err != nil {
		formatStr := "Failed to query proofs db, err: %+v"
		log.Printf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}
	defer rows.Close()

	if rows.Next() { // only one result
		if err != nil {
			formatStr := "Failed to query prover job status result, err: %+v"
			log.Printf(formatStr, err)
			return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
		}
		return getSuccessResult("proofs was generated successfully.")
	}
	return getErrorResult(pb.ResultCode_RESULT_BUSY, fmt.Sprintf("proof hasn't been ready, err: %+v", err))
}

func getSuccessResult(msg string) *pb.Result {
	return &pb.Result{Code: pb.ResultCode_RESULT_OK, Message: msg}
}

func getErrorResult(code pb.ResultCode, errorMsg string) *pb.Result {
	return &pb.Result{Code: code, Message: errorMsg}
}
