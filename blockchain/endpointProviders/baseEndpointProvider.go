package endpointProviders

import "fmt"

const (
	networkConfigEndpoint            = "network/config"
	networkEconomicsEndpoint         = "network/economics"
	ratingsConfigEndpoint            = "network/ratings"
	enableEpochsConfigEndpoint       = "network/enable-epochs"
	accountEndpoint                  = "address/%s"
	costTransactionEndpoint          = "transaction/cost"
	sendTransactionEndpoint          = "transaction/send"
	sendMultipleTransactionsEndpoint = "transaction/send-multiple"
	transactionStatusEndpoint        = "transaction/%s/status"
	transactionInfoEndpoint          = "transaction/%s"
	hyperBlockByNonceEndpoint        = "hyperblock/by-nonce/%d"
	hyperBlockByHashEndpoint         = "hyperblock/by-hash/%s"
	vmValuesEndpoint                 = "vm-values/query"
	genesisNodesConfigEndpoint       = "network/genesis-nodes"
	rawStartOfEpochMetaBlock         = "internal/raw/startofepoch/metablock/by-epoch/%d"
)

type baseEndpointProvider struct{}

// GetNetworkConfigEndpoint returns the network config endpoint
func (base *baseEndpointProvider) GetNetworkConfigEndpoint() string {
	return networkConfigEndpoint
}

// GetNetworkEconomicsEndpoint returns the network economics endpoint
func (base *baseEndpointProvider) GetNetworkEconomicsEndpoint() string {
	return networkEconomicsEndpoint
}

// GetRatingsConfigEndpoint returns the ratings config endpoint
func (base *baseEndpointProvider) GetRatingsConfigEndpoint() string {
	return ratingsConfigEndpoint
}

// GetEnableEpochsConfigEndpoint returns the enable epochs config endpoint
func (base *baseEndpointProvider) GetEnableEpochsConfigEndpoint() string {
	return enableEpochsConfigEndpoint
}

// GetAccountEndpoint returns the account endpoint
func (base *baseEndpointProvider) GetAccountEndpoint(addressAsBech32 string) string {
	return fmt.Sprintf(accountEndpoint, addressAsBech32)
}

// GetCostTransactionEndpoint returns the transaction cost endpoint
func (base *baseEndpointProvider) GetCostTransactionEndpoint() string {
	return costTransactionEndpoint
}

// GetSendTransactionEndpoint returns the send transaction endpoint
func (base *baseEndpointProvider) GetSendTransactionEndpoint() string {
	return sendTransactionEndpoint
}

// GetSendMultipleTransactionsEndpoint returns the send multiple transactions endpoint
func (base *baseEndpointProvider) GetSendMultipleTransactionsEndpoint() string {
	return sendMultipleTransactionsEndpoint
}

// GetTransactionStatusEndpoint returns the transaction status endpoint
func (base *baseEndpointProvider) GetTransactionStatusEndpoint(hexHash string) string {
	return fmt.Sprintf(transactionStatusEndpoint, hexHash)
}

// GetTransactionInfoEndpoint returns the transaction info endpoint
func (base *baseEndpointProvider) GetTransactionInfoEndpoint(hexHash string) string {
	return fmt.Sprintf(transactionInfoEndpoint, hexHash)
}

// GetHyperBlockByNonceEndpoint returns the hyper block by nonce endpoint
func (base *baseEndpointProvider) GetHyperBlockByNonceEndpoint(nonce uint64) string {
	return fmt.Sprintf(hyperBlockByNonceEndpoint, nonce)
}

// GetHyperBlockByHashEndpoint returns the hyper block by hash endpoint
func (base *baseEndpointProvider) GetHyperBlockByHashEndpoint(hexHash string) string {
	return fmt.Sprintf(hyperBlockByHashEndpoint, hexHash)
}

// GetVmValuesEndpoint returns the VM values endpoint
func (base *baseEndpointProvider) GetVmValuesEndpoint() string {
	return vmValuesEndpoint
}

// GetGenesisNodesConfigEndpoint returns the genesis nodes config endpoint
func (base *baseEndpointProvider) GetGenesisNodesConfigEndpoint() string {
	return genesisNodesConfigEndpoint
}

// GetRawStartOfEpochMetaBlock returns the raw start of epoch metablock endpoint
func (base *baseEndpointProvider) GetRawStartOfEpochMetaBlock(epoch uint32) string {
	return fmt.Sprintf(rawStartOfEpochMetaBlock, epoch)
}
