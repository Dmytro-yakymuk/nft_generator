package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/nfts"
	"generator/storage"
	"generator/whitelist"
)

// deployCmd deploys the contract of a collection from the command line, a fallback for the web page:
// the same deployment data (name, base URI, Merkle root of the whitelist, sale terms), signed with a key from the
// environment instead of MetaMask. The deployed address and block are saved to the collection.
var (
	deployCmd = &cobra.Command{
		Use:   "deploy-contract",
		Short: "deploys the contract of a collection with a key from the environment",
		RunE:  cmdDeployContract,
		// a failed deployment prints its reason, not the usage.
		SilenceUsage: true,
	}
	deployProject string
	deployKeyEnv  string
	deployRPC     string
	deployTipWei  int64
	deployDryRun  bool
	deployTerms   nfts.Terms
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVar(&deployProject, "project", "", "id of the collection")
	deployCmd.Flags().StringVar(&deployKeyEnv, "key-env", "DEPLOYER_PRIVATE_KEY", "environment variable with the hex private key of the collection owner's wallet")
	deployCmd.Flags().StringVar(&deployRPC, "rpc", "", "JSON-RPC of the network (default: the server RPC from the config)")
	deployCmd.Flags().Int64Var(&deployTipWei, "tip", 0, "priority fee per gas in wei (default: suggested by the node)")
	deployCmd.Flags().BoolVar(&deployDryRun, "dry-run", false, "estimate the gas and the cost without sending the transaction")
	deployCmd.Flags().StringVar(&deployTerms.Name, "name", "", "name of the contract (default: the collection name)")
	deployCmd.Flags().StringVar(&deployTerms.Price, "price", nfts.DefaultMintPrice, "price of one token, ETH")
	deployCmd.Flags().IntVar(&deployTerms.MaxPerMint, "max-per-mint", nfts.DefaultMaxPerMint, "maximum number of tokens in one transaction")
}

func cmdDeployContract(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	projectID, err := uuid.Parse(deployProject)
	if err != nil {
		return Error.New("--project must be the id of a collection: %v", err)
	}
	// the key is read from the environment, so it stays out of the shell history and the process list.
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(os.Getenv(deployKeyEnv)), "0x"))
	if err != nil {
		return Error.New("no private key in %s: %v", deployKeyEnv, err)
	}
	from := crypto.PubkeyToAddress(privateKey.PublicKey)

	cfg, err := readConfig()
	if err != nil {
		return Error.Wrap(err)
	}
	db, err := storage.New(cfg.Database)
	if err != nil {
		return Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	project, err := db.Collections().Get(ctx, projectID)
	if err != nil {
		return Error.Wrap(err)
	}
	owner, err := db.Accounts().Get(ctx, project.UserID)
	if err != nil {
		return Error.Wrap(err)
	}
	// the deployer becomes the owner of the contract, as when the owner deploys it from the web page.
	if owner.Wallet != from {
		return Error.New("the key of %s does not belong to the owner of the collection (%s)", from.Hex(), owner.Wallet.Hex())
	}

	collectionsService := collections.NewService(cfg.Collections.Config, db.Collections(), nil, nil, nil, nil)
	whitelistService := whitelist.NewService(cfg.Whitelist.Config, db.WhiteList(), collectionsService)
	data, err := nfts.NewService(cfg.NFTs.Config, collectionsService, whitelistService).
		GetDeployContractTransactionData(ctx, projectID, project.UserID, deployTerms)
	if err != nil {
		return Error.Wrap(err)
	}

	rpc := deployRPC
	if rpc == "" {
		rpc = cfg.Dashboard.RPCURL
	}
	if rpc == "" {
		rpc = cfg.Server.HTTP.Chain.RPCURL
	}
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		return Error.Wrap(err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return Error.Wrap(err)
	}
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return Error.Wrap(err)
	}
	tip := big.NewInt(deployTipWei)
	if deployTipWei == 0 {
		if tip, err = client.SuggestGasTipCap(ctx); err != nil {
			return Error.Wrap(err)
		}
	}
	// the fee cap leaves room for the base fee to double before the transaction is included.
	feeCap := new(big.Int).Add(new(big.Int).Mul(head.BaseFee, big.NewInt(2)), tip)
	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{From: from, Data: data, GasTipCap: tip, GasFeeCap: feeCap})
	if err != nil {
		return Error.Wrap(err)
	}
	balance, err := client.BalanceAt(ctx, from, nil)
	if err != nil {
		return Error.Wrap(err)
	}

	fmt.Printf("Collection:  %s (%s), %d tokens\n", project.Name, project.ID, project.NFTCount)
	fmt.Printf("Sale terms:  %s ETH per token, at most %d per transaction\n", deployTerms.Price, deployTerms.MaxPerMint)
	fmt.Printf("Network:     chain %s via %s\n", chainID, rpc)
	fmt.Printf("Deployer:    %s, balance %s ETH\n", from.Hex(), ether(balance))
	fmt.Printf("Gas:         %d, base fee %s gwei, tip %s gwei, at most %s ETH\n",
		gas, gwei(head.BaseFee), gwei(tip), ether(new(big.Int).Mul(new(big.Int).SetUint64(gas), feeCap)))
	if deployDryRun {
		return nil
	}

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return Error.Wrap(err)
	}
	tx, err := types.SignTx(types.NewTx(&types.DynamicFeeTx{
		ChainID: chainID, Nonce: nonce, GasTipCap: tip, GasFeeCap: feeCap, Gas: gas, Data: data,
	}), types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return Error.Wrap(err)
	}
	if err = client.SendTransaction(ctx, tx); err != nil {
		return Error.Wrap(err)
	}
	fmt.Printf("Transaction: %s, waiting for the receipt...\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return Error.Wrap(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return Error.New("the deployment failed in block %s", receipt.BlockNumber)
	}
	block, err := client.HeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return Error.Wrap(err)
	}
	price := new(big.Int).Add(block.BaseFee, tip)
	if price.Cmp(feeCap) > 0 {
		price = feeCap
	}

	if err = db.Collections().UpdateContractAddress(ctx, projectID, receipt.ContractAddress, receipt.BlockNumber.Uint64()); err != nil {
		return Error.Wrap(err)
	}
	fmt.Printf("Contract:    %s\n", receipt.ContractAddress.Hex())
	fmt.Printf("Block:       %s, gas used %d, cost %s ETH\n", receipt.BlockNumber, receipt.GasUsed,
		ether(new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), price)))
	fmt.Println("The address is saved to the collection: the dashboard and the mint site use it.")
	return nil
}

// ether formats wei as ETH with up to 9 decimals.
func ether(wei *big.Int) string {
	return strings.TrimRight(strings.TrimRight(new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 9), "0"), ".")
}

// gwei formats wei as gwei with up to 9 decimals.
func gwei(wei *big.Int) string {
	return strings.TrimRight(strings.TrimRight(new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e9)).Text('f', 9), "0"), ".")
}
