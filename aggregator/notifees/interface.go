package notifees

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// TxBuilder defines the component able to build & sign a transaction
type TxBuilder interface {
	ApplySignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNil() bool
}

// Proxy holds the primitive functions that the multiversx proxy engine supports & implements
// dependency inversion: blockchain package is considered inner business logic, this package is considered "plugin"
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	SendTransactions(ctx context.Context, txs []*transaction.FrontendTransaction) ([]string, error)
	IsInterfaceNil() bool
}

// TransactionNonceHandler defines the component able to apply nonce for a given FrontendTransaction
type TransactionNonceHandler interface {
	ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, tx *transaction.FrontendTransaction) error
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	IsInterfaceNil() bool
}
