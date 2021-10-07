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
// nonceTransactionsHandler should be terminated an collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandler struct {
	proxy       Proxy
	mutHandlers sync.Mutex
	handlers    map[string]*addressNonceHandler
	cancelFunc  func()
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
		proxy:    proxy,
		handlers: make(map[string]*addressNonceHandler),
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	nth.cancelFunc = cancelFunc
	go nth.resendTransactionsLoop(ctx, intervalToResend)

	return nth, nil
}

// GetNonce will return the nonce for the provided address
func (nth *nonceTransactionsHandler) GetNonce(address core.AddressHandler) (uint64, error) {
	if check.IfNil(address) {
		return 0, ErrNilAddress
	}

	anh := nth.getOrCreateAddressNonceHandler(address)

	return anh.getNonceUpdatingCurrent()
}

func (nth *nonceTransactionsHandler) getOrCreateAddressNonceHandler(address core.AddressHandler) *addressNonceHandler {
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

// SendTransaction will send and store the provided transaction
func (nth *nonceTransactionsHandler) SendTransaction(tx *data.Transaction) (string, error) {
	if tx == nil {
		return "", ErrNilTransaction
	}

	addrAsBech32 := tx.SndAddr
	addressHandler, err := data.NewAddressFromBech32String(addrAsBech32)
	if err != nil {
		return "", fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
	}

	anh := nth.getOrCreateAddressNonceHandler(addressHandler)
	sentHash, err := anh.sendTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("%w while sending transactions for address %s", err, addrAsBech32)
	}

	return sentHash, nil
}

func (nth *nonceTransactionsHandler) resendTransactionsLoop(ctx context.Context, intervalToResend time.Duration) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("finishing nonceTransactionsHandler.resendTransactionsLoop...")
			return
		case <-time.After(intervalToResend):
			nth.resendTransactions(ctx)
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

		err := anh.reSendTransactionsIfRequired()
		log.LogIfError(err)
	}
}

// ForceNonceReFetch will mark the addressNonceHandler to re-fetch its nonce from the blockchain account.
// This should be only used in a fallback plan, when some transactions are completely lost (or due to a bug, not even sent in first time)
func (nth *nonceTransactionsHandler) ForceNonceReFetch(address core.AddressHandler) error {
	if check.IfNil(address) {
		return ErrNilAddress
	}

	anh := nth.getOrCreateAddressNonceHandler(address)
	anh.markReFetchNonce()

	return nil
}

// Close finishes the transactions resend go routine
func (nth *nonceTransactionsHandler) Close() error {
	nth.cancelFunc()

	return nil
}
