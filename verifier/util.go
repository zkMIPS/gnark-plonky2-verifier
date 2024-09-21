package verifier

import (
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
	"github.com/consensys/gnark/std/math/uints"
	gl "github.com/succinctlabs/gnark-plonky2-verifier/goldilocks"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
	"math/big"
)

type ExampleVerifierCircuit struct {
	PublicInputsHash        frontend.Variable `gnark:",public"`
	PublicInputs            []gl.Variable
	Proof                   variables.Proof
	VerifierOnlyCircuitData variables.VerifierOnlyCircuitData

	// This is configuration for the circuit, it is a constant not a variable
	CommonCircuitData types.CommonCircuitData
}

func (c *ExampleVerifierCircuit) Define(api frontend.API) error {
	verifierChip := NewVerifierChip(api, c.CommonCircuitData)
	// Compute public inputs hash
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return fmt.Errorf("new uints api: %w", err)
	}
	publicInputs := make([]uints.U8, len(c.PublicInputs))
	for i, v := range c.PublicInputs {
		publicInputs[i] = uapi.ByteValueOf(v.Limb)
	}

	verifierChip.sha256Chip.Write(publicInputs)
	hash := verifierChip.sha256Chip.Sum()

	digest := make([]frontend.Variable, 32)
	msbBits := bits.ToBinary(api, hash[0].Val, bits.WithNbDigits(8))
	digest[0] = bits.FromBinary(api, msbBits[:5])
	for i := 1; i < len(hash); i++ {
		digest[i] = hash[i].Val
	}

	expectedPublicInputsHash := frontend.Variable(0)
	coeff := big.NewInt(1)
	for i := 31; i >= 0; i-- {
		expectedPublicInputsHash = api.Add(expectedPublicInputsHash, api.Mul(coeff, digest[i]))
		coeff.Lsh(coeff, 8)
	}
	api.AssertIsEqual(expectedPublicInputsHash, c.PublicInputsHash)

	verifierChip.Verify(c.Proof, c.PublicInputs, c.VerifierOnlyCircuitData)

	return nil
}
