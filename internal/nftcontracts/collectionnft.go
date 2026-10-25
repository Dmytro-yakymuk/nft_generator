// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nftcontracts

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

// CollectionNFTMetaData contains all meta data concerning the CollectionNFT contract.
var CollectionNFTMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot_\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"maxSupply_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPerMint_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_PER_MINT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRICE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROVENANCE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numberOfTokens\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_merkleProof\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mintWhitelist\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"presaleIsActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"saleIsActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI_\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"merkleRoot_\",\"type\":\"bytes32\"}],\"name\":\"setMerkleRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"newPresaleState\",\"type\":\"bool\"}],\"name\":\"setPresaleState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"provenance\",\"type\":\"string\"}],\"name\":\"setProvenance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"newState\",\"type\":\"bool\"}],\"name\":\"setSaleState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"whitelistClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040526008805461ffff191660011790555f600b5534801562000022575f80fd5b506040516200260b3803806200260b83398101604081905262000045916200023a565b60408051808201909152600381526213919560ea1b602082015284905f6200006e83826200034d565b5060016200007d82826200034d565b5050506200009a620000946200012460201b60201c565b62000128565b5f81118015620000aa5750828111155b620000fb5760405162461bcd60e51b815260206004820152601d60248201527f496e76616c6964206c696d697420706572207472616e73616374696f6e000000604482015260640160405180910390fd5b60096200010987826200034d565b50600a9490945560809190915260a0525060c0525062000419565b3390565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f8301126200019d575f80fd5b81516001600160401b0380821115620001ba57620001ba62000179565b604051601f8301601f19908116603f01168101908282118183101715620001e557620001e562000179565b816040528381526020925086602085880101111562000202575f80fd5b5f91505b8382101562000225578582018301518183018401529082019062000206565b5f602085830101528094505050505092915050565b5f805f805f8060c0878903121562000250575f80fd5b86516001600160401b038082111562000267575f80fd5b620002758a838b016200018d565b975060208901519650604089015191508082111562000292575f80fd5b50620002a189828a016200018d565b945050606087015192506080870151915060a087015190509295509295509295565b600181811c90821680620002d857607f821691505b602082108103620002f757634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200034857805f5260205f20601f840160051c81016020851015620003245750805b601f840160051c820191505b8181101562000345575f815560010162000330565b50505b505050565b81516001600160401b0381111562000369576200036962000179565b62000381816200037a8454620002c3565b84620002fd565b602080601f831160018114620003b7575f84156200039f5750858301515b5f19600386901b1c1916600185901b17855562000411565b5f85815260208120601f198616915b82811015620003e757888601518255948401946001909101908401620003c6565b50858210156200040557878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b60805160a05160c0516121b2620004595f395f81816102b3015261138101525f8181610474015261147f01525f818161036a01526113fd01526121b25ff3fe6080604052600436106101f1575f3560e01c8063715018a611610108578063b88d4fde1161009d578063db4bec441161006d578063db4bec4414610588578063e985e9c5146105b6578063eb8d2444146105d5578063f2fde38b146105f3578063ffe630b514610612575f80fd5b8063b88d4fde1461050c578063bce4d6ae1461052b578063c4e370951461054a578063c87b56dd14610569575f80fd5b806395d89b41116100d857806395d89b41146104b3578063a0712d68146104c7578063a22cb465146104da578063a6d612f9146104f9575f80fd5b8063715018a6146104305780637cb64759146104445780638d859f3e146104635780638da5cb5b14610496575f80fd5b80632eb4a7ab1161018957806342842e0e1161015957806342842e0e146103a057806355f804b3146103bf5780636352211e146103de5780636373a6b1146103fd57806370a0823114610411575f80fd5b80632eb4a7ab1461032b57806330f72cd41461034057806332cb6b0c146103595780633ccfd60b1461038c575f80fd5b806309d42b30116101c457806309d42b30146102a257806318160ddd146102e357806323b872dd146102f757806326987b6014610316575f80fd5b806301ffc9a7146101f557806306fdde0314610229578063081812fc1461024a578063095ea7b314610281575b5f80fd5b348015610200575f80fd5b5061021461020f366004611a83565b610631565b60405190151581526020015b60405180910390f35b348015610234575f80fd5b5061023d610682565b6040516102209190611aeb565b348015610255575f80fd5b50610269610264366004611afd565b610711565b6040516001600160a01b039091168152602001610220565b34801561028c575f80fd5b506102a061029b366004611b2f565b6107a9565b005b3480156102ad575f80fd5b506102d57f000000000000000000000000000000000000000000000000000000000000000081565b604051908152602001610220565b3480156102ee575f80fd5b50600b546102d5565b348015610302575f80fd5b506102a0610311366004611b57565b6108bd565b348015610321575f80fd5b506102d5600b5481565b348015610336575f80fd5b506102d5600a5481565b34801561034b575f80fd5b506008546102149060ff1681565b348015610364575f80fd5b506102d57f000000000000000000000000000000000000000000000000000000000000000081565b348015610397575f80fd5b506102a06108ee565b3480156103ab575f80fd5b506102a06103ba366004611b57565b61093c565b3480156103ca575f80fd5b506102a06103d9366004611c17565b610956565b3480156103e9575f80fd5b506102696103f8366004611afd565b610990565b348015610408575f80fd5b5061023d610a06565b34801561041c575f80fd5b506102d561042b366004611c5c565b610a92565b34801561043b575f80fd5b506102a0610b17565b34801561044f575f80fd5b506102a061045e366004611afd565b610b4a565b34801561046e575f80fd5b506102d57f000000000000000000000000000000000000000000000000000000000000000081565b3480156104a1575f80fd5b506006546001600160a01b0316610269565b3480156104be575f80fd5b5061023d610b79565b6102a06104d5366004611afd565b610b88565b3480156104e5575f80fd5b506102a06104f4366004611c84565b610c0a565b6102a0610507366004611cb5565b610c15565b348015610517575f80fd5b506102a0610526366004611d29565b610dc1565b348015610536575f80fd5b506102a0610545366004611da0565b610df3565b348015610555575f80fd5b506102a0610564366004611da0565b610e30565b348015610574575f80fd5b5061023d610583366004611afd565b610e74565b348015610593575f80fd5b506102146105a2366004611c5c565b600c6020525f908152604090205460ff1681565b3480156105c1575f80fd5b506102146105d0366004611db9565b610f4b565b3480156105e0575f80fd5b5060085461021490610100900460ff1681565b3480156105fe575f80fd5b506102a061060d366004611c5c565b610f78565b34801561061d575f80fd5b506102a061062c366004611c17565b611010565b5f6001600160e01b031982166380ac58cd60e01b148061066157506001600160e01b03198216635b5e139f60e01b145b8061067c57506301ffc9a760e01b6001600160e01b03198316145b92915050565b60605f805461069090611de1565b80601f01602080910402602001604051908101604052809291908181526020018280546106bc90611de1565b80156107075780601f106106de57610100808354040283529160200191610707565b820191905f5260205f20905b8154815290600101906020018083116106ea57829003601f168201915b5050505050905090565b5f818152600260205260408120546001600160a01b031661078e5760405162461bcd60e51b815260206004820152602c60248201527f4552433732313a20617070726f76656420717565727920666f72206e6f6e657860448201526b34b9ba32b73a103a37b5b2b760a11b60648201526084015b60405180910390fd5b505f908152600460205260409020546001600160a01b031690565b5f6107b382610990565b9050806001600160a01b0316836001600160a01b0316036108205760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b6064820152608401610785565b336001600160a01b038216148061083c575061083c8133610f4b565b6108ae5760405162461bcd60e51b815260206004820152603860248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f74206f7760448201527f6e6572206e6f7220617070726f76656420666f7220616c6c00000000000000006064820152608401610785565b6108b88383611046565b505050565b6108c733826110b3565b6108e35760405162461bcd60e51b815260040161078590611e19565b6108b8838383611188565b6006546001600160a01b031633146109185760405162461bcd60e51b815260040161078590611e6a565b60405133904780156108fc02915f818181858888f1935050505061093a575f80fd5b565b6108b883838360405180602001604052805f815250610dc1565b6006546001600160a01b031633146109805760405162461bcd60e51b815260040161078590611e6a565b600961098c8282611eea565b5050565b5f818152600260205260408120546001600160a01b03168061067c5760405162461bcd60e51b815260206004820152602960248201527f4552433732313a206f776e657220717565727920666f72206e6f6e657869737460448201526832b73a103a37b5b2b760b91b6064820152608401610785565b60078054610a1390611de1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a3f90611de1565b8015610a8a5780601f10610a6157610100808354040283529160200191610a8a565b820191905f5260205f20905b815481529060010190602001808311610a6d57829003601f168201915b505050505081565b5f6001600160a01b038216610afc5760405162461bcd60e51b815260206004820152602a60248201527f4552433732313a2062616c616e636520717565727920666f7220746865207a65604482015269726f206164647265737360b01b6064820152608401610785565b506001600160a01b03165f9081526003602052604090205490565b6006546001600160a01b03163314610b415760405162461bcd60e51b815260040161078590611e6a565b61093a5f611324565b6006546001600160a01b03163314610b745760405162461bcd60e51b815260040161078590611e6a565b600a55565b60606001805461069090611de1565b600854610100900460ff16610bdf5760405162461bcd60e51b815260206004820152601b60248201527f53616c65206d7573742062652061637469766520746f206d696e7400000000006044820152606401610785565b333214610bfe5760405162461bcd60e51b815260040161078590611faa565b610c0781611375565b50565b61098c338383611526565b60085460ff16610c675760405162461bcd60e51b815260206004820152601e60248201527f50726573616c65206d7573742062652061637469766520746f206d696e7400006044820152606401610785565b335f908152600c602052604090205460ff1615610cc65760405162461bcd60e51b815260206004820152601b60248201527f416464726573732068617320616c726561647920636c61696d656400000000006044820152606401610785565b333214610ce55760405162461bcd60e51b815260040161078590611faa565b6040516bffffffffffffffffffffffff193360601b1660208201525f90603401604051602081830303815290604052805190602001209050610d5d8484808060200260200160405190810160405280939291908181526020018383602002808284375f9201919091525050600a5491508490506115f3565b610d995760405162461bcd60e51b815260206004820152600d60248201526c24b73b30b634b210283937b7b360991b6044820152606401610785565b335f908152600c60205260409020805460ff19166001179055610dbb82611375565b50505050565b610dcb33836110b3565b610de75760405162461bcd60e51b815260040161078590611e19565b610dbb84848484611608565b6006546001600160a01b03163314610e1d5760405162461bcd60e51b815260040161078590611e6a565b6008805460ff1916911515919091179055565b6006546001600160a01b03163314610e5a5760405162461bcd60e51b815260040161078590611e6a565b600880549115156101000261ff0019909216919091179055565b5f818152600260205260409020546060906001600160a01b0316610ef25760405162461bcd60e51b815260206004820152602f60248201527f4552433732314d657461646174613a2055524920717565727920666f72206e6f60448201526e3732bc34b9ba32b73a103a37b5b2b760891b6064820152608401610785565b5f610efb61163b565b90505f815111610f195760405180602001604052805f815250610f44565b80610f238461164a565b604051602001610f34929190611fee565b6040516020818303038152906040525b9392505050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b6006546001600160a01b03163314610fa25760405162461bcd60e51b815260040161078590611e6a565b6001600160a01b0381166110075760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610785565b610c0781611324565b6006546001600160a01b0316331461103a5760405162461bcd60e51b815260040161078590611e6a565b600761098c8282611eea565b5f81815260046020526040902080546001600160a01b0319166001600160a01b038416908117909155819061107a82610990565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b5f818152600260205260408120546001600160a01b031661112b5760405162461bcd60e51b815260206004820152602c60248201527f4552433732313a206f70657261746f7220717565727920666f72206e6f6e657860448201526b34b9ba32b73a103a37b5b2b760a11b6064820152608401610785565b5f61113583610990565b9050806001600160a01b0316846001600160a01b031614806111705750836001600160a01b031661116584610711565b6001600160a01b0316145b8061118057506111808185610f4b565b949350505050565b826001600160a01b031661119b82610990565b6001600160a01b0316146112035760405162461bcd60e51b815260206004820152602960248201527f4552433732313a207472616e73666572206f6620746f6b656e2074686174206960448201526839903737ba1037bbb760b91b6064820152608401610785565b6001600160a01b0382166112655760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b6064820152608401610785565b61126f5f82611046565b6001600160a01b0383165f908152600360205260408120805460019290611297908490612030565b90915550506001600160a01b0382165f9081526003602052604081208054600192906112c4908490612043565b90915550505f8181526002602052604080822080546001600160a01b0319166001600160a01b0386811691821790925591518493918716917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b5f811180156113a457507f00000000000000000000000000000000000000000000000000000000000000008111155b6113fb5760405162461bcd60e51b815260206004820152602260248201527f546f6f206d616e7920746f6b656e7320696e206f6e65207472616e736163746960448201526137b760f11b6064820152608401610785565b7f000000000000000000000000000000000000000000000000000000000000000081600b5461142a9190612043565b11156114785760405162461bcd60e51b815260206004820181905260248201527f507572636861736520776f756c6420657863656564206d617820737570706c796044820152606401610785565b346114a3827f0000000000000000000000000000000000000000000000000000000000000000612056565b11156114f15760405162461bcd60e51b815260206004820152601f60248201527f45746865722076616c75652073656e74206973206e6f7420636f7272656374006044820152606401610785565b600b546114fe8282612043565b600b555f5b828110156108b85761151e336115198385612043565b611747565b600101611503565b816001600160a01b0316836001600160a01b0316036115875760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c6572000000000000006044820152606401610785565b6001600160a01b038381165f81815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b5f826115ff8584611760565b14949350505050565b611613848484611188565b61161f84848484611800565b610dbb5760405162461bcd60e51b81526004016107859061206d565b60606009805461069090611de1565b6060815f036116705750506040805180820190915260018152600360fc1b602082015290565b815f5b81156116995780611683816120bf565b91506116929050600a836120eb565b9150611673565b5f8167ffffffffffffffff8111156116b3576116b3611b90565b6040519080825280601f01601f1916602001820160405280156116dd576020820181803683370190505b5090505b8415611180576116f2600183612030565b91506116ff600a866120fe565b61170a906030612043565b60f81b81838151811061171f5761171f612111565b60200101906001600160f81b03191690815f1a905350611740600a866120eb565b94506116e1565b61098c828260405180602001604052805f8152506118fd565b5f81815b84518110156117f8575f85828151811061178057611780612111565b602002602001015190508083116117c25760408051602081018590529081018290526060016040516020818303038152906040528051906020012092506117ef565b60408051602081018390529081018490526060016040516020818303038152906040528051906020012092505b50600101611764565b509392505050565b5f6001600160a01b0384163b156118f257604051630a85bd0160e11b81526001600160a01b0385169063150b7a0290611843903390899088908890600401612125565b6020604051808303815f875af192505050801561187d575060408051601f3d908101601f1916820190925261187a91810190612161565b60015b6118d8573d8080156118aa576040519150601f19603f3d011682016040523d82523d5f602084013e6118af565b606091505b5080515f036118d05760405162461bcd60e51b81526004016107859061206d565b805181602001fd5b6001600160e01b031916630a85bd0160e11b149050611180565b506001949350505050565b611907838361192f565b6119135f848484611800565b6108b85760405162461bcd60e51b81526004016107859061206d565b6001600160a01b0382166119855760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f20616464726573736044820152606401610785565b5f818152600260205260409020546001600160a01b0316156119e95760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e746564000000006044820152606401610785565b6001600160a01b0382165f908152600360205260408120805460019290611a11908490612043565b90915550505f8181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b6001600160e01b031981168114610c07575f80fd5b5f60208284031215611a93575f80fd5b8135610f4481611a6e565b5f5b83811015611ab8578181015183820152602001611aa0565b50505f910152565b5f8151808452611ad7816020860160208601611a9e565b601f01601f19169290920160200192915050565b602081525f610f446020830184611ac0565b5f60208284031215611b0d575f80fd5b5035919050565b80356001600160a01b0381168114611b2a575f80fd5b919050565b5f8060408385031215611b40575f80fd5b611b4983611b14565b946020939093013593505050565b5f805f60608486031215611b69575f80fd5b611b7284611b14565b9250611b8060208501611b14565b9150604084013590509250925092565b634e487b7160e01b5f52604160045260245ffd5b5f67ffffffffffffffff80841115611bbe57611bbe611b90565b604051601f8501601f19908116603f01168101908282118183101715611be657611be6611b90565b81604052809350858152868686011115611bfe575f80fd5b858560208301375f602087830101525050509392505050565b5f60208284031215611c27575f80fd5b813567ffffffffffffffff811115611c3d575f80fd5b8201601f81018413611c4d575f80fd5b61118084823560208401611ba4565b5f60208284031215611c6c575f80fd5b610f4482611b14565b80358015158114611b2a575f80fd5b5f8060408385031215611c95575f80fd5b611c9e83611b14565b9150611cac60208401611c75565b90509250929050565b5f805f60408486031215611cc7575f80fd5b833567ffffffffffffffff80821115611cde575f80fd5b818601915086601f830112611cf1575f80fd5b813581811115611cff575f80fd5b8760208260051b8501011115611d13575f80fd5b6020928301989097509590910135949350505050565b5f805f8060808587031215611d3c575f80fd5b611d4585611b14565b9350611d5360208601611b14565b925060408501359150606085013567ffffffffffffffff811115611d75575f80fd5b8501601f81018713611d85575f80fd5b611d9487823560208401611ba4565b91505092959194509250565b5f60208284031215611db0575f80fd5b610f4482611c75565b5f8060408385031215611dca575f80fd5b611dd383611b14565b9150611cac60208401611b14565b600181811c90821680611df557607f821691505b602082108103611e1357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526031908201527f4552433732313a207472616e736665722063616c6c6572206973206e6f74206f6040820152701ddb995c881b9bdc88185c1c1c9bdd9959607a1b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b601f8211156108b857805f5260205f20601f840160051c81016020851015611ec45750805b601f840160051c820191505b81811015611ee3575f8155600101611ed0565b5050505050565b815167ffffffffffffffff811115611f0457611f04611b90565b611f1881611f128454611de1565b84611e9f565b602080601f831160018114611f4b575f8415611f345750858301515b5f19600386901b1c1916600185901b178555611fa2565b5f85815260208120601f198616915b82811015611f7957888601518255948401946001909101908401611f5a565b5085821015611f9657878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b60208082526024908201527f4e6f207472616e73616374696f6e2066726f6d20736d61727420636f6e7472616040820152636374732160e01b606082015260800190565b5f8351611fff818460208801611a9e565b835190830190612013818360208801611a9e565b01949350505050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561067c5761067c61201c565b8082018082111561067c5761067c61201c565b808202811582820484141761067c5761067c61201c565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b5f600182016120d0576120d061201c565b5060010190565b634e487b7160e01b5f52601260045260245ffd5b5f826120f9576120f96120d7565b500490565b5f8261210c5761210c6120d7565b500690565b634e487b7160e01b5f52603260045260245ffd5b6001600160a01b03858116825284166020820152604081018390526080606082018190525f9061215790830184611ac0565b9695505050505050565b5f60208284031215612171575f80fd5b8151610f4481611a6e56fea264697066735822122096ecacbc4bc1ca1b22ed8f653eb6f64efc6f654d4d1c54b1b7aee09f070dca2264736f6c63430008180033",
}

// CollectionNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use CollectionNFTMetaData.ABI instead.
var CollectionNFTABI = CollectionNFTMetaData.ABI

// CollectionNFTBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CollectionNFTMetaData.Bin instead.
var CollectionNFTBin = CollectionNFTMetaData.Bin

// DeployCollectionNFT deploys a new Ethereum contract, binding an instance of CollectionNFT to it.
func DeployCollectionNFT(auth *bind.TransactOpts, backend bind.ContractBackend, baseURI string, merkleRoot_ [32]byte, name_ string, maxSupply_ *big.Int, price_ *big.Int, maxPerMint_ *big.Int) (common.Address, *types.Transaction, *CollectionNFT, error) {
	parsed, err := CollectionNFTMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CollectionNFTBin), backend, baseURI, merkleRoot_, name_, maxSupply_, price_, maxPerMint_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CollectionNFT{CollectionNFTCaller: CollectionNFTCaller{contract: contract}, CollectionNFTTransactor: CollectionNFTTransactor{contract: contract}, CollectionNFTFilterer: CollectionNFTFilterer{contract: contract}}, nil
}

// CollectionNFT is an auto generated Go binding around an Ethereum contract.
type CollectionNFT struct {
	CollectionNFTCaller     // Read-only binding to the contract
	CollectionNFTTransactor // Write-only binding to the contract
	CollectionNFTFilterer   // Log filterer for contract events
}

// CollectionNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollectionNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollectionNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollectionNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollectionNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollectionNFTSession struct {
	Contract     *CollectionNFT    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CollectionNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollectionNFTCallerSession struct {
	Contract *CollectionNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// CollectionNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollectionNFTTransactorSession struct {
	Contract     *CollectionNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CollectionNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollectionNFTRaw struct {
	Contract *CollectionNFT // Generic contract binding to access the raw methods on
}

// CollectionNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollectionNFTCallerRaw struct {
	Contract *CollectionNFTCaller // Generic read-only contract binding to access the raw methods on
}

// CollectionNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollectionNFTTransactorRaw struct {
	Contract *CollectionNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollectionNFT creates a new instance of CollectionNFT, bound to a specific deployed contract.
func NewCollectionNFT(address common.Address, backend bind.ContractBackend) (*CollectionNFT, error) {
	contract, err := bindCollectionNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CollectionNFT{CollectionNFTCaller: CollectionNFTCaller{contract: contract}, CollectionNFTTransactor: CollectionNFTTransactor{contract: contract}, CollectionNFTFilterer: CollectionNFTFilterer{contract: contract}}, nil
}

// NewCollectionNFTCaller creates a new read-only instance of CollectionNFT, bound to a specific deployed contract.
func NewCollectionNFTCaller(address common.Address, caller bind.ContractCaller) (*CollectionNFTCaller, error) {
	contract, err := bindCollectionNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTCaller{contract: contract}, nil
}

