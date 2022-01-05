package fetchers

import "strings"

type baseFetcher struct{}

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
