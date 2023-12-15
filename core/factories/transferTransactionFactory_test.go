package factories

import (
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var (
	Dconf = &core.Config{
		ChainID:                      "D",
		MinGasLimit:                  50_000,
		GasLimitPerByte:              1_500,
		GasLimitESDTTransfer:         200_000,
		GasLimitESDTNFTTransfer:      200_000,
		GasLimitMultiESDTNFTTransfer: 200_000,
	}
)

func TestTransferTransactionFactory_CreateTransactionForNativeTokenTransferNoData(t *testing.T) {
	alice, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create address from bech32")

	bob, err := core.NewAddressFromBech32("erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx")
	require.NoError(t, err, "failed to create address from bech32")

	ttf := NewTransferTransactionFactory(Dconf, core.NewTokenComputer())
	transaction, err := ttf.CreateTransactionForNativeTokenTransfer(alice, bob, 1000000000000000000, "")
	require.NoError(t, err, "failed to create transaction")

	expectedTransaction := &core.Transaction{
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
		Value:    big.NewInt(1000000000000000000),
		ChainID:  "D",
		GasLimit: 50_000,
		Data:     []byte(""),
	}

	require.Equal(t, expectedTransaction, transaction)
}

func TestTransferTransactionFactory_CreateTransactionForNativeTokenTransferWithData(t *testing.T) {
	alice, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create address from bech32")

	bob, err := core.NewAddressFromBech32("erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx")
	require.NoError(t, err, "failed to create address from bech32")

	ttf := NewTransferTransactionFactory(Dconf, core.NewTokenComputer())
	transaction, err := ttf.CreateTransactionForNativeTokenTransfer(alice, bob, 1000000000000000000,
		"test data")
	require.NoError(t, err, "failed to create transaction")

	expectedTransaction := &core.Transaction{
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
		Value:    big.NewInt(1000000000000000000),
		ChainID:  "D",
		GasLimit: 63_500,
		Data:     []byte("test data"),
	}

	require.Equal(t, expectedTransaction, transaction)
}

func TestTransferTransactionFactory_CreateTransactionForESDTTokenTransfer(t *testing.T) {
	alice, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create address from bech32")

	bob, err := core.NewAddressFromBech32("erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx")
	require.NoError(t, err, "failed to create address from bech32")

	fooToken := core.Token{Identifier: "FOO-123456"}
	tt := core.TokenTransfer{Token: fooToken, Amount: big.NewInt(1000000)}

	ttf := NewTransferTransactionFactory(Dconf, core.NewTokenComputer())
	transaction, err := ttf.CreateTransactionForESDTTokenTransfer(alice, bob, []*core.TokenTransfer{&tt})
	require.NoError(t, err, "failed to create transaction")

	expectedTransaction := &core.Transaction{
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
		Value:    big.NewInt(0),
		ChainID:  "D",
		GasLimit: 410_000,
		Data:     []byte("ESDTTransfer@464f4f2d313233343536@0f4240"),
	}

	require.Equal(t, expectedTransaction, transaction)
}

func TestTransferTransactionFactory_CreateTransactionForNFTTransfer(t *testing.T) {
	alice, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create address from bech32")

	bob, err := core.NewAddressFromBech32("erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx")
	require.NoError(t, err, "failed to create address from bech32")

	nft := core.Token{Identifier: "NFT-123456", Nonce: 10}
	tt := core.TokenTransfer{Token: nft, Amount: big.NewInt(1)}

	ttf := NewTransferTransactionFactory(Dconf, core.NewTokenComputer())
	transaction, err := ttf.CreateTransactionForESDTTokenTransfer(alice, bob, []*core.TokenTransfer{&tt})
	require.NoError(t, err, "failed to create transaction")

	expectedTransaction := &core.Transaction{
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
		Value:    big.NewInt(0),
		ChainID:  "D",
		GasLimit: 1_210_500,
		Data:     []byte("ESDTNFTTransfer@4e46542d313233343536@0a@01@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8"),
	}

	require.Equal(t, expectedTransaction, transaction)
}

func TestTransferTransactionFactory_CreateTransactionForMultipleNFTTransfers(t *testing.T) {
	alice, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create address from bech32")

	bob, err := core.NewAddressFromBech32("erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx")
	require.NoError(t, err, "failed to create address from bech32")

	firstNFT := core.Token{Identifier: "NFT-123456", Nonce: 10}
	firstTransfer := core.TokenTransfer{Token: firstNFT, Amount: big.NewInt(1)}

	secondNFT := core.Token{Identifier: "TEST-987654", Nonce: 1}
	secondTransfer := core.TokenTransfer{Token: secondNFT, Amount: big.NewInt(1)}

	ttf := NewTransferTransactionFactory(Dconf, core.NewTokenComputer())
	transaction, err := ttf.CreateTransactionForESDTTokenTransfer(alice, bob, []*core.TokenTransfer{&firstTransfer,
		&secondTransfer})
	require.NoError(t, err, "failed to create transaction")

	expectedTransaction := &core.Transaction{
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Value:    big.NewInt(0),
		ChainID:  "D",
		GasLimit: 1_466_000,
		Data:     []byte("MultiESDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@02@4e46542d313233343536@0a@01@544553542d393837363534@01@01"),
	}

	require.Equal(t, expectedTransaction, transaction)
}
