package notifees

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsElrondNotifee() ArgsElrondNotifee {
	return ArgsElrondNotifee{
		Proxy:           &testsCommon.ProxyStub{},
		TxBuilder:       &testsCommon.TxBuilderStub{},
		TxNonceHandler:  &testsCommon.TxNonceHandlerStub{},
		ContractAddress: data.NewAddressFromBytes(bytes.Repeat([]byte{1}, 32)),
		PrivateKey:      &testsCommon.PrivateKeyStub{},
		BaseGasLimit:    1,
		GasLimitForEach: 1,
	}
}

func createMockArgsElrondNotifeeWithSomeRealComponents() ArgsElrondNotifee {
	proxy := &testsCommon.ProxyStub{
		GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				ChainID:     "test",
				MinGasLimit: 1000,
				MinGasPrice: 10,
			}, nil
		},
	}

	txBuilder, _ := builders.NewTxBuilder(blockchain.NewTxSigner())
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	skBytes, _ := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	sk, _ := keyGen.PrivateKeyFromByteArray(skBytes)

	return ArgsElrondNotifee{
		Proxy:           proxy,
		TxBuilder:       txBuilder,
		TxNonceHandler:  &testsCommon.TxNonceHandlerStub{},
		ContractAddress: data.NewAddressFromBytes(bytes.Repeat([]byte{1}, 32)),
		PrivateKey:      sk,
		BaseGasLimit:    2000,
		GasLimitForEach: 30,
	}
}

func createMockPriceChanges() []*aggregator.ArgsPriceChanged {
	return []*aggregator.ArgsPriceChanged{
		{
			Base:               "USD",
			Quote:              "ETH",
			DenominatedPrice:   380000,
			DenominationFactor: 100,
		},
		{
			Base:               "USD",
			Quote:              "BTC",
			DenominatedPrice:   47000000000,
			DenominationFactor: 1000000,
		},
	}
}

func TestNewElrondNotifee(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.Proxy = nil
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilProxy, err)
	})
	t.Run("nil tx builder should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.TxBuilder = nil
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilTxBuilder, err)
	})
	t.Run("nil tx nonce handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.TxNonceHandler = nil
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilTxNonceHandler, err)
	})
	t.Run("nil contract address should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.ContractAddress = nil
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilContractAddressHandler, err)
	})
	t.Run("invalid contract address should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.ContractAddress = data.NewAddressFromBytes([]byte("invalid"))
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidContractAddress, err)
	})
	t.Run("nil private key should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.PrivateKey = nil
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errNilPrivateKey, err)
	})
	t.Run("invalid base gas limit should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.BaseGasLimit = minGasLimit - 1
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidBaseGasLimit, err)
	})
	t.Run("invalid gas limit for each should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		args.GasLimitForEach = minGasLimit - 1
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, errInvalidGasLimitForEach, err)
	})
	t.Run("private key to byte array errors should error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifee()
		args.PrivateKey = &testsCommon.PrivateKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return nil, expectedErr
			},
		}
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, expectedErr, err)
	})
	t.Run("public key to byte array errors should error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifee()
		args.PrivateKey = &testsCommon.PrivateKeyStub{
			GeneratePublicCalled: func() crypto.PublicKey {
				return &testsCommon.PublicKeyStub{
					ToByteArrayCalled: func() ([]byte, error) {
						return nil, expectedErr
					},
				}
			},
		}
		en, err := NewElrondNotifee(args)

		assert.True(t, check.IfNil(en))
		assert.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifee()
		en, err := NewElrondNotifee(args)

		assert.False(t, check.IfNil(en))
		assert.Nil(t, err)
	})
}

func TestElrondNotifee_PriceChanged(t *testing.T) {
	t.Parallel()

	t.Run("get nonce errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 0, expectedErr
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("invalid price arguments", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 43, nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		priceChanges[0].Base = ""
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.True(t, errors.Is(err, builders.ErrInvalidValue))
	})
	t.Run("get network config errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.Proxy = &testsCommon.ProxyStub{
			GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
				return nil, expectedErr
			},
		}
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 43, nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("apply signature and generate transaction errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 43, nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				assert.Fail(t, "should have not called SendTransaction")
				return "", nil
			},
		}
		args.TxBuilder = &testsCommon.TxBuilderStub{
			ApplySignatureAndGenerateTxCalled: func(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error) {
				return nil, expectedErr
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("send transaction errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 43, nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				return "", expectedErr
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sentWasCalled := false
		args := createMockArgsElrondNotifeeWithSomeRealComponents()
		args.TxNonceHandler = &testsCommon.TxNonceHandlerStub{
			GetNonceCalled: func(ctx context.Context, address core.AddressHandler) (uint64, error) {
				return 43, nil
			},
			SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
				txDataStrings := []string{
					function,
					hex.EncodeToString([]byte("USD")),
					hex.EncodeToString([]byte("ETH")),
					hex.EncodeToString(big.NewInt(380000).Bytes()),
					hex.EncodeToString([]byte("USD")),
					hex.EncodeToString([]byte("BTC")),
					hex.EncodeToString(big.NewInt(47000000000).Bytes()),
				}
				txData := []byte(strings.Join(txDataStrings, "@"))

				assert.Equal(t, uint64(43), tx.Nonce)
				assert.Equal(t, "0", tx.Value)
				assert.Equal(t, "erd1qyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqsl6e0p7", tx.RcvAddr)
				assert.Equal(t, "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x", tx.SndAddr)
				assert.Equal(t, uint64(10), tx.GasPrice)
				assert.Equal(t, uint64(2060), tx.GasLimit)
				assert.Equal(t, txData, tx.Data)
				assert.Equal(t, "1bf31bd077f397dd3e8f1cb47f8211251d696b4d67003b22e638eea6e46200e8101d3e39305f3bf502a5a9838521ba61566bf7d04a7fe91e6d012a94cc6a810e", tx.Signature)
				assert.Equal(t, "test", tx.ChainID)
				assert.Equal(t, uint32(1), tx.Version)

				sentWasCalled = true

				return "hash", nil
			},
		}

		en, err := NewElrondNotifee(args)
		require.Nil(t, err)

		priceChanges := createMockPriceChanges()
		err = en.PriceChanged(context.Background(), priceChanges)
		assert.Nil(t, err)
		assert.True(t, sentWasCalled)
	})
}
