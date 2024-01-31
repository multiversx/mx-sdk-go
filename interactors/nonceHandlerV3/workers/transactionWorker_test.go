package workers

import (
	"context"
	"fmt"
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

	// We add roughly at the same time un-ordered transactions.
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 91})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 1})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 13})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 10})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 99})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 8})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: 7})

	// Verify that the results come in ordered.
	for _, n := range sortedNonces {
		require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(n, 10), Error: nil}, <-w.responsesChannels[n])
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
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[0]})
	w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[1]})

	// We add another transaction with a lower nonce after a while
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(2*time.Second, func() {
		w.AddTransaction(&transaction.FrontendTransaction{Nonce: nonces[2]})
		wg.Done()
	})

	// Verify that the transactions have been processed in the right order.
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[0], 10), Error: nil}, <-w.responsesChannels[nonces[0]])
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[1], 10), Error: nil}, <-w.responsesChannels[nonces[1]])

	// Wait for the scheduled transaction to finish. After that we verify that the transaction it has been processed.
	// Even though the nonce was lower than the first two.
	wg.Wait()
	require.Equal(t, &TransactionResponse{TxHash: strconv.FormatUint(nonces[2], 10), Error: nil}, <-w.responsesChannels[nonces[2]])
}

func TestMe(t *testing.T) {

	ticker := time.NewTicker(time.Second)
	i := 0
	for range ticker.C {
		fmt.Println(i)
		i++

		if i == 5 {
			break
		}
	}
}
