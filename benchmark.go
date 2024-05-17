package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	//bn254r1cs "github.com/consensys/gnark/constraint/bn254"
	"math/big"
	"os"
	"time"

	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/profile"
	"github.com/consensys/gnark/test"
)

var r1cs_circuit constraint.ConstraintSystem

var pk groth16.ProvingKey
var vk groth16.VerifyingKey

func init_circuit_keys(plonky2Circuit string, circuitPath string, pkPath string, vkPath string) {
	if r1cs_circuit != nil {
		return
	}

	// 使用os.Stat()函数检查文件是否存在
	_, err := os.Stat(circuitPath)

	// 检查错误
	if os.IsNotExist(err) {
		commonCircuitData, err := types.ReadCommonCircuitData("testdata/" + plonky2Circuit + "/common_circuit_data.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		proofWithPisData, err := types.ReadProofWithPublicInputs("testdata/" + plonky2Circuit + "/proof_with_public_inputs.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

		verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData("testdata/" + plonky2Circuit + "/verifier_only_circuit_data.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)

		circuit := verifier.ExampleVerifierCircuit{
			Proof:                   proofWithPis.Proof,
			PublicInputs:            proofWithPis.PublicInputs,
			VerifierOnlyCircuitData: verifierOnlyCircuitData,
			CommonCircuitData:       commonCircuitData,
		}

		var builder frontend.NewBuilder = r1cs.NewBuilder
		r1cs_circuit, err = frontend.Compile(ecc.BN254.ScalarField(), builder, &circuit)
		fR1CS, _ := os.Create(circuitPath)
		r1cs_circuit.WriteTo(fR1CS)
		fR1CS.Close()
	} else {
		fCircuit, err := os.Open(circuitPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		r1cs_circuit = groth16.NewCS(ecc.BN254)
		r1cs_circuit.ReadFrom(fCircuit)
		fCircuit.Close()
	}

	_, err = os.Stat(pkPath)
	if os.IsNotExist(err) {
		pk, vk, err = groth16.Setup(r1cs_circuit)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
		pk = groth16.NewProvingKey(ecc.BN254)
		vk = groth16.NewVerifyingKey(ecc.BN254)
		fPk, err := os.Open(pkPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pk.ReadFrom(fPk)

		fVk, err := os.Open(vkPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		vk.ReadFrom(fVk)
		defer fVk.Close()
	}
}

func runBenchmark(plonky2Circuit string, proofSystem string, profileCircuit bool, dummy bool, saveArtifacts bool) {
	commonCircuitData, err := types.ReadCommonCircuitData("testdata/" + plonky2Circuit + "/common_circuit_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proofWithPisData, err := types.ReadProofWithPublicInputs("testdata/" + plonky2Circuit + "/proof_with_public_inputs.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData("testdata/" + plonky2Circuit + "/verifier_only_circuit_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)

	circuit := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
		CommonCircuitData:       commonCircuitData,
	}

	var p *profile.Profile
	if profileCircuit {
		p = profile.Start()
	}

	var builder frontend.NewBuilder
	if proofSystem == "plonk" {
		builder = scs.NewBuilder
	} else if proofSystem == "groth16" {
		builder = r1cs.NewBuilder
	} else {
		fmt.Println("Please provide a valid proof system to benchmark, we only support plonk and groth16")
		os.Exit(1)
	}

	start := time.Now()
	fmt.Printf("frontend.Compile: %v\n", start)
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), builder, &circuit)
	fmt.Printf("frontend.Compile cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println("error in building circuit", err)
		os.Exit(1)
	}

	if profileCircuit {
		p.Stop()
		p.Top()
		println("r1cs.GetNbCoefficients(): ", r1cs.GetNbCoefficients())
		println("r1cs.GetNbConstraints(): ", r1cs.GetNbConstraints())
		println("r1cs.GetNbSecretVariables(): ", r1cs.GetNbSecretVariables())
		println("r1cs.GetNbPublicVariables(): ", r1cs.GetNbPublicVariables())
		println("r1cs.GetNbInternalVariables(): ", r1cs.GetNbInternalVariables())
	}

	if proofSystem == "plonk" {
		plonkProof(r1cs, plonky2Circuit, dummy, saveArtifacts)
	} else if proofSystem == "groth16" {
		groth16Proof(r1cs, plonky2Circuit, dummy, saveArtifacts)
	} else {
		panic("Please provide a valid proof system to benchmark, we only support plonk and groth16")
	}
}

func plonkProof(r1cs constraint.ConstraintSystem, circuitName string, dummy bool, saveArtifacts bool) {
	var pk plonk.ProvingKey
	var vk plonk.VerifyingKey
	var err error

	proofWithPisData, err := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	// Don't serialize the circuit for now, since it takes up too much memory
	// if saveArtifacts {
	// 	fR1CS, _ := os.Create("circuit")
	// 	r1cs.WriteTo(fR1CS)
	// 	fR1CS.Close()
	// }

	fmt.Println("Running circuit setup", time.Now())
	if dummy {
		panic("dummy setup not supported for plonk")
	} else {
		fmt.Println("Using real setup")
		srs, err := test.NewKZGSRS(r1cs)
		if err != nil {
			panic(err)
		}
		pk, vk, err = plonk.Setup(r1cs, srs)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if saveArtifacts {
		fPK, _ := os.Create("proving.key")
		pk.WriteTo(fPK)
		fPK.Close()

		if vk != nil {
			fVK, _ := os.Create("verifying.key")
			vk.WriteTo(fVK)
			fVK.Close()
		}

		fSolidity, _ := os.Create("proof.sol")
		err = vk.ExportSolidity(fSolidity)
	}

	fmt.Println("Generating witness", time.Now())
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()
	if saveArtifacts {
		fWitness, _ := os.Create("witness")
		witness.WriteTo(fWitness)
		fWitness.Close()
	}

	fmt.Println("Creating proof", time.Now())
	proof, err := plonk.Prove(r1cs, pk, witness)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if true {
		fProof, _ := os.Create("testdata/" + circuitName + "proof.proof")
		proof.WriteTo(fProof)
		fProof.Close()
	}

	if vk == nil {
		fmt.Println("vk is nil, means you're using dummy setup and we skip verification of proof")
		return
	}

	fmt.Println("Verifying proof", time.Now())
	err = plonk.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	const fpSize = 4 * 8
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()
	fmt.Printf("proofBytes: %v\n", proofBytes)
}

func groth16ProofWithCache(r1cs constraint.ConstraintSystem, circuitName string, saveArtifacts bool) {
	proofWithPisData, err := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}

	start := time.Now()
	fmt.Println("Generating witness", start)
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	fmt.Printf("frontend.NewWitness cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	publicWitness, _ := witness.Public()
	if saveArtifacts {
		fWitness, _ := os.Create("testdata/" + circuitName + "/witness")
		witness.WriteTo(fWitness)
		fWitness.Close()
	}

	start = time.Now()
	fmt.Println("Creating proof", start)
	proof, err := groth16.Prove(r1cs, pk, witness)
	fmt.Printf("groth16.Prove cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if saveArtifacts {
		fProof, _ := os.Create("testdata/" + circuitName + "/proof.proof")
		proof.WriteTo(fProof)
		fProof.Close()
	}

	if vk == nil {
		fmt.Println("vk is nil, means you're using dummy setup and we skip verification of proof")
		return
	}

	start = time.Now()
	fmt.Println("Verifying proof", start)
	err = groth16.Verify(proof, vk, publicWitness)
	fmt.Printf("groth16.Verify cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	const fpSize = 4 * 8
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	if saveArtifacts {
		fContractProof, _ := os.Create("testdata/" + circuitName + "/snark_proof_with_public_inputs.json")
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
		proofData, _ := json.Marshal(proofPublicData)
		fContractProof.Write(proofData)
		fContractProof.Close()
	}

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

	println("a[0] is ", a[0].String())
	println("a[1] is ", a[1].String())

	println("b[0][0] is ", b[0][0].String())
	println("b[0][1] is ", b[0][1].String())
	println("b[1][0] is ", b[1][0].String())
	println("b[1][1] is ", b[1][1].String())

	println("c[0] is ", c[0].String())
	println("c[1] is ", c[1].String())
}

func groth16Proof(r1cs constraint.ConstraintSystem, circuitName string, dummy bool, saveArtifacts bool) {
	var pk groth16.ProvingKey
	var vk groth16.VerifyingKey
	var err error

	proofWithPisData, err := types.ReadProofWithPublicInputs("testdata/" + circuitName + "/proof_with_public_inputs.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proofWithPis := variables.DeserializeProofWithPublicInputs(proofWithPisData)

	verifierOnlyCircuitRawData, err := types.ReadVerifierOnlyCircuitData("testdata/" + circuitName + "/verifier_only_circuit_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verifierOnlyCircuitData := variables.DeserializeVerifierOnlyCircuitData(verifierOnlyCircuitRawData)
	assignment := verifier.ExampleVerifierCircuit{
		Proof:                   proofWithPis.Proof,
		PublicInputs:            proofWithPis.PublicInputs,
		VerifierOnlyCircuitData: verifierOnlyCircuitData,
	}
	// Don't serialize the circuit for now, since it takes up too much memory
	//if saveArtifacts {
	//	fR1CS, _ := os.Create("testdata/" + "circuit")
	//	r1cs.WriteTo(fR1CS)
	//	fR1CS.Close()
	//}

	start := time.Now()
	fmt.Println("Running circuit setup", start)
	if dummy {
		fmt.Println("Using dummy setup")
		pk, err = groth16.DummySetup(r1cs)
	} else {
		fmt.Println("Using real setup")
		pk, vk, err = groth16.Setup(r1cs)
	}
	fmt.Printf("groth16.Setup cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if saveArtifacts {
		fPK, _ := os.Create("testdata/" + circuitName + "/proving.key")
		pk.WriteTo(fPK)
		fPK.Close()

		if vk != nil {
			fVK, _ := os.Create("testdata/" + circuitName + "/verifying.key")
			vk.WriteTo(fVK)
			fVK.Close()
		}

		fSolidity, _ := os.Create("testdata/" + circuitName + "/proof.sol")
		err = vk.ExportSolidity(fSolidity)
	}

	start = time.Now()
	fmt.Println("Generating witness", start)
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	fmt.Printf("frontend.NewWitness cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	publicWitness, _ := witness.Public()
	if saveArtifacts {
		fWitness, _ := os.Create("testdata/" + circuitName + "/witness")
		witness.WriteTo(fWitness)
		fWitness.Close()
	}

	start = time.Now()
	fmt.Println("Creating proof", start)
	proof, err := groth16.Prove(r1cs, pk, witness)
	fmt.Printf("groth16.Prove cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if saveArtifacts {
		fProof, _ := os.Create("testdata/" + circuitName + "/proof.proof")
		proof.WriteTo(fProof)
		fProof.Close()
	}

	if vk == nil {
		fmt.Println("vk is nil, means you're using dummy setup and we skip verification of proof")
		return
	}

	start = time.Now()
	fmt.Println("Verifying proof", start)
	err = groth16.Verify(proof, vk, publicWitness)
	fmt.Printf("groth16.Verify cost time: %v ms\n", time.Now().Sub(start).Milliseconds())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	println("a[0] is ", a[0].String())
	println("a[1] is ", a[1].String())

	println("b[0][0] is ", b[0][0].String())
	println("b[0][1] is ", b[0][1].String())
	println("b[1][0] is ", b[1][0].String())
	println("b[1][1] is ", b[1][1].String())

	println("c[0] is ", c[0].String())
	println("c[1] is ", c[1].String())

}

func runWithPrecomputedCircuit(plonky2Circuit string, saveArtifacts bool) {
	init_circuit_keys(plonky2Circuit, "testdata/circuit",
		"testdata/"+plonky2Circuit+"/proving.key",
		"testdata/"+plonky2Circuit+"/verifying.key")
	for i := 0; i < 10; i++ {
		groth16ProofWithCache(r1cs_circuit, plonky2Circuit, saveArtifacts)
	}
}

func main() {
	plonky2Circuit := flag.String("plonky2-circuit", "mips", "plonky2 circuit to benchmark")
	proofSystem := flag.String("proof-system", "groth16", "proof system to benchmark")
	profileCircuit := flag.Bool("profile", true, "profile the circuit")
	dummySetup := flag.Bool("dummy", false, "use dummy setup")
	saveArtifacts := flag.Bool("save", true, "save circuit artifacts")

	flag.Parse()

	if plonky2Circuit == nil || *plonky2Circuit == "" {
		fmt.Println("Please provide a plonky2 circuit to benchmark")
		os.Exit(1)
	}

	if *proofSystem == "plonk" {
		*dummySetup = false
	}

	fmt.Printf("Running benchmark for %s circuit with proof system %s\n", *plonky2Circuit, *proofSystem)
	fmt.Printf("Profiling: %t, DummySetup: %t, SaveArtifacts: %t\n", *profileCircuit, *dummySetup, *saveArtifacts)

	runWithPrecomputedCircuit(*plonky2Circuit, *saveArtifacts)
	//runBenchmark(*plonky2Circuit, *proofSystem, *profileCircuit, *dummySetup, *saveArtifacts)
}
