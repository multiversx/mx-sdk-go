package blockchain

import (
	"context"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// Proxy holds the primitive functions that the multiversx proxy engine supports & implements
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	SendTransactions(ctx context.Context, txs []*transaction.FrontendTransaction) ([]string, error)
	GetGuardianData(ctx context.Context, address core.AddressHandler) (*api.GuardianData, error)
	ExecuteVMQuery(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
	FilterLogs(ctx context.Context, filter *core.FilterQuery) ([]string, error)
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
	GetGuardianData(address string) string
	GetRestAPIEntityType() core.RestAPIEntityType
	GetValidatorsInfo(epoch uint32) string
	GetProcessedTransactionStatus(hexHash string) string
	GetESDTTokenData(addressAsBech32 string, tokenIdentifier string) string
	GetNFTTokenData(addressAsBech32 string, tokenIdentifier string, nonce uint64) string
	IsDataTrieMigrated(addressAsBech32 string) string
	GetBlockByNonce(shardID uint32, nonce uint64) string
	GetBlockByHash(shardID uint32, hash string) string
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
	IsInterfaceNil() bool
}

// BlockDataCache defines the methods required for a basic cache.
type BlockDataCache interface {
	Get(key []byte) (value interface{}, ok bool)
	Put(key []byte, value interface{}, sizeInBytes int) (evicted bool)
	IsInterfaceNil() bool
}
