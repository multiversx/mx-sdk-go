package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-go/state"
	"github.com/multiversx/mx-sdk-go/blockchain/factory"
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
	erdgoHttp "github.com/multiversx/mx-sdk-go/core/http"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	withResultsQueryParam = "?withResults=true"
)

// ArgsProxy is the DTO used in the multiversx proxy constructor
type ArgsProxy struct {
	ProxyURL            string
	Client              erdgoHttp.Client
	SameScState         bool
	ShouldBeSynced      bool
	FinalityCheck       bool
	AllowedDeltaToFinal int
	CacheExpirationTime time.Duration
	EntityType          erdgoCore.RestAPIEntityType
}

// proxy implements basic functions for interacting with a multiversx Proxy
type proxy struct {
	*baseProxy
	sameScState         bool
	shouldBeSynced      bool
	finalityCheck       bool
	allowedDeltaToFinal int
	finalityProvider    FinalityProvider
}

// NewProxy initializes and returns a proxy object
func NewProxy(args ArgsProxy) (*proxy, error) {
	err := checkArgsProxy(args)
	if err != nil {
		return nil, err
	}

	endpointProvider, err := factory.CreateEndpointProvider(args.EntityType)
	if err != nil {
		return nil, err
	}

	clientWrapper := erdgoHttp.NewHttpClientWrapper(args.Client, args.ProxyURL)
	baseArgs := argsBaseProxy{
		httpClientWrapper: clientWrapper,
		expirationTime:    args.CacheExpirationTime,
		endpointProvider:  endpointProvider,
	}
	baseProxyInstance, err := newBaseProxy(baseArgs)
	if err != nil {
		return nil, err
	}

	finalityProvider, err := factory.CreateFinalityProvider(baseProxyInstance, args.FinalityCheck)
	if err != nil {
		return nil, err
	}

	ep := &proxy{
		baseProxy:           baseProxyInstance,
		sameScState:         args.SameScState,
		shouldBeSynced:      args.ShouldBeSynced,
		finalityCheck:       args.FinalityCheck,
		allowedDeltaToFinal: args.AllowedDeltaToFinal,
		finalityProvider:    finalityProvider,
	}

	return ep, nil
}

func checkArgsProxy(args ArgsProxy) error {
	if args.FinalityCheck {
		if args.AllowedDeltaToFinal < erdgoCore.MinAllowedDeltaToFinal {
			return fmt.Errorf("%w, provided: %d, minimum: %d",
				ErrInvalidAllowedDeltaToFinal, args.AllowedDeltaToFinal, erdgoCore.MinAllowedDeltaToFinal)
		}
	}

	return nil
}

// ExecuteVMQuery retrieves data from existing SC trie through the use of a VM
func (ep *proxy) ExecuteVMQuery(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
	err := ep.checkFinalState(ctx, vmRequest.Address)
	if err != nil {
		return nil, err
	}

	jsonVMRequestWithOptionalParams := data.VmValueRequestWithOptionalParameters{
		VmValueRequest: vmRequest,
		SameScState:    ep.sameScState,
		ShouldBeSynced: ep.shouldBeSynced,
	}
	jsonVMRequest, err := json.Marshal(jsonVMRequestWithOptionalParams)
	if err != nil {
		return nil, err
	}

	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetVmValues(), jsonVMRequest)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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

func (ep *proxy) checkFinalState(ctx context.Context, address string) error {
	if !ep.finalityCheck {
		return nil
	}

	targetShardID, err := ep.GetShardOfAddress(ctx, address)
	if err != nil {
		return err
	}

	return ep.finalityProvider.CheckShardFinalization(ctx, targetShardID, uint64(ep.allowedDeltaToFinal))
}

// GetNetworkEconomics retrieves the network economics from the proxy
func (ep *proxy) GetNetworkEconomics(ctx context.Context) (*data.NetworkEconomics, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetNetworkEconomics())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetDefaultTransactionArguments(
	ctx context.Context,
	address erdgoCore.AddressHandler,
	networkConfigs *data.NetworkConfig,
) (transaction.FrontendTransaction, string, error) {
	if networkConfigs == nil {
		return transaction.FrontendTransaction{}, "", ErrNilNetworkConfigs
	}
	if check.IfNil(address) {
		return transaction.FrontendTransaction{}, "", ErrNilAddress
	}

	account, err := ep.GetAccount(ctx, address)
	if err != nil {
		return transaction.FrontendTransaction{}, "", err
	}

	return transaction.FrontendTransaction{
		Nonce:     account.Nonce,
		Value:     "",
		Receiver:  "",
		Sender:    address.AddressAsBech32String(),
		GasPrice:  networkConfigs.MinGasPrice,
		GasLimit:  networkConfigs.MinGasLimit,
		Data:      nil,
		Signature: "",
		ChainID:   networkConfigs.ChainID,
		Version:   networkConfigs.MinTransactionVersion,
		Options:   0,
	}, account.Balance, nil
}

