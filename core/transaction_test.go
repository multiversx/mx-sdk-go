package core

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestTransactionComputer_ComputeBytesForSigning(t *testing.T) {
	sender := "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th"
	receiver := "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx"

	tx := Transaction{
		Nonce:    90,
		Sender:   sender,
		Receiver: receiver,
		Value:    big.NewInt(0),
		GasLimit: 50000,
		GasPrice: 1000000000,
		ChainID:  "D",
		Version:  1,
	}

	tc := NewTransactionComputer()
	serializedTx, err := tc.ComputeBytesForSigning(tx)
	require.NoError(t, err, "failed to serialize transaction")

	fmt.Println(string(serializedTx))

	expectedResults := `{"sender":"erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th","receiver":"erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx","gas_limit":50000,"chain_id":"D","nonce":90,"value":0,"gas_price":1000000000,"version":1}`
	require.Equal(t, expectedResults, string(serializedTx))
}

func TestComputeTransactionHash(t *testing.T) {
	sig, err := hex.DecodeString("eaa9e4dfbd21695d9511e9754bde13e90c5cfb21748a339a79be11f744c71872e9fe8e73c6035c413f5f08eef09e5458e9ea6fc315ff4da0ab6d000b450b2a07")
	require.NoError(t, err, "failed to decode signature")
	tx := Transaction{
		Nonce:     17243,
		Value:     big.NewInt(1000000000000),
		Sender:    "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver:  "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		GasPrice:  1000000000,
		GasLimit:  100000,
		Data:      []byte("testtx"),
		ChainID:   "D",
		Version:   uint32(2),
		Signature: sig,
	}

	tc := NewTransactionComputer()
	hash, err := tc.ComputeTransactionHash(tx)
	require.NoError(t, err, "failed to compute transaction hash")
	require.Equal(t, "169b76b752b220a76a93aeebc462a1192db1dc2ec9d17e6b4d7b0dcc91792f03", hex.EncodeToString(hash))
}
