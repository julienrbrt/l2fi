package l2

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var _ L2 = (*ArbitrumClient)(nil)

// ArbitrumClient handles Arbitrum-specific logic
type ArbitrumClient struct {
	ethClient      *ethclient.Client
	arbitrumClient *ethclient.Client

	delayedInboxAddress common.Address
	inbox               *Inbox
}

// NewArbitrumClient creates a new ArbitrumClient instance.
func NewArbitrumClient(ethRPCURL, arbitrumRPCURL string, delayedInboxAddress string) (*ArbitrumClient, error) {
	client, err := ethclient.Dial(ethRPCURL)
	if err != nil {
		return nil, err
	}

	arbitrumClient, err := ethclient.Dial(arbitrumRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Arbitrum RPC: %w", err)
	}

	delayedAddr := common.HexToAddress(delayedInboxAddress)
	inbox, err := NewInbox(delayedAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind inbox: %w", err)
	}

	return &ArbitrumClient{
		ethClient:           client,
		arbitrumClient:      arbitrumClient,
		delayedInboxAddress: delayedAddr,
		inbox:               inbox,
	}, nil
}

// BuildForceInclusionTx builds a forced inclusion transaction for Arbitrum using CreateRetryableTicket.
// Returns a JSON string compatible with ethers.js for client-side signing.
func (a *ArbitrumClient) BuildForceInclusionTx(
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

	nonce, err := a.ethClient.PendingNonceAt(context.Background(), from)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := a.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	arbGasPrice, err := a.arbitrumClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get Arbitrum gas price: %w", err)
	}

	arbNonce, err := a.arbitrumClient.PendingNonceAt(context.Background(), from)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	tx, err := a.inbox.SendL1FundedUnsignedTransaction(
		&bind.TransactOpts{
			From:     from,
			Nonce:    big.NewInt(int64(nonce)),
			GasPrice: gasPrice,
			GasLimit: 200000,
			Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
				return t, nil // unsigned
			},
			NoSend: true,
			// Value:  value,
		},
		big.NewInt(int64(l2GasLimit)),
		arbGasPrice,
		big.NewInt(int64(arbNonce)),
		to,
		callData,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create retryable ticket: %w", err)
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
