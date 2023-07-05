package nonceHandlerV2

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleTransactionAddressNonceHandler_NewSingleTransactionAddressNonceHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		anh, err := NewSingleTransactionAddressNonceHandler(nil, nil)
		assert.Nil(t, anh)
		assert.Equal(t, interactors.ErrNilProxy, err)
	})
	t.Run("nil addressHandler", func(t *testing.T) {
		t.Parallel()

		anh, err := NewSingleTransactionAddressNonceHandler(&testsCommon.ProxyStub{}, nil)
		assert.Nil(t, anh)
		assert.Equal(t, interactors.ErrNilAddress, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		pubkey := make([]byte, 32)
		_, _ = rand.Read(pubkey)
		addressHandler := data.NewAddressFromBytes(pubkey)

		_, err := NewSingleTransactionAddressNonceHandler(&testsCommon.ProxyStub{}, addressHandler)
		assert.Nil(t, err)
	})
}

func TestSingleTransactionAddressNonceHandler_ApplyNonceAndGasPrice(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error should error", func(t *testing.T) {
		t.Parallel()

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

		tx := createDefaultTx()

		err := anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

		tx := createDefaultTx()

		err := anh.ApplyNonceAndGasPrice(context.Background(), &tx)
		require.Nil(t, err)
		require.Equal(t, blockchainNonce, tx.Nonce)
	})
}

func TestSingleTransactionAddressNonceHandler_fetchGasPriceIfRequired(t *testing.T) {
	t.Parallel()

	//proxy returns error should set invalid gasPrice(0)
	proxy := &testsCommon.ProxyStub{
		GetNetworkConfigCalled: func() (*data.NetworkConfig, error) {
			return nil, expectedErr
		},
	}
	anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
	anh.gasPrice = 100000
	anh.nonceUntilGasIncreased = 100

	anh.fetchGasPriceIfRequired(context.Background(), 101)
	require.Equal(t, uint64(0), anh.gasPrice)
}

func TestSingleTransactionAddressNonceHandler_DropTransactions(t *testing.T) {
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

	anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

	err := anh.ApplyNonceAndGasPrice(context.Background(), &tx)
	require.Nil(t, err)

	_, err = anh.SendTransaction(context.Background(), &tx)
	require.Nil(t, err)

	require.Equal(t, uint64(0), anh.nonceUntilGasIncreased)
	require.Equal(t, minGasPrice, anh.gasPrice)
	require.NotNil(t, anh.transaction)

	anh.DropTransactions()

	require.Equal(t, blockchainNonce, anh.nonceUntilGasIncreased)
	require.Equal(t, minGasPrice+1, anh.gasPrice)
	require.Nil(t, anh.transaction)
}

func TestSingleTransactionAddressNonceHandler_SendTransaction(t *testing.T) {
	t.Parallel()

	// proxy returns error should error
	proxy := &testsCommon.ProxyStub{
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			return "", expectedErr
		},
	}

	anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

	tx := createDefaultTx()

	_, err := anh.SendTransaction(context.Background(), &tx)
	require.Equal(t, expectedErr, err)
}

func TestSingleTransactionAddressNonceHandler_ReSendTransactionsIfRequired(t *testing.T) {
	t.Parallel()

	t.Run("no transaction to resend shall exit early with no error", func(t *testing.T) {
		t.Parallel()

		anh, _ := NewSingleTransactionAddressNonceHandler(&testsCommon.ProxyStub{}, testAddress)
		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
	})
	t.Run("proxy returns error shall error", func(t *testing.T) {
		t.Parallel()

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		tx := createDefaultTx()

		_, err := anh.SendTransaction(context.Background(), &tx)
		require.Nil(t, err)

		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, expectedErr, err)
	})
	t.Run("proxy returns error shall error", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
			SendTransactionCalled: func(txs *transaction.FrontendTransaction) (string, error) {
				return "", expectedErr
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce
		anh.transaction = &tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, expectedErr, err)
	})
	t.Run("anh.transaction.Nonce != account.Nonce", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce + 1}, nil
			},
			SendTransactionCalled: func(txs *transaction.FrontendTransaction) (string, error) {
				return "", expectedErr
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce
		anh.transaction = &tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
		require.Nil(t, anh.transaction)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
			SendTransactionCalled: func(txs *transaction.FrontendTransaction) (string, error) {
				return "hash", nil
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		tx := createDefaultTx()
		tx.Nonce = blockchainNonce
		anh.transaction = &tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
	})
}
