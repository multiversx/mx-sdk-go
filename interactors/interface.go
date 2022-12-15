package interactors

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain/cryptoProvider"
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
	ApplySignatureAndGenerateTx(cryptoHolder cryptoProvider.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error)
	IsInterfaceNil() bool
}

// AddressNonceHandler defines the component able to handler address nonces
type AddressNonceHandler interface {
	ApplyNonceAndGasPrice(ctx context.Context, txArgs *data.ArgCreateTransaction) error
	ReSendTransactionsIfRequired(ctx context.Context) error
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	DropTransactions()
	IsInterfaceNil() bool
}

// TransactionNonceHandlerV1 defines the component able to manage transaction nonces
type TransactionNonceHandlerV1 interface {
	GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	ForceNonceReFetch(address core.AddressHandler) error
	Close() error
	IsInterfaceNil() bool
}

// TransactionNonceHandlerV2 defines the component able to apply nonce for a given ArgCreateTransaction
type TransactionNonceHandlerV2 interface {
	ApplyNonceAndGasPrice(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	Close() error
	IsInterfaceNil() bool
}

// AddressNonceHandlerCreator defines the component able to create AddressNonceHandler instances
type AddressNonceHandlerCreator interface {
	Create(proxy Proxy, address core.AddressHandler) (AddressNonceHandler, error)
	IsInterfaceNil() bool
}
