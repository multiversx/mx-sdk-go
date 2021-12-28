package blockchain

import (
	"time"
)

func NewElrondProxyWithCacheWithHandlers(
	proxy Proxy,
	cacheExpiryDuration time.Duration,
	sinceTimeHandler func(t time.Time) time.Duration,
) *elrondProxyWithCache {
	return &elrondProxyWithCache{
		Proxy:               proxy,
		cacheExpiryDuration: cacheExpiryDuration,
		sinceTimeHandler:    sinceTimeHandler,
	}
}
