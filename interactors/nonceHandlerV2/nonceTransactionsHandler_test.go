package nonceHandlerV2

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNonceTransactionHandlerV2(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNonceTransactionsHandlerV2()
		args.Proxy = nil
		nth, err := NewNonceTransactionHandlerV2(args)
		require.Nil(t, nth)
		assert.Equal(t, interactors.ErrNilProxy, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNonceTransactionsHandlerV2()
		nth, err := NewNonceTransactionHandlerV2(args)
		require.NotNil(t, nth)
		require.Nil(t, err)

		require.Nil(t, nth.Close())
	})
}

func TestNonceTransactionsHandlerV2_GetNonce(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)
	numCalls := 0

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			numCalls++

			return &data.Account{
				Nonce: currentNonce,
			}, nil
		},
	}

	nth, _ := NewNonceTransactionHandlerV2(args)
	err := nth.ApplyNonceAndGasPrice(context.Background(), nil, nil)
	assert.Equal(t, interactors.ErrNilAddress, err)

	tx := transaction.FrontendTransaction{}
	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, &tx)
	assert.Nil(t, err)
	assert.Equal(t, currentNonce, tx.Nonce)

	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, &tx)
	assert.Nil(t, err)
	assert.Equal(t, currentNonce+1, tx.Nonce)

	assert.Equal(t, 2, numCalls)

	require.Nil(t, nth.Close())
}

func TestNonceTransactionsHandlerV2_SendMultipleTransactionsResendingEliminatingOne(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.FrontendTransaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionsCalled: func(txs []*transaction.FrontendTransaction) ([]string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = txs
			numCalls++
			hashes := make([]string, len(txs))

			return hashes, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.FrontendTransaction{tx}
			numCalls++

			return "", nil
		},
	}
	nth, _ := NewNonceTransactionHandlerV2(args)

	numTxs := 5
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.TODO(), txs[i])
		require.Nil(t, err)
	}

	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	numSentTransaction := 5
	numSentTransactions := 1
	assert.Equal(t, numSentTransaction+numSentTransactions, len(sentTransactions))
	for i := 0; i < numSentTransaction; i++ {
		assert.Equal(t, 1, len(sentTransactions[i]))
	}
	assert.Equal(t, numTxs-1, len(sentTransactions[numSentTransaction])) // resend
}

func TestNonceTransactionsHandlerV2_SendMultipleTransactionsResendingEliminatingAll(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.FrontendTransaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.FrontendTransaction{tx}
			numCalls++

			return "", nil
		},
	}
	numTxs := 5
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.Background(), txs[i])
		require.Nil(t, err)
	}

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, 5, len(sentTransactions))
	assert.Equal(t, 1, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendTransactionResendingEliminatingAll(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.FrontendTransaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.FrontendTransaction{tx}
			numCalls++

			return "", nil
		},
	}

	numTxs := 1
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))

	hash, err := nth.SendTransaction(context.Background(), txs[0])
	require.Nil(t, err)
	require.Equal(t, "", hash)

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, 1, len(sentTransactions))
	assert.Equal(t, numTxs, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendTransactionErrors(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	var errSent error

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			return "", errSent
		},
	}

	numTxs := 1
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactions(testAddress, numTxs, atomic.LoadUint64(&currentNonce))

	hash, err := nth.SendTransaction(context.Background(), nil)
	require.Equal(t, interactors.ErrNilTransaction, err)
	require.Equal(t, "", hash)

	errSent = errors.New("expected error")

	hash, err = nth.SendTransaction(context.Background(), txs[0])
	require.True(t, errors.Is(err, errSent))
	require.Equal(t, "", hash)
}

func createMockTransactions(addr core.AddressHandler, numTxs int, startNonce uint64) []*transaction.FrontendTransaction {
	txs := make([]*transaction.FrontendTransaction, 0, numTxs)
	addrAsBech32String, _ := addr.AddressAsBech32String()
	for i := 0; i < numTxs; i++ {
		tx := &transaction.FrontendTransaction{
			Nonce:     startNonce,
			Value:     "1",
			Receiver:  addrAsBech32String,
			Sender:    addrAsBech32String,
			GasPrice:  100000,
			GasLimit:  50000,
			Data:      nil,
			Signature: "sig",
			ChainID:   "3",
			Version:   1,
		}

		txs = append(txs, tx)
		startNonce++
	}

	return txs
}

