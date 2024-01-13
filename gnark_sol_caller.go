package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
	"log"
	"math/big"
	"os"
)


var GoerliId = big.NewInt(5)

func main() {
	deployAndCallVerifierContract()
	//verifyLocal()
}

func deployAndCallVerifierContract() {
	var network = ""
	client, err := ethclient.Dial(network)
	if err != nil {
		log.Fatalf("Failed to create eth client: %v", err)
	}
	var circuitName = "mips"

	unlockedKey, err := crypto.HexToECDSA("") //
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(unlockedKey, GoerliId)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	// Deploy Verifier Contract
	verifierAddr, tx, verifierContract, err := verifier.DeployContract(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy Verifier,err:[%+v]", err)
	}
	fmt.Printf("verifierAddress: %+v\n txHash: %+v\n", verifierAddr.String(), tx.Hash())
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to deploy Verifier when mining :%v", err)
	}

	proof := groth16.NewProof(ecc.BN254)
	fProof, _ := os.Open("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/"+ circuitName + "/proof.proof")
	proof.ReadFrom(fProof)
	defer fProof.Close()

	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk,_ := os.Open("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/"+ circuitName + "/verifying.key")
	vk.ReadFrom(fVk)
	defer fVk.Close()

	const fpSize = 4 * 8
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	// solidity contract inputs
	var proofInputs [8]*big.Int

	// proof.Ar, proof.Bs, proof.Krs
	for i := 0; i < 8; i++ {
		proofInputs[i] = new(big.Int).SetBytes(proofBytes[fpSize*i : fpSize*(i+1)])
	}

	proofWithPis := variables.DeserializeProofWithPublicInputs(types.ReadProofWithPublicInputs("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/" + circuitName + "/proof_with_public_inputs.json"))
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(types.ReadVerifierOnlyCircuitData("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/" + circuitName + "/verifier_only_circuit_data.json"))
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	err,bPublicWitness := groth16.GetBn254Witness(proof, vk, publicWitness)

	fmt.Printf("bPublicWitness len:%+v\n", len(bPublicWitness))

	// convert public inputs
	nbInputs := len(bPublicWitness)

	var nbPublicInputs = vk.NbPublicWitness()

	if nbInputs != nbPublicInputs {
		log.Fatalf("nbInputs != nbPublicInputs,nbInputs:{%+v} nbPublicInputs:{%+v}", nbInputs, nbPublicInputs)
	}

	var input [65]*big.Int
	for i := 0; i < nbInputs; i++ {
		input[i] = new(big.Int)
		bPublicWitness[i].BigInt(input[i])
		fmt.Printf("input[%v]:%s\n", i, input[i].String())
	}

	tx,err = verifierContract.VerifyProof(auth, proofInputs, input)
	if err != nil {
		log.Fatalf("Failed to VerifyProof,err:[%+v]", err)
	}
	fmt.Printf("verify proof txHash: %+v\n", tx.Hash())
}

func verifyLocal() {
	var circuitName = "mips"
	proofWithPis := variables.DeserializeProofWithPublicInputs(types.ReadProofWithPublicInputs("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/" + circuitName + "/proof_with_public_inputs.json"))
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(types.ReadVerifierOnlyCircuitData("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/" + circuitName + "/verifier_only_circuit_data.json"))
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk,_ := os.Open("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/"+ circuitName + "/verifying.key")
	vk.ReadFrom(fVk)
	defer fVk.Close()

	proof := groth16.NewProof(ecc.BN254)
	fProof, _ := os.Open("/Users/bj89200ml/Documents/golang_workspace/src/github.com/succinctlabs/gnark-plonky2-verifier/testdata/"+ circuitName + "/proof.proof")
	proof.ReadFrom(fProof)
	defer fProof.Close()

	fmt.Println("begin verify")
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("success")
}
