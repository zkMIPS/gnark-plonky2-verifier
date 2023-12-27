package verifier

import (
	"github.com/consensys/gnark/frontend"
	"github.com/succinctlabs/gnark-plonky2-verifier/challenger"
	"github.com/succinctlabs/gnark-plonky2-verifier/fri"
	gl "github.com/succinctlabs/gnark-plonky2-verifier/goldilocks"
	"github.com/succinctlabs/gnark-plonky2-verifier/plonk"
	"github.com/succinctlabs/gnark-plonky2-verifier/poseidon"
	"github.com/succinctlabs/gnark-plonky2-verifier/types"
	"github.com/succinctlabs/gnark-plonky2-verifier/variables"
)

type VerifierChip struct {
	api               frontend.API             `gnark:"-"`
	glChip            *gl.Chip                 `gnark:"-"`
	poseidonGlChip    *poseidon.GoldilocksChip `gnark:"-"`
	poseidonBN254Chip *poseidon.BN254Chip      `gnark:"-"`
	plonkChip         *plonk.PlonkChip         `gnark:"-"`
	friChip           *fri.Chip                `gnark:"-"`
	commonData        types.CommonCircuitData  `gnark:"-"`
}

func NewVerifierChip(api frontend.API, commonCircuitData types.CommonCircuitData) *VerifierChip {
	glChip := gl.New(api)
	friChip := fri.NewChip(api, &commonCircuitData, &commonCircuitData.ArithFriParams, &commonCircuitData.CpuFriParams, &commonCircuitData.LogicFriParams, &commonCircuitData.MemoryFriParams)
	plonkChip := plonk.NewPlonkChip(api, commonCircuitData)
	poseidonGlChip := poseidon.NewGoldilocksChip(api)
	poseidonBN254Chip := poseidon.NewBN254Chip(api)
	return &VerifierChip{
		api:               api,
		glChip:            glChip,
		poseidonGlChip:    poseidonGlChip,
		poseidonBN254Chip: poseidonBN254Chip,
		plonkChip:         plonkChip,
		friChip:           friChip,
		commonData:        commonCircuitData,
	}
}

func (c *VerifierChip) GetPublicInputsHash(publicInputs []gl.Variable) poseidon.GoldilocksHashOut {
	return c.poseidonGlChip.HashNoPad(publicInputs)
}

func (c *VerifierChip) GetChallenges(
	proof variables.Proof,
	publicInputsHash poseidon.GoldilocksHashOut,
	verifierData variables.VerifierOnlyCircuitData,
) variables.ProofChallenges {
	config := c.commonData.Config
	numChallenges := config.NumChallenges
	challenger := challenger.NewChip(c.api)

	var circuitDigest = verifierData.CircuitDigest

	challenger.ObserveBN254Hash(circuitDigest)
	challenger.ObserveHash(publicInputsHash)
	challenger.ObserveCap(proof.WiresCap)
	plonkBetas := challenger.GetNChallenges(numChallenges)
	plonkGammas := challenger.GetNChallenges(numChallenges)

	challenger.ObserveCap(proof.PlonkZsPartialProductsCap)
	plonkAlphas := challenger.GetNChallenges(numChallenges)

	challenger.ObserveCap(proof.QuotientPolysCap)
	plonkZeta := challenger.GetExtensionChallenge()

	challenger.ObserveOpenings(c.friChip.ToOpenings(proof.Openings))

	return variables.ProofChallenges{
		PlonkBetas:  plonkBetas,
		PlonkGammas: plonkGammas,
		PlonkAlphas: plonkAlphas,
		PlonkZeta:   plonkZeta,
		FriChallenges: challenger.GetFriChallenges(
			proof.OpeningProof.CommitPhaseMerkleCaps,
			proof.OpeningProof.FinalPoly,
			proof.OpeningProof.PowWitness,
			config.FriConfig,
		),
	}
}

func (c *VerifierChip) rangeCheckProof(proof variables.Proof) {
	// Need to verify the plonky2 proof's openings, openings proof (other than the sibling elements), fri's final poly, pow witness.

	// Note that this is NOT range checking the public inputs (first 32 elements should be no more than 8 bits and the last 4 elements should be no more than 64 bits).  Since this is currently being inputted via the smart contract,
	// we will assume that caller is doing that check.

	// Range check the proof's openings.
	for _, constant := range proof.Openings.Constants {
		c.glChip.RangeCheckQE(constant)
	}

	for _, plonkSigma := range proof.Openings.PlonkSigmas {
		c.glChip.RangeCheckQE(plonkSigma)
	}

	for _, wire := range proof.Openings.Wires {
		c.glChip.RangeCheckQE(wire)
	}

	for _, plonkZ := range proof.Openings.PlonkZs {
		c.glChip.RangeCheckQE(plonkZ)
	}

	for _, plonkZNext := range proof.Openings.PlonkZsNext {
		c.glChip.RangeCheckQE(plonkZNext)
	}

	for _, partialProduct := range proof.Openings.PartialProducts {
		c.glChip.RangeCheckQE(partialProduct)
	}

	for _, quotientPoly := range proof.Openings.QuotientPolys {
		c.glChip.RangeCheckQE(quotientPoly)
	}

	// Range check the openings proof.
	for _, queryRound := range proof.OpeningProof.QueryRoundProofs {
		for _, initialTreesElement := range queryRound.InitialTreesProof.EvalsProofs[0].Elements {
			c.glChip.RangeCheck(initialTreesElement)
		}

		for _, queryStep := range queryRound.Steps {
			for _, eval := range queryStep.Evals {
				c.glChip.RangeCheckQE(eval)
			}
		}
	}

	// Range check the fri's final poly.
	for _, coeff := range proof.OpeningProof.FinalPoly.Coeffs {
		c.glChip.RangeCheckQE(coeff)
	}

	// Range check the pow witness.
	c.glChip.RangeCheck(proof.OpeningProof.PowWitness)
}

