package nonceHandlerV2

import (
	"bytes"
	"context"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

// addressNonceHandler is the handler used for one address. It is able to handle the current
// nonce as max(current_stored_nonce, account_nonce). After each call of the getNonce function
// the current_stored_nonce is incremented. This will prevent "nonce too low in transaction"
// errors on the node interceptor. To prevent the "nonce too high in transaction" error,
// a retrial mechanism is implemented. This struct is able to store all sent transactions,
// having a function that sweeps the map in order to resend a transaction or remove them
// because they were executed. This struct is concurrent safe.
type addressNonceHandler struct {
	mut                 sync.RWMutex
	address             erdgoCore.AddressHandler
	proxy               interactors.Proxy
	computedNonceWasSet bool
	computedNonce       uint64
	lowestNonce         uint64
	transactions        map[uint64]*data.Transaction
}

func newAddressNonceHandler(proxy interactors.Proxy, address erdgoCore.AddressHandler) *addressNonceHandler {
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: make(map[uint64]*data.Transaction),
	}
}

func (anh *addressNonceHandler) ApplyNonce(ctx context.Context, txArgs *data.ArgCreateTransaction, checkForDuplicates bool) error {
	if checkForDuplicates && anh.isTxAlreadySent(txArgs) {
		// TODO: add gas comparation logic EN-11887
		return interactors.ErrTxAlreadySent
	}
	nonce, err := anh.getNonceUpdatingCurrent(ctx)
	if err != nil {
		return err
	}
	txArgs.Nonce = nonce
	return nil
}

func (anh *addressNonceHandler) getNonceUpdatingCurrent(ctx context.Context) (uint64, error) {
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return 0, err
	}

	if anh.lowestNonce > account.Nonce {
		return account.Nonce, nil
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
	minNonce := uint64(0)
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

func (anh *addressNonceHandler) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	anh.mut.Lock()
	anh.transactions[tx.Nonce] = tx
	anh.mut.Unlock()

	return anh.proxy.SendTransaction(ctx, tx)
}

func (anh *addressNonceHandler) DropTransactions() {
	anh.mut.Lock()
	anh.transactions = make(map[uint64]*data.Transaction)
	anh.computedNonceWasSet = false
	anh.mut.Unlock()
}

func (anh *addressNonceHandler) isTxAlreadySent(tx *data.ArgCreateTransaction) bool {
	anh.mut.RLock()
	defer anh.mut.RUnlock()
	for _, oldTx := range anh.transactions {
		isTheSameReceiverDataValue := oldTx.RcvAddr == tx.RcvAddr &&
			bytes.Equal(oldTx.Data, tx.Data) &&
			oldTx.Value == tx.Value
		if isTheSameReceiverDataValue {
			return true
		}
	}
	return false
}
