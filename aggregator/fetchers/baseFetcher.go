package fetchers

import "strings"

type baseFetcher struct{}

func (b *baseFetcher) normalizeQuoteName(quote string, fetcherName string) string {
	if strings.Contains(quote, quoteUSDFiat) {
		switch fetcherName {
		case binanceName, cryptocomName, hitbtcName, huobiName, okexName:
			return quoteUSDT
		default:
			return quoteUSDFiat
		}
	}
	return quote
}
