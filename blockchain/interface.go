package blockchain

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// Proxy holds the primitive functions that the elrond proxy engine supports & implements
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
