package nonceHandlerV3

import (
	"context"
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
	mut                    sync.Mutex
	address                sdkCore.AddressHandler
	proxy                  interactors.Proxy
	computedNonceWasSet    bool
	computedNonce          uint64
	gasPrice               uint64
	nonceUntilGasIncreased uint64
	transactionWorker      *workers.TransactionWorker
}

// NewAddressNonceHandlerV2 returns a new instance of a addressNonceHandler
func NewAddressNonceHandlerV2(proxy interactors.Proxy, address sdkCore.AddressHandler, pollingInterval time.Duration) (interactors.AddressNonceHandlerV2, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}

	anh := &addressNonceHandler{
		mut:               sync.Mutex{},
		address:           address,
		proxy:             proxy,
		transactionWorker: workers.NewTransactionWorker(proxy, pollingInterval),
	}

	return anh, nil
}

// ApplyNonceAndGasPrice will apply the computed nonce to the given FrontendTransaction
func (anh *addressNonceHandler) ApplyNonceAndGasPrice(ctx context.Context, txs ...*transaction.FrontendTransaction) error {
	for _, tx := range txs {
		nonce, err := anh.getNonceUpdatingCurrent(ctx)
		tx.Nonce = nonce
		if err != nil {
			return err
		}

		anh.fetchGasPriceIfRequired(ctx, nonce)
		tx.GasPrice = core.MaxUint64(anh.gasPrice, tx.GasPrice)
	}

	anh.computedNonceWasSet = false
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

func (anh *addressNonceHandler) fetchGasPriceIfRequired(ctx context.Context, nonce uint64) {
	if nonce == anh.nonceUntilGasIncreased+1 || anh.gasPrice == 0 {
		networkConfig, err := anh.proxy.GetNetworkConfig(ctx)

		anh.mut.Lock()
		defer anh.mut.Unlock()
		if err != nil {
			log.Error("%w: while fetching network config", err)
			anh.gasPrice = 0
			return
		}
		anh.gasPrice = networkConfig.MinGasPrice
	}
}

func (anh *addressNonceHandler) getNonceUpdatingCurrent(ctx context.Context) (uint64, error) {
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return 0, err
	}

	anh.mut.Lock()
	defer anh.mut.Unlock()

	if !anh.computedNonceWasSet {
		anh.computedNonce = account.Nonce
		anh.computedNonceWasSet = true

		return anh.computedNonce, nil
	}

	anh.computedNonce++

	return core.MaxUint64(anh.computedNonce, account.Nonce), nil
}
