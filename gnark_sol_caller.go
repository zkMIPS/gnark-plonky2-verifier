package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"text/template"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
)

// var GoerliId = big.NewInt(5)
// var GoerliNetwork = "https://eth-goerli.g.alchemy.com/v2/zKJf16XxhgdL6wMKT_NulOFfBfoT8YqE"
// var SepoliaId = big.NewInt(11155111)
// var SepoliaNetwork = "https://eth-sepolia.g.alchemy.com/v2/RH793ZL_pQkZb7KttcWcTlOjPrN0BjOW"

var ChainId *int64
var Network *string
var HexPrivaeKey *string // ("df4bc5647fdb9600ceb4943d4adff3749956a8512e5707716357b13d5ee687d9") // 0x21f59Cfb0d41FA2c0eeF0Fe1593F46f704C1Db50

func main() {
	ChainId = flag.Int64("chainId", 11155111, "chainId")
	Network = flag.String("network", "https://eth-sepolia.g.alchemy.com/v2/RH793ZL_pQkZb7KttcWcTlOjPrN0BjOW", "network")
	HexPrivaeKey = flag.String("privateKey", "df4bc5647fdb9600ceb4943d4adff3749956a8512e5707716357b13d5ee687d9", "privateKey")
	verifierAddr := flag.String("addr", "", "addr")
	if len(os.Args) < 2 {
		fmt.Println("expected 'deploy' or 'verify' or 'printvk' or 'all'  subcommands")
		os.Exit(1)
	}
	flag.CommandLine.Parse(os.Args[2:])
	switch os.Args[1] {
	case "deploy":
		deployVerifierContract()
	case "verify":
		callVerifierContract(*verifierAddr)
	case "printvk":
		PrintVk()
	case "all":
		deployAndCallVerifierContract()
	case "verifylocal":
		verifyLocal()
	case "generate":
		generateVerifySol()
	}
}

func PrintVk() {
	var circuitName = "mips"
	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk, _ := os.Open("testdata/" + circuitName + "/verifying.key")
	vk.ReadFrom(fVk)
	defer fVk.Close()
	groth16.PrintBn254Vk(vk)
}

