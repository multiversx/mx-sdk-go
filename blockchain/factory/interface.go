package factory

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	GetRestAPIEntityType() core.RestAPIEntityType
	IsInterfaceNil() bool
}

// EndpointProvider is able to return endpoint routes strings
type EndpointProvider interface {
	GetNetworkConfig() string
	GetNetworkEconomics() string
	GetRatingsConfig() string
	GetEnableEpochsConfig() string
	GetAccount(addressAsBech32 string) string
	GetAccountKeys(addressAsBech32 string) string
	GetCostTransaction() string
	GetSendTransaction() string
	GetSendMultipleTransactions() string
	GetTransactionStatus(hexHash string) string
	GetTransactionInfo(hexHash string) string
	GetHyperBlockByNonce(nonce uint64) string
	GetHyperBlockByHash(hexHash string) string
	GetVmValues() string
	GetGenesisNodesConfig() string
	GetRawStartOfEpochMetaBlock(epoch uint32) string
	GetNodeStatus(shardID uint32) string
	ShouldCheckShardIDForNodeStatus() bool
	GetRawBlockByHash(shardID uint32, hexHash string) string
	GetRawBlockByNonce(shardID uint32, nonce uint64) string
	GetRawMiniBlockByHash(shardID uint32, hexHash string, epoch uint32) string
	GetRestAPIEntityType() core.RestAPIEntityType
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
	IsInterfaceNil() bool
}
