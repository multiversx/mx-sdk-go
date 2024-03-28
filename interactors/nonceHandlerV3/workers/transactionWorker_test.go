package workers

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-sdk-go/testsCommon"
)

func TestTransactionWorker_AddTransaction(t *testing.T) {
	t.Parallel()
	sortedNonces := []uint64{1, 7, 8, 10, 13, 91, 99}
	proxy := &testsCommon.ProxyStub{
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			return strconv.FormatUint(tx.Nonce, 10), nil
		},
	}

	w := NewTransactionWorker(context.Background(), proxy, 2*time.Second)

	responseChannels := make([]<-chan *TransactionResponse, 7)

	// We add roughly at the same time un-ordered transactions.
	responseChannels[5] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 91})
	responseChannels[0] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 1})
	responseChannels[4] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 13})
	responseChannels[3] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 10})
	responseChannels[6] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 99})
	responseChannels[2] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 8})
	responseChannels[1] = w.AddTransaction(&transaction.FrontendTransaction{Nonce: 7})

	// Verify that the results come in ordered.
	for i, n := range sortedNonces {
		require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(n, 10), Error: nil}, <-responseChannels[i])
	}
}

func TestTransactionWorker_AddTransactionWithLowerNonceAfter(t *testing.T) {
	t.Parallel()
	nonces := []uint64{10, 11, 9}
	proxy := &testsCommon.ProxyStub{
		SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
			return strconv.FormatUint(tx.Nonce, 10), nil
		},
	}

	w := NewTransactionWorker(context.Background(), proxy, 1*time.Second)

	// We add two ordered by nonce transactions roughly at the same time.
	r1 := w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[0]})
	r2 := w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[1]})

	// We add another transaction with a lower nonce after a while
	var wg sync.WaitGroup
	r3 := make(<-chan *TransactionResponse, 1)
	wg.Add(1)
	time.AfterFunc(2*time.Second, func() {
		r3 = w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[2]})
		wg.Done()
	})

	// Verify that the transactions have been processed in the right order.
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[0], 10), Error: nil}, <-r1)
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[1], 10), Error: nil}, <-r2)

	// Wait for the scheduled transaction to finish. After that we verify that the transaction it has been processed.
	// Even though the nonce was lower than the first two.
	wg.Wait()
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[2], 10), Error: nil}, <-r3)
}
