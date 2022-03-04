package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

const (
	// endpoints
	networkConfigEndpoint            = "network/config"
	networkEconomicsEndpoint         = "network/economics"
	ratingsConfigEndpoint            = "network/ratings"
	enableEpochsConfigEndpoint       = "network/enable-epochs"
	accountEndpoint                  = "address/%s"
	costTransactionEndpoint          = "transaction/cost"
	sendTransactionEndpoint          = "transaction/send"
	sendMultipleTransactionsEndpoint = "transaction/send-multiple"
	getTransactionStatusEndpoint     = "transaction/%s/status"
	getTransactionInfoEndpoint       = "transaction/%s"
	getHyperBlockByNonceEndpoint     = "hyperblock/by-nonce/%v"
	getHyperBlockByHashEndpoint      = "hyperblock/by-hash/%s"
	getNetworkStatusEndpoint         = "network/status/%v"
	withResultsQueryParam            = "?withResults=true"
	vmValuesEndpoint                 = "vm-values/query"
	genesisNodesConfigEndpoint       = "node/genesisnodesconfig"

	getRawBlockByHashEndpoint     = "internal/%d/raw/block/by-hash/%s"
	getRawBlockByNonceEndpoint    = "internal/%d/raw/block/by-nonce/%d"
	getRawMiniBlockByHashEndpoint = "internal/%d/raw/miniblock/by-hash/%s/epoch/%d"
	getRawStartOfEpochMetaBlock   = "internal/raw/startofepoch/metablock/by-epoch/%d"
)

// HTTPClient is the interface we expect to call in order to do the HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// elrondProxy implements basic functions for interacting with an Elrond Proxy
type elrondProxy struct {
	proxyURL string
	client   HTTPClient
}

// NewElrondProxy initializes and returns an ElrondProxy object
func NewElrondProxy(url string, client HTTPClient) *elrondProxy {
	if check.IfNilReflect(client) {
		client = http.DefaultClient
	}

	ep := &elrondProxy{
		proxyURL: url,
		client:   client,
	}

	return ep
}

// ExecuteVMQuery retrieves data from existing SC trie through the use of a VM
func (ep *elrondProxy) ExecuteVMQuery(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
	jsonVMRequest, err := json.Marshal(vmRequest)
	if err != nil {
		return nil, err
	}

	buff, err := ep.PostHTTP(ctx, vmValuesEndpoint, jsonVMRequest)
	if err != nil {
		return nil, err
	}

	response := &data.ResponseVmValue{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response.Data, nil
}

// GetNetworkConfig retrieves the network configuration from the proxy
func (ep *elrondProxy) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	buff, err := ep.GetHTTP(ctx, networkConfigEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.NetworkConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetNetworkEconomics retrieves the network economics from the proxy
func (ep *elrondProxy) GetNetworkEconomics(ctx context.Context) (*data.NetworkEconomics, error) {
	buff, err := ep.GetHTTP(ctx, networkEconomicsEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.NetworkEconomicsResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Economics, nil
}

// GetDefaultTransactionArguments will prepare the transaction creation argument by querying the account's info
func (ep *elrondProxy) GetDefaultTransactionArguments(
	ctx context.Context,
	address erdgoCore.AddressHandler,
	networkConfigs *data.NetworkConfig,
) (data.ArgCreateTransaction, error) {
	if networkConfigs == nil {
		return data.ArgCreateTransaction{}, ErrNilNetworkConfigs
	}
	if check.IfNil(address) {
		return data.ArgCreateTransaction{}, ErrNilAddress
	}

	account, err := ep.GetAccount(ctx, address)
	if err != nil {
		return data.ArgCreateTransaction{}, err
	}

	return data.ArgCreateTransaction{
		Nonce:            account.Nonce,
		Value:            "",
		RcvAddr:          "",
		SndAddr:          address.AddressAsBech32String(),
		GasPrice:         networkConfigs.MinGasPrice,
		GasLimit:         networkConfigs.MinGasLimit,
		Data:             nil,
		Signature:        "",
		ChainID:          networkConfigs.ChainID,
		Version:          networkConfigs.MinTransactionVersion,
		Options:          0,
		AvailableBalance: account.Balance,
	}, nil
}

// GetAccount retrieves an account info from the network (nonce, balance)
func (ep *elrondProxy) GetAccount(ctx context.Context, address erdgoCore.AddressHandler) (*data.Account, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}
	endpoint := fmt.Sprintf(accountEndpoint, address.AddressAsBech32String())

	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	response := &data.AccountResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Account, nil
}

// SendTransaction broadcasts a transaction to the network and returns the txhash if successful
func (ep *elrondProxy) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	buff, err := ep.PostHTTP(ctx, sendTransactionEndpoint, jsonTx)
	if err != nil {
		return "", err
	}

	response := &data.SendTransactionResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.Data.TxHash, nil
}

// SendTransactions broadcasts the provided transactions to the network and returns the txhashes if successful
func (ep *elrondProxy) SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error) {
	jsonTx, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}
	buff, err := ep.PostHTTP(ctx, sendMultipleTransactionsEndpoint, jsonTx)
	if err != nil {
		return nil, err
	}

	response := &data.SendTransactionsResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return ep.postProcessSendMultipleTxsResult(response)
}

