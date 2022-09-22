package interactors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

const minimumIntervalToResend = time.Second

// nonceTransactionsHandler is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck transactions" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandler should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandler struct {
	proxy            Proxy
	mutHandlers      sync.Mutex
	handlers         map[string]*addressNonceHandler
	cancelFunc       func()
	intervalToResend time.Duration
}

// NewNonceTransactionHandler will create a new instance of the nonceTransactionsHandler. It requires a Proxy implementation
// and an interval at which the transactions sent are rechecked and eventually, resent.
func NewNonceTransactionHandler(proxy Proxy, intervalToResend time.Duration) (*nonceTransactionsHandler, error) {
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
func (nth *nonceTransactionsHandler) ApplyNonce(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
	if check.IfNil(address) {
		return ErrNilAddress
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
func (nth *nonceTransactionsHandler) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	if tx == nil {
		return "", ErrNilTransaction
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

func (nth *nonceTransactionsHandler) resendTransactionsLoop(ctx context.Context, intervalToResend time.Duration) {
	timer := time.NewTimer(intervalToResend)
	defer timer.Stop()

	for {
		timer.Reset(intervalToResend)

		select {
		case <-timer.C:
			nth.resendTransactions(ctx)
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandler.resendTransactionsLoop...")
			return
		}
	}
}

func (nth *nonceTransactionsHandler) resendTransactions(ctx context.Context) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	for _, anh := range nth.handlers {
		select {
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandler.resendTransactions...")
			return
		default:
		}

		resendCtx, cancel := context.WithTimeout(ctx, nth.intervalToResend)
		err := anh.ReSendTransactionsIfRequired(resendCtx)
		log.LogIfError(err)
		cancel()
	}
}

// DropTransactions will clean the addressNonceHandler cached transactions and will re-fetch its nonce from the blockchain account.
// This should be only used in a fallback plan, when some transactions are completely lost (or due to a bug, not even sent in first time)
func (nth *nonceTransactionsHandler) DropTransactions(address core.AddressHandler) error {
	if check.IfNil(address) {
		return ErrNilAddress
	}

	anh, err := nth.getOrCreateAddressNonceHandler(address)
	if err != nil {
		return err
	}
	anh.DropTransactions()

	return nil
}

// Close finishes the transactions resend go routine
func (nth *nonceTransactionsHandler) Close() error {
	nth.cancelFunc()

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nth *nonceTransactionsHandler) IsInterfaceNil() bool {
	return nth == nil
}
