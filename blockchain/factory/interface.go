package factory

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	IsInterfaceNil() bool
}

// EndpointProvider is able to return endpoint routes strings
type EndpointProvider interface {
	GetNetworkConfigEndpoint() string
	GetNetworkEconomicsEndpoint() string
	GetRatingsConfigEndpoint() string
	GetEnableEpochsConfigEndpoint() string
	GetAccountEndpoint(addressAsBech32 string) string
	GetCostTransactionEndpoint() string
	GetSendTransactionEndpoint() string
	GetSendMultipleTransactionsEndpoint() string
	GetTransactionStatusEndpoint(hexHash string) string
	GetTransactionInfoEndpoint(hexHash string) string
	GetHyperBlockByNonceEndpoint(nonce uint64) string
	GetHyperBlockByHashEndpoint(hexHash string) string
	GetVmValuesEndpoint() string
	GetGenesisNodesConfigEndpoint() string
	GetRawStartOfEpochMetaBlock(epoch uint32) string
	GetNodeStatusEndpoint(shardID uint32) string
	GetRawBlockByHashEndpoint(shardID uint32, hexHash string) string
	GetRawBlockByNonceEndpoint(shardID uint32, nonce uint64) string
	GetRawMiniBlockByHashEndpoint(shardID uint32, hexHash string, epoch uint32) string
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
	IsInterfaceNil() bool
}
