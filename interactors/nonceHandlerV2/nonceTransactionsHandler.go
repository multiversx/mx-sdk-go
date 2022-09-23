package interactors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

const minimumIntervalToResend = time.Second

var log = logger.GetOrCreate("elrond-sdk-erdgo/interactors/nonceHandlerV2")

// nonceTransactionsHandlerV2 is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck transactions" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandlerV2 should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandlerV2 struct {
	proxy              interactors.Proxy
	mutHandlers        sync.Mutex
	handlers           map[string]*addressNonceHandler
	cancelFunc         func()
	intervalToResend   time.Duration
}

// NewNonceTransactionHandlerV2 will create a new instance of the nonceTransactionsHandlerV2. It requires a Proxy implementation
// and an interval at which the transactions sent are rechecked and eventually, resent.
func NewNonceTransactionHandlerV2(proxy Proxy, intervalToResend time.Duration) (*nonceTransactionsHandler, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if intervalToResend < minimumIntervalToResend {
		return nil, fmt.Errorf("%w for intervalToResend in NewNonceTransactionHandler", ErrInvalidValue)
	}

	nth := &nonceTransactionsHandler{
		proxy:            proxy,
		handlers:         make(map[string]*addressNonceHandler),
		intervalToResend: intervalToResend,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	nth.cancelFunc = cancelFunc
	go nth.resendTransactionsLoop(ctx, intervalToResend)

	return nth, nil
}

// ApplyNonce will apply the nonce to the given ArgCreateTransaction
func (nth *nonceTransactionsHandlerV2) ApplyNonce(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
	if check.IfNil(address) {
		return interactors.ErrNilAddress
	}
	if txArgs == nil {
		return interactors.ErrNilArgCreateTransaction
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return err
	}

	return anh.ApplyNonce(ctx, txArgs)
}

func (nth *nonceTransactionsHandler) getOrCreateAddressNonceHandler(address core.AddressHandler) (*addressNonceHandler, error) {
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
func (nth *nonceTransactionsHandlerV2) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	if tx == nil {
		return "", interactors.ErrNilTransaction
	}

	addrAsBech32 := tx.SndAddr
	addressHandler, err := data.NewAddressFromBech32String(addrAsBech32)
	if err != nil {
		return "", fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
	}

	anh, err := nth.getOrCreateAddressNonceHandler(addressHandler)
	if err != nil {
		return "", err
	}

	sentHash, err := anh.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("%w while sending transaction for address %s", err, addrAsBech32)
	}

	return sentHash, nil
}

func (nth *nonceTransactionsHandlerV2) resendTransactionsLoop(ctx context.Context, intervalToResend time.Duration) {
	timer := time.NewTimer(intervalToResend)
	defer timer.Stop()

	for {
		timer.Reset(intervalToResend)

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
