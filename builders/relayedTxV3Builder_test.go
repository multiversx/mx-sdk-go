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

	innerTxs := []*transaction.FrontendTransaction{innerTx}

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		innerTxsCopy := copySlice(innerTxs)
		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions(innerTxsCopy)
		relayedV3Builder.SetRelayerAccount(relayerAcc)
		relayedV3Builder.SetNetworkConfig(netConfig)

		relayedTx, err := relayedV3Builder.Build()
		require.NoError(t, err)

		relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
		relayedTx.Signature = hex.EncodeToString(relayedTxSig)

		txJson, _ := json.Marshal(relayedTx)
		require.Equal(t,
			`{"nonce":37,"value":"0","receiver":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5","sender":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5",`+
				`"gasPrice":1000000000,"gasLimit":150000,"signature":"cbe00bc33a5742e613d0593879273d916dea1752314b34991c5d291e702b2e37a8b22ac4cad249cfaa62614facdd1566ac2e4736c52b06eabe7b767cfb767006",`+
				`"chainID":"T","version":1,"innerTransactions":[{"nonce":37,"value":"100000000","receiver":"erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",`+
				`"sender":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","gasPrice":1000000000,"gasLimit":50000,`+
				`"signature":"907f6dc73f2218c91180be9b027a513e92f669c36bc26300f90f5bf9d7729328eefd7e098e6fbcf34f6b3b13466ee6fb8918f956ca8efd543a8d2bdffc9d680f","chainID":"T","version":1,`+
				`"relayer":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5"}]}`,
			string(txJson),
		)
	})
	t.Run("nil inner txs should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrEmptyInnerTransactions, err)
		require.Nil(t, relayedTx)
	})
	t.Run("nil empty inner txs should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions([]*transaction.FrontendTransaction{})
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrEmptyInnerTransactions, err)
		require.Nil(t, relayedTx)
	})
	t.Run("empty inner tx signature should error", func(t *testing.T) {
		t.Parallel()

		innerTxsCopy := copySlice(innerTxs)
		innerTxsCopy[0].Signature = ""
		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions(innerTxsCopy)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilInnerTransactionSignature, err)
		require.Nil(t, relayedTx)
	})
	t.Run("empty inner tx relayer should error", func(t *testing.T) {
		t.Parallel()

		innerTxsCopy := copySlice(innerTxs)
		innerTxsCopy[0].Relayer = ""
		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions(innerTxsCopy)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrEmptyRelayerOnInnerTransaction, err)
		require.Nil(t, relayedTx)
	})
	t.Run("nil relayer account should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions(innerTxs)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilRelayerAccount, err)
		require.Nil(t, relayedTx)
	})
	t.Run("nil network config account should error", func(t *testing.T) {
		t.Parallel()

		relayedV3Builder := NewRelayedTxV3Builder()
		relayedV3Builder.SetInnerTransactions(innerTxs)
		relayedV3Builder.SetRelayerAccount(relayerAcc)
		relayedTx, err := relayedV3Builder.Build()
		require.Equal(t, ErrNilNetworkConfig, err)
		require.Nil(t, relayedTx)
	})
}

func copySlice(oldSlice []*transaction.FrontendTransaction) []*transaction.FrontendTransaction {
	newSlice := make([]*transaction.FrontendTransaction, 0, len(oldSlice))
	for _, sliceEntry := range oldSlice {
		entryCopy := *sliceEntry
		newSlice = append(newSlice, &entryCopy)
	}

	return newSlice
}
