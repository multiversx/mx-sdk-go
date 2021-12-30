package fetchers

import "github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) (aggregator.PriceFetcher, error) {
	if responseGetter == nil {
		return nil, errNilResponseGetter
	}

	return createFetcher(fetcherName, responseGetter)
}

func createFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) (aggregator.PriceFetcher, error) {
	switch fetcherName {
	case binanceName:
		return &binance{
			responseGetter,
			baseFetcher{},
		}, nil
	case bitfinexName:
		return &bitfinex{
			responseGetter,
			baseFetcher{},
		}, nil
	case cryptocomName:
		return &cryptocom{
			responseGetter,
			baseFetcher{},
		}, nil
	case geminiName:
		return &gemini{
			responseGetter,
			baseFetcher{},
		}, nil
	case hitbtcName:
		return &hitbtc{
			responseGetter,
			baseFetcher{},
		}, nil
	case huobiName:
		return &huobi{
			responseGetter,
			baseFetcher{},
		}, nil
	case krakenName:
		return &kraken{
			responseGetter,
			baseFetcher{},
		}, nil
	case okexName:
		return &okex{
			responseGetter,
			baseFetcher{},
		}, nil
	}
	return nil, errInvalidFetcherName
}
