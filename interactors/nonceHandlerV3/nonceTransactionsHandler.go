package nonceHandlerV3

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

var log = logger.GetOrCreate("mx-sdk-go/interactors/nonceHandlerV3")

// ArgsNonceTransactionsHandlerV3 is the argument DTO for a nonce transactions handler component
type ArgsNonceTransactionsHandlerV3 struct {
	Proxy            interactors.Proxy
	IntervalToResend time.Duration
}

// nonceTransactionsHandlerV3 is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck transactions" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandlerV3 should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandlerV3 struct {
	proxy            interactors.Proxy
	mutHandlers      sync.RWMutex
	handlers         map[string]interactors.AddressNonceHandler
	cancelFunc       func()
	intervalToResend time.Duration
}

// NewNonceTransactionHandlerV3 will create a new instance of the nonceTransactionsHandlerV3. It requires a Proxy implementation
// and an interval at which the transactions sent are rechecked and eventually, resent.
func NewNonceTransactionHandlerV3(args ArgsNonceTransactionsHandlerV3) (*nonceTransactionsHandlerV3, error) {
	if check.IfNil(args.Proxy) {
		return nil, interactors.ErrNilProxy
	}
	if args.IntervalToResend < minimumIntervalToResend {
		return nil, fmt.Errorf("%w for intervalToResend in NewNonceTransactionHandlerV2", interactors.ErrInvalidValue)
	}

	nth := &nonceTransactionsHandlerV3{
		proxy:            args.Proxy,
		handlers:         make(map[string]interactors.AddressNonceHandler),
		intervalToResend: args.IntervalToResend,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	nth.cancelFunc = cancelFunc
	go nth.resendTransactionsLoop(ctx)

	return nth, nil
}

// ApplyNonceAndGasPrice will apply the nonce to the given frontend transaction
func (nth *nonceTransactionsHandlerV3) ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error {
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

func (nth *nonceTransactionsHandlerV3) getOrCreateAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	anh := nth.getAddressNonceHandler(address)
	if !check.IfNil(anh) {
		return anh, nil
	}

	return nth.createAddressNonceHandler(address)
}

func (nth *nonceTransactionsHandlerV3) getAddressNonceHandler(address core.AddressHandler) interactors.AddressNonceHandler {
	nth.mutHandlers.RLock()
	defer nth.mutHandlers.RUnlock()

	anh, found := nth.handlers[string(address.AddressBytes())]
	if found {
		return anh
	}
	return nil
}

func (nth *nonceTransactionsHandlerV3) createAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	addressAsString := string(address.AddressBytes())
	anh, found := nth.handlers[addressAsString]
	if found {
		return anh, nil
	}

	anh, err := NewAddressNonceHandler(nth.proxy, address)
	if err != nil {
		return nil, err
	}
	nth.handlers[addressAsString] = anh

	return anh, nil
}

// SendTransaction will store and send the provided transaction
func (nth *nonceTransactionsHandlerV3) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if tx == nil {
		return "", interactors.ErrNilTransaction
	}

	// Work with a full copy of the provided transaction so the provided one can change without affecting this component.
	// Abnormal and unpredictable behaviors due to the resending mechanism are prevented this way
	txCopy := *tx

	addrAsBech32 := txCopy.Sender
	address, err := data.NewAddressFromBech32String(addrAsBech32)
	if err != nil {
		return "", fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return "", err
	}

	sentHash, err := anh.SendTransaction(ctx, &txCopy)
	if err != nil {
		return "", fmt.Errorf("%w while sending transaction for address %s", err, addrAsBech32)
	}

	return sentHash, nil
}

func (nth *nonceTransactionsHandlerV3) resendTransactionsLoop(ctx context.Context) {
	ticker := time.NewTicker(nth.intervalToResend)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nth.resendTransactions(ctx)
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV3.resendTransactionsLoop...")
			return
		}
	}
}

func (nth *nonceTransactionsHandlerV3) resendTransactions(ctx context.Context) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	for _, anh := range nth.handlers {
		select {
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandlerV3.resendTransactions...")
			return
		default:
		}

		resendCtx, cancel := context.WithTimeout(ctx, nth.intervalToResend)
		err := anh.ReSendTransactionsIfRequired(resendCtx)
		fmt.Println("i got cancelled once")
		log.LogIfError(err)
		cancel()
	}
}

// DropTransactions will clean the addressNonceHandler cached transactions. A little gas increase will be applied to the next transactions
// in order to also replace the transactions from the txPool.
// This should be only used in a fallback plan, when some transactions are completely lost (or due to a bug, not even sent in first time)
func (nth *nonceTransactionsHandlerV3) DropTransactions(address core.AddressHandler) error {
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
func (nth *nonceTransactionsHandlerV3) Close() error {
	nth.cancelFunc()

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nth *nonceTransactionsHandlerV3) IsInterfaceNil() bool {
	return nth == nil
}
