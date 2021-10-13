package blockchain

import (
	"bytes"
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
	getProofEndpoint                 = "proof/address/%s"
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
func (ep *elrondProxy) ExecuteVMQuery(vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
	jsonVMRequest, err := json.Marshal(vmRequest)
	if err != nil {
		return nil, err
	}

	buff, err := ep.PostHTTP(vmValuesEndpoint, jsonVMRequest)
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
func (ep *elrondProxy) GetNetworkConfig() (*data.NetworkConfig, error) {
	buff, err := ep.GetHTTP(networkConfigEndpoint)
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
func (ep *elrondProxy) GetNetworkEconomics() (*data.NetworkEconomics, error) {
	buff, err := ep.GetHTTP(networkEconomicsEndpoint)
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
	address erdgoCore.AddressHandler,
	networkConfigs *data.NetworkConfig,
) (data.ArgCreateTransaction, error) {
	if networkConfigs == nil {
		return data.ArgCreateTransaction{}, ErrNilNetworkConfigs
	}
	if check.IfNil(address) {
		return data.ArgCreateTransaction{}, ErrNilAddress
	}

	account, err := ep.GetAccount(address)
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
func (ep *elrondProxy) GetAccount(address erdgoCore.AddressHandler) (*data.Account, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}
	endpoint := fmt.Sprintf(accountEndpoint, address.AddressAsBech32String())

	buff, err := ep.GetHTTP(endpoint)
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
func (ep *elrondProxy) SendTransaction(tx *data.Transaction) (string, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	buff, err := ep.PostHTTP(sendTransactionEndpoint, jsonTx)
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
func (ep *elrondProxy) SendTransactions(txs []*data.Transaction) ([]string, error) {
	jsonTx, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}
	buff, err := ep.PostHTTP(sendMultipleTransactionsEndpoint, jsonTx)
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
func (ep *elrondProxy) GetTransactionStatus(hash string) (string, error) {
	endpoint := fmt.Sprintf(getTransactionStatusEndpoint, hash)
	buff, err := ep.GetHTTP(endpoint)
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
func (ep *elrondProxy) GetTransactionInfo(hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(hash, false)
}

// GetTransactionInfoWithResults retrieves a transaction's details from the network with events
func (ep *elrondProxy) GetTransactionInfoWithResults(hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(hash, true)
}

func (ep *elrondProxy) getTransactionInfo(hash string, withResults bool) (*data.TransactionInfo, error) {
	endpoint := fmt.Sprintf(getTransactionInfoEndpoint, hash)

	if withResults {
		endpoint += withResultsQueryParam
	}

	buff, err := ep.GetHTTP(endpoint)
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
func (ep *elrondProxy) RequestTransactionCost(tx *data.Transaction) (*data.TxCostResponseData, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	buff, err := ep.PostHTTP(costTransactionEndpoint, jsonTx)
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
func (ep *elrondProxy) GetLatestHyperBlockNonce() (uint64, error) {
	endpoint := fmt.Sprintf(getNetworkStatusEndpoint, core.MetachainShardId)
	buff, err := ep.GetHTTP(endpoint)
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
func (ep *elrondProxy) GetHyperBlockByNonce(nonce uint64) (*data.HyperBlock, error) {
	endpoint := fmt.Sprintf(getHyperBlockByNonceEndpoint, nonce)

	return ep.getHyperBlock(endpoint)
}

// GetHyperBlockByHash retrieves a hyper block's info by hash from the network
func (ep *elrondProxy) GetHyperBlockByHash(hash string) (*data.HyperBlock, error) {
	endpoint := fmt.Sprintf(getHyperBlockByHashEndpoint, hash)

	return ep.getHyperBlock(endpoint)
}

func (ep *elrondProxy) getHyperBlock(endpoint string) (*data.HyperBlock, error) {
	buff, err := ep.GetHTTP(endpoint)
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

func (ep *elrondProxy) GetMerkleProof(address string) (*data.ProofResponse, error) {
	endpoint := fmt.Sprintf(getProofEndpoint, address)
	buff, err := ep.GetHTTP(endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.ProofResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response, nil
}

func (ep *elrondProxy) GetHTTP(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ep.proxyURL, endpoint)
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

func (ep *elrondProxy) PostHTTP(endpoint string, data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ep.proxyURL, endpoint)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
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

// IsInterfaceNil returns true if there is no value under the interface
func (ep *elrondProxy) IsInterfaceNil() bool {
	return ep == nil
}
