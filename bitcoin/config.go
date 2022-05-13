package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
)

func init() {
	if err := chaincfg.Register(&EunoMainnetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&EunoTestnetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&EunoRegressionNetParams); err != nil {
		panic(err)
	}
}

// MainNetParams returns the chain configuration for mainnet
var EunoMainnetParams = chaincfg.Params{
	Net:         0xe9fdc490,
	DefaultPort: "46462",

	PubKeyHashAddrID: 0x21,
	ScriptHashAddrID: 0x11,
	PrivateKeyID:     0x9,

	HDPrivateKeyID: [4]byte{0x02, 0x21, 0x31, 0x2B},
	HDPublicKeyID:  [4]byte{0x02, 0x2D, 0x25, 0x33},

	HDCoinType: 119,
}

// TestnetParams returns the chain configuration for testnet
var EunoTestnetParams = chaincfg.Params{
	Net:         0xba657645,
	DefaultPort: "46464",

	PubKeyHashAddrID: 0x8B,
	ScriptHashAddrID: 0x13,
	PrivateKeyID:     0xEF,

	HDPrivateKeyID: [4]byte{0x3a, 0x80, 0x58, 0x37},
	HDPublicKeyID:  [4]byte{0x3a, 0x80, 0x61, 0xa0},

	HDCoinType: 1,
}

// RegressionNetParams returns the chain configuration for regression net
var EunoRegressionNetParams = chaincfg.Params{
	Net:         0xac7ecfa1,
	DefaultPort: "46466",
}
