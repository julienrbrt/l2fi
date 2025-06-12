package l2

import "math/big"

type L2 interface {
	BuildForceInclusionTx(
		fromAddress, toAddress, data string,
		value *big.Int,
		l2GasLimit uint64,
	) (string, error)
}