// NewCollectionNFTTransactor creates a new write-only instance of CollectionNFT, bound to a specific deployed contract.
func NewCollectionNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*CollectionNFTTransactor, error) {
	contract, err := bindCollectionNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTTransactor{contract: contract}, nil
}

// NewCollectionNFTFilterer creates a new log filterer instance of CollectionNFT, bound to a specific deployed contract.
func NewCollectionNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*CollectionNFTFilterer, error) {
	contract, err := bindCollectionNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTFilterer{contract: contract}, nil
}

// bindCollectionNFT binds a generic wrapper to an already deployed contract.
func bindCollectionNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CollectionNFTABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollectionNFT *CollectionNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollectionNFT.Contract.CollectionNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollectionNFT *CollectionNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollectionNFT.Contract.CollectionNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollectionNFT *CollectionNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollectionNFT.Contract.CollectionNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollectionNFT *CollectionNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollectionNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollectionNFT *CollectionNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollectionNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollectionNFT *CollectionNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollectionNFT.Contract.contract.Transact(opts, method, params...)
}

// MAXPERMINT is a free data retrieval call binding the contract method 0x09d42b30.
//
// Solidity: function MAX_PER_MINT() view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) MAXPERMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "MAX_PER_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPERMINT is a free data retrieval call binding the contract method 0x09d42b30.
//
// Solidity: function MAX_PER_MINT() view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) MAXPERMINT() (*big.Int, error) {
	return _CollectionNFT.Contract.MAXPERMINT(&_CollectionNFT.CallOpts)
}