func (ep *elrondProxy) postProcessSendMultipleTxsResult(response *data.SendTransactionsResponse) ([]string, error) {
	txHashes := make([]string, 0, len(response.Data.TxsHashes))
	indexes := make([]int, 0, len(response.Data.TxsHashes))
	for index := range response.Data.TxsHashes {
		indexes = append(indexes, index)
	}

	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	for _, idx := range indexes {
		txHashes = append(txHashes, response.Data.TxsHashes[idx])
	}

	return txHashes, nil
}

// GetTransactionStatus retrieves a transaction's status from the network
func (ep *elrondProxy) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	endpoint := fmt.Sprintf(getTransactionStatusEndpoint, hash)
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return "", err
	}

	response := &data.TransactionStatus{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.Data.Status, nil
}

// GetTransactionInfo retrieves a transaction's details from the network
func (ep *elrondProxy) GetTransactionInfo(ctx context.Context, hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, false)
}

// GetTransactionInfoWithResults retrieves a transaction's details from the network with events
func (ep *elrondProxy) GetTransactionInfoWithResults(ctx context.Context, hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, true)
}

func (ep *elrondProxy) getTransactionInfo(ctx context.Context, hash string, withResults bool) (*data.TransactionInfo, error) {
	endpoint := fmt.Sprintf(getTransactionInfoEndpoint, hash)

	if withResults {
		endpoint += withResultsQueryParam
	}

	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.TransactionInfo{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response, nil
}

// RequestTransactionCost retrieves how many gas a transaction will consume
func (ep *elrondProxy) RequestTransactionCost(ctx context.Context, tx *data.Transaction) (*data.TxCostResponseData, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	buff, err := ep.PostHTTP(ctx, costTransactionEndpoint, jsonTx)
	if err != nil {
		return nil, err
	}

	response := &data.ResponseTxCost{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response.Data, nil
}

// GetLatestHyperBlockNonce retrieves the latest hyper block (metachain) nonce from the network
func (ep *elrondProxy) GetLatestHyperBlockNonce(ctx context.Context) (uint64, error) {
	endpoint := fmt.Sprintf(getNetworkStatusEndpoint, core.MetachainShardId)
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return 0, err
	}

	response := &data.NetworkStatusResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}

	return response.Data.Status.Nonce, nil
}

// GetHyperBlockByNonce retrieves a hyper block's info by nonce from the network
func (ep *elrondProxy) GetHyperBlockByNonce(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
	endpoint := fmt.Sprintf(getHyperBlockByNonceEndpoint, nonce)

	return ep.getHyperBlock(ctx, endpoint)
}

