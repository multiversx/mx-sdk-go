package workers

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"

	"github.com/multiversx/mx-sdk-go/interactors"
)

var log = logger.GetOrCreate("mx-sdk-go/interactors/workers/transactionWorker")

// TransactionResponse wraps the results provided by the endpoint which will send the transaction in a struct.
type TransactionResponse struct {
	TxHash string
	Error  error
}

// TransactionQueueItem is a wrapper struct on the transaction itself that is used to encapsulate transactions in
// the priority queue.
type TransactionQueueItem struct {
	tx    *transaction.FrontendTransaction
	index int
}

// A transactionQueue implements heap.Interface and holds Items. Acts like a priority queue.
type transactionQueue []*TransactionQueueItem

// Push required by the heap.Interface
func (tq *transactionQueue) Push(x any) {
	n := len(*tq)
	item := x.(*TransactionQueueItem)
	item.index = n
	*tq = append(*tq, item)
}

// Pop required by the heap.Interface
func (tq *transactionQueue) Pop() any {
	old := *tq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*tq = old[0 : n-1]
	return item
}

// Len required by sort.Interface
func (tq transactionQueue) Len() int {
	return len(tq)
}

// Swap required by the sort.Interface
func (tq transactionQueue) Swap(a, b int) {
	tq[a], tq[b] = tq[b], tq[a]
	tq[a].index = a
	tq[b].index = b
}

// Less required by the sort.Interface
// Meaning that in the heap, the transaction with the lowest nonce has priority.
func (tq transactionQueue) Less(a, b int) bool {
	return tq[a].tx.Nonce < tq[b].tx.Nonce
}

// TransactionWorker handles all transaction stored inside a priority queue. The priority is given by the nonce, meaning
// that transactions with lower nonce will be sent first.
type TransactionWorker struct {
	mu sync.Mutex
	tq transactionQueue

	workerClosed      bool
	proxy             interactors.Proxy
	responsesChannels map[uint64]chan *TransactionResponse
}

// NewTransactionWorker creates a new instance of TransactionWorker.
func NewTransactionWorker(context context.Context, proxy interactors.Proxy, pollingInterval time.Duration) *TransactionWorker {
	tw := &TransactionWorker{
		mu:                sync.Mutex{},
		tq:                make(transactionQueue, 0),
		proxy:             proxy,
		responsesChannels: make(map[uint64]chan *TransactionResponse),
	}
	heap.Init(&tw.tq)

	tw.start(context, pollingInterval)
	return tw
}

// AddTransaction will add a transaction to the priority queue (heap) and will create a channel where the promised result
// will be broadcast on.
func (tw *TransactionWorker) AddTransaction(transaction *transaction.FrontendTransaction) <-chan *TransactionResponse {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	r := make(chan *TransactionResponse, 1)
	if tw.workerClosed {
		r <- &TransactionResponse{TxHash: "", Error: interactors.ErrWorkerClosed}
		return r
	}

	tw.responsesChannels[transaction.Nonce] = r

	heap.Push(&tw.tq, &TransactionQueueItem{tx: transaction})
	return r
}

// start will spawn a goroutine tasked with iterating all the transactions inside the priority queue. The priority is
// given by the nonce, meaning that transaction with lower nonce will be sent first.
func (tw *TransactionWorker) start(ctx context.Context, pollingInterval time.Duration) {
	ticker := time.NewTicker(pollingInterval)

	go func() {
		chMoreWorkToDo := make(chan struct{}, 1)
		for {
			select {
			case <-ctx.Done():
				log.Info("context cancelled - transaction worker has stopped")
				tw.closeAllChannels(ctx)
				return
			case <-ticker.C:
				tw.processNextTransaction(ctx, chMoreWorkToDo)
			case <-chMoreWorkToDo:
				tw.processNextTransaction(ctx, chMoreWorkToDo)
			}
		}
	}()
}

func (tw *TransactionWorker) processNextTransaction(ctx context.Context, chMoreWorkToDo chan struct{}) {
	tx := tw.nextTransaction()
	if tx == nil {
		return
	}

	// Retrieve channel where the response will be broadcast on.
	r := tw.retrieveChannel(tx.Nonce)

	// Send the transaction and forward the response on the channel promised.
	txHash, err := tw.proxy.SendTransaction(ctx, tx)
	r <- &TransactionResponse{TxHash: txHash, Error: err}
	chMoreWorkToDo <- struct{}{}
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
	return nextTransaction.(*TransactionQueueItem).tx
}

func (tw *TransactionWorker) closeAllChannels(ctx context.Context) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	for _, ch := range tw.responsesChannels {
		ch <- &TransactionResponse{TxHash: "", Error: ctx.Err()}
	}
	tw.workerClosed = true
}

func (tw *TransactionWorker) retrieveChannel(nonce uint64) chan *TransactionResponse {
	// We retrieve the channel where we will send the response.
	// Everytime a transaction is added to the queue, such a channel is created and placed in a map.
	// After retrieving it, delete the entry from the map that stores all of them.
	tw.mu.Lock()
	r := tw.responsesChannels[nonce]
	delete(tw.responsesChannels, nonce)
	tw.mu.Unlock()

	return r
}
