package fetchers

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) (aggregator.PriceFetcher, error) {
	if responseGetter == nil {
		return nil, errNilResponseGetter
	}

	return createFetcher(fetcherName, responseGetter)
}

func createFetcher(fetcherName string, responseGetter aggregator.ResponseGetter) (aggregator.PriceFetcher, error) {
	switch fetcherName {
	case BinanceName:
		return &binance{
			responseGetter,
			baseFetcher{},
		}, nil
	case BitfinexName:
		return &bitfinex{
			responseGetter,
			baseFetcher{},
		}, nil
	case CryptocomName:
		return &cryptocom{
			responseGetter,
			baseFetcher{},
		}, nil
	case GeminiName:
		return &gemini{
			responseGetter,
			baseFetcher{},
		}, nil
	case HitbtcName:
		return &hitbtc{
			responseGetter,
			baseFetcher{},
		}, nil
	case HuobiName:
		return &huobi{
			responseGetter,
			baseFetcher{},
		}, nil
	case KrakenName:
		return &kraken{
			responseGetter,
			baseFetcher{},
		}, nil
	case OkexName:
		return &okex{
			responseGetter,
			baseFetcher{},
		}, nil
	}
	return nil, fmt.Errorf("%w, fetcherName %s", errInvalidFetcherName, fetcherName)
}
