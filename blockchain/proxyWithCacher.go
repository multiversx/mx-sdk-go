package blockchain

import (
	"context"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/blockchain")

const minimumCachingInterval = time.Second

type elrondProxyWithCache struct {
	Proxy
	mut                 sync.RWMutex
	fetchedConfigs      *data.NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
}

// NewElrondProxyWithCache will create an elrond proxy with cache instance
func NewElrondProxyWithCache(proxy Proxy, expirationTime time.Duration) (*elrondProxyWithCache, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if expirationTime < minimumCachingInterval {
		return nil, ErrInvalidCacherDuration
	}
	return &elrondProxyWithCache{
		Proxy:               proxy,
		cacheExpiryDuration: expirationTime,
		sinceTimeHandler:    since,
	}, nil
}

func since(t time.Time) time.Duration {
	return time.Since(t)
}

// GetNetworkConfig will return the cached network configs fetching new values and saving them if necessary
func (proxy *elrondProxyWithCache) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.RLock()
	cachedConfigs := proxy.getCachedConfigs()
	proxy.mut.RUnlock()

	if cachedConfigs != nil {
		return cachedConfigs, nil
	}

	return proxy.cacheConfigs(ctx)
}

func (proxy *elrondProxyWithCache) getCachedConfigs() *data.NetworkConfig {
	if proxy.sinceTimeHandler(proxy.lastFetchedTime) > proxy.cacheExpiryDuration {
		return nil
	}

	return proxy.fetchedConfigs
}

func (proxy *elrondProxyWithCache) cacheConfigs(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.Lock()
	defer proxy.mut.Unlock()

	// maybe another parallel running go routine already did the fetching
	cachedConfig := proxy.getCachedConfigs()
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	log.Debug("Network config not cached. caching...")
	configs, err := proxy.Proxy.GetNetworkConfig(ctx)
	if err != nil {
		return nil, err
	}

	proxy.lastFetchedTime = time.Now()
	proxy.fetchedConfigs = configs

	return configs, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *elrondProxyWithCache) IsInterfaceNil() bool {
	return proxy == nil
}
