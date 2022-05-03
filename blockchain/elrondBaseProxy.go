package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/blockchain")

const (
	minimumCachingInterval = time.Second
)

type argsElrondBaseProxy struct {
	expirationTime    time.Duration
	httpClientWrapper httpClientWrapper
	endpointProvider  EndpointProvider
}

type elrondBaseProxy struct {
	httpClientWrapper
	mut                 sync.RWMutex
	fetchedConfigs      *data.NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
	endpointProvider    EndpointProvider
}

// newElrondBaseProxy will create a base elrond proxy with cache instance
func newElrondBaseProxy(args argsElrondBaseProxy) (*elrondBaseProxy, error) {
	err := checkArgsBaseProxy(args)
	if err != nil {
		return nil, err
	}

	return &elrondBaseProxy{
		httpClientWrapper:   args.httpClientWrapper,
		cacheExpiryDuration: args.expirationTime,
		endpointProvider:    args.endpointProvider,
		sinceTimeHandler:    since,
	}, nil
}

func checkArgsBaseProxy(args argsElrondBaseProxy) error {
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
func (proxy *elrondBaseProxy) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.RLock()
	cachedConfigs := proxy.getCachedConfigs()
	proxy.mut.RUnlock()

	if cachedConfigs != nil {
		return cachedConfigs, nil
	}

	return proxy.cacheConfigs(ctx)
}

func (proxy *elrondBaseProxy) getCachedConfigs() *data.NetworkConfig {
	if proxy.sinceTimeHandler(proxy.lastFetchedTime) > proxy.cacheExpiryDuration {
		return nil
	}

	return proxy.fetchedConfigs
}

func (proxy *elrondBaseProxy) cacheConfigs(ctx context.Context) (*data.NetworkConfig, error) {
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
func (proxy *elrondBaseProxy) getNetworkConfigFromSource(ctx context.Context) (*data.NetworkConfig, error) {
	buff, code, err := proxy.GetHTTP(ctx, proxy.endpointProvider.GetNetworkConfigEndpoint())
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
func (proxy *elrondBaseProxy) GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error) {
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
func (proxy *elrondBaseProxy) GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
	endpoint := proxy.endpointProvider.GetNodeStatusEndpoint(shardID)
	buff, code, err := proxy.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &data.NetworkStatusResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Status, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *elrondBaseProxy) IsInterfaceNil() bool {
	return proxy == nil
}
