package builders

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func TestRelayedTxV2Builder(t *testing.T) {
	t.Parallel()

	netConfig := &data.NetworkConfig{
		ChainID:               "T",
		MinTransactionVersion: 1,
		GasPerDataByte:        1500,
	}

	relayerAcc, relayerPrivKey := getAccount(t, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, testInnerSenderMnemonic)

	innerTx := &data.Transaction{
		Nonce:    innerSenderAcc.Nonce,
		Value:    "0",
		RcvAddr:  "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		SndAddr:  innerSenderAcc.Address,
		GasPrice: netConfig.MinGasPrice,
		GasLimit: 0,
		Data:     []byte("getContractConfig"),
		ChainID:  netConfig.ChainID,
		Version:  netConfig.MinTransactionVersion,
		Options:  0,
	}

	innerTxSig := signTx(t, innerSenderPrivKey, innerTx)
	innerTx.Signature = hex.EncodeToString(innerTxSig)

	txJson, _ := json.Marshal(innerTx)
	require.Equal(t,
		`{"nonce":37,"value":"0","receiver":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u","sender":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","gasPrice":0,"gasLimit":0,"data":"Z2V0Q29udHJhY3RDb25maWc=","signature":"e578eb15076d9bf42f7610c17125ef88d38422a927e80075d11790b618ec106d4dec7da0d312801e0df28f5c8885d648844f54e7dc83c3e7fd71f60998867103","chainID":"T","version":1}`,
		string(txJson),
	)

	relayedV2Builder := NewRelayedTxV2Builder()
	relayedV2Builder.SetInnerTransaction(innerTx)
	relayedV2Builder.SetRelayerAccount(relayerAcc)
	relayedV2Builder.SetNetworkConfig(netConfig)
	relayedV2Builder.SetGasLimitNeededForInnerTransaction(60_000_000)

	relayedTx, err := relayedV2Builder.Build()
	require.NoError(t, err)

	relayedTx.GasPrice = netConfig.MinGasPrice

	relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
	relayedTx.Signature = hex.EncodeToString(relayedTxSig)

	txJson, _ = json.Marshal(relayedTx)
	require.Equal(t,
		`{"nonce":37,"value":"0","receiver":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","sender":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5","gasPrice":0,"gasLimit":60364500,"data":"cmVsYXllZFR4VjJAMDAwMDAwMDAwMDAwMDAwMDAwMDEwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAyZmZmZkAyNUA2NzY1NzQ0MzZmNmU3NDcyNjE2Mzc0NDM2ZjZlNjY2OTY3QGU1NzhlYjE1MDc2ZDliZjQyZjc2MTBjMTcxMjVlZjg4ZDM4NDIyYTkyN2U4MDA3NWQxMTc5MGI2MThlYzEwNmQ0ZGVjN2RhMGQzMTI4MDFlMGRmMjhmNWM4ODg1ZDY0ODg0NGY1NGU3ZGM4M2MzZTdmZDcxZjYwOTk4ODY3MTAz","signature":"6c4f404274d499857be292c59bfca9c29f1ea5ead9544346fadf9ab2944310775aa5c0954cc289d916c6c8cadec128614903f73b2bcd6c0266f6d651383d7207","chainID":"T","version":1}`,
		string(txJson),
	)
}
