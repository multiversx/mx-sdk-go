package nonceHandlerV3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"

	sdkCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/interactors/nonceHandlerV3/workers"
)

// addressNonceHandler is the handler used for one address. It is able to handle the current
// nonce as max(current_stored_nonce, account_nonce). After each call of the getNonce function
// the current_stored_nonce is incremented. This will prevent "nonce too low in transaction"
// errors on the node interceptor. To prevent the "nonce too high in transaction" error,
// a retrial mechanism is implemented. This struct is able to store all sent transactions,
// having a function that sweeps the map in order to resend a transaction or remove them
// because they were executed. This struct is concurrent safe.
type addressNonceHandler struct {
	mut               sync.Mutex
	address           sdkCore.AddressHandler
	proxy             interactors.Proxy
	gasPrice          uint64
	transactionWorker *workers.TransactionWorker
	cancelFunc        func()
}

// NewAddressNonceHandlerV3 returns a new instance of a addressNonceHandler
func NewAddressNonceHandlerV3(proxy interactors.Proxy, address sdkCore.AddressHandler, pollingInterval time.Duration) (*addressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	anh := &addressNonceHandler{
		mut:               sync.Mutex{},
		address:           address,
		proxy:             proxy,
		transactionWorker: workers.NewTransactionWorker(ctx, proxy, pollingInterval),
		cancelFunc:        cancelFunc,
	}

	return anh, nil
}

// ApplyNonceAndGasPrice will apply the computed nonce to the given FrontendTransaction
func (anh *addressNonceHandler) ApplyNonceAndGasPrice(ctx context.Context, txs ...*transaction.FrontendTransaction) error {
	nonce, err := anh.fetchNonce(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch nonce: %w", err)
	}

	for i, tx := range txs {
		anh.mut.Lock()
		tx.Nonce = nonce + uint64(i)
		anh.mut.Unlock()

		anh.applyGasPriceIfRequired(ctx, tx)
	}

	return nil
}

// SendTransaction will save and propagate a transaction to the network
func (anh *addressNonceHandler) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	ch := anh.transactionWorker.AddTransaction(tx)

	select {
	case response := <-ch:
		return response.TxHash, response.Error

	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (anh *addressNonceHandler) IsInterfaceNil() bool {
	return anh == nil
}

// Close will cancel all related processes..
func (anh *addressNonceHandler) Close() {
	anh.cancelFunc()
}

func (anh *addressNonceHandler) applyGasPriceIfRequired(ctx context.Context, tx *transaction.FrontendTransaction) {
	anh.mut.Lock()
	defer anh.mut.Unlock()

	if anh.gasPrice == 0 {
		networkConfig, err := anh.proxy.GetNetworkConfig(ctx)

		if err != nil {
			log.Error("%w: while fetching network config", err)
			anh.gasPrice = 0
			return
		}
		anh.gasPrice = networkConfig.MinGasPrice
	}
	tx.GasPrice = core.MaxUint64(anh.gasPrice, tx.GasPrice)
}

func (anh *addressNonceHandler) fetchNonce(ctx context.Context) (uint64, error) {
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return 0, err
	}

	return account.Nonce, nil
}
