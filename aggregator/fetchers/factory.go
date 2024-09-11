package fetchers

import (
	"fmt"

	"github.com/multiversx/mx-sdk-go/aggregator"
)

// XExchangeTokensPair defines a base-quote pair of ids used by XExchange
type XExchangeTokensPair struct {
	Base  string
	Quote string
}

// NewPriceFetcher returns a new price fetcher of the type provided
func NewPriceFetcher(fetcherName string, responseGetter aggregator.ResponseGetter, graphqlGetter aggregator.GraphqlGetter, xExchangeTokensMap map[string]XExchangeTokensPair, config EVMGasPriceFetcherConfig) (aggregator.PriceFetcher, error) {
	if responseGetter == nil {
		return nil, errNilResponseGetter
	}
	if graphqlGetter == nil {
		return nil, errNilGraphqlGetter
	}
	if xExchangeTokensMap == nil && fetcherName == XExchangeName {
		return nil, errNilXExchangeTokensMap
	}

	return createFetcher(fetcherName, responseGetter, graphqlGetter, xExchangeTokensMap, config)
}

func createFetcher(fetcherName string, responseGetter aggregator.ResponseGetter, graphqlGetter aggregator.GraphqlGetter, xExchangeTokensMap map[string]XExchangeTokensPair, config EVMGasPriceFetcherConfig) (aggregator.PriceFetcher, error) {
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
	case OkxName:
		return &okx{
			ResponseGetter: responseGetter,
			baseFetcher:    newBaseFetcher(),
		}, nil
	case XExchangeName:
		return &xExchange{
			GraphqlGetter:      graphqlGetter,
			baseFetcher:        newBaseFetcher(),
			xExchangeTokensMap: xExchangeTokensMap,
		}, nil
	case EVMGasPriceStation:
		return &evmGasPriceFetcher{
			ResponseGetter: responseGetter,
			config:         config,
			baseFetcher:    newBaseFetcher(),
		}, nil
	}
	return nil, fmt.Errorf("%w, fetcherName %s", errInvalidFetcherName, fetcherName)
}
