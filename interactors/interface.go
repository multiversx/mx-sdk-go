package interactors

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// Proxy holds the primitive functions that the elrond proxy engine supports & implements
// dependency inversion: blockchain package is considered inner business logic, this package is considered "plugin"
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error)
	IsInterfaceNil() bool
}

// TxBuilder defines the component able to build & sign a transaction
type TxBuilder interface {
	ApplySignatureAndGenerateTx(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error)
	IsInterfaceNil() bool
}

type AddressNonceHandler interface {
	ApplyNonce(ctx context.Context, txArgs *data.ArgCreateTransaction, checkForDuplicates bool) error
	ReSendTransactionsIfRequired(ctx context.Context) error
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	DropTransactions()
}