// MAXPERMINT is a free data retrieval call binding the contract method 0x09d42b30.
//
// Solidity: function MAX_PER_MINT() view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) MAXPERMINT() (*big.Int, error) {
	return _CollectionNFT.Contract.MAXPERMINT(&_CollectionNFT.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) MAXSUPPLY() (*big.Int, error) {
	return _CollectionNFT.Contract.MAXSUPPLY(&_CollectionNFT.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _CollectionNFT.Contract.MAXSUPPLY(&_CollectionNFT.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x8d859f3e.
//
// Solidity: function PRICE() view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) PRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "PRICE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRICE is a free data retrieval call binding the contract method 0x8d859f3e.
//
// Solidity: function PRICE() view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) PRICE() (*big.Int, error) {
	return _CollectionNFT.Contract.PRICE(&_CollectionNFT.CallOpts)
}

// PRICE is a free data retrieval call binding the contract method 0x8d859f3e.
//
// Solidity: function PRICE() view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) PRICE() (*big.Int, error) {
	return _CollectionNFT.Contract.PRICE(&_CollectionNFT.CallOpts)
}

// PROVENANCE is a free data retrieval call binding the contract method 0x6373a6b1.
//
// Solidity: function PROVENANCE() view returns(string)
func (_CollectionNFT *CollectionNFTCaller) PROVENANCE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "PROVENANCE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PROVENANCE is a free data retrieval call binding the contract method 0x6373a6b1.
//
// Solidity: function PROVENANCE() view returns(string)
func (_CollectionNFT *CollectionNFTSession) PROVENANCE() (string, error) {
	return _CollectionNFT.Contract.PROVENANCE(&_CollectionNFT.CallOpts)
}

// PROVENANCE is a free data retrieval call binding the contract method 0x6373a6b1.
//
// Solidity: function PROVENANCE() view returns(string)
func (_CollectionNFT *CollectionNFTCallerSession) PROVENANCE() (string, error) {
	return _CollectionNFT.Contract.PROVENANCE(&_CollectionNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CollectionNFT.Contract.BalanceOf(&_CollectionNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CollectionNFT.Contract.BalanceOf(&_CollectionNFT.CallOpts, owner)
}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) CurrentIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "currentIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) CurrentIndex() (*big.Int, error) {
	return _CollectionNFT.Contract.CurrentIndex(&_CollectionNFT.CallOpts)
}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) CurrentIndex() (*big.Int, error) {
	return _CollectionNFT.Contract.CurrentIndex(&_CollectionNFT.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CollectionNFT.Contract.GetApproved(&_CollectionNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CollectionNFT.Contract.GetApproved(&_CollectionNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CollectionNFT *CollectionNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CollectionNFT *CollectionNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CollectionNFT.Contract.IsApprovedForAll(&_CollectionNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CollectionNFT *CollectionNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CollectionNFT.Contract.IsApprovedForAll(&_CollectionNFT.CallOpts, owner, operator)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_CollectionNFT *CollectionNFTCaller) MerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "merkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_CollectionNFT *CollectionNFTSession) MerkleRoot() ([32]byte, error) {
	return _CollectionNFT.Contract.MerkleRoot(&_CollectionNFT.CallOpts)
}

// MerkleRoot is a free data retrieval call binding the contract method 0x2eb4a7ab.
//
// Solidity: function merkleRoot() view returns(bytes32)
func (_CollectionNFT *CollectionNFTCallerSession) MerkleRoot() ([32]byte, error) {
	return _CollectionNFT.Contract.MerkleRoot(&_CollectionNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CollectionNFT *CollectionNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CollectionNFT *CollectionNFTSession) Name() (string, error) {
	return _CollectionNFT.Contract.Name(&_CollectionNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CollectionNFT *CollectionNFTCallerSession) Name() (string, error) {
	return _CollectionNFT.Contract.Name(&_CollectionNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollectionNFT *CollectionNFTCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollectionNFT *CollectionNFTSession) Owner() (common.Address, error) {
	return _CollectionNFT.Contract.Owner(&_CollectionNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollectionNFT *CollectionNFTCallerSession) Owner() (common.Address, error) {
	return _CollectionNFT.Contract.Owner(&_CollectionNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CollectionNFT.Contract.OwnerOf(&_CollectionNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CollectionNFT *CollectionNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CollectionNFT.Contract.OwnerOf(&_CollectionNFT.CallOpts, tokenId)
}

// PresaleIsActive is a free data retrieval call binding the contract method 0x30f72cd4.
//
// Solidity: function presaleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTCaller) PresaleIsActive(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "presaleIsActive")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PresaleIsActive is a free data retrieval call binding the contract method 0x30f72cd4.
//
// Solidity: function presaleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTSession) PresaleIsActive() (bool, error) {
	return _CollectionNFT.Contract.PresaleIsActive(&_CollectionNFT.CallOpts)
}

// PresaleIsActive is a free data retrieval call binding the contract method 0x30f72cd4.
//
// Solidity: function presaleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTCallerSession) PresaleIsActive() (bool, error) {
	return _CollectionNFT.Contract.PresaleIsActive(&_CollectionNFT.CallOpts)
}

// SaleIsActive is a free data retrieval call binding the contract method 0xeb8d2444.
//
// Solidity: function saleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTCaller) SaleIsActive(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "saleIsActive")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SaleIsActive is a free data retrieval call binding the contract method 0xeb8d2444.
//
// Solidity: function saleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTSession) SaleIsActive() (bool, error) {
	return _CollectionNFT.Contract.SaleIsActive(&_CollectionNFT.CallOpts)
}

// SaleIsActive is a free data retrieval call binding the contract method 0xeb8d2444.
//
// Solidity: function saleIsActive() view returns(bool)
func (_CollectionNFT *CollectionNFTCallerSession) SaleIsActive() (bool, error) {
	return _CollectionNFT.Contract.SaleIsActive(&_CollectionNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CollectionNFT *CollectionNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CollectionNFT *CollectionNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CollectionNFT.Contract.SupportsInterface(&_CollectionNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CollectionNFT *CollectionNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CollectionNFT.Contract.SupportsInterface(&_CollectionNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CollectionNFT *CollectionNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CollectionNFT *CollectionNFTSession) Symbol() (string, error) {
	return _CollectionNFT.Contract.Symbol(&_CollectionNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CollectionNFT *CollectionNFTCallerSession) Symbol() (string, error) {
	return _CollectionNFT.Contract.Symbol(&_CollectionNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CollectionNFT *CollectionNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CollectionNFT *CollectionNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CollectionNFT.Contract.TokenURI(&_CollectionNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CollectionNFT *CollectionNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CollectionNFT.Contract.TokenURI(&_CollectionNFT.CallOpts, tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CollectionNFT *CollectionNFTCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CollectionNFT *CollectionNFTSession) TotalSupply() (*big.Int, error) {
	return _CollectionNFT.Contract.TotalSupply(&_CollectionNFT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CollectionNFT *CollectionNFTCallerSession) TotalSupply() (*big.Int, error) {
	return _CollectionNFT.Contract.TotalSupply(&_CollectionNFT.CallOpts)
}

// WhitelistClaimed is a free data retrieval call binding the contract method 0xdb4bec44.
//
// Solidity: function whitelistClaimed(address ) view returns(bool)
func (_CollectionNFT *CollectionNFTCaller) WhitelistClaimed(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _CollectionNFT.contract.Call(opts, &out, "whitelistClaimed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WhitelistClaimed is a free data retrieval call binding the contract method 0xdb4bec44.
//
// Solidity: function whitelistClaimed(address ) view returns(bool)
func (_CollectionNFT *CollectionNFTSession) WhitelistClaimed(arg0 common.Address) (bool, error) {
	return _CollectionNFT.Contract.WhitelistClaimed(&_CollectionNFT.CallOpts, arg0)
}

// WhitelistClaimed is a free data retrieval call binding the contract method 0xdb4bec44.
//
// Solidity: function whitelistClaimed(address ) view returns(bool)
func (_CollectionNFT *CollectionNFTCallerSession) WhitelistClaimed(arg0 common.Address) (bool, error) {
	return _CollectionNFT.Contract.WhitelistClaimed(&_CollectionNFT.CallOpts, arg0)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.Approve(&_CollectionNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.Approve(&_CollectionNFT.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 numberOfTokens) payable returns()
func (_CollectionNFT *CollectionNFTTransactor) Mint(opts *bind.TransactOpts, numberOfTokens *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "mint", numberOfTokens)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 numberOfTokens) payable returns()
func (_CollectionNFT *CollectionNFTSession) Mint(numberOfTokens *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.Mint(&_CollectionNFT.TransactOpts, numberOfTokens)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 numberOfTokens) payable returns()
func (_CollectionNFT *CollectionNFTTransactorSession) Mint(numberOfTokens *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.Mint(&_CollectionNFT.TransactOpts, numberOfTokens)
}

// MintWhitelist is a paid mutator transaction binding the contract method 0xa6d612f9.
//
// Solidity: function mintWhitelist(bytes32[] _merkleProof, uint256 amount) payable returns()
func (_CollectionNFT *CollectionNFTTransactor) MintWhitelist(opts *bind.TransactOpts, _merkleProof [][32]byte, amount *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "mintWhitelist", _merkleProof, amount)
}

// MintWhitelist is a paid mutator transaction binding the contract method 0xa6d612f9.
//
// Solidity: function mintWhitelist(bytes32[] _merkleProof, uint256 amount) payable returns()
func (_CollectionNFT *CollectionNFTSession) MintWhitelist(_merkleProof [][32]byte, amount *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.MintWhitelist(&_CollectionNFT.TransactOpts, _merkleProof, amount)
}

// MintWhitelist is a paid mutator transaction binding the contract method 0xa6d612f9.
//
// Solidity: function mintWhitelist(bytes32[] _merkleProof, uint256 amount) payable returns()
func (_CollectionNFT *CollectionNFTTransactorSession) MintWhitelist(_merkleProof [][32]byte, amount *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.MintWhitelist(&_CollectionNFT.TransactOpts, _merkleProof, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollectionNFT *CollectionNFTTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollectionNFT *CollectionNFTSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollectionNFT.Contract.RenounceOwnership(&_CollectionNFT.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollectionNFT *CollectionNFTTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollectionNFT.Contract.RenounceOwnership(&_CollectionNFT.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SafeTransferFrom(&_CollectionNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SafeTransferFrom(&_CollectionNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_CollectionNFT *CollectionNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_CollectionNFT *CollectionNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SafeTransferFrom0(&_CollectionNFT.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SafeTransferFrom0(&_CollectionNFT.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CollectionNFT *CollectionNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetApprovalForAll(&_CollectionNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetApprovalForAll(&_CollectionNFT.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetBaseURI(opts *bind.TransactOpts, baseURI_ string) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setBaseURI", baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_CollectionNFT *CollectionNFTSession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetBaseURI(&_CollectionNFT.TransactOpts, baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI_) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetBaseURI(baseURI_ string) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetBaseURI(&_CollectionNFT.TransactOpts, baseURI_)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 merkleRoot_) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetMerkleRoot(opts *bind.TransactOpts, merkleRoot_ [32]byte) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setMerkleRoot", merkleRoot_)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 merkleRoot_) returns()
func (_CollectionNFT *CollectionNFTSession) SetMerkleRoot(merkleRoot_ [32]byte) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetMerkleRoot(&_CollectionNFT.TransactOpts, merkleRoot_)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x7cb64759.
//
// Solidity: function setMerkleRoot(bytes32 merkleRoot_) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetMerkleRoot(merkleRoot_ [32]byte) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetMerkleRoot(&_CollectionNFT.TransactOpts, merkleRoot_)
}

// SetPresaleState is a paid mutator transaction binding the contract method 0xbce4d6ae.
//
// Solidity: function setPresaleState(bool newPresaleState) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetPresaleState(opts *bind.TransactOpts, newPresaleState bool) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setPresaleState", newPresaleState)
}

// SetPresaleState is a paid mutator transaction binding the contract method 0xbce4d6ae.
//
// Solidity: function setPresaleState(bool newPresaleState) returns()
func (_CollectionNFT *CollectionNFTSession) SetPresaleState(newPresaleState bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetPresaleState(&_CollectionNFT.TransactOpts, newPresaleState)
}

// SetPresaleState is a paid mutator transaction binding the contract method 0xbce4d6ae.
//
// Solidity: function setPresaleState(bool newPresaleState) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetPresaleState(newPresaleState bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetPresaleState(&_CollectionNFT.TransactOpts, newPresaleState)
}

// SetProvenance is a paid mutator transaction binding the contract method 0xffe630b5.
//
// Solidity: function setProvenance(string provenance) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetProvenance(opts *bind.TransactOpts, provenance string) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setProvenance", provenance)
}

// SetProvenance is a paid mutator transaction binding the contract method 0xffe630b5.
//
// Solidity: function setProvenance(string provenance) returns()
func (_CollectionNFT *CollectionNFTSession) SetProvenance(provenance string) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetProvenance(&_CollectionNFT.TransactOpts, provenance)
}

// SetProvenance is a paid mutator transaction binding the contract method 0xffe630b5.
//
// Solidity: function setProvenance(string provenance) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetProvenance(provenance string) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetProvenance(&_CollectionNFT.TransactOpts, provenance)
}

// SetSaleState is a paid mutator transaction binding the contract method 0xc4e37095.
//
// Solidity: function setSaleState(bool newState) returns()
func (_CollectionNFT *CollectionNFTTransactor) SetSaleState(opts *bind.TransactOpts, newState bool) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "setSaleState", newState)
}

// SetSaleState is a paid mutator transaction binding the contract method 0xc4e37095.
//
// Solidity: function setSaleState(bool newState) returns()
func (_CollectionNFT *CollectionNFTSession) SetSaleState(newState bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetSaleState(&_CollectionNFT.TransactOpts, newState)
}

// SetSaleState is a paid mutator transaction binding the contract method 0xc4e37095.
//
// Solidity: function setSaleState(bool newState) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) SetSaleState(newState bool) (*types.Transaction, error) {
	return _CollectionNFT.Contract.SetSaleState(&_CollectionNFT.TransactOpts, newState)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.TransferFrom(&_CollectionNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CollectionNFT.Contract.TransferFrom(&_CollectionNFT.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollectionNFT *CollectionNFTTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollectionNFT *CollectionNFTSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollectionNFT.Contract.TransferOwnership(&_CollectionNFT.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollectionNFT *CollectionNFTTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollectionNFT.Contract.TransferOwnership(&_CollectionNFT.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CollectionNFT *CollectionNFTTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollectionNFT.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CollectionNFT *CollectionNFTSession) Withdraw() (*types.Transaction, error) {
	return _CollectionNFT.Contract.Withdraw(&_CollectionNFT.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_CollectionNFT *CollectionNFTTransactorSession) Withdraw() (*types.Transaction, error) {
	return _CollectionNFT.Contract.Withdraw(&_CollectionNFT.TransactOpts)
}

// CollectionNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CollectionNFT contract.
type CollectionNFTApprovalIterator struct {
	Event *CollectionNFTApproval // Event containing the contract specifics and raw log

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
func (it *CollectionNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollectionNFTApproval)
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
		it.Event = new(CollectionNFTApproval)
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
func (it *CollectionNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollectionNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollectionNFTApproval represents a Approval event raised by the CollectionNFT contract.
type CollectionNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*CollectionNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CollectionNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTApprovalIterator{contract: _CollectionNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CollectionNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CollectionNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollectionNFTApproval)
				if err := _CollectionNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) ParseApproval(log types.Log) (*CollectionNFTApproval, error) {
	event := new(CollectionNFTApproval)
	if err := _CollectionNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollectionNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the CollectionNFT contract.
type CollectionNFTApprovalForAllIterator struct {
	Event *CollectionNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *CollectionNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollectionNFTApprovalForAll)
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
		it.Event = new(CollectionNFTApprovalForAll)
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
func (it *CollectionNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollectionNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollectionNFTApprovalForAll represents a ApprovalForAll event raised by the CollectionNFT contract.
type CollectionNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CollectionNFT *CollectionNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*CollectionNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CollectionNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTApprovalForAllIterator{contract: _CollectionNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CollectionNFT *CollectionNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *CollectionNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CollectionNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollectionNFTApprovalForAll)
				if err := _CollectionNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CollectionNFT *CollectionNFTFilterer) ParseApprovalForAll(log types.Log) (*CollectionNFTApprovalForAll, error) {
	event := new(CollectionNFTApprovalForAll)
	if err := _CollectionNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollectionNFTOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CollectionNFT contract.
type CollectionNFTOwnershipTransferredIterator struct {
	Event *CollectionNFTOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CollectionNFTOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollectionNFTOwnershipTransferred)
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
		it.Event = new(CollectionNFTOwnershipTransferred)
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
func (it *CollectionNFTOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollectionNFTOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollectionNFTOwnershipTransferred represents a OwnershipTransferred event raised by the CollectionNFT contract.
type CollectionNFTOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollectionNFT *CollectionNFTFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CollectionNFTOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollectionNFT.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTOwnershipTransferredIterator{contract: _CollectionNFT.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollectionNFT *CollectionNFTFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CollectionNFTOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollectionNFT.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollectionNFTOwnershipTransferred)
				if err := _CollectionNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollectionNFT *CollectionNFTFilterer) ParseOwnershipTransferred(log types.Log) (*CollectionNFTOwnershipTransferred, error) {
	event := new(CollectionNFTOwnershipTransferred)
	if err := _CollectionNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollectionNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CollectionNFT contract.
type CollectionNFTTransferIterator struct {
	Event *CollectionNFTTransfer // Event containing the contract specifics and raw log

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
func (it *CollectionNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollectionNFTTransfer)
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
		it.Event = new(CollectionNFTTransfer)
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
func (it *CollectionNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollectionNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollectionNFTTransfer represents a Transfer event raised by the CollectionNFT contract.
type CollectionNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*CollectionNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CollectionNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CollectionNFTTransferIterator{contract: _CollectionNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CollectionNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _CollectionNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollectionNFTTransfer)
				if err := _CollectionNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CollectionNFT *CollectionNFTFilterer) ParseTransfer(log types.Log) (*CollectionNFTTransfer, error) {
	event := new(CollectionNFTTransfer)
	if err := _CollectionNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
