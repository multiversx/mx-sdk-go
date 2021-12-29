package fetchers

import "github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) (aggregator.PriceFetcher, error) {
	if !isValidFetcherName(fetcherName) {
		return nil, errInvalidFetcherName
	}
	if responseGetter == nil {
		return nil, errNilResponseGetter
	}

	return createFetcher(fetcherName, responseGetter), nil
}

func isValidFetcherName(name string) bool {
	for _, fetcherName := range knownFetchers {
		if name == fetcherName {
			return true
		}
	}
	return false
}

func createFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) aggregator.PriceFetcher {
	switch fetcherName {
	case binanceName:
		return &binance{
			responseGetter,
			baseFetcher{},
		}
	case bitfinexName:
		return &bitfinex{
			responseGetter,
			baseFetcher{},
		}
	case cryptocomName:
		return &cryptocom{
			responseGetter,
			baseFetcher{},
		}
	case geminiName:
		return &gemini{
			responseGetter,
			baseFetcher{},
		}

	case hitbtcName:
		return &hitbtc{
			responseGetter,
			baseFetcher{},
		}
	case huobiName:
		return &huobi{
			responseGetter,
			baseFetcher{},
		}
	case krakenName:
		return &kraken{
			responseGetter,
			baseFetcher{},
		}
	case okexName:
		return &okex{
			responseGetter,
			baseFetcher{},
		}
	}
	return nil
}
