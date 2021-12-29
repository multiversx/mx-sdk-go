package aggregator

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinance_FunctionalTesting(t *testing.T) {
	t.Skip("this test should be run only when doing debugging work on the component")

	t.Parallel()

	bin := &binance{
		ResponseGetter: &HttpResponseGetter{},
	}
	ethTicker := "ETH"
	price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
	require.Nil(t, err)
	fmt.Printf("price between %s and %s is: %v\n", ethTicker, QuoteUSDFiat, price)
	require.True(t, price > 0)
}

func TestBinance_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	t.Run("response getter errors should error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		bin := &binance{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					return expectedError
				},
			},
		}
		assert.False(t, check.IfNil(bin))

		ethTicker := "ETH"
		price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, expectedError, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("empty string for price should error", func(t *testing.T) {
		t.Parallel()

		bin := &binance{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*binancePriceRequest)
					cast.Price = ""
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(bin))

		ethTicker := "ETH"
		price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, ErrInvalidResponseData, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("invalid string for price should error", func(t *testing.T) {
		t.Parallel()

		bin := &binance{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*binancePriceRequest)
					cast.Price = "not a number"
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(bin))

		ethTicker := "ETH"
		price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.NotNil(t, err)
		require.Equal(t, float64(0), price)
		require.IsType(t, err, &strconv.NumError{})
	})
	t.Run("negative price should error", func(t *testing.T) {
		t.Parallel()

		bin := &binance{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*binancePriceRequest)
					cast.Price = "-1"
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(bin))

		ethTicker := "ETH"
		price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, ErrInvalidResponseData, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		bin := &binance{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*binancePriceRequest)
					cast.Price = "4714.05000000"
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(bin))

		ethTicker := "ETH"
		price, err := bin.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Nil(t, err)
		require.Equal(t, 4714.05, price)
		assert.Equal(t, "Binance", bin.Name())
	})
}
