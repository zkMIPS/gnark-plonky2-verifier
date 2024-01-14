// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PairingG1Point is an auto generated low-level Go binding around an user-defined struct.
type PairingG1Point struct {
	X *big.Int
	Y *big.Int
}

// PairingG2Point is an auto generated low-level Go binding around an user-defined struct.
type PairingG2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// VerifierProof is an auto generated low-level Go binding around an user-defined struct.
type VerifierProof struct {
	A PairingG1Point
	B PairingG2Point
	C PairingG1Point
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"Value\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"VerifyEvent\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structPairing.G1Point\",\"name\":\"a\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structPairing.G2Point\",\"name\":\"b\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structPairing.G1Point\",\"name\":\"c\",\"type\":\"tuple\"}],\"internalType\":\"structVerifier.Proof\",\"name\":\"proof\",\"type\":\"tuple\"},{\"internalType\":\"uint256[65]\",\"name\":\"input\",\"type\":\"uint256[65]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proof_commitment\",\"type\":\"uint256[2]\"}],\"name\":\"verifyTx\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"r\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506123388061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063bc232a001461002d575b5f80fd5b61004760048036038101906100429190612007565b61005d565b6040516100549190612074565b60405180910390f35b5f8061006a8486856100ba565b036100af577fe78b29f09ce8ea10f38c2cd1016403d7a96bc8a14f6c1daa0da40c836bd2067e3360405161009e91906120cc565b60405180910390a1600190506100b3565b5f90505b9392505050565b5f807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000190505f6100e8610295565b9050806080015151600160416100fe9190612112565b14610107575f80fd5b5f60405180604001604052805f81526020015f81525090505f5b60418110156101b5578388826041811061013e5761013d612145565b5b60200201511061014c575f80fd5b6101a08261019b85608001516001856101659190612112565b8151811061017657610175612145565b5b60200260200101518b856041811061019157610190612145565b5b6020020151611497565b611534565b915080806101ad90612172565b915050610121565b505f6040518060400160405280875f600281106101d5576101d4612145565b5b60200201518152602001876001600281106101f3576101f2612145565b5b602002015181525090506102258284608001515f8151811061021857610217612145565b5b6020026020010151611534565b91506102318282611534565b9150610275875f01518860200151610248856115f6565b866040015161025a8c604001516115f6565b886060015161026b8a5f01516115f6565b8a60200151611694565b61028657600194505050505061028e565b5f9450505050505b9392505050565b61029d611bce565b60405180604001604052807f02dc9c089bb0e34249b78da4103a5f437d1e9ce588d84927ee128ff2c7d6b0e981526020017f214914fbe215f043d62807f1451b90a5687b2e86d0734cdc766bffe396812838815250815f0181905250604051806040016040528060405180604001604052807f2e292b2f34e273052e2534f429f883a322b538df17464a2d442d1e2347cfeb0b81526020017f28ee14ca1da9dde73a52ec46577ed31e7baa70caa00a47cf50648f0aeebf3991815250815260200160405180604001604052807f094b8a7016f537dce1c265654e3934fc1e4de099693dbc454390775d82d3a46b81526020017f16f3ac6d4ec5c4c3574d8b502b50c2850ebd7e3e4d66a70fac8c8dff242f88a88152508152508160200181905250604051806040016040528060405180604001604052807f21fd60ade9c55c0469cc37b6972f3002bb6999a4df660bd81eb5803fd94336e081526020017f175c079896d2d96406c3a3b1751d74523865517a22b35c5231f69c4adc6d3ea2815250815260200160405180604001604052807f031eec713d6674c7f1316d941b78f949636d9204acf0ffe572c9f1ecb9b6aa0281526020017f1426618f61cb58516f757d74bdfdd2935ffec1bcdf7309c484257a422e81abf98152508152508160400181905250604051806040016040528060405180604001604052807f2d4a3b6f3c106d4c990743d2ba4030dd1c4bb1f6011b044e670e80e6e5c0b74681526020017f0fe1a0ebe6e0fc9c9915a12ba2b5ffd6d8e3f5779489fb4b9da35c6c56497c6d815250815260200160405180604001604052807f0daf583d10ed0322df6772ea2182eba74ad7650f9c93b50ae0afef88165b355e81526020017f13a17ae636012f33fd7d5f65e5db43371ad09d75e43af0fa68ac68b4ad2e45628152508152508160600181905250604267ffffffffffffffff81111561056357610562611cfc565b5b60405190808252806020026020018201604052801561059c57816020015b610589611c15565b8152602001906001900390816105815790505b50816080018190525060405180604001604052807f0693b413605cdb94c07381639b61ad44b354424c6742a56543b4b3ebc39f4a6681526020017f2d5dfa40f53e7c083fe370aae326521da184e379f50f08859f82b90fc5d9b15281525081608001515f8151811061061157610610612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160018151811061064957610648612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160028151811061068157610680612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516003815181106106b9576106b8612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516004815181106106f1576106f0612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160058151811061072957610728612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160068151811061076157610760612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160078151811061079957610798612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516008815181106107d1576107d0612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160098151811061080957610808612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600a8151811061084157610840612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600b8151811061087957610878612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600c815181106108b1576108b0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600d815181106108e9576108e8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600e8151811061092157610920612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600f8151811061095957610958612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160108151811061099157610990612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516011815181106109c9576109c8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601281518110610a0157610a00612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601381518110610a3957610a38612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601481518110610a7157610a70612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601581518110610aa957610aa8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601681518110610ae157610ae0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601781518110610b1957610b18612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601881518110610b5157610b50612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601981518110610b8957610b88612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601a81518110610bc157610bc0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601b81518110610bf957610bf8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601c81518110610c3157610c30612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601d81518110610c6957610c68612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601e81518110610ca157610ca0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601f81518110610cd957610cd8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602081518110610d1157610d10612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602181518110610d4957610d48612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602281518110610d8157610d80612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602381518110610db957610db8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602481518110610df157610df0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602581518110610e2957610e28612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602681518110610e6157610e60612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602781518110610e9957610e98612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602881518110610ed157610ed0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602981518110610f0957610f08612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602a81518110610f4157610f40612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602b81518110610f7957610f78612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602c81518110610fb157610fb0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602d81518110610fe957610fe8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602e8151811061102157611020612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602f8151811061105957611058612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160308151811061109157611090612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516031815181106110c9576110c8612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160328151811061110157611100612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160338151811061113957611138612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160348151811061117157611170612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516035815181106111a9576111a8612145565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516036815181106111e1576111e0612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160378151811061121957611218612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160388151811061125157611250612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160398151811061128957611288612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603a815181106112c1576112c0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603b815181106112f9576112f8612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603c8151811061133157611330612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603d8151811061136957611368612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603e815181106113a1576113a0612145565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603f815181106113d9576113d8612145565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160408151811061141157611410612145565b5b602002602001018190525060405180604001604052807f06b54b73ed7d5281f1b6d2e8990108f5ed669750f01dcfadb61a6520c974021681526020017f271ff51660b024be18683877fe329abc1420a22d6477606ca8930ff002855c91815250816080015160418151811061148957611488612145565b5b602002602001018190525090565b61149f611c15565b6114a7611c2d565b835f0151815f600381106114be576114bd612145565b5b6020020181815250508360200151816001600381106114e0576114df612145565b5b60200201818152505082816002600381106114fe576114fd612145565b5b6020020181815250505f60608360808460076107d05a03fa9050805f810361152257fe5b508061152c575f80fd5b505092915050565b61153c611c15565b611544611c4f565b835f0151815f6004811061155b5761155a612145565b5b60200201818152505083602001518160016004811061157d5761157c612145565b5b602002018181525050825f01518160026004811061159e5761159d612145565b5b6020020181815250508260200151816003600481106115c0576115bf612145565b5b6020020181815250505f60608360c08460066107d05a03fa9050805f81036115e457fe5b50806115ee575f80fd5b505092915050565b6115fe611c15565b5f7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4790505f835f015114801561163757505f8360200151145b156116595760405180604001604052805f81526020015f81525091505061168f565b6040518060400160405280845f0151815260200182856020015161167d91906121e6565b836116889190612216565b8152509150505b919050565b5f80600467ffffffffffffffff8111156116b1576116b0611cfc565b5b6040519080825280602002602001820160405280156116ea57816020015b6116d7611c15565b8152602001906001900390816116cf5790505b5090505f600467ffffffffffffffff81111561170957611708611cfc565b5b60405190808252806020026020018201604052801561174257816020015b61172f611c71565b8152602001906001900390816117275790505b5090508a825f8151811061175957611758612145565b5b6020026020010181905250888260018151811061177957611778612145565b5b6020026020010181905250868260028151811061179957611798612145565b5b602002602001018190525084826003815181106117b9576117b8612145565b5b602002602001018190525089815f815181106117d8576117d7612145565b5b602002602001018190525087816001815181106117f8576117f7612145565b5b6020026020010181905250858160028151811061181857611817612145565b5b6020026020010181905250838160038151811061183857611837612145565b5b602002602001018190525061184d828261185d565b9250505098975050505050505050565b5f815183511461186b575f80fd5b5f835190505f60068261187e9190612249565b90505f8167ffffffffffffffff81111561189b5761189a611cfc565b5b6040519080825280602002602001820160405280156118c95781602001602082028036833780820191505090505b5090505f5b83811015611b47578681815181106118e9576118e8612145565b5b60200260200101515f0151825f6006846119039190612249565b61190d9190612112565b8151811061191e5761191d612145565b5b60200260200101818152505086818151811061193d5761193c612145565b5b6020026020010151602001518260016006846119599190612249565b6119639190612112565b8151811061197457611973612145565b5b60200260200101818152505085818151811061199357611992612145565b5b60200260200101515f01516001600281106119b1576119b0612145565b5b60200201518260026006846119c69190612249565b6119d09190612112565b815181106119e1576119e0612145565b5b602002602001018181525050858181518110611a00576119ff612145565b5b60200260200101515f01515f60028110611a1d57611a1c612145565b5b6020020151826003600684611a329190612249565b611a3c9190612112565b81518110611a4d57611a4c612145565b5b602002602001018181525050858181518110611a6c57611a6b612145565b5b602002602001015160200151600160028110611a8b57611a8a612145565b5b6020020151826004600684611aa09190612249565b611aaa9190612112565b81518110611abb57611aba612145565b5b602002602001018181525050858181518110611ada57611ad9612145565b5b6020026020010151602001515f60028110611af857611af7612145565b5b6020020151826005600684611b0d9190612249565b611b179190612112565b81518110611b2857611b27612145565b5b6020026020010181815250508080611b3f90612172565b9150506118ce565b50611b50611c97565b5f602082602086026020860160086107d05a03fa905080611ba6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b9d906122e4565b60405180910390fd5b5f825f60018110611bba57611bb9612145565b5b602002015114159550505050505092915050565b6040518060a00160405280611be1611c15565b8152602001611bee611c71565b8152602001611bfb611c71565b8152602001611c08611c71565b8152602001606081525090565b60405180604001604052805f81526020015f81525090565b6040518060600160405280600390602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b6040518060400160405280611c84611cb9565b8152602001611c91611cb9565b81525090565b6040518060200160405280600190602082028036833780820191505090505090565b6040518060400160405280600290602082028036833780820191505090505090565b5f604051905090565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611d3282611cec565b810181811067ffffffffffffffff82111715611d5157611d50611cfc565b5b80604052505050565b5f611d63611cdb565b9050611d6f8282611d29565b919050565b5f819050919050565b611d8681611d74565b8114611d90575f80fd5b50565b5f81359050611da181611d7d565b92915050565b5f60408284031215611dbc57611dbb611ce8565b5b611dc66040611d5a565b90505f611dd584828501611d93565b5f830152506020611de884828501611d93565b60208301525092915050565b5f80fd5b5f67ffffffffffffffff821115611e1257611e11611cfc565b5b602082029050919050565b5f80fd5b5f611e33611e2e84611df8565b611d5a565b90508060208402830185811115611e4d57611e4c611e1d565b5b835b81811015611e765780611e628882611d93565b845260208401935050602081019050611e4f565b5050509392505050565b5f82601f830112611e9457611e93611df4565b5b6002611ea1848285611e21565b91505092915050565b5f60808284031215611ebf57611ebe611ce8565b5b611ec96040611d5a565b90505f611ed884828501611e80565b5f830152506040611eeb84828501611e80565b60208301525092915050565b5f6101008284031215611f0d57611f0c611ce8565b5b611f176060611d5a565b90505f611f2684828501611da7565b5f830152506040611f3984828501611eaa565b60208301525060c0611f4d84828501611da7565b60408301525092915050565b5f67ffffffffffffffff821115611f7357611f72611cfc565b5b602082029050919050565b5f611f90611f8b84611f59565b611d5a565b90508060208402830185811115611faa57611fa9611e1d565b5b835b81811015611fd35780611fbf8882611d93565b845260208401935050602081019050611fac565b5050509392505050565b5f82601f830112611ff157611ff0611df4565b5b6041611ffe848285611f7e565b91505092915050565b5f805f610960848603121561201f5761201e611ce4565b5b5f61202c86828701611ef7565b93505061010061203e86828701611fdd565b92505061092061205086828701611e80565b9150509250925092565b5f8115159050919050565b61206e8161205a565b82525050565b5f6020820190506120875f830184612065565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6120b68261208d565b9050919050565b6120c6816120ac565b82525050565b5f6020820190506120df5f8301846120bd565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61211c82611d74565b915061212783611d74565b925082820190508082111561213f5761213e6120e5565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f61217c82611d74565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036121ae576121ad6120e5565b5b600182019050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6121f082611d74565b91506121fb83611d74565b92508261220b5761220a6121b9565b5b828206905092915050565b5f61222082611d74565b915061222b83611d74565b9250828203905081811115612243576122426120e5565b5b92915050565b5f61225382611d74565b915061225e83611d74565b925082820261226c81611d74565b91508282048414831517612283576122826120e5565b5b5092915050565b5f82825260208201905092915050565b7f6e6f0000000000000000000000000000000000000000000000000000000000005f82015250565b5f6122ce60028361228a565b91506122d98261229a565b602082019050919050565b5f6020820190508181035f8301526122fb816122c2565b905091905056fea2646970667358221220ee4aa53c2404d38561311c34bab28bfdc49c3645f67ae15bffb662fc7acc60ce64736f6c63430008150033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// VerifyTx is a paid mutator transaction binding the contract method 0xbc232a00.
//
// Solidity: function verifyTx(((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)) proof, uint256[65] input, uint256[2] proof_commitment) returns(bool r)
func (_Contract *ContractTransactor) VerifyTx(opts *bind.TransactOpts, proof VerifierProof, input [65]*big.Int, proof_commitment [2]*big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "verifyTx", proof, input, proof_commitment)
}

