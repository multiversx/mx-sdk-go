package builders

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func TestRelayedTxV3Builder(t *testing.T) {
	t.Parallel()

	netConfig := &data.NetworkConfig{
		ChainID:               "T",
		MinTransactionVersion: 1,
		GasPerDataByte:        1500,
		MinGasLimit:           50000,
		MinGasPrice:           1000000000,
	}

	relayerAcc, relayerPrivKey := getAccount(t, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, testInnerSenderMnemonic)

	innerTx := &transaction.FrontendTransaction{
		Nonce:    innerSenderAcc.Nonce,
		Value:    "100000000",
		Receiver: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Sender:   innerSenderAcc.Address,
		GasPrice: netConfig.MinGasPrice,
		GasLimit: netConfig.MinGasLimit,
		ChainID:  netConfig.ChainID,
		Version:  netConfig.MinTransactionVersion,
		Relayer:  relayerAcc.Address,
	}

	innerTxSig := signTx(t, innerSenderPrivKey, innerTx)
	innerTx.Signature = hex.EncodeToString(innerTxSig)

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransaction(innerTx)
		relayedV3Builder.SetRelayerAccount(relayerAcc)
		relayedV3Builder.SetNetworkConfig(netConfig)

		relayedTx, err := relayedV3Builder.Build()
		require.NoError(t, err)

		relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
		relayedTx.Signature = hex.EncodeToString(relayedTxSig)

		txJson, _ := json.Marshal(relayedTx)
		require.Equal(t,
			`{"nonce":37,"value":"0","receiver":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","sender":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5",`+
				`"gasPrice":1000000000,"gasLimit":100000,"signature":"4de727262dcdb8118d9df83cd1376ff0c920e4b8cd0e26f5041dca543aa8522def65e2f84ea6e85ee4a72a585749fd7896d879258a57792599ef81066b432c00",`+
				`"chainID":"T","version":1,"innerTransaction":{"nonce":37,"value":"100000000","receiver":"erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",`+
				`"sender":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","gasPrice":1000000000,"gasLimit":50000,`+
				`"signature":"907f6dc73f2218c91180be9b027a513e92f669c36bc26300f90f5bf9d7729328eefd7e098e6fbcf34f6b3b13466ee6fb8918f956ca8efd543a8d2bdffc9d680f","chainID":"T","version":1,`+
				`"relayer":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5"}}`,
			string(txJson),
		)
	})
	t.Run("nil inner tx should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilInnerTransaction, err)
		require.Nil(t, relayedTx)
	})
	t.Run("empty inner tx signature should error", func(t *testing.T) {
		t.Parallel()

		innerTxCopy := *innerTx
		innerTxCopy.Signature = ""
		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransaction(&innerTxCopy)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilInnerTransactionSignature, err)
		require.Nil(t, relayedTx)
	})
	t.Run("empty inner tx relayer should error", func(t *testing.T) {
		t.Parallel()

		innerTxCopy := *innerTx
		innerTxCopy.Relayer = ""
		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransaction(&innerTxCopy)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrEmptyRelayerOnInnerTransaction, err)
		require.Nil(t, relayedTx)
	})
	t.Run("nil relayer account should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransaction(innerTx)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilRelayerAccount, err)
		require.Nil(t, relayedTx)
	})
	t.Run("nil network config account should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransaction(innerTx)
		relayedV3Builder.SetRelayerAccount(relayerAcc)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilNetworkConfig, err)
		require.Nil(t, relayedTx)
	})
}
