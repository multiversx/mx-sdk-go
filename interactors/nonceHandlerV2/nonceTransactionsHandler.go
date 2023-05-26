package nonceHandlerV2

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

var log = logger.GetOrCreate("mx-sdk-go/interactors/nonceHandlerV2")

// ArgsNonceTransactionsHandlerV2 is the argument DTO for a nonce transactions handler component
type ArgsNonceTransactionsHandlerV2 struct {
	Proxy            interactors.Proxy
	IntervalToResend time.Duration
	Creator          interactors.AddressNonceHandlerCreator
}

// nonceTransactionsHandlerV2 is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck transactions" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandlerV2 should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandlerV2 struct {
	proxy            interactors.Proxy
	mutHandlers      sync.RWMutex
	creator          interactors.AddressNonceHandlerCreator
	handlers         map[string]interactors.AddressNonceHandler
	cancelFunc       func()
	intervalToResend time.Duration
}

// NewNonceTransactionHandlerV2 will create a new instance of the nonceTransactionsHandlerV2. It requires a Proxy implementation
// and an interval at which the transactions sent are rechecked and eventually, resent.
func NewNonceTransactionHandlerV2(args ArgsNonceTransactionsHandlerV2) (*nonceTransactionsHandlerV2, error) {
	if check.IfNil(args.Proxy) {
		return nil, interactors.ErrNilProxy
	}
	if args.IntervalToResend < minimumIntervalToResend {
		return nil, fmt.Errorf("%w for intervalToResend in NewNonceTransactionHandlerV2", interactors.ErrInvalidValue)
	}
	if check.IfNil(args.Creator) {
		return nil, interactors.ErrNilAddressNonceHandlerCreator
	}

	nth := &nonceTransactionsHandlerV2{
		proxy:            args.Proxy,
		handlers:         make(map[string]interactors.AddressNonceHandler),
		intervalToResend: args.IntervalToResend,
		creator:          args.Creator,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	nth.cancelFunc = cancelFunc
	go nth.resendTransactionsLoop(ctx)

	return nth, nil
}

// ApplyNonceAndGasPrice will apply the nonce to the given ArgCreateTransaction
func (nth *nonceTransactionsHandlerV2) ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error {
	if check.IfNil(address) {
		return interactors.ErrNilAddress
	}
	if tx == nil {
		return interactors.ErrNilTransaction
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return err
	}

	return anh.ApplyNonceAndGasPrice(ctx, tx)
}

func (nth *nonceTransactionsHandlerV2) getOrCreateAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	anh := nth.getAddressNonceHandler(address)
	if !check.IfNil(anh) {
		return anh, nil
	}

	return nth.createAddressNonceHandler(address)
}

func (nth *nonceTransactionsHandlerV2) getAddressNonceHandler(address core.AddressHandler) interactors.AddressNonceHandler {
	nth.mutHandlers.RLock()
	defer nth.mutHandlers.RUnlock()

	anh, found := nth.handlers[string(address.AddressBytes())]
	if found {
		return anh
	}
	return nil
}

func (nth *nonceTransactionsHandlerV2) createAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	addressAsString := string(address.AddressBytes())
	anh, found := nth.handlers[addressAsString]
	if found {
		return anh, nil
	}
	anh, err := nth.creator.Create(nth.proxy, address)
	if err != nil {
		return nil, err
	}
	nth.handlers[addressAsString] = anh

	return anh, nil
}

// SendTransaction will store and send the provided transaction
func (nth *nonceTransactionsHandlerV2) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if tx == nil {
		return "", interactors.ErrNilTransaction
	}

	addrAsBech32 := tx.Sender
	address, err := data.NewAddressFromBech32String(addrAsBech32)
	if err != nil {
		return "", fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return "", err
	}

	sentHash, err := anh.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("%w while sending transaction for address %s", err, addrAsBech32)
	}

	return sentHash, nil
}

func (nth *nonceTransactionsHandlerV2) resendTransactionsLoop(ctx context.Context) {
	timer := time.NewTimer(nth.intervalToResend)
	defer timer.Stop()

	for {
		timer.Reset(nth.intervalToResend)

		select {
		case <-timer.C:
			nth.resendTransactions(ctx)
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV2.resendTransactionsLoop...")
			return
		}
	}
}

func (nth *nonceTransactionsHandlerV2) resendTransactions(ctx context.Context) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	for _, anh := range nth.handlers {
		select {
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV2.resendTransactions...")
			return
		default:
		}

		resendCtx, cancel := context.WithTimeout(ctx, nth.intervalToResend)
		err := anh.ReSendTransactionsIfRequired(resendCtx)
		log.LogIfError(err)
		cancel()
	}
}

// DropTransactions will clean the addressNonceHandler cached transactions. A little gas increase will be applied to the next transactions
// in order to also replace the transactions from the txPool.
// This should be only used in a fallback plan, when some transactions are completely lost (or due to a bug, not even sent in first time)
func (nth *nonceTransactionsHandlerV2) DropTransactions(address core.AddressHandler) error {
	if check.IfNil(address) {
		return interactors.ErrNilAddress
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return err
	}
	anh.DropTransactions()

	return nil
}

// Close finishes the transactions resend go routine
func (nth *nonceTransactionsHandlerV2) Close() error {
	nth.cancelFunc()

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nth *nonceTransactionsHandlerV2) IsInterfaceNil() bool {
	return nth == nil
}
