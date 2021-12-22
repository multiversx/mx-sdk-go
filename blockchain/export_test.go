package blockchain

import (
	"sync"
	"time"
)

func NewElrondProxyWithCacheWithHandlers(
	proxy Proxy,
	cacheExpiryDuration time.Duration,
	sinceTimeHandler func(t time.Time) time.Duration,
) *elrondProxyWithCache {
	return &elrondProxyWithCache{
		Proxy:               proxy,
		mut:                 sync.RWMutex{},
		fetchedConfigs:      nil,
		lastFetchedTime:     time.Time{},
		cacheExpiryDuration: cacheExpiryDuration,
		sinceTimeHandler:    sinceTimeHandler,
	}
}
