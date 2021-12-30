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

func Test_FunctionalTesting(t *testing.T) {
	t.Parallel()

	for _, f := range knownFetchers {
		fetcherName := f
		t.Run("Test_FunctionalTesting_"+fetcherName, func(t *testing.T) {
			t.Skip("this test should be run only when doing debugging work on the component")

			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &aggregator.HttpResponseGetter{})
			ethTicker := "ETH"
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			require.Nil(t, err)
			fmt.Printf("price between %s and %s is: %v\n", ethTicker, quoteUSDFiat, price)
			require.True(t, price > 0)
		})
	}
}

func Test_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	ethTicker := "ETH"
	pair := ethTicker + quoteUSDFiat

	for _, f := range knownFetchers {
		fetcherName := f

		t.Run("response getter errors should error"+fetcherName, func(t *testing.T) {
			t.Parallel()

			expectedError := errors.New("expected error")
			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "", pair, expectedError),
			})

			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, expectedError, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("empty string for price should error"+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "", pair, nil),
			})
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("negative price should error"+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "-1", pair, nil),
			})
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("invalid string for price should error"+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "not a number", pair, nil),
			})
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.NotNil(t, err)
			require.Equal(t, float64(0), price)
			require.IsType(t, err, &strconv.NumError{})
		})
		t.Run("should work eth-usd"+fetcherName, func(t *testing.T) {
			t.Parallel()

			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", pair, nil),
			})
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if err == errShouldSkipTest {
				return
			}
			require.Nil(t, err)
			require.Equal(t, 4714.05, price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
		t.Run("should work btc-usd"+fetcherName, func(t *testing.T) {
			t.Parallel()

			btcTicker := "BTC"
			btcUsdPair := btcTicker + quoteUSDFiat
			fetcher, _ := NewPriceFetcher(fetcherName, &mock.HttpResponseGetterStub{
				GetCalled: getFuncGetCalled(fetcherName, "4714.05000000", btcUsdPair, nil),
			})
			assert.False(t, check.IfNil(fetcher))

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
	case binanceName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*binancePriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case bitfinexName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*bitfinexPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case cryptocomName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*cryptocomPriceRequest)
			var err error
			cast.Result.Data.Price, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	case geminiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*geminiPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case hitbtcName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*hitbtcPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case huobiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*huobiPriceRequest)
			var err error
			cast.Ticker.Price, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	case krakenName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*krakenPriceRequest)
			cast.Result = map[string]krakenPricePair{
				pair: {[]string{returnPrice, ""}},
			}
			return returnErr
		}
	case okexName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*okexPriceRequest)
			cast.Data = []okexTicker{{returnPrice}}
			return returnErr
		}
	}

	return nil
}
