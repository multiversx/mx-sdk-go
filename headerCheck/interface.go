package headerCheck

import (
	"context"

	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// Proxy holds the behaviour needed for header verifier in order to interact with proxy
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetRatingsConfig(ctx context.Context) (*data.RatingsConfig, error)
	GetEnableEpochsConfig(ctx context.Context) (*data.EnableEpochsConfig, error)
	GetNonceAtEpochStart(ctx context.Context, shardId uint32) (uint64, error)
	GetRawMiniBlockByHash(ctx context.Context, shardId uint32, hash string, epoch uint32) ([]byte, error)
	GetRawBlockByNonce(ctx context.Context, shardId uint32, nonce uint64) ([]byte, error)
	GetRawBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error)
	GetRawStartOfEpochMetaBlock(ctx context.Context, epoch uint32) ([]byte, error)
	GetGenesisNodesConfig(ctx context.Context) (*data.GenesisNodesConfig, error)
	IsInterfaceNil() bool
}

// RawHeaderHandler holds the behaviour needed to handler raw header data from proxy
type RawHeaderHandler interface {
	GetMetaBlockByHash(ctx context.Context, hash string) (coreData.MetaHeaderHandler, error)
	GetShardBlockByHash(ctx context.Context, shardId uint32, hash string) (coreData.HeaderHandler, error)
	GetValidatorsInfoPerEpoch(ctx context.Context, epoch uint32) ([]*state.ShardValidatorInfo, []byte, error)
	IsInterfaceNil() bool
}

// HeaderVerifier defines the functions needed for verifying headers
type HeaderVerifier interface {
	VerifyHeaderByHash(ctx context.Context, shardId uint32, hash string) (bool, error)
	IsInterfaceNil() bool
}

// HeaderSigVerifierHandler defines the functions needed to verify headers signature
type HeaderSigVerifierHandler interface {
	VerifySignature(header coreData.HeaderHandler) error
	IsInterfaceNil() bool
}
