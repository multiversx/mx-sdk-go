package workflows

import (
	"context"
	"encoding/json"
	"math/big"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("workflows")

// WalletTrackerArgs is the argument DTO for the NewWalletTracker constructor function
type WalletTrackerArgs struct {
	TrackableAddressesProvider TrackableAddressesProvider
	Proxy                      ProxyHandler
	NonceHandler               LastProcessedNonceHandler
	CheckInterval              time.Duration
	MinimumBalance             *big.Int
}

// walletTracker is able to track a set of addresses by storing those that received a greater-than-specified
// amount of EGLD. It does this by parsing hyper block by hyper block and checking each transaction
type walletTracker struct {
	accumulator                *addressesAccumulator
	trackableAddressesProvider TrackableAddressesProvider
	proxy                      ProxyHandler
	nonceHandler               LastProcessedNonceHandler
	checkInterval              time.Duration
	cancelFunc                 func()
	minimumBalance             *big.Int

	mutHandlers                  sync.RWMutex
	handlerNewDepositTransaction func(transaction data.TransactionOnNetwork)
}

// NewWalletTracker will create a new walletTracker instance. It automatically starts an inner
// processLoop go routine that can be stopped by calling the Close method
func NewWalletTracker(args WalletTrackerArgs) (*walletTracker, error) {
	if check.IfNil(args.TrackableAddressesProvider) {
		return nil, ErrNilTrackableAddressesProvider
	}
	if check.IfNil(args.Proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(args.NonceHandler) {
		return nil, ErrNilLastProcessedNonceHandler
	}
	if args.MinimumBalance == nil {
		return nil, ErrNilMinimumBalance
	}

	wt := &walletTracker{
		accumulator:                newAddressesAccumulator(),
		trackableAddressesProvider: args.TrackableAddressesProvider,
		proxy:                      args.Proxy,
		nonceHandler:               args.NonceHandler,
		checkInterval:              args.CheckInterval,
		minimumBalance:             args.MinimumBalance,
	}

	var ctx context.Context
	ctx, wt.cancelFunc = context.WithCancel(context.Background())
	go wt.processLoop(ctx)

	return wt, nil
}

func (wt *walletTracker) processLoop(ctx context.Context) {
	log.Debug("walletTracker.processLoop started")

	timer := time.NewTimer(wt.checkInterval)
	defer timer.Stop()

	for {
		timer.Reset(wt.checkInterval)

		select {
		case <-timer.C:
			err := wt.fetchAndProcessHyperBlocks(ctx)
			log.LogIfError(err)
		case <-ctx.Done():
			log.Debug("terminating walletTracker.processLoop...")
			return
		}
	}
}

func (wt *walletTracker) fetchAndProcessHyperBlocks(ctx context.Context) error {
	lastProcessedNonce := wt.nonceHandler.GetLastProcessedNonce()
	networkNonce, err := wt.proxy.GetLatestHyperBlockNonce(ctx)
	if err != nil {
		return err
	}

	for nonce := lastProcessedNonce + 1; nonce <= networkNonce; nonce++ {
		err = wt.fetchAndProcessHyperBlock(ctx, nonce)
		if err != nil {
			return err
		}

		wt.nonceHandler.ProcessedNonce(nonce)
	}

	return nil
}

func (wt *walletTracker) fetchAndProcessHyperBlock(ctx context.Context, nonce uint64) error {
	block, err := wt.proxy.GetHyperBlockByNonce(ctx, nonce)
	if err != nil {
		return err
	}

	wt.processHyperBlock(block)

	log.Debug("processed hyper block", "nonce", nonce, "hash", block.Hash, "num txs", block.NumTxs)

	return nil
}

func (wt *walletTracker) processHyperBlock(block *data.HyperBlock) {
	for _, transaction := range block.Transactions {
		err := wt.processTransaction(transaction)
		if err != nil {
			transactionString, _ := json.Marshal(&transaction)
			log.Warn("error processing transaction, ignoring",
				"transaction", transactionString, "error", err)
		}
	}
}

func (wt *walletTracker) processTransaction(transaction data.TransactionOnNetwork) error {
	value, ok := big.NewInt(0).SetString(transaction.Value, 10)
	if !ok {
		return ErrInvalidTransactionValue
	}

	if !wt.trackableAddressesProvider.IsTrackableAddresses(transaction.Receiver) {
		return nil
	}

	if value.Cmp(wt.minimumBalance) < 0 {
		// transaction has a very small value transfer (possible attack vector as someone
		// can trigger millions of these transactions as to consume the owner's balance through fees)
		return nil
	}

	wt.notifyNewDepositTransactionFound(transaction)
	wt.accumulator.push(transaction.Receiver)

	return nil
}

func (wt *walletTracker) notifyNewDepositTransactionFound(transaction data.TransactionOnNetwork) {
	wt.mutHandlers.RLock()
	defer wt.mutHandlers.RUnlock()

	if wt.handlerNewDepositTransaction != nil {
		wt.handlerNewDepositTransaction(transaction)
	}
}

// SetHandlerForNewDepositTransactionFound will set the handler that will get notified each time a new deposit
// transaction is found on a hyper block
func (wt *walletTracker) SetHandlerForNewDepositTransactionFound(handler func(tx data.TransactionOnNetwork)) {
	if handler == nil {
		return
	}

	wt.mutHandlers.Lock()
	wt.handlerNewDepositTransaction = handler
	wt.mutHandlers.Unlock()
}

// GetLatestTrackedAddresses returns the accumulated addresses that contained changed balances
func (wt *walletTracker) GetLatestTrackedAddresses() []string {
	return wt.accumulator.pop()
}

// Close will close the process loop go routine
func (wt *walletTracker) Close() error {
	wt.cancelFunc()

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wt *walletTracker) IsInterfaceNil() bool {
	return wt == nil
}
