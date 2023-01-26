package blockchain

import (
	"context"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// Proxy holds the primitive functions that the multiversx proxy engine supports & implements
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error)
	ExecuteVMQuery(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
	IsInterfaceNil() bool
}

type httpClientWrapper interface {
	GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
	IsInterfaceNil() bool
}

// EndpointProvider is able to return endpoint routes strings
type EndpointProvider interface {
	GetNetworkConfig() string
	GetNetworkEconomics() string
	GetRatingsConfig() string
	GetEnableEpochsConfig() string
	GetAccount(addressAsBech32 string) string
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
	GetValidatorsInfo(epoch uint32) string
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
	IsInterfaceNil() bool
}
