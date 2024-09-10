package blockchain

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
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
	sdkCore "github.com/multiversx/mx-sdk-go/core"
	sdkHttp "github.com/multiversx/mx-sdk-go/core/http"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	withResultsQueryParam = "?withResults=true"
)

var (
	// MaximumBlocksDelta is the maximum allowed delta between the final block and the current block
	MaximumBlocksDelta uint64 = 500
)

// ArgsProxy is the DTO used in the multiversx proxy constructor
type ArgsProxy struct {
	ProxyURL               string
	Client                 sdkHttp.Client
	SameScState            bool
	ShouldBeSynced         bool
	FinalityCheck          bool
	AllowedDeltaToFinal    int
	CacheExpirationTime    time.Duration
	EntityType             sdkCore.RestAPIEntityType
	FilterQueryBlockCacher BlockDataCache
}

// proxy implements basic functions for interacting with a multiversx Proxy
type proxy struct {
	*baseProxy
	sameScState            bool
	shouldBeSynced         bool
	finalityCheck          bool
	allowedDeltaToFinal    int
	finalityProvider       FinalityProvider
	filterQueryBlockCacher BlockDataCache
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

	clientWrapper := sdkHttp.NewHttpClientWrapper(args.Client, args.ProxyURL)
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

	cacher := args.FilterQueryBlockCacher
	if cacher == nil {
		cacher = &DisabledBlockDataCache{}
	}

	ep := &proxy{
		baseProxy:              baseProxyInstance,
		sameScState:            args.SameScState,
		shouldBeSynced:         args.ShouldBeSynced,
		finalityCheck:          args.FinalityCheck,
		allowedDeltaToFinal:    args.AllowedDeltaToFinal,
		finalityProvider:       finalityProvider,
		filterQueryBlockCacher: cacher,
	}

	return ep, nil
}

func checkArgsProxy(args ArgsProxy) error {
	if args.FinalityCheck {
		if args.AllowedDeltaToFinal < sdkCore.MinAllowedDeltaToFinal {
			return fmt.Errorf("%w, provided: %d, minimum: %d",
				ErrInvalidAllowedDeltaToFinal, args.AllowedDeltaToFinal, sdkCore.MinAllowedDeltaToFinal)
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
	address sdkCore.AddressHandler,
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

	addressAsBech32String, err := address.AddressAsBech32String()
	if err != nil {
		return transaction.FrontendTransaction{}, "", err
	}

	return transaction.FrontendTransaction{
		Nonce:     account.Nonce,
		Value:     "",
		Receiver:  "",
		Sender:    addressAsBech32String,
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
func (ep *proxy) GetAccount(ctx context.Context, address sdkCore.AddressHandler) (*data.Account, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}

	addressAsBech32, err := address.AddressAsBech32String()
	if err != nil {
		return nil, err
	}

	err = ep.checkFinalState(ctx, addressAsBech32)
	if err != nil {
		return nil, err
	}

	endpoint := ep.endpointProvider.GetAccount(addressAsBech32)

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
	address sdkCore.AddressHandler,
	tokenIdentifier string,
	queryOptions api.AccountQueryOptions, // TODO: provide AccountQueryOptions on all accounts-related getters
) (*data.ESDTFungibleTokenData, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}

	addressAsBech32String, err := address.AddressAsBech32String()
	if err != nil {
		return nil, err
	}

	endpoint := ep.endpointProvider.GetESDTTokenData(addressAsBech32String, tokenIdentifier)
	endpoint = sdkCore.BuildUrlWithAccountQueryOptions(endpoint, queryOptions)
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
	address sdkCore.AddressHandler,
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

	addressAsBech32String, err := address.AddressAsBech32String()
	if err != nil {
		return nil, err
	}

	endpoint := ep.endpointProvider.GetNFTTokenData(addressAsBech32String, tokenIdentifier, nonce)
	endpoint = sdkCore.BuildUrlWithAccountQueryOptions(endpoint, queryOptions)
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

// GetGuardianData retrieves guardian data from proxy
func (ep *proxy) GetGuardianData(ctx context.Context, address sdkCore.AddressHandler) (*api.GuardianData, error) {
	if check.IfNil(address) {
		return nil, ErrNilAddress
	}
	if !address.IsValid() {
		return nil, ErrInvalidAddress
	}
	bech32Address, err := address.AddressAsBech32String()
	if err != nil {
		return nil, err
	}

	err = ep.checkFinalState(ctx, bech32Address)
	if err != nil {
		return nil, err
	}

	endpoint := ep.endpointProvider.GetGuardianData(bech32Address)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.GuardianDataResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.GuardianData, nil
}

// IsDataTrieMigrated returns true if the data trie of the given account is migrated
func (ep *proxy) IsDataTrieMigrated(ctx context.Context, address sdkCore.AddressHandler) (bool, error) {
	if check.IfNil(address) {
		return false, ErrNilAddress
	}

	bech32Address, err := address.AddressAsBech32String()
	if err != nil {
		return false, err
	}

	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.IsDataTrieMigrated(bech32Address))
	if err != nil || code != http.StatusOK {
		return false, createHTTPStatusError(code, err)
	}

	response := &data.IsDataTrieMigratedResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return false, err
	}
	if response.Error != "" {
		return false, errors.New(response.Error)
	}

	isMigrated, ok := response.Data["isMigrated"]
	if !ok {
		return false, errors.New("isMigrated key not found in response map")
	}

	return isMigrated, nil
}

// GetBlockBytesByNonce retrieves bytes of a block with a specific nonce
func (ep *proxy) GetBlockBytesByNonce(ctx context.Context, shardID uint32, nonce uint64) ([]byte, error) {
	endpoint := ep.endpointProvider.GetBlockByNonce(shardID, nonce)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}
	return buff, nil
}

