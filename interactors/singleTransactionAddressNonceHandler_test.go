package interactors

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleTransactionAddressNonceHandler_NewSingleTransactionAddressNonceHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		anh, err := NewSingleTransactionAddressNonceHandler(nil, nil)
		assert.Nil(t, anh)
		assert.Equal(t, ErrNilProxy, err)
	})
	t.Run("nil addressHandler", func(t *testing.T) {
		t.Parallel()

		anh, err := NewSingleTransactionAddressNonceHandler(&testsCommon.ProxyStub{}, nil)
		assert.Nil(t, anh)
		assert.Equal(t, ErrNilAddress, err)
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

func TestSingleTransactionAddressNonceHandler_ApplyNonce(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error should error", func(t *testing.T) {

		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

		txArgs := createTxArgs()

		err := anh.ApplyNonce(context.Background(), &txArgs)
		require.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {

		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

		txArgs := createTxArgs()

		err := anh.ApplyNonce(context.Background(), &txArgs)
		require.Nil(t, err)
		require.Equal(t, blockchainNonce, txArgs.Nonce)
	})
}

func TestSingleTransactionAddressNonceHandler_fetchGasPriceIfRequired(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error should set invalid gasPrice(0)", func(t *testing.T) {
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
	})
}

func TestSingleTransactionAddressNonceHandler_DropTransactions(t *testing.T) {
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

	anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

	err := anh.ApplyNonce(context.Background(), &txArgs)

	tx := createTx(txArgs.GasPrice, txArgs)
	_, err = anh.SendTransaction(context.Background(), tx)
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

	t.Run("proxy returns error should error", func(t *testing.T) {
		proxy := &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *data.Transaction) (string, error) {
				return "", expectedErr
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)

		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)

		_, err := anh.SendTransaction(context.Background(), tx)
		require.Equal(t, expectedErr, err)
	})
}

func TestSingleTransactionAddressNonceHandler_ReSendTransactionsIfRequired(t *testing.T) {
	t.Parallel()

	t.Run("no transaction to resend shall exit early with no error", func(t *testing.T) {
		anh, _ := NewSingleTransactionAddressNonceHandler(&testsCommon.ProxyStub{}, testAddress)
		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
	})
	t.Run("proxy returns error shall error", func(t *testing.T) {
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return nil, expectedErr
			},
		}

		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)

		_, err := anh.SendTransaction(context.Background(), tx)
		require.Nil(t, err)

		err = anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, expectedErr, err)
	})
	t.Run("proxy returns error shall error", func(t *testing.T) {
		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
			SendTransactionCalled: func(txs *data.Transaction) (string, error) {
				return "", expectedErr
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce
		anh.transaction = tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Equal(t, expectedErr, err)
	})
	t.Run("anh.transaction.Nonce != account.Nonce", func(t *testing.T) {
		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce + 1}, nil
			},
			SendTransactionCalled: func(txs *data.Transaction) (string, error) {
				return "", expectedErr
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce
		anh.transaction = tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
		require.Nil(t, anh.transaction)
	})
	t.Run("should work", func(t *testing.T) {
		blockchainNonce := uint64(100)
		proxy := &testsCommon.ProxyStub{
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{Nonce: blockchainNonce}, nil
			},
			SendTransactionCalled: func(txs *data.Transaction) (string, error) {
				return "hash", nil
			},
		}
		anh, _ := NewSingleTransactionAddressNonceHandler(proxy, testAddress)
		txArgs := createTxArgs()
		tx := createTx(txArgs.GasPrice, txArgs)
		tx.Nonce = blockchainNonce
		anh.transaction = tx

		err := anh.ReSendTransactionsIfRequired(context.Background())
		require.Nil(t, err)
	})
}
