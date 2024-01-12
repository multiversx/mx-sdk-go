package nonceHandlerV3

import (
	"context"
	"fmt"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"

	sdkCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// addressNonceHandler is the handler used for one address. It is able to handle the current
// nonce as max(current_stored_nonce, account_nonce). After each call of the getNonce function
// the current_stored_nonce is incremented. This will prevent "nonce too low in transaction"
// errors on the node interceptor. To prevent the "nonce too high in transaction" error,
// a retrial mechanism is implemented. This struct is able to store all sent transactions,
// having a function that sweeps the map in order to resend a transaction or remove them
// because they were executed. This struct is concurrent safe.
type addressNonceHandler struct {
	mut                    sync.RWMutex
	address                sdkCore.AddressHandler
	proxy                  interactors.Proxy
	computedNonceWasSet    bool
	computedNonce          uint64
	lowestNonce            uint64
	gasPrice               uint64
	nonceUntilGasIncreased uint64
	transactions           *TransactionQueueHandler

	hashChannel chan string
}

// NewAddressNonceHandler returns a new instance of a addressNonceHandler
func NewAddressNonceHandler(proxy interactors.Proxy, address sdkCore.AddressHandler) (interactors.AddressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: NewTransactionQueueHandler(),

		hashChannel: make(chan string),
	}, nil
}

// ApplyNonceAndGasPrice will apply the computed nonce to the given FrontendTransaction
func (anh *addressNonceHandler) ApplyNonceAndGasPrice(ctx context.Context, tx *transaction.FrontendTransaction) error {
	oldTx := anh.getOlderTxWithSameNonce(tx)
	if oldTx != nil {
		err := anh.handleTxWithSameNonce(oldTx, tx)
		if err != nil {
			return err
		}
	}

	nonce, err := anh.getNonceUpdatingCurrent(ctx)
	tx.Nonce = nonce
	if err != nil {
		return err
	}

	anh.fetchGasPriceIfRequired(ctx, nonce)
	tx.GasPrice = core.MaxUint64(anh.gasPrice, tx.GasPrice)
	return nil
}

func (anh *addressNonceHandler) handleTxWithSameNonce(oldTx *transaction.FrontendTransaction, tx *transaction.FrontendTransaction) error {
	if oldTx.GasPrice < tx.GasPrice {
		return nil
	}

	if oldTx.GasPrice == tx.GasPrice && oldTx.GasPrice < anh.gasPrice {
		return nil
	}

	return interactors.ErrTxWithSameNonceAndGasPriceAlreadySent
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

	if anh.lowestNonce > account.Nonce {
		return account.Nonce, interactors.ErrGapNonce
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

// ReSendTransactionsIfRequired will resend the cached transactions that still have a nonce greater that the one fetched from the blockchain
func (anh *addressNonceHandler) ReSendTransactionsIfRequired(ctx context.Context) error {
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return err
	}

	anh.mut.Lock()
	if account.Nonce == anh.computedNonce {
		anh.lowestNonce = anh.computedNonce
		anh.transactions = NewTransactionQueueHandler()
		anh.mut.Unlock()

		return nil
	}

	//resendableTxs := make([]*transaction.FrontendTransaction, 0, anh.transactions.tq.Len())
	//minNonce := anh.computedNonce

	addressAsBech32String, err := anh.address.AddressAsBech32String()
	if err != nil {
		return fmt.Errorf("failed to retrieve address as bech32: %w", err)
	}

	for anh.transactions.Len() > 0 {
		t := anh.transactions.NextTransaction()

		if t.Nonce > account.Nonce {
			hash, err := anh.proxy.SendTransaction(ctx, t)
			if err != nil {
				log.Error("failed to send transaction", "error", err.Error())
			}
			anh.hashChannel <- hash
			log.Info(fmt.Sprintf("successfully resent transaction with nonce %d for address %q", t.Nonce, addressAsBech32String), "hash", hash)
		}
	}

	//anh.transactions.NextTransaction()
	//for txNonce, tx := range anh.transactions {
	//	if txNonce <= account.Nonce {
	//		delete(anh.transactions, txNonce)
	//		continue
	//	}
	//	minNonce = core.MinUint64(txNonce, minNonce)
	//	resendableTxs = append(resendableTxs, tx)
	//}
	//anh.lowestNonce = minNonce
	anh.mut.Unlock()

	//if len(resendableTxs) == 0 {
	//	return nil
	//}
	//
	//hashes, err := anh.proxy.SendTransactions(ctx, resendableTxs)
	//if err != nil {
	//	return err
	//}
	//
	//addressAsBech32String, err := anh.address.AddressAsBech32String()
	//if err != nil {
	//	return err
	//}
	//
	//log.Info("resent transactions", "address", addressAsBech32String, "total txs", len(resendableTxs), "received hashes", len(hashes))

	return nil
}

// SendTransaction will save and propagate a transaction to the network
func (anh *addressNonceHandler) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	anh.mut.Lock()
	anh.transactions.AddTransaction(tx)
	anh.mut.Unlock()

	hash := <-anh.hashChannel

	return hash, nil
}

// DropTransactions will delete the cached transactions and will try to replace the current transactions from the pool using more gas price
func (anh *addressNonceHandler) DropTransactions() {
	//anh.mut.Lock()
	//anh.transactions = make(map[uint64]*transaction.FrontendTransaction)
	//anh.computedNonceWasSet = false
	//anh.gasPrice++
	//anh.nonceUntilGasIncreased = anh.computedNonce
	//anh.mut.Unlock()
}

func (anh *addressNonceHandler) getOlderTxWithSameNonce(tx *transaction.FrontendTransaction) *transaction.FrontendTransaction {
	anh.mut.RLock()
	defer anh.mut.RUnlock()

	return anh.transactions.SearchForNonce(tx.Nonce)
}

// IsInterfaceNil returns true if there is no value under the interface
func (anh *addressNonceHandler) IsInterfaceNil() bool {
	return anh == nil
}