// GetBlockBytesByHash retrieves bytes of a block with a specific hash
func (ep *proxy) GetBlockBytesByHash(ctx context.Context, shardID uint32, hash string) ([]byte, error) {
	endpoint := ep.endpointProvider.GetBlockByHash(shardID, hash)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}
	return buff, nil
}

// FilterLogs retrieves logs from the network and filters them based on the provided filter
func (ep *proxy) FilterLogs(ctx context.Context, filter *sdkCore.FilterQuery) ([]*transaction.Events, error) {
	status, err := ep.GetNetworkStatus(ctx, filter.ShardID)
	if err != nil {
		return nil, err
	}

	fromBlock, toBlock, err := ep.computeFromToBlocksForFilter(ctx, filter, status.Nonce)
	if err != nil {
		return nil, err
	}

	matchingEvents := make([]*transaction.Events, 0, toBlock-fromBlock+1)
	for blockNum := fromBlock; blockNum <= toBlock; blockNum++ {
		blockLogs, err := ep.getLogsFromBlock(ctx, filter.ShardID, blockNum, filter)
		if err != nil {
			return nil, err
		}

		matchingEvents = append(matchingEvents, blockLogs...)
	}

	return matchingEvents, nil
}

func (ep *proxy) computeFromToBlocksForFilter(ctx context.Context, filter *sdkCore.FilterQuery, latestBlock uint64) (uint64, uint64, error) {
	if filter.BlockHash != nil {
		blockNum, err := ep.getBlockNumberByHash(ctx, filter.ShardID, filter.BlockHash)
		if err != nil {
			return 0, 0, err
		}
		return blockNum, blockNum, nil
	}

	return resolveBlockRange(filter, latestBlock)
}

func resolveBlockRange(filter *sdkCore.FilterQuery, latestBlock uint64) (uint64, uint64, error) {
	var genesisBlock uint64 = 0

	if filter.ToBlock.HasValue && filter.ToBlock.Value > latestBlock {
		return 0, 0, errors.New("toBlock is greater than the latest block")
	}

	// Check if both fromBlock and toBlock are set
	if filter.FromBlock.HasValue && filter.ToBlock.HasValue {
		if filter.ToBlock.Value-filter.FromBlock.Value <= MaximumBlocksDelta {
			return filter.FromBlock.Value, filter.ToBlock.Value, nil
		}
		return 0, 0, errors.New("invalid block range: too many blocks to process")
	}

	// Check if only fromBlock is set
	if filter.FromBlock.HasValue {
		toBlock := latestBlock // Set toBlock to latestBlock
		if toBlock-filter.FromBlock.Value <= MaximumBlocksDelta {
			return filter.FromBlock.Value, toBlock, nil
		}
		return 0, 0, errors.New("invalid block range: too many blocks to process")
	}

	// Check if only toBlock is set
	if filter.ToBlock.HasValue {
		fromBlock := genesisBlock // Set fromBlock to genesisBlock
		if filter.ToBlock.Value-fromBlock <= MaximumBlocksDelta {
			return fromBlock, filter.ToBlock.Value, nil
		}
		return 0, 0, errors.New("invalid block range: too many blocks to process")
	}

	return 0, 0, errors.New("no block range specified")
}

