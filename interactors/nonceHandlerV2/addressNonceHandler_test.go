package nonceHandlerV2

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAddress, _ = data.NewAddressFromBech32String("erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht")
var expectedErr = errors.New("expected error")

func TestAddressNonceHandler_NewAddressNonceHandlerWithPrivateAccess(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		anh, err := NewAddressNonceHandlerWithPrivateAccess(nil, nil)
		assert.Nil(t, anh)
		assert.Equal(t, interactors.ErrNilProxy, err)
	})
	t.Run("nil addressHandler", func(t *testing.T) {
		t.Parallel()

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, nil)
		assert.Nil(t, anh)
		assert.Equal(t, interactors.ErrNilAddress, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		pubkey := make([]byte, 32)
		_, _ = rand.Read(pubkey)
		addressHandler := data.NewAddressFromBytes(pubkey)

		_, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, addressHandler)
		assert.Nil(t, err)
	})
}

func TestAddressNonceHandler_ApplyNonceAndGasPrice(t *testing.T) {
	t.Parallel()
	t.Run("tx already sent; oldTx.GasPrice == txArgs.GasPrice == anh.gasPrice", func(t *testing.T) {
		t.Parallel()

		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)

		anh.gasPrice = txArgs.GasPrice
		err = anh.ApplyNonceAndGasPrice(context.Background(), &txArgs)
		require.Equal(t, interactors.ErrTxAlreadySent, err)
	})
	t.Run("tx already sent; oldTx.GasPrice < txArgs.GasPrice", func(t *testing.T) {
		t.Parallel()

		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice-1, txArgs)

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)

		anh.gasPrice = txArgs.GasPrice
		err = anh.ApplyNonceAndGasPrice(context.Background(), &txArgs)
		require.Nil(t, err)
	})
	t.Run("oldTx.GasPrice == txArgs.GasPrice && oldTx.GasPrice < anh.gasPrice", func(t *testing.T) {
		t.Parallel()

		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)

		anh.gasPrice = txArgs.GasPrice + 1
		err = anh.ApplyNonceAndGasPrice(context.Background(), &txArgs)
		require.Nil(t, err)
	})
}

func TestAddressNonceHandler_getNonceUpdatingCurrent(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error shall return error", func(t *testing.T) {
		t.Parallel()

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		nonce, err := anh.getNonceUpdatingCurrent(context.Background())
		require.Equal(t, expectedErr, err)
		require.Equal(t, uint64(0), nonce)
	})
	t.Run("gap nonce detected", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
		}

		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		anh.lowestNonce = blockchainNonce + 1

		nonce, err := anh.getNonceUpdatingCurrent(context.Background())
		require.Equal(t, interactors.ErrGapNonce, err)
		require.Equal(t, nonce, blockchainNonce)
	})
	t.Run("when computedNonce already set, getNonceUpdatingCurrent shall increase it", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
		}

		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		anh.computedNonceWasSet = true
		computedNonce := uint64(105)
		anh.computedNonce = computedNonce

		nonce, err := anh.getNonceUpdatingCurrent(context.Background())
		require.Nil(t, err)
		require.Equal(t, nonce, computedNonce+1)
	})
	t.Run("getNonceUpdatingCurrent returns error should error", func(t *testing.T) {
		t.Parallel()

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		txArgs := createTxArgs()

		err := anh.ApplyNonceAndGasPrice(context.Background(), &txArgs)
		require.Equal(t, expectedErr, err)
	})
}

func TestAddressNonceHandler_DropTransactions(t *testing.T) {
	t.Parallel()

	txArgs := createTxArgs()

	blockchainNonce := uint64(100)
	minGasPrice := uint64(10)
	proxy := &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			return &data.Account{Nonce: blockchainNonce}, nil
		},
		GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
			return &data.NetworkConfig{MinGasPrice: minGasPrice}, nil
		},
	}

	anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)

	err := anh.ApplyNonceAndGasPrice(context.Background(), &txArgs)
	require.Nil(t, err)

	tx := createTx(txArgs.GasPrice, txArgs)
	_, err = anh.SendTransaction(context.Background(), tx)
	require.Nil(t, err)

	require.True(t, anh.computedNonceWasSet)
	require.Equal(t, blockchainNonce, anh.computedNonce)
	require.Equal(t, uint64(0), anh.nonceUntilGasIncreased)
	require.Equal(t, minGasPrice, anh.gasPrice)
	require.Equal(t, 1, len(anh.transactions))

	anh.DropTransactions()

	require.False(t, anh.computedNonceWasSet)
	require.Equal(t, blockchainNonce, anh.nonceUntilGasIncreased)
	require.Equal(t, minGasPrice+1, anh.gasPrice)
	require.Equal(t, 0, len(anh.transactions))
}

