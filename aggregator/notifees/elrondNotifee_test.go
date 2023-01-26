package notifees

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/aggregator"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

func createMockArgsMultiversXNotifee() ArgsMultiversXNotifee {
	return ArgsMultiversXNotifee{
		Proxy:           &testsCommon.ProxyStub{},
		TxBuilder:       &testsCommon.TxBuilderStub{},
		TxNonceHandler:  &testsCommon.TxNonceHandlerV2Stub{},
		ContractAddress: data.NewAddressFromBytes(bytes.Repeat([]byte{1}, 32)),
		CryptoHolder:    &testsCommon.CryptoComponentsHolderStub{},
		BaseGasLimit:    1,
		GasLimitForEach: 1,
	}
}

func createMockArgsMultiversXNotifeeWithSomeRealComponents() ArgsMultiversXNotifee {
	proxy := &testsCommon.ProxyStub{
		GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				ChainID:     "test",
				MinGasLimit: 1000,
				MinGasPrice: 10,
			}, nil
		},
	}

	skBytes, _ := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, skBytes)
	txBuilder, _ := builders.NewTxBuilder(cryptoProvider.NewSigner())

	return ArgsMultiversXNotifee{
		Proxy:           proxy,
		TxBuilder:       txBuilder,
		TxNonceHandler:  &testsCommon.TxNonceHandlerV2Stub{},
		ContractAddress: data.NewAddressFromBytes(bytes.Repeat([]byte{1}, 32)),
		CryptoHolder:    holder,
		BaseGasLimit:    2000,
		GasLimitForEach: 30,
	}
}

func createMockPriceChanges() []*aggregator.ArgsPriceChanged {
	return []*aggregator.ArgsPriceChanged{
		{
			Base:             "USD",
			Quote:            "ETH",
			DenominatedPrice: 380000,
			Decimals:         2,
			Timestamp:        200,
		},
		{
			Base:             "USD",
			Quote:            "BTC",
			DenominatedPrice: 47000000000,
			Decimals:         6,
			Timestamp:        300,
		},
	}
}

func TestNewMultiversXNotifee(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.Proxy = nil
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilProxy, err)
	})
	t.Run("nil tx builder should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.TxBuilder = nil
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilTxBuilder, err)
	})
	t.Run("nil tx nonce handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.TxNonceHandler = nil
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilTxNonceHandler, err)
	})
	t.Run("nil contract address should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.ContractAddress = nil
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilContractAddressHandler, err)
	})
	t.Run("invalid contract address should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.ContractAddress = data.NewAddressFromBytes([]byte("invalid"))
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidContractAddress, err)
	})
	t.Run("nil cryptoHlder should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.CryptoHolder = nil
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, builders.ErrNilCryptoComponentsHolder, err)
	})
	t.Run("invalid base gas limit should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.BaseGasLimit = minGasLimit - 1
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidBaseGasLimit, err)
	})
	t.Run("invalid gas limit for each should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		args.GasLimitForEach = minGasLimit - 1
		en, err := NewMultiversXNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidGasLimitForEach, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifee()
		en, err := NewMultiversXNotifee(args)

		assert.False(t, check.IfNil(en))
		assert.Nil(t, err)
	})
}

func TestMultiversXNotifee_PriceChanged(t *testing.T) {
	t.Parallel()

	t.Run("get nonce errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				return expectedErr
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("invalid price arguments", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				txArgs.Nonce = 43
				return nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		priceChanges[0].Base = ""
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.True(t, errors.Is(err, builders.ErrInvalidValue))
	})
	t.Run("get network config errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.Proxy = &testsCommon.ProxyStub{
			GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
				return nil, expectedErr
			},
		}
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				txArgs.Nonce = 43
				return nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("apply signature and generate transaction errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				txArgs.Nonce = 43
				return nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}
		args.TxBuilder = &testsCommon.TxBuilderStub{
			ApplySignatureAndGenerateTxCalled: func(cryptoHolder core.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error) {
				return nil, expectedErr
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("send transaction errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				txArgs.Nonce = 43
				return nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				return "", expectedErr
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		priceChanges := createMockPriceChanges()
		sentWasCalled := false
		args := createMockArgsMultiversXNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerV2Stub{
			ApplyNonceAndGasPriceCalled: func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
				txArgs.Nonce = 43
				return nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				txDataStrings := []string{
					function,
					hex.EncodeToString([]byte(priceChanges[0].Base)),
					hex.EncodeToString([]byte(priceChanges[0].Quote)),
					hex.EncodeToString(big.NewInt(priceChanges[0].Timestamp).Bytes()),
					hex.EncodeToString(big.NewInt(int64(priceChanges[0].DenominatedPrice)).Bytes()),
					hex.EncodeToString(big.NewInt(int64(priceChanges[0].Decimals)).Bytes()),
					hex.EncodeToString([]byte(priceChanges[1].Base)),
					hex.EncodeToString([]byte(priceChanges[1].Quote)),
					hex.EncodeToString(big.NewInt(priceChanges[1].Timestamp).Bytes()),
					hex.EncodeToString(big.NewInt(int64(priceChanges[1].DenominatedPrice)).Bytes()),
					hex.EncodeToString(big.NewInt(int64(priceChanges[1].Decimals)).Bytes()),
				}
				txData := []byte(strings.Join(txDataStrings, "@"))

				assert.Equal(t, uint64(43), tx.Nonce)
				assert.Equal(t, "0", tx.Value)
				assert.Equal(t, "erd1qyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqsl6e0p7", tx.RcvAddr)
				assert.Equal(t, "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x", tx.SndAddr)
				assert.Equal(t, uint64(10), tx.GasPrice)
				assert.Equal(t, uint64(2060), tx.GasLimit)
				assert.Equal(t, txData, tx.Data)
				assert.Equal(t, "test", tx.ChainID)
				assert.Equal(t, uint32(1), tx.Version)
				assert.Equal(t, uint32(0), tx.Options)

				sentWasCalled = true

				return "hash", nil
			},
		}

		en, err := NewMultiversXNotifee(args)
		require.Nil(t, err)

		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Nil(t, err)
		assert.True(t, sentWasCalled)
	})
}
