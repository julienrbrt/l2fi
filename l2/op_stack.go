package l2

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum-optimism/optimism/op-node/bindings"
)

var _ L2 = (*OpStackClient)(nil)

// OpStackClient handles Optimism-specific logic
type OpStackClient struct {
	client                *ethclient.Client
	optimismPortal        *bindings.OptimismPortal
	optimismPortalAddress common.Address
}

// NewOpStackClient creates a new OpStackClient instance.
func NewOpStackClient(rpcURL string, optimismPortalAddress string) (*OpStackClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	portalAddr := common.HexToAddress(optimismPortalAddress)
	optimismPortal, err := bindings.NewOptimismPortal(
		portalAddr,
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OptimismPortal binding: %w", err)
	}

	return &OpStackClient{
		client:                client,
		optimismPortal:        optimismPortal,
		optimismPortalAddress: portalAddr,
	}, nil
}

// BuildForceInclusionTx builds a forced inclusion transaction for the Optimism stack.
// It returns a JSON string compatible with ethers.js for client-side signing.
func (o *OpStackClient) BuildForceInclusionTx(
	fromAddress string,
	toAddress string,
	data string,
	value *big.Int,
	l2GasLimit uint64,
) (string, error) {
	if !common.IsHexAddress(fromAddress) {
		return "", fmt.Errorf("invalid from address: %s", fromAddress)
	}
	if !common.IsHexAddress(toAddress) {
		return "", fmt.Errorf("invalid to address: %s", toAddress)
	}
	if value == nil || value.Sign() < 0 {
		return "", fmt.Errorf("value must be non-negative and non-nil")
	}

	from := common.HexToAddress(fromAddress)
	to := common.HexToAddress(toAddress)

	callData := []byte{}
	if data != "0x" && data != "" {
		callData = common.FromHex(data)
	}

	nonce, err := o.client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := o.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// create forced inclusion tx
	tx, err := o.optimismPortal.DepositTransaction(
		&bind.TransactOpts{
			From:     from,
			Nonce:    big.NewInt(int64(nonce)),
			GasPrice: gasPrice,
			GasLimit: 200000,
			Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
				// transaction is unsigned; must be signed by the client
				return t, nil
			},
			NoSend: true,
		},
		to,
		value,
		l2GasLimit,
		false,
		callData,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create forced inclusion transaction: %w", err)
	}

	ethersTransaction := map[string]any{
		"to":       tx.To().Hex(),
		"data":     fmt.Sprintf("0x%x", tx.Data()),
		"value":    fmt.Sprintf("0x%x", tx.Value()),
		"gasLimit": fmt.Sprintf("0x%x", tx.Gas()),
		"gasPrice": fmt.Sprintf("0x%x", tx.GasPrice()),
		"nonce":    fmt.Sprintf("0x%x", tx.Nonce()),
	}

	txData, err := json.Marshal(ethersTransaction)
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %w", err)
	}

	return string(txData), nil
}
