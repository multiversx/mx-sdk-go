package fetchers

import "strings"

type baseFetcher struct{}

func (b *baseFetcher) updateQuoteIfNeeded(quote *string, fetcherName string) {
	if strings.Contains(*quote, QuoteUSDFiat) {
		switch fetcherName {
		case binanceName, cryptocomName, hitbtcName, huobiName, okexName:
			*quote = QuoteUSDT
		default:
			*quote = QuoteUSDFiat
		}
	}
}
