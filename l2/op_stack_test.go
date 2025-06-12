package l2_test

import (
	"math/big"
	"testing"

	"github.com/julienrbrt/l2fi/l2"
	"github.com/stretchr/testify/require"
)

func TestOptimismForceInclusion_BuildForceInclusionTx(t *testing.T) {
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
	}

	for _, tt := range tests {
		fromAddress := "0x742d35cC6600C91D844c4c0C7A7D26b2b39d3c07"

		t.Run(tt.name, func(t *testing.T) {
			var (
				rpcURL        = "https://mainnet.optimism.io" // Example RPC URL
				portalAddress = "0x420"
			)

			o, err := l2.NewOpStackClient(rpcURL, portalAddress)
			require.NoError(t, err)

			got, err := o.BuildForceInclusionTx(fromAddress, tt.toAddress, tt.data, tt.amount, tt.gasLimit)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
