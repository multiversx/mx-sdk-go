package notifees

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxBuilder defines the component able to build & sign a transaction
type TxBuilder interface {
	ApplyUserSignatureAndGenerateTx(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error)
	IsInterfaceNil() bool
}

// Proxy holds the primitive functions that the elrond proxy engine supports & implements
// dependency inversion: blockchain package is considered inner business logic, this package is considered "plugin"
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error)
	IsInterfaceNil() bool
}

// TransactionNonceHandler is able to handle the
type TransactionNonceHandler interface {
	GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	IsInterfaceNil() bool
}
