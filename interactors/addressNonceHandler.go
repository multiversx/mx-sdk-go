package interactors

import (
	"bytes"
	"context"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

//TODO EN-13182: create a baseAddressNonceHandler component that can remove the duplicate code as much as possible from the
// addressNonceHandler and singleTransactionAddressNonceHandler

// addressNonceHandler is the handler used for one address. It is able to handle the current
// nonce as max(current_stored_nonce, account_nonce). After each call of the getNonce function
// the current_stored_nonce is incremented. This will prevent "nonce too low in transaction"
// errors on the node interceptor. To prevent the "nonce too high in transaction" error,
// a retrial mechanism is implemented. This struct is able to store all sent transactions,
// having a function that sweeps the map in order to resend a transaction or remove them
// because they were executed. This struct is concurrent safe.
type addressNonceHandler struct {
	mut                    sync.RWMutex
	address                erdgoCore.AddressHandler
	proxy                  Proxy
	computedNonceWasSet    bool
	computedNonce          uint64
	lowestNonce            uint64
	gasPrice               uint64
	nonceUntilGasIncreased uint64
	transactions           map[uint64]*data.Transaction
}

// NewAddressNonceHandler returns a new instance of a addressNonceHandler
func NewAddressNonceHandler(proxy Proxy, address erdgoCore.AddressHandler) (*addressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: make(map[uint64]*data.Transaction),
	}, nil
}

// ApplyNonce will apply the computed nonce to the given ArgCreateTransaction
func (anh *addressNonceHandler) ApplyNonce(ctx context.Context, txArgs *data.ArgCreateTransaction) error {
	oldTx, alreadyExists := anh.isTxAlreadySent(txArgs)
	if alreadyExists {
		err := anh.handleTxAlreadyExists(oldTx, txArgs)
		if err != nil {
			return err
		}
	}

	nonce, err := anh.getNonceUpdatingCurrent(ctx)
	txArgs.Nonce = nonce
	if err != nil {
		return err
	}

	anh.fetchGasPriceIfRequired(ctx, nonce)
	txArgs.GasPrice = core.MaxUint64(anh.gasPrice, txArgs.GasPrice)
	return nil
}

func (anh *addressNonceHandler) handleTxAlreadyExists(oldTx *data.Transaction, txArgs *data.ArgCreateTransaction) error {
	if oldTx.GasPrice < txArgs.GasPrice {
		return nil
	}

	if oldTx.GasPrice == txArgs.GasPrice && oldTx.GasPrice < anh.gasPrice {
		return nil
	}

	return ErrTxAlreadySent
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
		return account.Nonce, ErrGapNonce
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
		anh.transactions = make(map[uint64]*data.Transaction)
		anh.mut.Unlock()

		return nil
	}

	resendableTxs := make([]*data.Transaction, 0, len(anh.transactions))
	minNonce := anh.computedNonce
	for txNonce, tx := range anh.transactions {
		if txNonce <= account.Nonce {
			delete(anh.transactions, txNonce)
			continue
		}
		minNonce = core.MinUint64(txNonce, minNonce)
		resendableTxs = append(resendableTxs, tx)
	}
	anh.lowestNonce = minNonce
	anh.mut.Unlock()

	if len(resendableTxs) == 0 {
		return nil
	}

	hashes, err := anh.proxy.SendTransactions(ctx, resendableTxs)
	if err != nil {
		return err
	}

	log.Debug("resent transactions", "address", anh.address.AddressAsBech32String(), "total txs", len(resendableTxs), "received hashes", len(hashes))

	return nil
}

// SendTransaction will save and propagate a transaction to the network
func (anh *addressNonceHandler) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	anh.mut.Lock()
	anh.transactions[tx.Nonce] = tx
	anh.mut.Unlock()

	return anh.proxy.SendTransaction(ctx, tx)
}

// DropTransactions will delete the cached transactions and will try to replace the current transactions from the pool using more gas price
func (anh *addressNonceHandler) DropTransactions() {
	anh.mut.Lock()
	anh.transactions = make(map[uint64]*data.Transaction)
	anh.computedNonceWasSet = false
	anh.gasPrice++
	anh.nonceUntilGasIncreased = anh.computedNonce
	anh.mut.Unlock()
}

func (anh *addressNonceHandler) isTxAlreadySent(tx *data.ArgCreateTransaction) (*data.Transaction, bool) {
	anh.mut.RLock()
	defer anh.mut.RUnlock()
	for _, oldTx := range anh.transactions {
		isTheSameReceiverDataValue := oldTx.RcvAddr == tx.RcvAddr &&
			bytes.Equal(oldTx.Data, tx.Data) &&
			oldTx.Value == tx.Value
		if isTheSameReceiverDataValue {
			return oldTx, true
		}
	}
	return nil, false
}
