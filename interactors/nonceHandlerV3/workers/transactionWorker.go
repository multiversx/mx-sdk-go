package workers

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"

	"github.com/multiversx/mx-sdk-go/interactors"
)

// TransactionResponse wraps the results provided by the endpoint which will send the transaction in a struct.
type TransactionResponse struct {
	TxHash string
	Error  error
}

type TransactionQueue interface {
	Push(newTx interface{})
	Pop() interface{}
	Len() int
}

type transactionQueue []*transaction.FrontendTransaction

// Push required by the heap.Interface
func (tq *transactionQueue) Push(newTx interface{}) {
	tx := newTx.(*transaction.FrontendTransaction)
	*tq = append(*tq, tx)
}

// Pop required by the heap.Interface
func (tq *transactionQueue) Pop() interface{} {
	currentQueue := *tq
	n := len(currentQueue)
	tx := currentQueue[n-1]
	*tq = currentQueue[0 : n-1]
	return tx
}

// Len required by sort.Interface
func (tq *transactionQueue) Len() int {
	return len(*tq)
}

// Swap required by the sort.Interface
func (tq transactionQueue) Swap(a, b int) {
	tq[a], tq[b] = tq[b], tq[a]
}

// Less required by the sort.Interface
// Meaning that in the heap, the transaction with the lowest nonce has priority.
func (tq transactionQueue) Less(a, b int) bool {
	return tq[a].Nonce < tq[b].Nonce
}

// TransactionWorker handles all transaction stored inside a priority queue. The priority is given by the nonce, meaning
// that transactions with lower nonce will be sent first.
type TransactionWorker struct {
	mu sync.Mutex
	tq transactionQueue

	proxy             interactors.Proxy
	responsesChannels map[uint64]chan *TransactionResponse
}

// NewTransactionWorker creates a new instance of TransactionWorker.
func NewTransactionWorker(proxy interactors.Proxy, pollingInterval time.Duration) *TransactionWorker {
	tw := &TransactionWorker{
		mu:                sync.Mutex{},
		tq:                make(transactionQueue, 0),
		proxy:             proxy,
		responsesChannels: make(map[uint64]chan *TransactionResponse),
	}
	heap.Init(&tw.tq)

	tw.start(pollingInterval)
	return tw
}

// AddTransaction will add a transaction to the priority queue (heap) and will create a channel where the promised result
// will be broadcast on.
func (tw *TransactionWorker) AddTransaction(transaction *transaction.FrontendTransaction) <-chan *TransactionResponse {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	r := make(chan *TransactionResponse)
	tw.responsesChannels[transaction.Nonce] = r

	heap.Push(&tw.tq, transaction)
	return r
}

// start will spawn a goroutine tasked with iterating all the transactions inside the priority queue. The priority is
// given by the nonce, meaning that transaction with lower nonce will be sent first.
func (tw *TransactionWorker) start(pollingInterval time.Duration) {
	go func() {
		for {
			// Retrieve the transaction in the queue with the lowest nonce.
			tx := tw.nextTransaction()

			// If there are no transaction in the queue, the result will be nil.
			// That means there are no transactions to send.
			if tx == nil {
				time.Sleep(pollingInterval)
				continue
			}

			// We retrieve the channel where we will send the response.
			// Everytime a transaction is added to the queue, such a channel is created and placed in a map.
			tw.mu.Lock()
			r := tw.responsesChannels[tx.Nonce]
			tw.mu.Unlock()

			// Send the transaction and forward the response on the channel promised.
			txHash, err := tw.proxy.SendTransaction(context.TODO(), tx)
			r <- &TransactionResponse{TxHash: txHash, Error: err}
		}
	}()
}

// nextTransaction will return the transaction stored in the priority queue (heap) with the lowest nonce.
// If there aren't any transaction, the result will be nil.
func (tw *TransactionWorker) nextTransaction() *transaction.FrontendTransaction {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.tq.Len() == 0 {
		return nil
	}

	nextTransaction := heap.Pop(&tw.tq)
	return nextTransaction.(*transaction.FrontendTransaction)
}
