package blockchain

import (
	"context"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type elrondProxyWithCache struct {
	Proxy
	mut                 sync.RWMutex
	fetchedConfigs      *data.NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
}

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

	//TODO log debug here
	configs, err := proxy.Proxy.GetNetworkConfig(ctx)
	if err != nil {
		return nil, err
	}

	proxy.lastFetchedTime = time.Now()
	proxy.fetchedConfigs = configs

	return configs, nil
}

func since(t time.Time) time.Duration {
	return time.Since(t)
}

func (proxy *elrondProxyWithCache) IsInterfaceNil() bool {
	// TODO implement me
	panic("implement me")
}
