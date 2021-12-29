package fetchers

import "strings"

type baseFetcher struct{}

func (b *baseFetcher) updateQuoteIfNeeded(quote *string, fetcherName string) {
	if strings.Contains(*quote, quoteUSDFiat) {
		switch fetcherName {
		case binanceName, cryptocomName, hitbtcName, huobiName, okexName:
			*quote = quoteUSDT
		default:
			*quote = quoteUSDFiat
		}
	}
}
