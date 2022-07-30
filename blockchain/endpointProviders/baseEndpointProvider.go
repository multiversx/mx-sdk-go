package endpointProviders

import "fmt"

const (
	networkConfig            = "network/config"
	networkEconomics         = "network/economics"
	ratingsConfig            = "network/ratings"
	enableEpochsConfig       = "network/enable-epochs"
	account                  = "address/%s"
	accountKeys              = "address/%s/keys"
	costTransaction          = "transaction/cost"
	sendTransaction          = "transaction/send"
	sendMultipleTransactions = "transaction/send-multiple"
	transactionStatus        = "transaction/%s/status"
	transactionInfo          = "transaction/%s"
	hyperBlockByNonce        = "hyperblock/by-nonce/%d"
	hyperBlockByHash         = "hyperblock/by-hash/%s"
	vmValues                 = "vm-values/query"
	genesisNodesConfig       = "network/genesis-nodes"
	rawStartOfEpochMetaBlock = "internal/raw/startofepoch/metablock/by-epoch/%d"
)

type baseEndpointProvider struct{}

// GetNetworkConfig returns the network config endpoint
func (base *baseEndpointProvider) GetNetworkConfig() string {
	return networkConfig
}

// GetNetworkEconomics returns the network economics endpoint
func (base *baseEndpointProvider) GetNetworkEconomics() string {
	return networkEconomics
}

// GetRatingsConfig returns the ratings config endpoint
func (base *baseEndpointProvider) GetRatingsConfig() string {
	return ratingsConfig
}

// GetEnableEpochsConfig returns the enable epochs config endpoint
func (base *baseEndpointProvider) GetEnableEpochsConfig() string {
	return enableEpochsConfig
}

// GetAccount returns the account endpoint
func (base *baseEndpointProvider) GetAccount(addressAsBech32 string) string {
	return fmt.Sprintf(account, addressAsBech32)
}

// GetAccountKeys retrieves all key-value pairs stored under a given account
func (base *baseEndpointProvider) GetAccountKeys(addressAsBech32 string) string {
	return fmt.Sprintf(accountKeys, addressAsBech32)
}

// GetCostTransaction returns the transaction cost endpoint
func (base *baseEndpointProvider) GetCostTransaction() string {
	return costTransaction
}

// GetSendTransaction returns the send transaction endpoint
func (base *baseEndpointProvider) GetSendTransaction() string {
	return sendTransaction
}

// GetSendMultipleTransactions returns the send multiple transactions endpoint
func (base *baseEndpointProvider) GetSendMultipleTransactions() string {
	return sendMultipleTransactions
}

// GetTransactionStatus returns the transaction status endpoint
func (base *baseEndpointProvider) GetTransactionStatus(hexHash string) string {
	return fmt.Sprintf(transactionStatus, hexHash)
}

// GetTransactionInfo returns the transaction info endpoint
func (base *baseEndpointProvider) GetTransactionInfo(hexHash string) string {
	return fmt.Sprintf(transactionInfo, hexHash)
}

// GetHyperBlockByNonce returns the hyper block by nonce endpoint
func (base *baseEndpointProvider) GetHyperBlockByNonce(nonce uint64) string {
	return fmt.Sprintf(hyperBlockByNonce, nonce)
}

// GetHyperBlockByHash returns the hyper block by hash endpoint
func (base *baseEndpointProvider) GetHyperBlockByHash(hexHash string) string {
	return fmt.Sprintf(hyperBlockByHash, hexHash)
}

// GetVmValues returns the VM values endpoint
func (base *baseEndpointProvider) GetVmValues() string {
	return vmValues
}

// GetGenesisNodesConfig returns the genesis nodes config endpoint
func (base *baseEndpointProvider) GetGenesisNodesConfig() string {
	return genesisNodesConfig
}

// GetRawStartOfEpochMetaBlock returns the raw start of epoch metablock endpoint
func (base *baseEndpointProvider) GetRawStartOfEpochMetaBlock(epoch uint32) string {
	return fmt.Sprintf(rawStartOfEpochMetaBlock, epoch)
}
