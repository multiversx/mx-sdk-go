package interactors

import (
	"context"

	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// Proxy holds the primitive functions that the elrond proxy engine supports & implements
// dependency inversion: blockchain package is considered inner business logic, this package is considered "plugin"
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetRatingsConfig(ctx context.Context) (*data.RatingsConfig, error)
	GetEnableEpochsConfig(ctx context.Context) (*data.EnableEpochsConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error)

	GetNonceAtEpochStart(ctx context.Context, shardId uint32) (uint64, error)
	GetCurrentEpoch(ctx context.Context, shardId uint32) (uint64, error)
	GetRawMiniBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error)
	GetRawBlockByNonce(ctx context.Context, shardId uint32, nonce uint64) ([]byte, error)
	GetRawBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error)

	IsInterfaceNil() bool
}

// TxBuilder defines the component able to build & sign a transaction
type TxBuilder interface {
	ApplySignatureAndGenerateTx(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error)
	IsInterfaceNil() bool
}

// HeaderVerifierHandler defines the behaviour of a header verifier instance
type HeaderVerifierHandler interface {
	IsInCache(epoch uint32) bool
	SetNodesConfigPerEpoch(validatorsInfo []*state.ShardValidatorInfo, epoch uint32, randomness []byte) error
	VerifyHeader(header coreData.HeaderHandler) bool
}
