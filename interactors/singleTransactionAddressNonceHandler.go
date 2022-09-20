package interactors

import (
	"context"
	"sync"

	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type singleTransactionAddressNonceHandler struct {
	mut           sync.RWMutex
	address       erdgoCore.AddressHandler
	transaction   *data.Transaction
	computedNonce uint64
	proxy         Proxy
}

func (anh *singleTransactionAddressNonceHandler) ApplyNonce(ctx context.Context, txArgs *data.ArgCreateTransaction, checkForDuplicates bool) error {
	nonce, err := anh.getNonce(ctx)
	if err != nil {
		return err
	}
	txArgs.Nonce = nonce

	return nil
}

func (anh *singleTransactionAddressNonceHandler) getNonce(ctx context.Context) (uint64, error) {
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return 0, err
	}
	anh.computedNonce = account.Nonce

	return account.Nonce, nil
}

func (anh *singleTransactionAddressNonceHandler) ReSendTransactionsIfRequired(ctx context.Context) error {
	if anh.transaction == nil {
		return nil
	}
	account, err := anh.proxy.GetAccount(ctx, anh.address)
	if err != nil {
		return err
	}

	if anh.transaction.Nonce != account.Nonce {
		anh.DropTransactions()
	}

	hash, err := anh.proxy.SendTransaction(ctx, anh.transaction)
	if err != nil {
		return err
	}

	log.Debug("resent transaction", "address", anh.address.AddressAsBech32String(), "hash", hash)

	return nil
}

func (anh *singleTransactionAddressNonceHandler) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	anh.mut.Lock()
	anh.transaction = tx
	anh.mut.Unlock()

	return anh.proxy.SendTransaction(ctx, tx)
}

func (anh *singleTransactionAddressNonceHandler) DropTransactions() {
	anh.mut.Lock()
	defer anh.mut.Unlock()

	anh.transaction = nil
}
