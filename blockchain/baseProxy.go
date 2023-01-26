package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

var log = logger.GetOrCreate("mx-sdk-go/blockchain")

const (
	minimumCachingInterval = time.Second
)

type argsBaseProxy struct {
	expirationTime    time.Duration
	httpClientWrapper httpClientWrapper
	endpointProvider  EndpointProvider
}

type baseProxy struct {
	httpClientWrapper
	mut                 sync.RWMutex
	fetchedConfigs      *data.NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
	endpointProvider    EndpointProvider
}

// newBaseProxy will create a base multiversx proxy with cache instance
func newBaseProxy(args argsBaseProxy) (*baseProxy, error) {
	err := checkArgsBaseProxy(args)
	if err != nil {
		return nil, err
	}

	return &baseProxy{
		httpClientWrapper:   args.httpClientWrapper,
		cacheExpiryDuration: args.expirationTime,
		endpointProvider:    args.endpointProvider,
		sinceTimeHandler:    since,
	}, nil
}

func checkArgsBaseProxy(args argsBaseProxy) error {
	if args.expirationTime < minimumCachingInterval {
		return fmt.Errorf("%w, provided: %v, minimum: %v", ErrInvalidCacherDuration, args.expirationTime, minimumCachingInterval)
	}
	if check.IfNil(args.httpClientWrapper) {
		return ErrNilHTTPClientWrapper
	}
	if check.IfNil(args.endpointProvider) {
		return ErrNilEndpointProvider
	}

	return nil
}

func since(t time.Time) time.Duration {
	return time.Since(t)
}

// GetNetworkConfig will return the cached network configs fetching new values and saving them if necessary
func (proxy *baseProxy) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.RLock()
	cachedConfigs := proxy.getCachedConfigs()
	proxy.mut.RUnlock()

	if cachedConfigs != nil {
		return cachedConfigs, nil
	}

	return proxy.cacheConfigs(ctx)
}

func (proxy *baseProxy) getCachedConfigs() *data.NetworkConfig {
	if proxy.sinceTimeHandler(proxy.lastFetchedTime) > proxy.cacheExpiryDuration {
		return nil
	}

	return proxy.fetchedConfigs
}

func (proxy *baseProxy) cacheConfigs(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.Lock()
	defer proxy.mut.Unlock()

	// maybe another parallel running go routine already did the fetching
	cachedConfig := proxy.getCachedConfigs()
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	log.Debug("Network config not cached. caching...")
	configs, err := proxy.getNetworkConfigFromSource(ctx)
	if err != nil {
		return nil, err
	}

	proxy.lastFetchedTime = time.Now()
	proxy.fetchedConfigs = configs

	return configs, nil
}

// getNetworkConfigFromSource retrieves the network configuration from the proxy
func (proxy *baseProxy) getNetworkConfigFromSource(ctx context.Context) (*data.NetworkConfig, error) {
	buff, code, err := proxy.GetHTTP(ctx, proxy.endpointProvider.GetNetworkConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
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

// GetShardOfAddress returns the shard ID of a provided address by using a shardCoordinator object and querying the
// network config route
func (proxy *baseProxy) GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error) {
	addr, err := data.NewAddressFromBech32String(bech32Address)
	if err != nil {
		return 0, err
	}

	networkConfigs, err := proxy.GetNetworkConfig(ctx)
	if err != nil {
		return 0, err
	}

	shardCoordinatorInstance, err := NewShardCoordinator(networkConfigs.NumShardsWithoutMeta, 0)
	if err != nil {
		return 0, err
	}

	return shardCoordinatorInstance.ComputeShardId(addr)
}

// GetNetworkStatus will return the network status of a provided shard
func (proxy *baseProxy) GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
	endpoint := proxy.endpointProvider.GetNodeStatus(shardID)
	buff, code, err := proxy.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	endpointProviderType := proxy.endpointProvider.GetRestAPIEntityType()
	switch endpointProviderType {
	case core.Proxy:
		return proxy.getNetworkStatus(buff, shardID)
	case core.ObserverNode:
		return proxy.getNodeStatus(buff, shardID)
	}

	return &data.NetworkStatus{}, ErrInvalidEndpointProvider
}

func (proxy *baseProxy) getNetworkStatus(buff []byte, shardID uint32) (*data.NetworkStatus, error) {
	response := &data.NetworkStatusResponse{}
	err := json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	err = proxy.checkReceivedNodeStatus(response.Data.Status, shardID)
	if err != nil {
		return nil, err
	}

	return response.Data.Status, nil
}

func (proxy *baseProxy) getNodeStatus(buff []byte, shardID uint32) (*data.NetworkStatus, error) {
	response := &data.NodeStatusResponse{}
	err := json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	err = proxy.checkReceivedNodeStatus(response.Data.Status, shardID)
	if err != nil {
		return nil, err
	}

	return response.Data.Status, nil
}

func (proxy *baseProxy) checkReceivedNodeStatus(networkStatus *data.NetworkStatus, shardID uint32) error {
	if networkStatus == nil {
		return fmt.Errorf("%w, requested from %d", ErrNilNetworkStatus, shardID)
	}
	if !proxy.endpointProvider.ShouldCheckShardIDForNodeStatus() {
		return nil
	}
	if networkStatus.ShardID == shardID {
		return nil
	}

	return fmt.Errorf("%w, requested from %d, got response from %d", ErrShardIDMismatch, shardID, networkStatus.ShardID)
}

// GetRestAPIEntityType returns the REST API entity type that this implementation works with
func (proxy *baseProxy) GetRestAPIEntityType() core.RestAPIEntityType {
	return proxy.endpointProvider.GetRestAPIEntityType()
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *baseProxy) IsInterfaceNil() bool {
	return proxy == nil
}
