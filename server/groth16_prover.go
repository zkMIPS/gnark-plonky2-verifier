package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/profile"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

type ProverJobStatus int

const (
	Idle ProverJobStatus = iota
	InProgress
	Done
)

type ProverJobType string

const (
	SingleProof ProverJobType = "SINGLE_PROOF"
	AggregatedProof ProverJobType = "AGGREGATED_PROOF"
)

const SingleProofJobPriority = 1
const AggregatedProofJobPriority = 0

type ProverInputResponse struct {
	JobId             int
	SnarkProofRequest *pb.FinalProofRequest
}

type Groth16ProofResult struct {
	ProofId           string
	ComputedRequestId string
	ProofBytes        []byte
	Err               error
}

func proverWorkCycle(workerName string, interval uint64, proverTimeout uint64, heartBeat uint64) {
	logger().Infof("Running worker cycle")
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		logger().Infof("Prover cycle started.")

		job, err := getJob(workerName)
		if err != nil {
			logger().Errorf("Failed to get prover job,err :%+v", err)
		}
		if err != nil || job.SnarkProofRequest == nil {
			continue
		}

		ch := make(chan Groth16ProofResult)

		go computeProof(job, workerName, ch, heartBeat)

		select {
		case result := <-ch:
			if result.Err != nil {
				logger().Errorf("Failed to compute groth16 proof,err: %+v", result.Err)
			} else {
				logger().Infof("groth16 proof was computed.")
				err := storeProof(job, result)
				if err != nil {
					logger().Infof("Fialed to store proof,err: %+v", err)
				}
			}
		case <-time.After(time.Duration(proverTimeout) * time.Second):
			logger().Warnf("Prover timeout.")
		}
	}
}

func getJob(proverName string) (ProverInputResponse, error) {
	logger().Infof("Request stark proof to prove from worker: %s", proverName)
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return resp, rollbackErr
		}
		return resp, err
	}
	defer rows.Close()

	if rows.Next() { // if and only if one result
		var id int
		var jobData string
		err := rows.Scan(&id, &jobData)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return resp, rollbackErr
			}
			return resp, err
		}

		_, err = db.Exec("UPDATE prover_job_queue "+
			"SET job_status=?,updated_at=now(),updated_by='server_give_job' WHERE id=?", InProgress, id)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return resp, rollbackErr
			}
			return resp, err
		}

		proofReq := &pb.FinalProofRequest{}

		if err := jsonpb.UnmarshalString(jobData, proofReq); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return resp, rollbackErr
			}
			return resp, err
		}

		resp = ProverInputResponse{JobId: id, SnarkProofRequest: proofReq}
	}

	if err := tx.Commit(); err != nil {
		return resp, fmt.Errorf("failed to commit transaction for loadIdleProverJobFromQueue")
	}

	return resp, nil
}

func computeProof(job ProverInputResponse, proverName string, ch chan Groth16ProofResult, heartBeat uint64) {
	go func() {
		for {
			time.Sleep(time.Duration(heartBeat) * time.Millisecond)
			recordProverIsWorking(job.JobId, proverName)
		}
	}()

	res := Groth16ProofResult{
		ProofId:           job.SnarkProofRequest.ProofId,
		ComputedRequestId: job.SnarkProofRequest.ComputedRequestId,
		ProofBytes:        []byte{},
		Err:               nil,
	}

	// remove useless slash
	cleanInputDir := filepath.Clean(job.SnarkProofRequest.InputDir)

	commonCircuitData, err := types.ReadCommonCircuitData(cleanInputDir + "/common_circuit_data.json")
	if err != nil {
		res.Err = err
		ch <- res
		return
	}

	proofWithPisData, err := types.ReadProofWithPublicInputs(cleanInputDir + "/proof_with_public_inputs.json")
	if err != nil {
		res.Err = err
		ch <- res
		return
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData(cleanInputDir + "/verifier_only_circuit_data.json")
	if err != nil {
		res.Err = err
		ch <- res
		return
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)

	circuit := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
		CommonCircuitData:       commonCircuitData,
	}

	var p *profile.Profile
	if *profileCircuit {
		p = profile.Start()
	}

	var builder frontend.NewBuilder = r1cs.NewBuilder

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), builder, &circuit)
	if err != nil {
		res.Err = err
		ch <- res
		return
	}

	if *profileCircuit {
		p.Stop()
		p.Top()
		logger().Infof("r1cs.GetNbCoefficients(): %v", r1cs.GetNbCoefficients())
		logger().Infof("r1cs.GetNbConstraints(): %v", r1cs.GetNbConstraints())
		logger().Infof("r1cs.GetNbSecretVariables(): %v", r1cs.GetNbSecretVariables())
		logger().Infof("r1cs.GetNbPublicVariables(): %v", r1cs.GetNbPublicVariables())
		logger().Infof("r1cs.GetNbInternalVariables(): %v", r1cs.GetNbInternalVariables())
	}

	bytes, err := generateGroth16Proof(r1cs, cleanInputDir, job.SnarkProofRequest.OutputPath)
	if err != nil {
		res.Err = err
		ch <- res
		return
	}

	res.ProofBytes = bytes
	ch <- res
}