// GetAccount retrieves an account info from the network (nonce, balance)
func (ep *proxy) GetAccount(ctx context.Context, address erdgoCore.AddressHandler) (*data.Account, error) {
	err := ep.checkFinalState(ctx, address.AddressAsBech32String())
	if err != nil {
		return nil, err
	}

	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}
	endpoint := ep.endpointProvider.GetAccount(address.AddressAsBech32String())

	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetSendTransaction(), jsonTx)
	if err != nil {
		return "", createHTTPStatusError(code, err)
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
func (ep *proxy) SendTransactions(ctx context.Context, txs []*transaction.FrontendTransaction) ([]string, error) {
	jsonTx, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetSendMultipleTransactions(), jsonTx)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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

func (ep *proxy) postProcessSendMultipleTxsResult(response *data.SendTransactionsResponse) ([]string, error) {
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
func (ep *proxy) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	endpoint := ep.endpointProvider.GetTransactionStatus(hash)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return "", createHTTPStatusError(code, err)
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
func (ep *proxy) GetTransactionInfo(ctx context.Context, hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, false)
}

// GetTransactionInfoWithResults retrieves a transaction's details from the network with events
func (ep *proxy) GetTransactionInfoWithResults(ctx context.Context, hash string) (*data.TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, true)
}

func (ep *proxy) getTransactionInfo(ctx context.Context, hash string, withResults bool) (*data.TransactionInfo, error) {
	endpoint := ep.endpointProvider.GetTransactionInfo(hash)
	if withResults {
		endpoint += withResultsQueryParam
	}

	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) RequestTransactionCost(ctx context.Context, tx *transaction.FrontendTransaction) (*data.TxCostResponseData, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetCostTransaction(), jsonTx)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetLatestHyperBlockNonce(ctx context.Context) (uint64, error) {
	response, err := ep.GetNetworkStatus(ctx, core.MetachainShardId)
	if err != nil {
		return 0, err
	}

	return response.Nonce, nil
}

// GetHyperBlockByNonce retrieves a hyper block's info by nonce from the network
func (ep *proxy) GetHyperBlockByNonce(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
	endpoint := ep.endpointProvider.GetHyperBlockByNonce(nonce)

	return ep.getHyperBlock(ctx, endpoint)
}

// GetHyperBlockByHash retrieves a hyper block's info by hash from the network
func (ep *proxy) GetHyperBlockByHash(ctx context.Context, hash string) (*data.HyperBlock, error) {
	endpoint := ep.endpointProvider.GetHyperBlockByHash(hash)

	return ep.getHyperBlock(ctx, endpoint)
}

func (ep *proxy) getHyperBlock(ctx context.Context, endpoint string) (*data.HyperBlock, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetRawBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawBlockByHash(shardId, hash)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawBlockByNonce retrieves a raw block by hash from the network
func (ep *proxy) GetRawBlockByNonce(ctx context.Context, shardId uint32, nonce uint64) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawBlockByNonce(shardId, nonce)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawStartOfEpochMetaBlock retrieves a raw block by hash from the network
func (ep *proxy) GetRawStartOfEpochMetaBlock(ctx context.Context, epoch uint32) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawStartOfEpochMetaBlock(epoch)

	return ep.getRawBlock(ctx, endpoint)
}

func (ep *proxy) getRawBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetRawMiniBlockByHash(ctx context.Context, shardId uint32, hash string, epoch uint32) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawMiniBlockByHash(shardId, hash, epoch)

	return ep.getRawMiniBlock(ctx, endpoint)
}

func (ep *proxy) getRawMiniBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetNonceAtEpochStart(ctx context.Context, shardId uint32) (uint64, error) {
	response, err := ep.GetNetworkStatus(ctx, shardId)
	if err != nil {
		return 0, err
	}

	return response.NonceAtEpochStart, nil
}

// GetRatingsConfig retrieves the ratings configuration from the proxy
func (ep *proxy) GetRatingsConfig(ctx context.Context) (*data.RatingsConfig, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetRatingsConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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
func (ep *proxy) GetEnableEpochsConfig(ctx context.Context) (*data.EnableEpochsConfig, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetEnableEpochsConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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

// GetGenesisNodesPubKeys retrieves genesis nodes configuration from proxy
func (ep *proxy) GetGenesisNodesPubKeys(ctx context.Context) (*data.GenesisNodes, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetGenesisNodesConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.GenesisNodesResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Nodes, nil
}

// GetValidatorsInfoByEpoch retrieves the validators info by epoch
func (ep *proxy) GetValidatorsInfoByEpoch(ctx context.Context, epoch uint32) ([]*state.ShardValidatorInfo, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetValidatorsInfo(epoch))
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.ValidatorsInfoResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.ValidatorsInfo, nil
}

// GetESDTTokenData returns the address' fungible token data
func (ep *proxy) GetESDTTokenData(
	ctx context.Context,
	address erdgoCore.AddressHandler,
	tokenIdentifier string,
	queryOptions api.AccountQueryOptions, // TODO: provide AccountQueryOptions on all accounts-related getters
) (*data.ESDTFungibleTokenData, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}

	endpoint := ep.endpointProvider.GetESDTTokenData(address.AddressAsBech32String(), tokenIdentifier)
	endpoint = erdgoCore.BuildUrlWithAccountQueryOptions(endpoint, queryOptions)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.ESDTFungibleResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.TokenData, nil
}

// GetNFTTokenData returns the address' NFT/SFT/MetaESDT token data
func (ep *proxy) GetNFTTokenData(
	ctx context.Context,
	address erdgoCore.AddressHandler,
	tokenIdentifier string,
	nonce uint64,
	queryOptions api.AccountQueryOptions, // TODO: provide AccountQueryOptions on all accounts-related getters
) (*data.ESDTNFTTokenData, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}

	endpoint := ep.endpointProvider.GetNFTTokenData(address.AddressAsBech32String(), tokenIdentifier, nonce)
	endpoint = erdgoCore.BuildUrlWithAccountQueryOptions(endpoint, queryOptions)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.ESDTNFTResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.TokenData, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ep *proxy) IsInterfaceNil() bool {
	return ep == nil
}