// VerifyTx is a paid mutator transaction binding the contract method 0xbc232a00.
//
// Solidity: function verifyTx(((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)) proof, uint256[65] input, uint256[2] proof_commitment) returns(bool r)
func (_Contract *ContractSession) VerifyTx(proof VerifierProof, input [65]*big.Int, proof_commitment [2]*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.VerifyTx(&_Contract.TransactOpts, proof, input, proof_commitment)
}

// VerifyTx is a paid mutator transaction binding the contract method 0xbc232a00.
//
// Solidity: function verifyTx(((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)) proof, uint256[65] input, uint256[2] proof_commitment) returns(bool r)
func (_Contract *ContractTransactorSession) VerifyTx(proof VerifierProof, input [65]*big.Int, proof_commitment [2]*big.Int) (*types.Transaction, error) {
	return _Contract.Contract.VerifyTx(&_Contract.TransactOpts, proof, input, proof_commitment)
}

// ContractValueIterator is returned from FilterValue and is used to iterate over the raw logs and unpacked data for Value events raised by the Contract contract.
type ContractValueIterator struct {
	Event *ContractValue // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractValueIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractValue)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractValue)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractValueIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractValueIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractValue represents a Value event raised by the Contract contract.
type ContractValue struct {
	X   *big.Int
	Y   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterValue is a free log retrieval operation binding the contract event 0xd0df8930e73a69258b2c5f54b88f056f4a3594e30a976d7e9d02f45cb0c8d72f.
//
// Solidity: event Value(uint256 x, uint256 y)
func (_Contract *ContractFilterer) FilterValue(opts *bind.FilterOpts) (*ContractValueIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Value")
	if err != nil {
		return nil, err
	}
	return &ContractValueIterator{contract: _Contract.contract, event: "Value", logs: logs, sub: sub}, nil
}

// WatchValue is a free log subscription operation binding the contract event 0xd0df8930e73a69258b2c5f54b88f056f4a3594e30a976d7e9d02f45cb0c8d72f.
//
// Solidity: event Value(uint256 x, uint256 y)
func (_Contract *ContractFilterer) WatchValue(opts *bind.WatchOpts, sink chan<- *ContractValue) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Value")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractValue)
				if err := _Contract.contract.UnpackLog(event, "Value", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValue is a log parse operation binding the contract event 0xd0df8930e73a69258b2c5f54b88f056f4a3594e30a976d7e9d02f45cb0c8d72f.
//
// Solidity: event Value(uint256 x, uint256 y)
func (_Contract *ContractFilterer) ParseValue(log types.Log) (*ContractValue, error) {
	event := new(ContractValue)
	if err := _Contract.contract.UnpackLog(event, "Value", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractVerifyEventIterator is returned from FilterVerifyEvent and is used to iterate over the raw logs and unpacked data for VerifyEvent events raised by the Contract contract.
type ContractVerifyEventIterator struct {
	Event *ContractVerifyEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractVerifyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractVerifyEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractVerifyEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractVerifyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractVerifyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractVerifyEvent represents a VerifyEvent event raised by the Contract contract.
type ContractVerifyEvent struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterVerifyEvent is a free log retrieval operation binding the contract event 0xe78b29f09ce8ea10f38c2cd1016403d7a96bc8a14f6c1daa0da40c836bd2067e.
//
// Solidity: event VerifyEvent(address user)
func (_Contract *ContractFilterer) FilterVerifyEvent(opts *bind.FilterOpts) (*ContractVerifyEventIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "VerifyEvent")
	if err != nil {
		return nil, err
	}
	return &ContractVerifyEventIterator{contract: _Contract.contract, event: "VerifyEvent", logs: logs, sub: sub}, nil
}

// WatchVerifyEvent is a free log subscription operation binding the contract event 0xe78b29f09ce8ea10f38c2cd1016403d7a96bc8a14f6c1daa0da40c836bd2067e.
//
// Solidity: event VerifyEvent(address user)
func (_Contract *ContractFilterer) WatchVerifyEvent(opts *bind.WatchOpts, sink chan<- *ContractVerifyEvent) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "VerifyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractVerifyEvent)
				if err := _Contract.contract.UnpackLog(event, "VerifyEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVerifyEvent is a log parse operation binding the contract event 0xe78b29f09ce8ea10f38c2cd1016403d7a96bc8a14f6c1daa0da40c836bd2067e.
//
// Solidity: event VerifyEvent(address user)
func (_Contract *ContractFilterer) ParseVerifyEvent(log types.Log) (*ContractVerifyEvent, error) {
	event := new(ContractVerifyEvent)
	if err := _Contract.contract.UnpackLog(event, "VerifyEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