func generateGroth16Proof(r1cs constraint.ConstraintSystem, inputDir string, outputPath string) ([]byte, error) {
	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey
	var err error

	proofWithPisRawdata, err := types.ReadProofWithPublicInputs(inputDir + "/proof_with_public_inputs.json")
	if err != nil {
		return nil, err
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisRawdata)
	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData(inputDir + "/verifier_only_circuit_data.json")
	if err != nil {
		return nil, err
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	logger().Infof("Running circuit setup: %v", time.Now())
	logger().Infof("Using real setup")
	pk, vk, err = groth16.Setup(r1cs)

	if err != nil {
		return nil, err
	}

	fPK, _ := os.Create(inputDir + "/proving.key")
	pk.WriteTo(fPK)
	fPK.Close()

	if vk != nil {
		fVK, _ := os.Create(inputDir + "/verifying.key")
		vk.WriteTo(fVK)
		fVK.Close()
	}

	logger().Infof("Generating witness: %v", time.Now())
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	fWitness, _ := os.Create(inputDir + "/witness")
	witness.WriteTo(fWitness)
	fWitness.Close()

	logger().Infof("Creating proof: %v", time.Now())
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return nil, err
	}

	fProof, _ := os.Create(outputPath)
	proof.WriteTo(fProof)
	fProof.Close()

	if vk == nil {
		return nil, fmt.Errorf("vk is nil, means you're using dummy setup and we skip verification of proof")
	}

	logger().Infof("Verifying proof: %v", time.Now())
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return nil, err
	}

	const fpSize = 4 * 8
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	var (
		a [2]*big.Int
		b [2][2]*big.Int
		c [2]*big.Int
	)

	// proof.Ar, proof.Bs, proof.Krs
	a[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	a[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	b[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	b[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	b[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	b[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	c[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	c[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	logger().Infof("a[0] is %s", a[0].String())
	logger().Infof("a[1] is %s", a[1].String())

	logger().Infof("b[0][0] is %s", b[0][0].String())
	logger().Infof("b[0][1] is %s", b[0][1].String())
	logger().Infof("b[1][0] is %s", b[1][0].String())
	logger().Infof("b[1][1] is %s", b[1][1].String())

	logger().Infof("c[0] is %s", c[0].String())
	logger().Infof("c[1] is %s", c[1].String())

	return proofBytes, nil
}

func recordProverIsWorking(jobId int, proverName string) error {
	updateQuery := "UPDATE prover_job_queue SET updated_at = now(),updated_by = ? WHERE id = ?"
	_, err := db.Exec(updateQuery, proverName, jobId)
	if err != nil {
		logger().Errorf("Failed to recordProverIsWorking,err: %+v", err)
		return err
	}

	return nil
}

func storeProof(job ProverInputResponse, proof Groth16ProofResult) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateQuery := "UPDATE prover_job_queue SET updated_at=now(),job_status=?,updated_by='server_finish_job' WHERE id = ? AND job_type = ?"
	rows, err := db.Exec(updateQuery, Done, job.JobId, SingleProof)
	if err != nil {
		return err
	}

	updatedRows, err := rows.RowsAffected()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	if updatedRows != 1 {
		return fmt.Errorf("missing job for stored proof")
	}

	insertQuery := "INSERT INTO proofs (proof_id, computed_request_id, proof) VALUES (?, ?, ?)"
	_, err = db.Exec(insertQuery, proof.ProofId, proof.ComputedRequestId, proof.ProofBytes)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for storeProof")
	}

	return nil
}

func checkStarkProofExists(req *pb.FinalProofRequest) error {
	cleanInputDir := filepath.Clean(req.InputDir)

	commonCircuitFile := cleanInputDir + "/common_circuit_data.json"
	_, err := os.Stat(commonCircuitFile)
	if err != nil {
		return fmt.Errorf("file:%s not found", commonCircuitFile)
	}

	proofWithPublicInputsFile := cleanInputDir + "/proof_with_public_inputs.json"
	_, err = os.Stat(proofWithPublicInputsFile)
	if err != nil {
		return fmt.Errorf("file:%s not found", proofWithPublicInputsFile)
	}

	verifierOnlyCircuitDataFile := cleanInputDir + "/verifier_only_circuit_data.json"
	_, err = os.Stat(verifierOnlyCircuitDataFile)
	if err != nil {
		return fmt.Errorf("file:%s not found", verifierOnlyCircuitDataFile)
	}

	return nil
}

func addProverJobToQueue(ctx context.Context, req *pb.FinalProofRequest) *pb.Result {
	marshaller := jsonpb.Marshaler{}
	jobData, err := marshaller.MarshalToString(req)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}
	err = checkStarkProofExists(req)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}
	insertQuery := "INSERT INTO prover_job_queue (job_status, job_priority, job_type, updated_by, proof_id, computed_request_id, job_data) VALUES(?,?,?,?,?,?,?)"
	_, err = db.Exec(insertQuery, Idle, SingleProofJobPriority, SingleProof, "server_add_job", req.ProofId, req.ComputedRequestId, jobData)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}

	return getSuccessResult("get task successfully.")
}

func queryProverJobStatus(req *pb.GetTaskResultRequest) *pb.Result {
	query := "SELECT * FROM proofs WHERE proof_id = ? and computed_request_id = ?"

	rows, err := db.Query(query, req.ProofId, req.ComputedRequestId)
	if err != nil {
		formatStr := "Failed to query proofs db, err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_RESULT_ERROR, fmt.Sprintf(formatStr, err))
	}
	defer rows.Close()

	if rows.Next() { // if and only if one result
		if err != nil {
			formatStr := "Failed to query prover job status result, err: %+v"
			logger().Errorf(formatStr, err)
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
