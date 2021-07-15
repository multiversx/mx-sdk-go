package blockchain

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// ProxyHandler will handle different implementations of proxy client
type ProxyHandler interface {
	ExecuteVMQuery(vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
	GetNetworkConfig() (*data.NetworkConfig, error)
	GetNetworkEconomics() (*data.NetworkEconomics, error)
	GetAccount(address core.AddressHandler) (*data.Account, error)
	SendTransaction(tx *data.Transaction) (string, error)
	RequestTransactionCost(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatus(hash string) (string, error)
	GetTransactionInfo(hash string) (*data.TransactionInfo, error)
	GetTransactionInfoWithResults(hash string) (*data.TransactionInfo, error)
	GetLatestHyperblockNonce() (uint64, error)
	GetHyperblockByNonce(nonce uint64) (*data.Hyperblock, error)
	GetHyperblockByHash(hash string) (*data.Hyperblock, error)
}
