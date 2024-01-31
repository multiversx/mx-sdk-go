package nonceHandlerV3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"
	"golang.org/x/sync/errgroup"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
)

const minimumIntervalToResend = time.Second

var log = logger.GetOrCreate("mx-sdk-go/interactors/nonceHandlerV3")

// ArgsNonceTransactionsHandlerV3 is the argument DTO for a nonce workers handler component
type ArgsNonceTransactionsHandlerV3 struct {
	Proxy           interactors.Proxy
	PollingInterval time.Duration
}

// nonceTransactionsHandlerV3 is the handler used for an unlimited number of addresses.
// It basically contains a map of addressNonceHandler, creating new entries on the first
// access of a provided address. This struct delegates all the operations on the right
// instance of addressNonceHandler. It also starts a go routine that will periodically
// try to resend "stuck workers" and to clean the inner state. The recommended resend
// interval is 1 minute. The Close method should be called whenever the current instance of
// nonceTransactionsHandlerV3 should be terminated and collected by the GC.
// This struct is concurrent safe.
type nonceTransactionsHandlerV3 struct {
	proxy           interactors.Proxy
	mutHandlers     sync.RWMutex
	handlers        map[string]interactors.AddressNonceHandlerV3
	pollingInterval time.Duration
}

// NewNonceTransactionHandlerV3 will create a new instance of the nonceTransactionsHandlerV3. It requires a Proxy implementation
// and an interval at which the workers sent are rechecked and eventually, resent.
func NewNonceTransactionHandlerV3(args ArgsNonceTransactionsHandlerV3) (*nonceTransactionsHandlerV3, error) {
	if check.IfNil(args.Proxy) {
		return nil, interactors.ErrNilProxy
	}
	if args.PollingInterval < minimumIntervalToResend {
		return nil, fmt.Errorf("%w for pollingInterval in NewNonceTransactionHandlerV2", interactors.ErrInvalidValue)
	}

	nth := &nonceTransactionsHandlerV3{
		proxy:           args.Proxy,
		handlers:        make(map[string]interactors.AddressNonceHandlerV3),
		pollingInterval: args.PollingInterval,
	}

	return nth, nil
}

// ApplyNonceAndGasPrice will apply the nonce to the given frontend transaction
func (nth *nonceTransactionsHandlerV3) ApplyNonceAndGasPrice(ctx context.Context, tx ...*transaction.FrontendTransaction) error {
	if tx == nil {
		return interactors.ErrNilTransaction
	}

	mapAddressTransactions := nth.filterTransactionsBySenderAddress(tx)

	for addressRawString, transactions := range mapAddressTransactions {
		address, err := data.NewAddressFromBech32String(addressRawString)
		if err != nil {
			return err
		}
		anh, err := nth.getOrCreateAddressNonceHandler(address)
		if err != nil {
			return err
		}

		err = anh.ApplyNonceAndGasPrice(ctx, transactions...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (nth *nonceTransactionsHandlerV3) getOrCreateAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandlerV3, error) {
	anh := nth.getAddressNonceHandler(address)
	if !check.IfNil(anh) {
		return anh, nil
	}

	return nth.createAddressNonceHandler(address)
}

func (nth *nonceTransactionsHandlerV3) getAddressNonceHandler(address core.AddressHandler) interactors.AddressNonceHandlerV3 {
	nth.mutHandlers.RLock()
	defer nth.mutHandlers.RUnlock()

	anh, found := nth.handlers[string(address.AddressBytes())]
	if found {
		return anh
	}
	return nil
}

func (nth *nonceTransactionsHandlerV3) createAddressNonceHandler(address core.AddressHandler) (interactors.AddressNonceHandlerV3, error) {
	nth.mutHandlers.Lock()
	defer nth.mutHandlers.Unlock()

	addressAsString := string(address.AddressBytes())
	anh, found := nth.handlers[addressAsString]
	if found {
		return anh, nil
	}

	anh, err := NewAddressNonceHandlerV3(nth.proxy, address, nth.pollingInterval)
	if err != nil {
		return nil, err
	}
	nth.handlers[addressAsString] = anh

	return anh, nil
}

func (nth *nonceTransactionsHandlerV3) filterTransactionsBySenderAddress(transactions []*transaction.FrontendTransaction) map[string][]*transaction.FrontendTransaction {
	filterMap := make(map[string][]*transaction.FrontendTransaction)
	for _, tx := range transactions {
		if _, ok := filterMap[tx.Sender]; !ok {
			transactionsPerAddress := make([]*transaction.FrontendTransaction, 0)
			transactionsPerAddress = append(transactionsPerAddress, tx)
			filterMap[tx.Sender] = transactionsPerAddress
		} else {
			filterMap[tx.Sender] = append(filterMap[tx.Sender], tx)
		}
	}

	return filterMap
}

// SendTransactions will store and send the provided transaction
func (nth *nonceTransactionsHandlerV3) SendTransactions(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	sentHashes := make([]string, len(txs))
	for i, tx := range txs {
		if tx == nil {
			return nil, interactors.ErrNilTransaction
		}

		// Work with a full copy of the provided transaction so the provided one can change without affecting this component.
		// Abnormal and unpredictable behaviors due to the resending mechanism are prevented this way
		txCopy := *tx

		addrAsBech32 := txCopy.Sender
		address, err := data.NewAddressFromBech32String(addrAsBech32)
		if err != nil {
			return nil, fmt.Errorf("%w while creating address handler for string %s", err, addrAsBech32)
		}

		anh, err := nth.getOrCreateAddressNonceHandler(address)
		if err != nil {
			return nil, err
		}

		i := i
		g.Go(func() error {
			sentHash, err := anh.SendTransaction(ctx, &txCopy)
			if err != nil {
				return fmt.Errorf("%w while sending transaction for address %s", err, addrAsBech32)
			}

			sentHashes[i] = sentHash
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return sentHashes, err
	}

	return sentHashes, nil
}

// Close will cancel all related processes.
func (nth *nonceTransactionsHandlerV3) Close() {
	nth.mutHandlers.RLock()
	defer nth.mutHandlers.RUnlock()
	for _, handler := range nth.handlers {
		handler.Close()
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (nth *nonceTransactionsHandlerV3) IsInterfaceNil() bool {
	return nth == nil
}
