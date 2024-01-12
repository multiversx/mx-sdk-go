package nonceHandlerV2

import (
	"container/heap"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

type TransactionQueue []*transaction.FrontendTransaction

// Push required by the heap.Interface
func (tq *TransactionQueue) Push(newTx interface{}) {
	tx := newTx.(*transaction.FrontendTransaction)
	*tq = append(*tq, tx)
}

// Pop required by the heap.Interface
func (tq *TransactionQueue) Pop() interface{} {
	currentQueue := *tq
	n := len(currentQueue)
	tx := currentQueue[n-1]
	*tq = currentQueue[0 : n-1]
	return tx
}

// Len required by sort.Interface
func (tq *TransactionQueue) Len() int {
	return len(*tq)
}

// Swap required by the sort.Interface
func (tq TransactionQueue) Swap(a, b int) {
	tq[a], tq[b] = tq[b], tq[a]
}

// Less required by the sort.Interface
// we flip the comparer (to greater than) because we need
// the comparer to sort by highest prio, not lowest
func (tq TransactionQueue) Less(a, b int) bool {
	return tq[a].Nonce > tq[b].Nonce
}

type TransactionQueueHandler struct {
	tq TransactionQueue
}

func NewTransactionQueueHandler() *TransactionQueueHandler {
	th := &TransactionQueueHandler{tq: make(TransactionQueue, 0)}
	heap.Init(&th.tq)
	return th
}

func (th *TransactionQueueHandler) AddTransaction(transaction *transaction.FrontendTransaction) {
	heap.Push(&th.tq, transaction)
}

func (th *TransactionQueueHandler) NextTransaction() *transaction.FrontendTransaction {
	if th.tq.Len() == 0 {
		return nil
	}

	nextTransaction := heap.Pop(&th.tq)
	return nextTransaction.(*transaction.FrontendTransaction)
}

func (th *TransactionQueueHandler) Len() int {
	return th.tq.Len()
}

func (th *TransactionQueueHandler) SearchForNonce(nonce uint64) *transaction.FrontendTransaction {
	if th.tq.Len() == 0 {
		return nil
	}

	for i, _ := range th.tq {
		if th.tq[i].Nonce == nonce {
			return th.tq[i]
		}
	}

	return nil
}