// getBlockNumberByHash retrieves the block number associated with the given block hash
func (ep *proxy) getBlockNumberByHash(ctx context.Context, shardID uint32, blockHash []byte) (uint64, error) {
	blockHashStr := hex.EncodeToString(blockHash)
	buff, err := ep.GetBlockBytesByHash(ctx, shardID, blockHashStr)
	if err != nil {
		return 0, err
	}

	var response data.BlockResponse
	if err := json.Unmarshal(buff, &response); err != nil {
		return 0, err
	}
	blockNonce := response.Data.Block.Nonce

	// Cache the raw response bytes
	if len(buff) > 0 {
		cacheKey := make([]byte, 8)
		binary.BigEndian.PutUint64(cacheKey, blockNonce)
		ep.filterQueryBlockCacher.Put(cacheKey, buff, len(buff))
	}

	return blockNonce, nil
}

// getLogsFromBlock retrieves logs from a specific block and filters them
func (ep *proxy) getLogsFromBlock(ctx context.Context, shardID uint32, blockNum uint64, filter *sdkCore.FilterQuery) ([]*transaction.Events, error) {
	buff, err := getBlockBytesByNonce(ctx, ep, shardID, blockNum)
	if err != nil {
		return nil, err
	}

	var response data.BlockResponse
	if err := json.Unmarshal(buff, &response); err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return extractMatchingEvents(response, filter), nil
}

func getBlockBytesByNonce(ctx context.Context, ep *proxy, shardID uint32, nonce uint64) ([]byte, error) {
	cacheKey := make([]byte, 8)
	binary.BigEndian.PutUint64(cacheKey, nonce)

	cachedResponse, found := ep.filterQueryBlockCacher.Get(cacheKey)
	if found {
		cachedBuff, ok := cachedResponse.([]byte)
		if ok && len(cachedBuff) > 0 {
			return cachedBuff, nil
		}
	}

	buff, err := ep.GetBlockBytesByNonce(ctx, shardID, nonce)
	if err != nil {
		return nil, err
	}

	// Cache the raw response bytes
	ep.filterQueryBlockCacher.Put(cacheKey, buff, len(buff))

	return buff, nil
}

func extractMatchingEvents(response data.BlockResponse, filter *sdkCore.FilterQuery) []*transaction.Events {
	var matchingEvents []*transaction.Events
	for _, miniblock := range response.Data.Block.MiniBlocks {
		for _, tx := range miniblock.Transactions {
			if tx.Logs == nil {
				continue
			}
			for _, event := range tx.Logs.Events {
				if matchesFilter(filter, event) {
					matchingEvents = append(matchingEvents, event)
				}
			}
		}
	}
	return matchingEvents
}

func matchesFilter(filter *sdkCore.FilterQuery, event *transaction.Events) bool {
	// Check if the event's address matches any of the filter addresses (if set)
	if len(filter.Addresses) > 0 && !contains(filter.Addresses, event.Address) {
		return false
	}

	// Check if the event's topics match the filter topics
	if len(filter.Topics) > 0 && !topicsMatch(filter.Topics, event.Topics) {
		return false
	}

	return true
}

func contains(addresses []string, address string) bool {
	for _, a := range addresses {
		if a == address {
			return true
		}
	}
	return false
}

func topicsMatch(filterTopics [][]byte, eventTopics [][]byte) bool {
	if len(filterTopics) > len(eventTopics) {
		return false
	}

	for i, filterTopic := range filterTopics {
		if !bytes.Equal(filterTopic, eventTopics[i]) {
			return false
		}
	}

	return true
}

// IsInterfaceNil returns true if there is no value under the interface
func (ep *proxy) IsInterfaceNil() bool {
	return ep == nil
}
