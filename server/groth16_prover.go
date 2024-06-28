package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/profile"
	pb "github.com/succinctlabs/gnark-plonky2-verifier/proto/prover/v1"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
)

type ProverJobStatus int

const (
	Idle ProverJobStatus = iota
	InProgress
	Done
)

type ProverJobType string

const (
	SingleProof     ProverJobType = "SINGLE_PROOF"
	AggregatedProof ProverJobType = "AGGREGATED_PROOF"
)

const SingleProofJobPriority = 1
const AggregatedProofJobPriority = 0

type ProverInputResponse struct {
	JobId             int
	SnarkProofRequest *OriginalFinalProofRequest
}

type Groth16ProofResult struct {
	ProofId           string
	ComputedRequestId string
	ProofBytes        []byte
	Err               error
}

type OriginalFinalProofRequest struct {
	ChainId           uint64 `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Timestamp         uint64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ProofId           string `protobuf:"bytes,3,opt,name=proof_id,json=proofId,proto3" json:"proof_id,omitempty"`
	ComputedRequestId string `protobuf:"bytes,4,opt,name=computed_request_id,json=computedRequestId,proto3" json:"computed_request_id,omitempty"`
	// There are three files in the folder
	// common_circuit_data.json
	// verifier_only_circuit_data.json
	// proof_with_public_inputs.json
	InputDir string `protobuf:"bytes,5,opt,name=input_dir,json=inputDir,proto3" json:"input_dir,omitempty"`
	// The file path for storing the results
	OutputPath string `protobuf:"bytes,6,opt,name=output_path,json=outputPath,proto3" json:"output_path,omitempty"`
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

		start := time.Now()
		go record_metrics("snark::compute_proof", func() {
			computeProof(job, workerName, ch, heartBeat)
		})

		select {
		case result := <-ch:
			if result.Err != nil {
				logger().Errorf("Failed to compute groth16 proof,err: %+v", result.Err)
			} else {
				logger().Infof("computeProof cost time: %v ms", time.Now().Sub(start).Milliseconds())
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

		originalReq := OriginalFinalProofRequest{}
		if err = json.Unmarshal([]byte(jobData), &originalReq); err != nil {
			return resp, err
		}

		resp = ProverInputResponse{JobId: id, SnarkProofRequest: &originalReq}
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

	if r1csCircuit != nil && pk != nil && vk != nil { // has cache
		bytes, err := generateGroth16ProofWithCache(r1csCircuit, cleanInputDir, job.SnarkProofRequest.OutputPath)
		if err != nil {
			res.Err = err
			ch <- res
			return
		}

		res.ProofBytes = bytes
		ch <- res
		return
	}

	// without cache
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

	start := time.Now()
	logger().Infof("frontend.Compile: %v", start)
	r1csCircuit, err = frontend.Compile(ecc.BN254.ScalarField(), builder, &circuit)
	logger().Infof("frontend.Compile cost time: %v ms", time.Now().Sub(start).Milliseconds())
	if err != nil {
		res.Err = err
		ch <- res
		return
	}

	fR1CS, _ := os.Create(*cacheDir + "/circuit")
	r1csCircuit.WriteTo(fR1CS)
	fR1CS.Close()

	if *profileCircuit {
		p.Stop()
		p.Top()
		logger().Infof("r1cs.GetNbCoefficients(): %v", r1csCircuit.GetNbCoefficients())
		logger().Infof("r1cs.GetNbConstraints(): %v", r1csCircuit.GetNbConstraints())
		logger().Infof("r1cs.GetNbSecretVariables(): %v", r1csCircuit.GetNbSecretVariables())
		logger().Infof("r1cs.GetNbPublicVariables(): %v", r1csCircuit.GetNbPublicVariables())
		logger().Infof("r1cs.GetNbInternalVariables(): %v", r1csCircuit.GetNbInternalVariables())
	}

	bytes, err := generateGroth16Proof(r1csCircuit, cleanInputDir, job.SnarkProofRequest.OutputPath)
	if err != nil {
		res.Err = err
		ch <- res
		return
	}

	res.ProofBytes = bytes
	ch <- res
}

func getWitness(inputDir string) (verifier.ExampleVerifierCircuit, error) {
	var err error

	proofWithPisRawdata, err := types.ReadProofWithPublicInputs(inputDir + "/proof_with_public_inputs.json")
	if err != nil {
		return verifier.ExampleVerifierCircuit{}, err
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisRawdata)
	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData(inputDir + "/verifier_only_circuit_data.json")
	if err != nil {
		return verifier.ExampleVerifierCircuit{}, err
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}
	return assignment, nil
}

func generateProof(inputDir string, outputPath string, assignment verifier.ExampleVerifierCircuit, r1cs constraint.ConstraintSystem) ([]byte, error) {
	fPK, _ := os.Create(inputDir + "/proving.key")
	pk.WriteTo(fPK)
	fPK.Close()

	if vk != nil {
		fVK, _ := os.Create(inputDir + "/verifying.key")
		vk.WriteTo(fVK)
		fVK.Close()
	}

	start := time.Now()
	logger().Infof("Generating witness: %v", start)
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	logger().Infof("frontend.NewWitness cost time: %v ms", time.Now().Sub(start).Milliseconds())
	publicWitness, _ := witness.Public()

	fWitness, _ := os.Create(inputDir + "/witness")
	witness.WriteTo(fWitness)
	fWitness.Close()

	start = time.Now()
	logger().Infof("Creating proof: %v", start)
	proof, err := groth16.Prove(r1cs, pk, witness)
	logger().Infof("groth16.Prove cost time: %v ms", time.Now().Sub(start).Milliseconds())
	if err != nil {
		return nil, err
	}

	fProof, _ := os.Create(outputPath)
	proof.WriteTo(fProof)
	fProof.Close()

	if vk == nil {
		return nil, fmt.Errorf("vk is nil, means you're using dummy setup and we skip verification of proof")
	}

	start = time.Now()
	logger().Infof("Verifying proof: %v", start)
	err = groth16.Verify(proof, vk, publicWitness)
	logger().Infof("groth16.Verify cost time: %v ms", time.Now().Sub(start).Milliseconds())
	if err != nil {
		return nil, err
	}

	_, bPublicWitness, _, _ := groth16.GetBn254Witness(proof, vk, publicWitness)
	nbInputs := len(bPublicWitness)

	type ProofPublicData struct {
		Proof         groth16.Proof
		PublicWitness []string
	}
	proofPublicData := ProofPublicData{
		Proof:         proof,
		PublicWitness: make([]string, nbInputs),
	}
	for i := 0; i < nbInputs; i++ {
		input := new(big.Int)
		bPublicWitness[i].BigInt(input)
		proofPublicData.PublicWitness[i] = input.String()
	}
	proofBytes, _ := json.Marshal(proofPublicData)
	logger().Infof("proof.json %s", string(proofBytes))

	return proofBytes, nil
}

func generateGroth16ProofWithCache(r1cs constraint.ConstraintSystem, inputDir string, outputPath string) ([]byte, error) {
	assignment, err := getWitness(inputDir)
	if err != nil {
		return nil, err
	}

	return generateProof(inputDir, outputPath, assignment, r1cs)
}

func generateGroth16Proof(r1cs constraint.ConstraintSystem, inputDir string, outputPath string) ([]byte, error) {
	assignment, err := getWitness(inputDir)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	logger().Infof("Running circuit setup: %v", time.Now())
	logger().Infof("Using real setup")
	pk, vk, err = groth16.Setup(r1cs)
	logger().Infof("groth16.Setup cost time: %v ms", time.Now().Sub(start).Milliseconds())

	if err != nil {
		return nil, err
	}

	fPK, _ := os.Create(*cacheDir + "/proving.key")
	pk.WriteTo(fPK)
	fPK.Close()

	if vk != nil {
		fVK, _ := os.Create(*cacheDir + "/verifying.key")
		vk.WriteTo(fVK)
		fVK.Close()
	}

	return generateProof(inputDir, outputPath, assignment, r1cs)
}

func recordProverIsWorking(jobId int, proverName string) error {
	updateQuery := "UPDATE prover_job_queue SET updated_at = now(),updated_by = ? WHERE id = ? AND job_status != ?"
	_, err := db.Exec(updateQuery, proverName, jobId, Done)
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

func saveStarkProofIntoDisk(commonCircuitData []byte, verifierOnlyCircuitData []byte,
	proofWithPublicInputs []byte, starkProofDir string, outputDir string) error {

	err := os.MkdirAll(starkProofDir+"/aggregate", 0755)
	if err != nil {
		return err
	}

	err = saveJsonFile(starkProofDir+"/common_circuit_data.json", commonCircuitData)
	if err != nil {
		return err
	}

	err = saveJsonFile(starkProofDir+"/verifier_only_circuit_data.json", verifierOnlyCircuitData)
	if err != nil {
		return err
	}

	err = saveJsonFile(starkProofDir+"/proof_with_public_inputs.json", proofWithPublicInputs)
	if err != nil {
		return err
	}

	// create output path parent dir
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func saveJsonFile(filePath string, byteData []byte) error {
	var jsonData interface{}
	err := json.Unmarshal(byteData, &jsonData)
	if err != nil {
		return err
	}

	// jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	// if err != nil {
	// 	return err
	// }

	err = os.WriteFile(filePath, byteData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addProverJobToQueue(ctx context.Context, req *pb.FinalProofRequest) *pb.Result {
	starkProofDir := filepath.Clean(*inputParentDir) + "/" + req.ProofId
	outputDir := starkProofDir + "/final"
	originalReq := OriginalFinalProofRequest{
		ChainId:           req.ChainId,
		Timestamp:         req.Timestamp,
		ProofId:           req.ProofId,
		ComputedRequestId: req.ComputedRequestId,
		InputDir:          starkProofDir + "/aggregate",
		OutputPath:        outputDir + "/output",
	}
	jobData, err := json.Marshal(originalReq)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
	}
	err = saveStarkProofIntoDisk(
		req.CommonCircuitData,
		req.VerifierOnlyCircuitData,
		req.ProofWithPublicInputs,
		originalReq.InputDir,
		outputDir,
	)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
	}
	insertQuery := "INSERT INTO prover_job_queue (job_status, job_priority, job_type, updated_by, proof_id, computed_request_id, job_data) VALUES(?,?,?,?,?,?,?)"
	_, err = db.Exec(insertQuery, Idle, SingleProofJobPriority, SingleProof, "server_add_job", req.ProofId, req.ComputedRequestId, jobData)
	if err != nil {
		formatStr := "Failed to addProverJobToQueue,err: %+v"
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
	}

	return getSuccessResult("get task successfully.")
}

func queryProverJobStatus(req *pb.GetTaskResultRequest) *pb.Result {
	formatStr := "Failed to query proofs db, err: %+v"
	query := "SELECT proof FROM proofs WHERE proof_id = ? and computed_request_id = ?"

	rows, err := db.Query(query, req.ProofId, req.ComputedRequestId)
	if err != nil {
		logger().Errorf(formatStr, err)
		return getErrorResult(pb.ResultCode_INVALID_PARAMETER, fmt.Sprintf(formatStr, err))
	}
	defer rows.Close()

	if rows.Next() { // if and only if one result
		var proofResult string
		rows.Scan(&proofResult)
		return getSuccessResult(proofResult)

		// query = "SELECT job_data FROM prover_job_queue WHERE proof_id = ? and computed_request_id = ?"

		// rows, err = db.Query(query, req.ProofId, req.ComputedRequestId)
		// if err != nil {
		// 	logger().Errorf(formatStr, err)
		// 	return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
		// }
		// if rows.Next() {
		// 	var jobData string
		// 	err = rows.Scan(&jobData)
		// 	if err != nil {
		// 		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
		// 	}
		// 	proofReq := OriginalFinalProofRequest{}
		// 	if err = json.Unmarshal([]byte(jobData), &proofReq); err != nil {
		// 		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
		// 	}

		// 	proofBytes, err := os.ReadFile(proofReq.OutputPath)
		// 	if err != nil {
		// 		return getErrorResult(pb.ResultCode_INTERNAL_ERROR, fmt.Sprintf(formatStr, err))
		// 	}

		// 	return getSuccessResult(hex.EncodeToString(proofBytes))
		// }
		// return getErrorResult(pb.ResultCode_BUSY, fmt.Sprintf("proof hasn't been ready, err: %+v", err))
	}
	return getErrorResult(pb.ResultCode_BUSY, fmt.Sprintf("proof hasn't been ready, err: %+v", err))
}

func getSuccessResult(msg string) *pb.Result {
	return &pb.Result{Code: pb.ResultCode_OK, Message: msg}
}

func getErrorResult(code pb.ResultCode, errorMsg string) *pb.Result {
	return &pb.Result{Code: code, Message: errorMsg}
}
