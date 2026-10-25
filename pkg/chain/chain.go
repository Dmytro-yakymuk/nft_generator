// Package chain describes the Ethereum network the platform works with.
package chain

// Config describes the network in the format of the wallet_addEthereumChain request,
// so the web pages can switch MetaMask to it or add it to the wallet.
type Config struct {
	ChainID          string `json:"chainId"`
	ChainName        string `json:"chainName"`
	RPCURL           string `json:"rpcUrl"`
	BlockExplorerURL string `json:"blockExplorerUrl"`
}
