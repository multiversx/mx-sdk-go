package nonceHandlerV2

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAddressAsBech32String = "erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht"
var testAddress, _ = data.NewAddressFromBech32String(testAddressAsBech32String)
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

		tx := createDefaultTx()

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), &tx)
		require.Nil(t, err)

		anh.gasPrice = tx.GasPrice
		err = anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Equal(t, interactors.ErrTxAlreadySent, err)
	})
	t.Run("tx already sent; oldTx.GasPrice < txArgs.GasPrice", func(t *testing.T) {
		t.Parallel()

		tx := createDefaultTx()
		initialGasPrice := tx.GasPrice
		tx.GasPrice--

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), &tx)
		require.Nil(t, err)

		anh.gasPrice = initialGasPrice
		err = anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Nil(t, err)
	})
	t.Run("oldTx.GasPrice == txArgs.GasPrice && oldTx.GasPrice < anh.gasPrice", func(t *testing.T) {
		t.Parallel()

		tx := createDefaultTx()
		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), &tx)
		require.Nil(t, err)

		anh.gasPrice = tx.GasPrice + 1
		err = anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Nil(t, err)
	})
	t.Run("same transaction but with different nonce should work", func(t *testing.T) {
		t.Parallel()

		tx1 := createDefaultTx()
		tx2 := createDefaultTx()
		tx2.Nonce++

		anh, err := NewAddressNonceHandlerWithPrivateAccess(&testsCommon.ProxyStub{}, testAddress)
		require.Nil(t, err)

		_, err = anh.SendTransaction(context.Background(), &tx1)
		require.Nil(t, err)

		err = anh.ApplyNonceAndGasPrice(context.Background(), &tx2)
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
		tx := createDefaultTx()

		err := anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Equal(t, expectedErr, err)
	})
}

func TestAddressNonceHandler_DropTransactions(t *testing.T) {
	t.Parallel()

	tx := createDefaultTx()

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

	err := anh.ApplyNonceAndGasPrice(context.Background(), &tx)
	require.Nil(t, err)

	_, err = anh.SendTransaction(context.Background(), &tx)
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
			SendTransactionsCalled: func(txs []*transaction.FrontendTransaction) ([]string, error) {
				return make([]string, 0), expectedErr
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce
		_, err := anh.SendTransaction(context.Background(), &tx)
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
		tx := createDefaultTx()
		_, err := anh.SendTransaction(context.Background(), &tx)
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
		tx := createDefaultTx()
		_, err := anh.SendTransaction(context.Background(), &tx)
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
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce + 1
		_, err := anh.SendTransaction(context.Background(), &tx)
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
			SendTransactionsCalled: func(txs []*transaction.FrontendTransaction) ([]string, error) {
				return make([]string, 0), nil
			},
		}
		anh, _ := NewAddressNonceHandlerWithPrivateAccess(proxy, testAddress)
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce
		_, err := anh.SendTransaction(context.Background(), &tx)
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

func createDefaultTx() transaction.FrontendTransaction {
	return transaction.FrontendTransaction{
		Value:    "1",
		Receiver: testAddressAsBech32String,
		Sender:   testAddressAsBech32String,
		GasPrice: 100000,
		GasLimit: 50000,
		Data:     nil,
		ChainID:  "3",
		Version:  1,
	}
}
