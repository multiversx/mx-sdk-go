package nonceHandlerV1

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
)

const minimumIntervalToResend = time.Second

var log = logger.GetOrCreate("mx-sdk-go/interactors/nonceHandlerV1")

// nonceTransactionsHandlerV1 is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck transactions" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandlerV1 should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandlerV1 struct {
	proxy              interactors.Proxy
	mutHandlers        sync.Mutex
	handlers           map[string]*addressNonceHandler
	checkForDuplicates bool
	cancelFunc         func()
	intervalToResend   time.Duration
}

// NewNonceTransactionHandlerV1 will create a new instance of the nonceTransactionsHandlerV1. It requires a Proxy implementation
// and an interval at which the transactions sent are rechecked and eventually, resent.
// checkForDuplicates set as true will prevent sending a transaction with the same receiver, value and data.
func NewNonceTransactionHandlerV1(proxy interactors.Proxy, intervalToResend time.Duration, checkForDuplicates bool) (*nonceTransactionsHandlerV1, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if intervalToResend < minimumIntervalToResend {
		return nil, fmt.Errorf("%w for intervalToResend in NewNonceTransactionHandlerV1", interactors.ErrInvalidValue)
	}

	nth := &nonceTransactionsHandlerV1{
		proxy:              proxy,
		handlers:           make(map[string]*addressNonceHandler),
		intervalToResend:   intervalToResend,
		checkForDuplicates: checkForDuplicates,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	nth.cancelFunc = cancelFunc
	go nth.resendTransactionsLoop(ctx, intervalToResend)

	return nth, nil
}

// GetNonce will return the nonce for the provided address
func (nth *nonceTransactionsHandlerV1) GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error) {
	if check.IfNil(address) {
		return 0, interactors.ErrNilAddress
	}

	anh := nth.getOrCreateAddressNonceHandler(address)

	return anh.getNonceUpdatingCurrent(ctx)
}

func (nth *nonceTransactionsHandlerV1) getOrCreateAddressNonceHandler(address core.AddressHandler) *addressNonceHandler {
	nth.mutHandlers.Lock()
	addressAsString := string(address.AddressBytes())
	anh, found := nth.handlers[addressAsString]
	if !found {
		anh = newAddressNonceHandler(nth.proxy, address)
		nth.handlers[addressAsString] = anh
	}
	nth.mutHandlers.Unlock()

	return anh
}

// SendTransaction will store and send the provided transaction
func (nth *nonceTransactionsHandlerV1) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if tx == nil {
		return "", interactors.ErrNilTransaction
	}

	addrAsBech32 := tx.Sender
	addressHandler, err := data.NewAddressFromBech32String(addrAsBech32)
	if err != nil {
		return "", fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
	}

	anh := nth.getOrCreateAddressNonceHandler(addressHandler)
	if nth.checkForDuplicates && anh.isTxAlreadySent(tx) {
		// TODO: add gas comparation logic EN-11887
		anh.decrementComputedNonce()
		return "", interactors.ErrTxAlreadySent
	}
	sentHash, err := anh.sendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("%w while sending transaction for address %s", err, addrAsBech32)
	}

	return sentHash, nil
}

func (nth *nonceTransactionsHandlerV1) resendTransactionsLoop(ctx context.Context, intervalToResend time.Duration) {
	timer := time.NewTimer(intervalToResend)
	defer timer.Stop()

	for {
		timer.Reset(intervalToResend)

		select {
		case <-timer.C:
			nth.resendTransactions(ctx)
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV1.resendTransactionsLoop...")
			return
		}
	}
}

func (nth *nonceTransactionsHandlerV1) resendTransactions(ctx context.Context) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	for _, anh := range nth.handlers {
		select {
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV1.resendTransactions...")
			return
		default:
		}

		resendCtx, cancel := context.WithTimeout(ctx, nth.intervalToResend)
		err := anh.reSendTransactionsIfRequired(resendCtx)
		log.LogIfError(err)
		cancel()
	}
}

// ForceNonceReFetch will mark the addressNonceHandler to re-fetch its nonce from the blockchain account.
// This should be only used in a fallback plan, when some transactions are completely lost (or due to a bug, not even sent in first time)
func (nth *nonceTransactionsHandlerV1) ForceNonceReFetch(address core.AddressHandler) error {
	if check.IfNil(address) {
		return interactors.ErrNilAddress
	}

	anh := nth.getOrCreateAddressNonceHandler(address)
	anh.markReFetchNonce()

	return nil
}

// Close finishes the transactions resend go routine
func (nth *nonceTransactionsHandlerV1) Close() error {
	nth.cancelFunc()

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nth *nonceTransactionsHandlerV1) IsInterfaceNil() bool {
	return nth == nil
}
