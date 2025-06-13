package l2_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/julienrbrt/l2fi/l2"
)

func TestArbitrumForceInclusion_BuildForceInclusionTx(t *testing.T) {
	tests := []struct {
		name      string
		toAddress string
		data      string
		amount    *big.Int
		gasLimit  uint64
		wantErr   bool
	}{
		{
			name:      "basic transaction",
			toAddress: "0x742d35cC6600C91D844c4c0C7A7D26b2b39d3c07",
			data:      "0x",
			amount:    big.NewInt(1000),
			gasLimit:  200000,
			wantErr:   false,
		},
		{
			name:      "transaction with data",
			toAddress: "0x742d35cC6600C91D844c4c0C7A7D26b2b39d3c07",
			data:      "0xa9059cbb000000000000000000000000742d35cc6600c91d844c4c0c7a7d26b2b39d3c070000000000000000000000000000000000000000000000000000000000000064",
			amount:    big.NewInt(0),
			gasLimit:  200000,
			wantErr:   false,
		},
		{
			name:      "invalid to address",
			toAddress: "0x123",
			data:      "0x",
			amount:    big.NewInt(1),
			gasLimit:  21000,
			wantErr:   true,
		},
		{
			name:      "negative value",
			toAddress: "0x742d35cC6600C91D844c4c0C7A7D26b2b39d3c07",
			data:      "0x",
			amount:    big.NewInt(-1),
			gasLimit:  21000,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		fromAddress := "0x742d35cC6600C91D844c4c0C7A7D26b2b39d3c07"
		t.Run(tt.name, func(t *testing.T) {
			var (
				rpcURL           = "https://arb1.arbitrum.io/rpc"
				delayedInboxAddr = "0x4Fb6e0c9c2c8e0b2e0b2e0b2e0b2e0b2e0b2e0b2" // dummy address
			)
			a, err := l2.NewArbitrumClient(rpcURL, delayedInboxAddr)
			require.NoError(t, err)

			got, err := a.BuildForceInclusionTx(fromAddress, tt.toAddress, tt.data, tt.amount, tt.gasLimit)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