// GetHyperBlockByHash retrieves a hyper block's info by hash from the network
func (ep *elrondProxy) GetHyperBlockByHash(ctx context.Context, hash string) (*data.HyperBlock, error) {
	endpoint := fmt.Sprintf(getHyperBlockByHashEndpoint, hash)

	return ep.getHyperBlock(ctx, endpoint)
}

func (ep *elrondProxy) getHyperBlock(ctx context.Context, endpoint string) (*data.HyperBlock, error) {
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.HyperBlockResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response.Data.HyperBlock, nil
}

// GetRawBlockByHash retrieves a raw block by hash from the network
func (ep *elrondProxy) GetRawBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error) {
	endpoint := fmt.Sprintf(getRawBlockByHashEndpoint, shardId, hash)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawBlockByNonce retrieves a raw block by hash from the network
func (ep *elrondProxy) GetRawBlockByNonce(ctx context.Context, shardId uint32, nonce uint64) ([]byte, error) {
	endpoint := fmt.Sprintf(getRawBlockByNonceEndpoint, shardId, nonce)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawStartOfEpochMetaBlock retrieves a raw block by hash from the network
func (ep *elrondProxy) GetRawStartOfEpochMetaBlock(ctx context.Context, epoch uint32) ([]byte, error) {
	endpoint := fmt.Sprintf(getRawStartOfEpochMetaBlock, epoch)

	return ep.getRawBlock(ctx, endpoint)
}

func (ep *elrondProxy) getRawBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.RawBlockRespone{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Block, nil
}

// GetRawMiniBlockByHash retrieves a raw block by hash from the network
func (ep *elrondProxy) GetRawMiniBlockByHash(ctx context.Context, shardId uint32, hash string, epoch uint32) ([]byte, error) {
	endpoint := fmt.Sprintf(getRawMiniBlockByHashEndpoint, shardId, hash, epoch)

	return ep.getRawMiniBlock(ctx, endpoint)
}

func (ep *elrondProxy) getRawMiniBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.RawMiniBlockRespone{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.MiniBlock, nil
}

// GetNonceAtEpochStart retrieves the start of epoch nonce from hyper block (metachain)
func (ep *elrondProxy) GetNonceAtEpochStart(ctx context.Context, shardId uint32) (uint64, error) {
	endpoint := fmt.Sprintf(getNetworkStatusEndpoint, shardId)
	buff, err := ep.GetHTTP(ctx, endpoint)
	if err != nil {
		return 0, err
	}

	response := &data.NetworkStatusResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return 0, err
	}
	if response.Error != "" {
		return 0, errors.New(response.Error)
	}

	return response.Data.Status.NonceAtEpochStart, nil
}

// GetHTTP does a GET method operation on the specified endpoint
func (ep *elrondProxy) GetHTTP(ctx context.Context, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ep.proxyURL, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := ep.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// PostHTTP does a POST method operation on the specified endpoint with the provided raw data bytes
func (ep *elrondProxy) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ep.proxyURL, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "")
	response, err := ep.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	return ioutil.ReadAll(response.Body)
}

// GetRatingsConfig retrieves the ratings configuration from the proxy
func (ep *elrondProxy) GetRatingsConfig(ctx context.Context) (*data.RatingsConfig, error) {
	buff, err := ep.GetHTTP(ctx, ratingsConfigEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.RatingsConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetEnableEpochsConfig retrieves the ratings configuration from the proxy
func (ep *elrondProxy) GetEnableEpochsConfig(ctx context.Context) (*data.EnableEpochsConfig, error) {
	buff, err := ep.GetHTTP(ctx, enableEpochsConfigEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.EnableEpochsConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetGenesisNodesConfig
func (ep *elrondProxy) GetGenesisNodesConfig(ctx context.Context) (*data.GenesisNodesConfig, error) {
	buff, err := ep.GetHTTP(ctx, genesisNodesConfigEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.GenesisNodesConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ep *elrondProxy) IsInterfaceNil() bool {
	return ep == nil
}
