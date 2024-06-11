package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"text/template"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	groth16_bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	verifierAddr := flag.String("addr", "0x012ef3e31BA2664163bD039535889aE7bE9E7E86", "addr")
	outputDir := flag.String("outputDir", "hardhat/contracts", "outputDir")
	proofPath := flag.String("proofPath", "./hardhat/test/snark_proof_with_public_inputs.json", "proofPath")
	if len(os.Args) < 2 {
		log.Printf("expected 'verify' or 'generate'  subcommands")
		os.Exit(1)
	}
	flag.CommandLine.Parse(os.Args[2:])
	switch os.Args[1] {
	case "verify":
		callSnarkVerifierContract(*verifierAddr, *proofPath)
	case "generate":
		generateVerifySol(*outputDir)
	}
}

func generateVerifySol(outputDir string) {
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
	fSol, _ := os.Create(filepath.Join(outputDir, "verifier.sol"))
	_, err = fSol.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fSol.Close()
	log.Println("success")
}

func callSnarkVerifierContract(addr string, proofPath string) {

	flag.Parse()

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
	auth.GasLimit = 1000000

	contractAddr := common.HexToAddress(addr)
	verifierContract, _ := verifier.NewContract(contractAddr, client)

	jsonFile, err := os.Open(proofPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	type ProofPublicData struct {
		Proof struct {
			Ar struct {
				X string
				Y string
			}
			Krs struct {
				X string
				Y string
			}
			Bs struct {
				X struct {
					A0 string
					A1 string
				}
				Y struct {
					A0 string
					A1 string
				}
			}
			Commitments []struct {
				X string
				Y string
			}
		}
		PublicWitness []string
	}
	proofPublicData := ProofPublicData{}
	err = json.Unmarshal(byteValue, &proofPublicData)
	if err != nil {
		log.Fatal(err)
	}

	var input [65]*big.Int
	for i := 0; i < len(proofPublicData.PublicWitness); i++ {
		input[i], _ = new(big.Int).SetString(proofPublicData.PublicWitness[i], 0)
	}

	var vp = verifier.VerifierProof{}
	vp.A.X, _ = new(big.Int).SetString(proofPublicData.Proof.Ar.X, 0)
	vp.A.Y, _ = new(big.Int).SetString(proofPublicData.Proof.Ar.Y, 0)

	vp.B.X[0], _ = new(big.Int).SetString(proofPublicData.Proof.Bs.X.A0, 0)
	vp.B.X[1], _ = new(big.Int).SetString(proofPublicData.Proof.Bs.X.A1, 0)
	vp.B.Y[0], _ = new(big.Int).SetString(proofPublicData.Proof.Bs.Y.A0, 0)
	vp.B.Y[1], _ = new(big.Int).SetString(proofPublicData.Proof.Bs.Y.A1, 0)

	vp.C.X, _ = new(big.Int).SetString(proofPublicData.Proof.Krs.X, 0)
	vp.C.Y, _ = new(big.Int).SetString(proofPublicData.Proof.Krs.Y, 0)

	var proofCommitment [2]*big.Int
	proofCommitment[0], _ = new(big.Int).SetString(proofPublicData.Proof.Commitments[0].X, 0)
	proofCommitment[1], _ = new(big.Int).SetString(proofPublicData.Proof.Commitments[0].Y, 0)

	tx, err := verifierContract.VerifyTx(auth, vp, input, proofCommitment)
	if err != nil {
		log.Fatalf("Failed to VerifyProof,err:[%+v]", err)
	}
	log.Printf("verify proof txHash: %+v\n", tx.Hash())
}
