package workflows

import (
	"context"

	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TrackableAddressesProvider defines the behavior for a component that can determine if one address is tracked or not
type TrackableAddressesProvider interface {
	IsTrackableAddresses(addressAsBech32 string) bool
	PrivateKeyOfBech32Address(addressAsBech32 string) []byte
	IsInterfaceNil() bool
}

// LastProcessedNonceHandler will keep track of the last processed hyper block nonce.
// fraction of hyper blocks sets after an application restart
type LastProcessedNonceHandler interface {
	ProcessedNonce(nonce uint64)
	GetLastProcessedNonce() uint64
	IsInterfaceNil() bool
}

// ProxyHandler defines the behavior of a proxy handler that can process requests
type ProxyHandler interface {
	GetLatestHyperBlockNonce(ctx context.Context) (uint64, error)
	GetHyperBlockByNonce(ctx context.Context, nonce uint64) (*data.HyperBlock, error)
	GetDefaultTransactionArguments(ctx context.Context, address erdgoCore.AddressHandler, networkConfigs *data.NetworkConfig) (data.ArgCreateTransaction, error)
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	IsInterfaceNil() bool
}

// TransactionInteractor defines the transaction interactor behavior used in workflows
type TransactionInteractor interface {
	AddTransaction(tx *data.Transaction)
	ApplyUserSignatureAndGenerateTx(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error)
	IsInterfaceNil() bool
}
