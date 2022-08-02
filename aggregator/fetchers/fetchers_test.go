package fetchers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errShouldSkipTest = errors.New("should skip test")

func createMockMap() map[string]MaiarTokensPair {
	return map[string]MaiarTokensPair{
		"ETH-USD": {
			Base:  "WEGLD-bd4d79", // for tests only until we have an ETH id
			Quote: "USDC-c76f1f",
		},
		"EGLD-USD": {
			Base:  "WEGLD-bd4d79",
			Quote: "USDC-c76f1f",
		},
		"BTC-USD": {
			Base:  "BTC-test1",
			Quote: "USD-test1",
		},
	}
}

func Test_FunctionalTesting(t *testing.T) {
	t.Parallel()

	for f := range ImplementedFetchers {
		fetcherName := f
		t.Run("Test_FunctionalTesting_"+fetcherName, func(t *testing.T) {
			t.Skip("this test should be run only when doing debugging work on the component")

			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &aggregator.HttpResponseGetter{}, createMockMap())
			ethTicker := "ETH"
			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			require.Nil(t, err)
			fmt.Printf("price between %s and %s is: %v from %s\n", ethTicker, quoteUSDFiat, price, fetcherName)
			require.True(t, price > 0)
		})
	}
}

func Test_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	ethTicker := "ETH"
	pair := ethTicker + quoteUSDFiat

	for f := range ImplementedFetchers {
		fetcherName := f

		t.Run("response getter errors should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			expectedError := errors.New("expected error")
			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "", pair, expectedError),
			}, createMockMap())

			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, expectedError, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("empty string for price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("negative price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "-1", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("invalid string for price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "not a number", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.NotNil(t, err)
			require.Equal(t, float64(0), price)
			require.IsType(t, err, &strconv.NumError{})
		})
		t.Run("maiar: missing key from map should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			if fetcherName != MaiarName {
				return
			}

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			missingTicker := "missing ticker"
			fetcher.AddPair(missingTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), missingTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			assert.Equal(t, errInvalidPair, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("pair not added should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, aggregator.ErrPairNotSupported, err)
			require.Equal(t, float64(0), price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
		t.Run("should work eth-usd "+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", pair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Nil(t, err)
			require.Equal(t, 4714.05, price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
		t.Run("should work btc-usd "+fetcherName, func(t *testing.T) {
			t.Parallel()

			btcTicker := "BTC"
			btcUsdPair := btcTicker + quoteUSDFiat
			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", btcUsdPair, nil),
			}, createMockMap())
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(btcTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), btcTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Nil(t, err)
			require.Equal(t, 4714.05, price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
	}
}

func getFuncGetCalled(name, returnPrice, pair string, returnErr error) func(ctx context.Context, url string, response interface{}) error {
	switch name {
	case BinanceName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*binancePriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case BitfinexName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*bitfinexPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case CryptocomName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*cryptocomPriceRequest)
			var err error
			cast.Result.Data.Price, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	case GeminiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*geminiPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case HitbtcName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*hitbtcPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case HuobiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*huobiPriceRequest)
			var err error
			cast.Ticker.Price, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	case KrakenName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*krakenPriceRequest)
			cast.Result = map[string]krakenPricePair{
				pair: {[]string{returnPrice, ""}},
			}
			return returnErr
		}
	case OkexName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*okexPriceRequest)
			cast.Data = []okexTicker{{returnPrice}}
			return returnErr
		}
	case MaiarName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*maiarPriceRequest)
			var err error
			cast.BasePrice, err = strconv.ParseFloat(returnPrice, 64)
			cast.QuotePrice = 1
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	}

	return nil
}
