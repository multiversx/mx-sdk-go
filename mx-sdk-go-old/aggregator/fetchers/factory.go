package fetchers

import (
	"fmt"

	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/aggregator"
)

// XExchangeTokensPair defines a base-quote pair of ids used by XExchange
type XExchangeTokensPair struct {
	Base  string
	Quote string
}

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(fetcherName string, responseGetter aggregator.ResponseGetter, graphqlGetter aggregator.GraphqlGetter, xExchangeTokensMap map[string]XExchangeTokensPair) (aggregator.PriceFetcher, error) {
	if responseGetter == nil {
		return nil, errNilResponseGetter
	}
	if graphqlGetter == nil {
		return nil, errNilGraphqlGetter
	}
	if xExchangeTokensMap == nil && fetcherName == XExchangeName {
		return nil, errNilXExchangeTokensMap
	}

	return createFetcher(fetcherName, responseGetter, graphqlGetter, xExchangeTokensMap)
}

func createFetcher(fetcherName string, responseGetter aggregator.ResponseGetter, graphqlGetter aggregator.GraphqlGetter, xExchangeTokensMap map[string]XExchangeTokensPair) (aggregator.PriceFetcher, error) {
	switch fetcherName {
	case BinanceName:
		return &binance{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case BitfinexName:
		return &bitfinex{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case CryptocomName:
		return &cryptocom{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case GeminiName:
		return &gemini{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case HitbtcName:
		return &hitbtc{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case HuobiName:
		return &huobi{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case KrakenName:
		return &kraken{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case OkexName:
		return &okex{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case XExchangeName:
		return &xExchange{
			GraphqlGetter:      graphqlGetter,
			baseFetcher:        newBaseFetcher(),
			xExchangeTokensMap: xExchangeTokensMap,
		}, nil
	}
	return nil, fmt.Errorf("%w, fetcherName %s", errInvalidFetcherName, fetcherName)
}
