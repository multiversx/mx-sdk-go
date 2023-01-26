package interactors

import (
	"context"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/data"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/interactors")

const defaultTimeBetweenBunches = time.Second

type transactionInteractor struct {
	Proxy
	TxBuilder
	mutTxAccumulator      sync.RWMutex
	mutTimeBetweenBunches sync.RWMutex
	timeBetweenBunches    time.Duration
	txAccumulator         []*data.Transaction
}

// NewTransactionInteractor will create an interactor that extends the proxy functionality with some transaction-oriented functionality
func NewTransactionInteractor(proxy Proxy, txBuilder TxBuilder) (*transactionInteractor, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(txBuilder) {
		return nil, ErrNilTxBuilder
	}

	return &transactionInteractor{
		Proxy:              proxy,
		TxBuilder:          txBuilder,
		timeBetweenBunches: defaultTimeBetweenBunches,
	}, nil
}

// SetTimeBetweenBunches sets the time between bunch sends
func (ti *transactionInteractor) SetTimeBetweenBunches(timeBetweenBunches time.Duration) {
	ti.mutTimeBetweenBunches.Lock()
	ti.timeBetweenBunches = timeBetweenBunches
	ti.mutTimeBetweenBunches.Unlock()
}

// AddTransaction will add the provided transaction in the transaction accumulator
func (ti *transactionInteractor) AddTransaction(tx *data.Transaction) {
	if tx == nil {
		return
	}

	ti.mutTxAccumulator.Lock()
	ti.txAccumulator = append(ti.txAccumulator, tx)
	ti.mutTxAccumulator.Unlock()
}

// PopAccumulatedTransactions will return the whole accumulated contents emptying the accumulator
func (ti *transactionInteractor) PopAccumulatedTransactions() []*data.Transaction {
	ti.mutTxAccumulator.Lock()
	result := make([]*data.Transaction, len(ti.txAccumulator))
	copy(result, ti.txAccumulator)
	ti.txAccumulator = make([]*data.Transaction, 0)
	ti.mutTxAccumulator.Unlock()

	return result
}

// SendTransactionsAsBunch will send all stored transactions as bunches
func (ti *transactionInteractor) SendTransactionsAsBunch(ctx context.Context, bunchSize int) ([]string, error) {
	if bunchSize <= 0 {
		return nil, ErrInvalidValue
	}

	ti.mutTimeBetweenBunches.RLock()
	timeBetweenBunches := ti.timeBetweenBunches
	ti.mutTimeBetweenBunches.RUnlock()

	transactions := ti.PopAccumulatedTransactions()
	allHashes := make([]string, 0)
	for bunchIndex := 0; len(transactions) > 0; bunchIndex++ {
		var bunch []*data.Transaction

		log.Debug("sending bunch", "index", bunchIndex)

		if len(transactions) > bunchSize {
			bunch = transactions[0:bunchSize]
			transactions = transactions[bunchSize:]
		} else {
			bunch = transactions
			transactions = make([]*data.Transaction, 0)
		}

		hashes, err := ti.Proxy.SendTransactions(ctx, bunch)
		if err != nil {
			return nil, err
		}

		allHashes = append(allHashes, hashes...)

		time.Sleep(timeBetweenBunches)
	}

	return allHashes, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ti *transactionInteractor) IsInterfaceNil() bool {
	return ti == nil
}
