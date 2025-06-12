package l2

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

var _ L2 = (*ArbitrumClient)(nil)

// ArbitrumClient handles Arbitrum-specific logic
type ArbitrumClient struct {
	client *ethclient.Client
}

// NewArbitrumClient creates a new ArbitrumClient instance.
func NewArbitrumClient(rpcURL string, delayedInboxAddress string) (*ArbitrumClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	return &ArbitrumClient{
		client: client,
	}, nil
}

func (a *ArbitrumClient) BuildForceInclusionTx(fromAddress, toAddress, data string, value *big.Int, gasLimit uint64) (string, error) {
	return "", errors.New("Arbitrum stack is currently disabled")
}
