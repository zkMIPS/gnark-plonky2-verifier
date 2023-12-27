package types

import (
	"encoding/json"
	"io"
	"os"

	"github.com/succinctlabs/gnark-plonky2-verifier/plonk/gates"
)

type CommonCircuitDataRaw struct {
	Config struct {
		NumWires                uint64 `json:"num_wires"`
		NumRoutedWires          uint64 `json:"num_routed_wires"`
		NumConstants            uint64 `json:"num_constants"`
		UseBaseArithmeticGate   bool   `json:"use_base_arithmetic_gate"`
		SecurityBits            uint64 `json:"security_bits"`
		NumChallenges           uint64 `json:"num_challenges"`
		ZeroKnowledge           bool   `json:"zero_knowledge"`
		MaxQuotientDegreeFactor uint64 `json:"max_quotient_degree_factor"`
		FriConfig               struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"fri_config"`
	} `json:"config"`
	ArithFriParams struct {
		Config struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"config"`
		Hiding             bool     `json:"hiding"`
		DegreeBits         uint64   `json:"degree_bits"`
		ReductionArityBits []uint64 `json:"reduction_arity_bits"`
	} `json:"arith_fri_params"`
	CpuFriParams struct {
		Config struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"config"`
		Hiding             bool     `json:"hiding"`
		DegreeBits         uint64   `json:"degree_bits"`
		ReductionArityBits []uint64 `json:"reduction_arity_bits"`
	} `json:"cpu_fri_params"`
	LogicFriParams struct {
		Config struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"config"`
		Hiding             bool     `json:"hiding"`
		DegreeBits         uint64   `json:"degree_bits"`
		ReductionArityBits []uint64 `json:"reduction_arity_bits"`
	} `json:"logic_fri_params"`
	MemoryFriParams struct {
		Config struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"config"`
		Hiding             bool     `json:"hiding"`
		DegreeBits         uint64   `json:"degree_bits"`
		ReductionArityBits []uint64 `json:"reduction_arity_bits"`
	} `json:"memory_fri_params"`
	Gates         []string `json:"gates"`
	SelectorsInfo struct {
		SelectorIndices []uint64 `json:"selector_indices"`
		Groups          []struct {
			Start uint64 `json:"start"`
			End   uint64 `json:"end"`
		} `json:"groups"`
	} `json:"selectors_info"`
	QuotientDegreeFactor uint64   `json:"quotient_degree_factor"`
	NumGateConstraints   uint64   `json:"num_gate_constraints"`
	NumConstants         uint64   `json:"num_constants"`
	NumPublicInputs      uint64   `json:"num_public_inputs"`
	KIs                  []uint64 `json:"k_is"`
	NumPartialProducts   uint64   `json:"num_partial_products"`
}

func ReadCommonCircuitData(path string) CommonCircuitData {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()
	rawBytes, _ := io.ReadAll(jsonFile)

	var raw CommonCircuitDataRaw
	err = json.Unmarshal(rawBytes, &raw)
	if err != nil {
		panic(err)
	}

	var commonCircuitData CommonCircuitData
	commonCircuitData.Config.NumWires = raw.Config.NumWires
	commonCircuitData.Config.NumRoutedWires = raw.Config.NumRoutedWires
	commonCircuitData.Config.NumConstants = raw.Config.NumConstants
	commonCircuitData.Config.UseBaseArithmeticGate = raw.Config.UseBaseArithmeticGate
	commonCircuitData.Config.SecurityBits = raw.Config.SecurityBits
	commonCircuitData.Config.NumChallenges = raw.Config.NumChallenges
	commonCircuitData.Config.ZeroKnowledge = raw.Config.ZeroKnowledge
	commonCircuitData.Config.MaxQuotientDegreeFactor = raw.Config.MaxQuotientDegreeFactor

	commonCircuitData.Config.FriConfig.RateBits = raw.Config.FriConfig.RateBits
	commonCircuitData.Config.FriConfig.CapHeight = raw.Config.FriConfig.CapHeight
	commonCircuitData.Config.FriConfig.ProofOfWorkBits = raw.Config.FriConfig.ProofOfWorkBits
	commonCircuitData.Config.FriConfig.NumQueryRounds = raw.Config.FriConfig.NumQueryRounds

	commonCircuitData.ArithFriParams.DegreeBits = raw.ArithFriParams.DegreeBits
	commonCircuitData.ArithFriParams.Config.RateBits = raw.ArithFriParams.Config.RateBits
	commonCircuitData.ArithFriParams.Config.CapHeight = raw.ArithFriParams.Config.CapHeight
	commonCircuitData.ArithFriParams.Config.ProofOfWorkBits = raw.ArithFriParams.Config.ProofOfWorkBits
	commonCircuitData.ArithFriParams.Config.NumQueryRounds = raw.ArithFriParams.Config.NumQueryRounds
	commonCircuitData.ArithFriParams.ReductionArityBits = raw.ArithFriParams.ReductionArityBits

	commonCircuitData.CpuFriParams.DegreeBits = raw.CpuFriParams.DegreeBits
	commonCircuitData.CpuFriParams.Config.RateBits = raw.CpuFriParams.Config.RateBits
	commonCircuitData.CpuFriParams.Config.CapHeight = raw.CpuFriParams.Config.CapHeight
	commonCircuitData.CpuFriParams.Config.ProofOfWorkBits = raw.CpuFriParams.Config.ProofOfWorkBits
	commonCircuitData.CpuFriParams.Config.NumQueryRounds = raw.CpuFriParams.Config.NumQueryRounds
	commonCircuitData.CpuFriParams.ReductionArityBits = raw.CpuFriParams.ReductionArityBits

	commonCircuitData.LogicFriParams.DegreeBits = raw.LogicFriParams.DegreeBits
	commonCircuitData.LogicFriParams.Config.RateBits = raw.LogicFriParams.Config.RateBits
	commonCircuitData.LogicFriParams.Config.CapHeight = raw.LogicFriParams.Config.CapHeight
	commonCircuitData.LogicFriParams.Config.ProofOfWorkBits = raw.LogicFriParams.Config.ProofOfWorkBits
	commonCircuitData.LogicFriParams.Config.NumQueryRounds = raw.LogicFriParams.Config.NumQueryRounds
	commonCircuitData.LogicFriParams.ReductionArityBits = raw.LogicFriParams.ReductionArityBits

	commonCircuitData.MemoryFriParams.DegreeBits = raw.MemoryFriParams.DegreeBits
	commonCircuitData.MemoryFriParams.Config.RateBits = raw.MemoryFriParams.Config.RateBits
	commonCircuitData.MemoryFriParams.Config.CapHeight = raw.MemoryFriParams.Config.CapHeight
	commonCircuitData.MemoryFriParams.Config.ProofOfWorkBits = raw.MemoryFriParams.Config.ProofOfWorkBits
	commonCircuitData.MemoryFriParams.Config.NumQueryRounds = raw.MemoryFriParams.Config.NumQueryRounds
	commonCircuitData.MemoryFriParams.ReductionArityBits = raw.MemoryFriParams.ReductionArityBits

	commonCircuitData.GateIds = raw.Gates

	selectorGroupStart := []uint64{}
	selectorGroupEnd := []uint64{}
	for _, group := range raw.SelectorsInfo.Groups {
		selectorGroupStart = append(selectorGroupStart, group.Start)
		selectorGroupEnd = append(selectorGroupEnd, group.End)
	}

	commonCircuitData.SelectorsInfo = *gates.NewSelectorsInfo(
		raw.SelectorsInfo.SelectorIndices,
		selectorGroupStart,
		selectorGroupEnd,
	)

	commonCircuitData.QuotientDegreeFactor = raw.QuotientDegreeFactor
	commonCircuitData.NumGateConstraints = raw.NumGateConstraints
	commonCircuitData.NumConstants = raw.NumConstants
	commonCircuitData.NumPublicInputs = raw.NumPublicInputs
	commonCircuitData.KIs = raw.KIs
	commonCircuitData.NumPartialProducts = raw.NumPartialProducts

	return commonCircuitData
}