func generateVerifySol() {
	tmpl, err := template.ParseFiles("verifier/verifier.sol.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	type VerifyingKeyConfig struct {
		Alpha     string
		Beta      string
		Gamma     string
		Delta     string
		Gamma_abc string
	}

	var config VerifyingKeyConfig
	var circuitName = "mips"
	var vkBN254 = groth16.NewVerifyingKey(ecc.BN254)

	fVk, _ := os.Open("testdata/" + circuitName + "/verifying.key")
	vkBN254.ReadFrom(fVk)
	defer fVk.Close()

	vk := vkBN254.(*groth16_bn254.VerifyingKey)

	config.Alpha = fmt.Sprint("Pairing.G1Point(uint256(", vk.G1.Alpha.X.String(), "), uint256(", vk.G1.Alpha.Y.String(), "))")
	config.Beta = fmt.Sprint("Pairing.G2Point([uint256(", vk.G2.Beta.X.A0.String(), "), uint256(", vk.G2.Beta.X.A1.String(), ")], [uint256(", vk.G2.Beta.Y.A0.String(), "), uint256(", vk.G2.Beta.Y.A1.String(), ")])")
	config.Gamma = fmt.Sprint("Pairing.G2Point([uint256(", vk.G2.Gamma.X.A0.String(), "), uint256(", vk.G2.Gamma.X.A1.String(), ")], [uint256(", vk.G2.Gamma.Y.A0.String(), "), uint256(", vk.G2.Gamma.Y.A1.String(), ")])")
	config.Delta = fmt.Sprint("Pairing.G2Point([uint256(", vk.G2.Delta.X.A0.String(), "), uint256(", vk.G2.Delta.X.A1.String(), ")], [uint256(", vk.G2.Delta.Y.A0.String(), "), uint256(", vk.G2.Delta.Y.A1.String(), ")])")
	config.Gamma_abc = fmt.Sprint("vk.gamma_abc = new Pairing.G1Point[](", len(vk.G1.K), ");\n")
	for k, v := range vk.G1.K {
		config.Gamma_abc += fmt.Sprint("        vk.gamma_abc[", k, "] = Pairing.G1Point(uint256(", v.X.String(), "), uint256(", v.Y.String(), "));\n")
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	fSol, _ := os.Create("testdata/" + circuitName + "/verifier.sol")
	_, err = fSol.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fSol.Close()
	fmt.Println("success")
}

func callVerifierContract(addr string) {

	flag.Parse()

	client, err := ethclient.Dial(*Network)
	if err != nil {
		log.Fatalf("Failed to create eth client: %v", err)
	}
	var circuitName = "mips"

	unlockedKey, err := crypto.HexToECDSA(*HexPrivaeKey)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(unlockedKey, big.NewInt(*ChainId))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.GasLimit = 1000000

	contractAddr := common.HexToAddress(addr)
	verifierContract, _ := verifier.NewContract(contractAddr, client)

	proof := groth16.NewProof(ecc.BN254)
	fProof, _ := os.Open("testdata/" + circuitName + "/proof.proof")
	proof.ReadFrom(fProof)
	defer fProof.Close()

	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk, _ := os.Open("testdata/" + circuitName + "/verifying.key")
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

	proofInputs[2], proofInputs[3] = proofInputs[3], proofInputs[2]
	proofInputs[4], proofInputs[5] = proofInputs[5], proofInputs[4]

	for i := 0; i < 8; i++ {
		fmt.Printf("proofInputs[%v]:%s\n", i, proofInputs[i].String())
	}
	publicInput, _ := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	proofWithPis := variables.DeserializeProofWithPublicInputs(publicInput)
	circuitData, _ := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(circuitData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	err, bPublicWitness, commitmentX, commitmentY := groth16.GetBn254Witness(proof, vk, publicWitness)

	fmt.Printf("bPublicWitness len:%+v\n", len(bPublicWitness))
	// fmt.Printf("bPublicWitness last:%+v\n", bPublicWitness)

	// convert public inputs
	nbInputs := len(bPublicWitness)

	var nbPublicInputs = vk.NbPublicWitness()

	if nbInputs != nbPublicInputs {
		log.Fatalf("nbInputs != nbPublicInputs,nbInputs:{%+v} nbPublicInputs:{%+v}", nbInputs, nbPublicInputs)
	}

	type ProofPublicData struct {
		Proof         groth16.Proof
		PublicWitness []string
	}

	proofPublicData := ProofPublicData{
		Proof:         proof,
		PublicWitness: make([]string, nbInputs),
	}

	var input [65]*big.Int
	for i := 0; i < nbInputs; i++ {
		input[i] = new(big.Int)
		bPublicWitness[i].BigInt(input[i])
		proofPublicData.PublicWitness[i] = input[i].String()
		fmt.Printf("input[%v]:%s\n", i, input[i].String())
	}

	jproofPublicData, _ := json.Marshal(proofPublicData)
	fmt.Printf("proofPublicData json: %s\n", string(jproofPublicData))

	var vp = verifier.VerifierProof{
		A: verifier.PairingG1Point{
			X: proofInputs[0],
			Y: proofInputs[1],
		},
		B: verifier.PairingG2Point{
			X: [2]*big.Int{proofInputs[2], proofInputs[3]},
			Y: [2]*big.Int{proofInputs[4], proofInputs[5]},
		},
		C: verifier.PairingG1Point{
			X: proofInputs[6],
			Y: proofInputs[7],
		},
	}

	var proofCommitment [2]*big.Int
	proofCommitment[0] = new(big.Int)
	commitmentX.BigInt(proofCommitment[0])
	proofCommitment[1] = new(big.Int)
	commitmentY.BigInt(proofCommitment[1])

	fmt.Printf("proofCommitmentX:%s,proofCommitmentY:%s\n", proofCommitment[0].String(), proofCommitment[1].String())

	tx, err := verifierContract.VerifyTx(auth, vp, input, proofCommitment)
	if err != nil {
		log.Fatalf("Failed to VerifyProof,err:[%+v]", err)
	}
	fmt.Printf("verify proof txHash: %+v\n", tx.Hash())
}

func deployVerifierContract() {
	client, err := ethclient.Dial(*Network)
	if err != nil {
		log.Fatalf("Failed to create eth client: %v", err)
	}

	unlockedKey, err := crypto.HexToECDSA(*HexPrivaeKey)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(unlockedKey, big.NewInt(*ChainId))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	// Deploy Verifier Contract
	verifierAddr, tx, _, err := verifier.DeployContract(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy Verifier,err:[%+v]", err)
	}
	fmt.Printf("verifierAddress: %+v\n txHash: %+v\n", verifierAddr.String(), tx.Hash())
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to deploy Verifier when mining :%v", err)
	}
}

func deployAndCallVerifierContract() {
	client, err := ethclient.Dial(*Network)
	if err != nil {
		log.Fatalf("Failed to create eth client: %v", err)
	}
	var circuitName = "mips"

	unlockedKey, err := crypto.HexToECDSA(*HexPrivaeKey)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(unlockedKey, big.NewInt(*ChainId))
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
	fProof, _ := os.Open("testdata/" + circuitName + "/proof.proof")
	proof.ReadFrom(fProof)
	defer fProof.Close()

	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk, _ := os.Open("testdata/" + circuitName + "/verifying.key")
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

	proofInputs[2], proofInputs[3] = proofInputs[3], proofInputs[2]
	proofInputs[4], proofInputs[5] = proofInputs[5], proofInputs[4]

	for i := 0; i < 8; i++ {
		fmt.Printf("proofInputs[%v]:%s\n", i, proofInputs[i].String())
	}

	publicInput, _ := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	proofWithPis := variables.DeserializeProofWithPublicInputs(publicInput)
	circuitData, _ := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(circuitData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	err, bPublicWitness, commitmentX, commitmentY := groth16.GetBn254Witness(proof, vk, publicWitness)

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

	var vp = verifier.VerifierProof{
		A: verifier.PairingG1Point{
			X: proofInputs[0],
			Y: proofInputs[1],
		},
		B: verifier.PairingG2Point{
			X: [2]*big.Int{proofInputs[2], proofInputs[3]},
			Y: [2]*big.Int{proofInputs[4], proofInputs[5]},
		},
		C: verifier.PairingG1Point{
			X: proofInputs[6],
			Y: proofInputs[7],
		},
	}

	var proofCommitment [2]*big.Int
	proofCommitment[0] = new(big.Int)
	commitmentX.BigInt(proofCommitment[0])
	proofCommitment[1] = new(big.Int)
	commitmentY.BigInt(proofCommitment[1])

	fmt.Printf("proofCommitmentX:%s,proofCommitmentY:%s\n", proofCommitment[0].String(), proofCommitment[1].String())

	tx, err = verifierContract.VerifyTx(auth, vp, input, proofCommitment)
	if err != nil {
		log.Fatalf("Failed to VerifyProof,err:[%+v]", err)
	}
	fmt.Printf("verify proof txHash: %+v\n", tx.Hash())
}

func verifyLocal() {
	var circuitName = "mips"
	publicInput, _ := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	proofWithPis := variables.DeserializeProofWithPublicInputs(publicInput)
	circuitData, _ := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(circuitData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	var vk = groth16.NewVerifyingKey(ecc.BN254)

	fVk, _ := os.Open("testdata/" + circuitName + "/verifying.key")
	vk.ReadFrom(fVk)
	defer fVk.Close()

	proof := groth16.NewProof(ecc.BN254)
	fProof, _ := os.Open("testdata/" + circuitName + "/proof.proof")
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
