package fetchers

import (
	"fmt"
	"strings"
	"sync"
)

type baseFetcher struct {
	knownPairs    map[string]struct{}
	knownPairsMut sync.RWMutex
}

func newBaseFetcher() baseFetcher {
	return baseFetcher{
		knownPairs:    make(map[string]struct{}),
		knownPairsMut: sync.RWMutex{},
	}
}

func (b *baseFetcher) normalizeQuoteName(quote string, fetcherName string) string {
	if strings.Contains(quote, quoteUSDFiat) {
		switch fetcherName {
		case BinanceName, CryptocomName, HitbtcName, HuobiName, OkexName:
			return quoteUSDT
		default:
			return quoteUSDFiat
		}
	}
	return quote
}

// AddPair adds the specified base-quote pair to the internal cache
func (b *baseFetcher) AddPair(base, quote string) {
	key := b.getPairKey(base, quote)

	b.knownPairsMut.Lock()
	b.knownPairs[key] = struct{}{}
	b.knownPairsMut.Unlock()
}

func (b *baseFetcher) hasPair(base, quote string) bool {
	key := b.getPairKey(base, quote)

	b.knownPairsMut.RLock()
	_, ok := b.knownPairs[key]
	b.knownPairsMut.RUnlock()

	return ok
}

func (b *baseFetcher) getPairKey(base, quote string) string {
	return fmt.Sprintf("%s-%s", base, quote)
}