func TestAddressNonceHandler_ReSendTransactionsIfRequired(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error shall error", func(t *testing.T) {
		t.Parallel()

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, expectedErr, err)
	})
	t.Run("proxy returns error shall error", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce - 1}, nil
			},
			SendTransactionsCalled: func(txs []*data.Transaction) ([]string, error) {
				return make([]string, 0), expectedErr
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce
		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		require.Equal(t, 1, len(anh.transactions))

		anh.computedNonce = blockchainNonce

		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, 1, len(anh.transactions))
		require.Equal(t, expectedErr, err)
	})
	t.Run("account.Nonce == anh.computedNonce", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		require.Equal(t, 1, len(anh.transactions))

		anh.computedNonce = blockchainNonce
		anh.lowestNonce = 80
		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, anh.computedNonce, anh.lowestNonce)
		require.Equal(t, 0, len(anh.transactions))
		require.Nil(t, err)
	})
	t.Run("len(anh.transactions) == 0", func(t *testing.T) {
		t.Parallel()

		anh, _ := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		require.Equal(t, 1, len(anh.transactions))

		anh.computedNonce = 100
		anh.lowestNonce = 80
		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, anh.computedNonce, anh.lowestNonce)
		require.Equal(t, 0, len(anh.transactions))
		require.Nil(t, err)
	})
	t.Run("lowestNonce should be recalculated each time", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce - 1}, nil
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce + 1
		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		require.Equal(t, 1, len(anh.transactions))

		anh.computedNonce = blockchainNonce + 2
		anh.lowestNonce = blockchainNonce
		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, blockchainNonce+1, anh.lowestNonce)
		require.Equal(t, 1, len(anh.transactions))
		require.Nil(t, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce - 1}, nil
			},
			SendTransactionsCalled: func(txs []*data.Transaction) ([]string, error) {
				return make([]string, 0), nil
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce
		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		require.Equal(t, 1, len(anh.transactions))

		anh.computedNonce = blockchainNonce

		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, 1, len(anh.transactions))
		require.Nil(t, err)
	})
}

func TestAddressNonceHandler_fetchGasPriceIfRequired(t *testing.T) {
	t.Parallel()

	// proxy returns error should set invalid gasPrice(0)
	proxy := &testsCommon.ProxyStub{
		GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
			return nil, expectedErr
		},
	}
	anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
	anh.gasPrice = 100000
	anh.nonceUntilGasIncreased = 100

	anh.fetchGasPriceIfRequired(context.Background(), 101)
	require.Equal(t, uint64(0), anh.gasPrice)
}

func createTxArgs() data.ArgCreateTransaction {
	return data.ArgCreateTransaction{
		Value:    "1",
		RcvAddr:  testAddress.AddressAsBech32String(),
		SndAddr:  testAddress.AddressAsBech32String(),
		GasPrice: 100000,
		GasLimit: 50000,
		Data:     nil,
		ChainID:  "3",
		Version:  1,
	}
}

func createTx(gasPrice uint64, txArgs data.ArgCreateTransaction) *data.Transaction {
	return &data.Transaction{
		Nonce:     txArgs.Nonce,
		Value:     txArgs.Value,
		RcvAddr:   txArgs.RcvAddr,
		SndAddr:   txArgs.SndAddr,
		GasPrice:  gasPrice,
		GasLimit:  txArgs.GasLimit,
		Data:      txArgs.Data,
		Signature: "sig",
		ChainID:   txArgs.ChainID,
		Version:   txArgs.Version,
	}
}