func (c *VerifierChip) Verify(
	arithProof variables.Proof,
	cpuProof variables.Proof,
	logicProof variables.Proof,
	memoryProof variables.Proof,
	publicInputs []gl.Variable,
	verifierData variables.VerifierOnlyCircuitData,
) {
	c.rangeCheckProof(arithProof)
	c.rangeCheckProof(cpuProof)
	c.rangeCheckProof(logicProof)
	c.rangeCheckProof(memoryProof)

	// Generate the parts of the witness that is for the plonky2 proof input
	//publicInputsHash := c.GetPublicInputsHash(publicInputs)
	// Generate the parts of the witness that is for the plonky2 proof input
	publicInputsHash := c.GetPublicInputsHash(publicInputs)

	//arith
	arithProofChallenges := c.GetChallenges(arithProof, publicInputsHash, verifierData)
	c.plonkChip.Verify(arithProofChallenges, arithProof.Openings, publicInputsHash, c.commonData.ArithFriParams.DegreeBits)
	//cpu
	cpuProofChallenges := c.GetChallenges(cpuProof, publicInputsHash, verifierData)
	c.plonkChip.Verify(cpuProofChallenges, cpuProof.Openings, publicInputsHash, c.commonData.CpuFriParams.DegreeBits)
	//logic
	logicProofChallenges := c.GetChallenges(logicProof, publicInputsHash, verifierData)
	c.plonkChip.Verify(logicProofChallenges, logicProof.Openings, publicInputsHash, c.commonData.LogicFriParams.DegreeBits)
	//memory
	memoryProofChallenges := c.GetChallenges(memoryProof, publicInputsHash, verifierData)
	c.plonkChip.Verify(memoryProofChallenges, memoryProof.Openings, publicInputsHash, c.commonData.MemoryFriParams.DegreeBits)

	//arith
	arithInitialMerkleCaps := []variables.FriMerkleCap{
		verifierData.ConstantSigmasCap,
		arithProof.WiresCap,
		arithProof.PlonkZsPartialProductsCap,
		arithProof.QuotientPolysCap,
	}
	c.friChip.ArithVerifyFriProof(
		c.friChip.GetInstance(arithProofChallenges.PlonkZeta, c.commonData.ArithFriParams.DegreeBits),
		c.friChip.ToOpenings(arithProof.Openings),
		&arithProofChallenges.FriChallenges,
		arithInitialMerkleCaps,
		&arithProof.OpeningProof,
	)
	//cpu
	cpuInitialMerkleCaps := []variables.FriMerkleCap{
		verifierData.ConstantSigmasCap,
		cpuProof.WiresCap,
		cpuProof.PlonkZsPartialProductsCap,
		cpuProof.QuotientPolysCap,
	}
	c.friChip.CpuVerifyFriProof(
		c.friChip.GetInstance(cpuProofChallenges.PlonkZeta, c.commonData.CpuFriParams.DegreeBits),
		c.friChip.ToOpenings(cpuProof.Openings),
		&cpuProofChallenges.FriChallenges,
		cpuInitialMerkleCaps,
		&cpuProof.OpeningProof,
	)
	//logic
	logicInitialMerkleCaps := []variables.FriMerkleCap{
		verifierData.ConstantSigmasCap,
		logicProof.WiresCap,
		logicProof.PlonkZsPartialProductsCap,
		logicProof.QuotientPolysCap,
	}
	c.friChip.LogicVerifyFriProof(
		c.friChip.GetInstance(logicProofChallenges.PlonkZeta, c.commonData.LogicFriParams.DegreeBits),
		c.friChip.ToOpenings(logicProof.Openings),
		&logicProofChallenges.FriChallenges,
		logicInitialMerkleCaps,
		&logicProof.OpeningProof,
	)
	//memory
	memoryInitialMerkleCaps := []variables.FriMerkleCap{
		verifierData.ConstantSigmasCap,
		memoryProof.WiresCap,
		memoryProof.PlonkZsPartialProductsCap,
		memoryProof.QuotientPolysCap,
	}
	c.friChip.MemoryVerifyFriProof(
		c.friChip.GetInstance(memoryProofChallenges.PlonkZeta, c.commonData.MemoryFriParams.DegreeBits),
		c.friChip.ToOpenings(memoryProof.Openings),
		&memoryProofChallenges.FriChallenges,
		memoryInitialMerkleCaps,
		&memoryProof.OpeningProof,
	)
}