func TestNonceTransactionsHandlerV2_SendTransactionsWithGetNonce(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	mutSentTransactions := sync.Mutex{}
	numCalls := 0
	sentTransactions := make(map[int][]*transaction.FrontendTransaction)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			mutSentTransactions.Lock()
			defer mutSentTransactions.Unlock()

			sentTransactions[numCalls] = []*transaction.FrontendTransaction{tx}
			numCalls++

			return "", nil
		},
	}

	numTxs := 5
	nth, _ := NewNonceTransactionHandlerV2(args)
	txs := createMockTransactionsWithGetNonce(t, testAddress, 5, nth)
	for i := 0; i < numTxs; i++ {
		_, err := nth.SendTransaction(context.Background(), txs[i])
		require.Nil(t, err)
	}

	atomic.AddUint64(&currentNonce, uint64(numTxs))
	time.Sleep(time.Second * 3)
	_ = nth.Close()

	mutSentTransactions.Lock()
	defer mutSentTransactions.Unlock()

	//no resend operation was made because all transactions were executed (nonce was incremented)
	assert.Equal(t, numTxs, len(sentTransactions))
	assert.Equal(t, 1, len(sentTransactions[0]))
}

func TestNonceTransactionsHandlerV2_SendDuplicateTransactions(t *testing.T) {
	currentNonce := uint64(664)

	numCalls := 0

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			require.LessOrEqual(t, numCalls, 1)
			currentNonce++
			return "", nil
		},
	}
	nth, _ := NewNonceTransactionHandlerV2(args)

	tx := &transaction.FrontendTransaction{
		Value:    "1",
		Receiver: testAddressAsBech32String,
		Sender:   testAddressAsBech32String,
		GasPrice: 100000,
		GasLimit: 50000,
		Data:     nil,
		ChainID:  "3",
		Version:  1,
	}
	err := nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Nil(t, err)

	_, err = nth.SendTransaction(context.Background(), tx)
	require.Nil(t, err)
	acc, _ := nth.getOrCreateAddressNonceHandler(testAddress)
	accWithPrivateAccess, ok := acc.(*addressNonceHandler)
	require.True(t, ok)

	// after sending first tx, nonce shall increase
	require.Equal(t, accWithPrivateAccess.computedNonce+1, currentNonce)

	// trying to apply nonce for the same tx, NonceTransactionHandler shall return ErrTxAlreadySent
	// and computedNonce shall not increase
	tx.Nonce = 0
	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Equal(t, err, interactors.ErrTxAlreadySent)
	require.Equal(t, tx.Nonce, uint64(0))
	require.Equal(t, accWithPrivateAccess.computedNonce+1, currentNonce)
}

func createMockTransactionsWithGetNonce(
	tb testing.TB,
	addr core.AddressHandler,
	numTxs int,
	nth *nonceTransactionsHandlerV2,
) []*transaction.FrontendTransaction {
	txs := make([]*transaction.FrontendTransaction, 0, numTxs)
	addrAsBech32String, _ := addr.AddressAsBech32String()
	for i := 0; i < numTxs; i++ {
		tx := &transaction.FrontendTransaction{}
		err := nth.ApplyNonceAndGasPrice(context.Background(), addr, tx)
		require.Nil(tb, err)

		tx.Value = "1"
		tx.Receiver = addrAsBech32String
		tx.Sender = addrAsBech32String
		tx.GasLimit = 50000
		tx.Data = nil
		tx.Signature = "sig"
		tx.ChainID = "3"
		tx.Version = 1

		txs = append(txs, tx)
	}

	return txs
}

func TestNonceTransactionsHandlerV2_ForceNonceReFetch(t *testing.T) {
	t.Parallel()

	currentNonce := uint64(664)

	args := createMockArgsNonceTransactionsHandlerV2()
	args.Proxy = &testsCommon.ProxyStub{
		GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
			addressAsBech32String, _ := address.AddressAsBech32String()
			if addressAsBech32String != testAddressAsBech32String {
				return nil, errors.New("unexpected address")
			}

			return &data.Account{
				Nonce: atomic.LoadUint64(&currentNonce),
			}, nil
		},
	}

	nth, _ := NewNonceTransactionHandlerV2(args)
	tx := &transaction.FrontendTransaction{}
	_ = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	_ = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	err := nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	require.Nil(t, err)
	assert.Equal(t, atomic.LoadUint64(&currentNonce)+2, tx.Nonce)

	err = nth.DropTransactions(nil)
	assert.Equal(t, interactors.ErrNilAddress, err)

	err = nth.DropTransactions(testAddress)
	assert.Nil(t, err)

	err = nth.ApplyNonceAndGasPrice(context.Background(), testAddress, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, atomic.LoadUint64(&currentNonce), tx.Nonce)
}

func createMockArgsNonceTransactionsHandlerV2() ArgsNonceTransactionsHandlerV2 {
	return ArgsNonceTransactionsHandlerV2{
		Proxy:            &testsCommon.ProxyStub{},
		IntervalToResend: time.Second * 2,
	}
}
