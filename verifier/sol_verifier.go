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
	_ = abi.ConvertType
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
	Bin: "0x6080604052348015600e575f80fd5b506122e58061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063bc232a001461002d575b5f80fd5b61004760048036038101906100429190611ffb565b61005d565b6040516100549190612068565b60405180910390f35b5f8061006a8486856100ba565b036100af577fe78b29f09ce8ea10f38c2cd1016403d7a96bc8a14f6c1daa0da40c836bd2067e3360405161009e91906120c0565b60405180910390a1600190506100b3565b5f90505b9392505050565b5f807f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000190505f6100e861028f565b9050806080015151600160416100fe9190612106565b14610107575f80fd5b5f60405180604001604052805f81526020015f81525090505f5b60418110156101af578388826041811061013e5761013d612139565b5b60200201511061014c575f80fd5b6101a08261019b85608001516001856101659190612106565b8151811061017657610175612139565b5b60200260200101518b856041811061019157610190612139565b5b6020020151611491565b61152e565b91508080600101915050610121565b505f6040518060400160405280875f600281106101cf576101ce612139565b5b60200201518152602001876001600281106101ed576101ec612139565b5b6020020151815250905061021f8284608001515f8151811061021257610211612139565b5b602002602001015161152e565b915061022b828261152e565b915061026f875f01518860200151610242856115f0565b86604001516102548c604001516115f0565b88606001516102658a5f01516115f0565b8a6020015161168e565b610280576001945050505050610288565b5f9450505050505b9392505050565b610297611bc2565b60405180604001604052807f0ad0241a8532daca1a4e2c518bbfedd80b619aa43b147aac9ac549762f39974f81526020017f2697f5bb55787983681e3cac8f5ec390a7dc296a5b932255d1ddf73f634658e2815250815f0181905250604051806040016040528060405180604001604052807f2498349f28ec0e20cac6bf813af5939e03565077c266b76740c784a295b8d19f81526020017f15bada385f2a9c3d39840ab83305d413fe851f4c2f1d204b0b2440c892954640815250815260200160405180604001604052807f14c5ad67713431676a4d4adb680633cfa52b9ad0abc7fa922c59e84603e518ab81526020017f266a5ddd2b74c8170964b5e0e88592105d71fc5742637a7fd8690935f4c5f2058152508152508160200181905250604051806040016040528060405180604001604052807f08eb9cf73d2f8ad7c29965f7f57dbab073fc36564c0949dfef956a6aae2050ef81526020017f2c4f9a8cb56ac6a61b1a51bdb595e51b9561a6af6b330f978559d73abc3a090d815250815260200160405180604001604052807f27400506bb4e918510c3219db6934fc691538783c2cdc362a3f4fd0de95d563181526020017f1d37cf99c477a451e01762cbcd67ee97d1e1f5a0e5ee6908f707f24b6047955d8152508152508160400181905250604051806040016040528060405180604001604052807f0988cf33a9fe0a8a6d4f0d34a7ecddb996505f74ada15b8a71c5ed6833f8db0681526020017f2e38b6c4bfe7a4115df17fac89e31e02df7ed645252d512c499302e94eebe36f815250815260200160405180604001604052807f201b4d8800b376e0460f5e3bc06f198e4b467b6b8cc33d9de3421d15547ebffa81526020017f1fe396a865f30cf5d2a4a9292582dfbe564a8c4fceead809ee347d4e305c91848152508152508160600181905250604267ffffffffffffffff81111561055d5761055c611cf0565b5b60405190808252806020026020018201604052801561059657816020015b610583611c09565b81526020019060019003908161057b5790505b50816080018190525060405180604001604052807f0902c314d07eee8df5c07448f931c7084ecb3dfb787a3de0e1551d4f1397da3981526020017f160c48bcfaa43808bbe0dcf0ff1793a6281275eecf1909613980d7d4f4503e8581525081608001515f8151811061060b5761060a612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160018151811061064357610642612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160028151811061067b5761067a612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516003815181106106b3576106b2612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516004815181106106eb576106ea612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160058151811061072357610722612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160068151811061075b5761075a612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160078151811061079357610792612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516008815181106107cb576107ca612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160098151811061080357610802612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600a8151811061083b5761083a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600b8151811061087357610872612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600c815181106108ab576108aa612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600d815181106108e3576108e2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600e8151811061091b5761091a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151600f8151811061095357610952612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160108151811061098b5761098a612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516011815181106109c3576109c2612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516012815181106109fb576109fa612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601381518110610a3357610a32612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601481518110610a6b57610a6a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601581518110610aa357610aa2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601681518110610adb57610ada612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601781518110610b1357610b12612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601881518110610b4b57610b4a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601981518110610b8357610b82612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601a81518110610bbb57610bba612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601b81518110610bf357610bf2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601c81518110610c2b57610c2a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601d81518110610c6357610c62612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601e81518110610c9b57610c9a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151601f81518110610cd357610cd2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602081518110610d0b57610d0a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602181518110610d4357610d42612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602281518110610d7b57610d7a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602381518110610db357610db2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602481518110610deb57610dea612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602581518110610e2357610e22612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602681518110610e5b57610e5a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602781518110610e9357610e92612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602881518110610ecb57610eca612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602981518110610f0357610f02612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602a81518110610f3b57610f3a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602b81518110610f7357610f72612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602c81518110610fab57610faa612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602d81518110610fe357610fe2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602e8151811061101b5761101a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151602f8151811061105357611052612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160308151811061108b5761108a612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516031815181106110c3576110c2612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516032815181106110fb576110fa612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160338151811061113357611132612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160348151811061116b5761116a612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516035815181106111a3576111a2612139565b5b602002602001018190525060405180604001604052805f81526020015f81525081608001516036815181106111db576111da612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160378151811061121357611212612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160388151811061124b5761124a612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160398151811061128357611282612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603a815181106112bb576112ba612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603b815181106112f3576112f2612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603c8151811061132b5761132a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603d8151811061136357611362612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603e8151811061139b5761139a612139565b5b602002602001018190525060405180604001604052805f81526020015f8152508160800151603f815181106113d3576113d2612139565b5b602002602001018190525060405180604001604052805f81526020015f815250816080015160408151811061140b5761140a612139565b5b602002602001018190525060405180604001604052807f1105188506db2f91e15a080919cfce406ac7d3819e94fbed8c092202ce7dc5e981526020017f2905cc90f3e3c7852468dc168604f6c9430952e8ca4c44e3b472028a39e0a77c815250816080015160418151811061148357611482612139565b5b602002602001018190525090565b611499611c09565b6114a1611c21565b835f0151815f600381106114b8576114b7612139565b5b6020020181815250508360200151816001600381106114da576114d9612139565b5b60200201818152505082816002600381106114f8576114f7612139565b5b6020020181815250505f60608360808460076107d05a03fa9050805f810361151c57fe5b5080611526575f80fd5b505092915050565b611536611c09565b61153e611c43565b835f0151815f6004811061155557611554612139565b5b60200201818152505083602001518160016004811061157757611576612139565b5b602002018181525050825f01518160026004811061159857611597612139565b5b6020020181815250508260200151816003600481106115ba576115b9612139565b5b6020020181815250505f60608360c08460066107d05a03fa9050805f81036115de57fe5b50806115e8575f80fd5b505092915050565b6115f8611c09565b5f7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4790505f835f015114801561163157505f8360200151145b156116535760405180604001604052805f81526020015f815250915050611689565b6040518060400160405280845f015181526020018285602001516116779190612193565b8361168291906121c3565b8152509150505b919050565b5f80600467ffffffffffffffff8111156116ab576116aa611cf0565b5b6040519080825280602002602001820160405280156116e457816020015b6116d1611c09565b8152602001906001900390816116c95790505b5090505f600467ffffffffffffffff81111561170357611702611cf0565b5b60405190808252806020026020018201604052801561173c57816020015b611729611c65565b8152602001906001900390816117215790505b5090508a825f8151811061175357611752612139565b5b6020026020010181905250888260018151811061177357611772612139565b5b6020026020010181905250868260028151811061179357611792612139565b5b602002602001018190525084826003815181106117b3576117b2612139565b5b602002602001018190525089815f815181106117d2576117d1612139565b5b602002602001018190525087816001815181106117f2576117f1612139565b5b6020026020010181905250858160028151811061181257611811612139565b5b6020026020010181905250838160038151811061183257611831612139565b5b60200260200101819052506118478282611857565b9250505098975050505050505050565b5f8151835114611865575f80fd5b5f835190505f60068261187891906121f6565b90505f8167ffffffffffffffff81111561189557611894611cf0565b5b6040519080825280602002602001820160405280156118c35781602001602082028036833780820191505090505b5090505f5b83811015611b3b578681815181106118e3576118e2612139565b5b60200260200101515f0151825f6006846118fd91906121f6565b6119079190612106565b8151811061191857611917612139565b5b60200260200101818152505086818151811061193757611936612139565b5b60200260200101516020015182600160068461195391906121f6565b61195d9190612106565b8151811061196e5761196d612139565b5b60200260200101818152505085818151811061198d5761198c612139565b5b60200260200101515f01516001600281106119ab576119aa612139565b5b60200201518260026006846119c091906121f6565b6119ca9190612106565b815181106119db576119da612139565b5b6020026020010181815250508581815181106119fa576119f9612139565b5b60200260200101515f01515f60028110611a1757611a16612139565b5b6020020151826003600684611a2c91906121f6565b611a369190612106565b81518110611a4757611a46612139565b5b602002602001018181525050858181518110611a6657611a65612139565b5b602002602001015160200151600160028110611a8557611a84612139565b5b6020020151826004600684611a9a91906121f6565b611aa49190612106565b81518110611ab557611ab4612139565b5b602002602001018181525050858181518110611ad457611ad3612139565b5b6020026020010151602001515f60028110611af257611af1612139565b5b6020020151826005600684611b0791906121f6565b611b119190612106565b81518110611b2257611b21612139565b5b60200260200101818152505080806001019150506118c8565b50611b44611c8b565b5f602082602086026020860160086107d05a03fa905080611b9a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b9190612291565b60405180910390fd5b5f825f60018110611bae57611bad612139565b5b602002015114159550505050505092915050565b6040518060a00160405280611bd5611c09565b8152602001611be2611c65565b8152602001611bef611c65565b8152602001611bfc611c65565b8152602001606081525090565b60405180604001604052805f81526020015f81525090565b6040518060600160405280600390602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b6040518060400160405280611c78611cad565b8152602001611c85611cad565b81525090565b6040518060200160405280600190602082028036833780820191505090505090565b6040518060400160405280600290602082028036833780820191505090505090565b5f604051905090565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b611d2682611ce0565b810181811067ffffffffffffffff82111715611d4557611d44611cf0565b5b80604052505050565b5f611d57611ccf565b9050611d638282611d1d565b919050565b5f819050919050565b611d7a81611d68565b8114611d84575f80fd5b50565b5f81359050611d9581611d71565b92915050565b5f60408284031215611db057611daf611cdc565b5b611dba6040611d4e565b90505f611dc984828501611d87565b5f830152506020611ddc84828501611d87565b60208301525092915050565b5f80fd5b5f67ffffffffffffffff821115611e0657611e05611cf0565b5b602082029050919050565b5f80fd5b5f611e27611e2284611dec565b611d4e565b90508060208402830185811115611e4157611e40611e11565b5b835b81811015611e6a5780611e568882611d87565b845260208401935050602081019050611e43565b5050509392505050565b5f82601f830112611e8857611e87611de8565b5b6002611e95848285611e15565b91505092915050565b5f60808284031215611eb357611eb2611cdc565b5b611ebd6040611d4e565b90505f611ecc84828501611e74565b5f830152506040611edf84828501611e74565b60208301525092915050565b5f6101008284031215611f0157611f00611cdc565b5b611f0b6060611d4e565b90505f611f1a84828501611d9b565b5f830152506040611f2d84828501611e9e565b60208301525060c0611f4184828501611d9b565b60408301525092915050565b5f67ffffffffffffffff821115611f6757611f66611cf0565b5b602082029050919050565b5f611f84611f7f84611f4d565b611d4e565b90508060208402830185811115611f9e57611f9d611e11565b5b835b81811015611fc75780611fb38882611d87565b845260208401935050602081019050611fa0565b5050509392505050565b5f82601f830112611fe557611fe4611de8565b5b6041611ff2848285611f72565b91505092915050565b5f805f610960848603121561201357612012611cd8565b5b5f61202086828701611eeb565b93505061010061203286828701611fd1565b92505061092061204486828701611e74565b9150509250925092565b5f8115159050919050565b6120628161204e565b82525050565b5f60208201905061207b5f830184612059565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6120aa82612081565b9050919050565b6120ba816120a0565b82525050565b5f6020820190506120d35f8301846120b1565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61211082611d68565b915061211b83611d68565b9250828201905080821115612133576121326120d9565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61219d82611d68565b91506121a883611d68565b9250826121b8576121b7612166565b5b828206905092915050565b5f6121cd82611d68565b91506121d883611d68565b92508282039050818111156121f0576121ef6120d9565b5b92915050565b5f61220082611d68565b915061220b83611d68565b925082820261221981611d68565b915082820484148315176122305761222f6120d9565b5b5092915050565b5f82825260208201905092915050565b7f6e6f0000000000000000000000000000000000000000000000000000000000005f82015250565b5f61227b600283612237565b915061228682612247565b602082019050919050565b5f6020820190508181035f8301526122a88161226f565b905091905056fea26469706673582212204106e55329fc667a5af666cf534bd72d2f5a7c983a6086079d180af71107687564736f6c63430008190033",
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
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
