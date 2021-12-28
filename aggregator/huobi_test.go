package aggregator

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHuobi_FunctionalTesting(t *testing.T) {
	t.Skip("this test should be run only when doing debugging work on the component")

	t.Parallel()

	huo := &huobi{
		ResponseGetter: &HttpResponseGetter{},
	}
	ethTicker := "ETH"
	price, err := huo.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
	require.Nil(t, err)
	fmt.Printf("price between %s and %s is: %v\n", ethTicker, QuoteUSDFiat, price)
	require.True(t, price > 0)
}

func TestHuobi_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	t.Run("response getter errors should error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		huo := &huobi{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					return expectedError
				},
			},
		}
		assert.False(t, check.IfNil(huo))

		ethTicker := "ETH"
		price, err := huo.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, expectedError, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("invalid value for price should error", func(t *testing.T) {
		t.Parallel()

		huo := &huobi{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*huobiPriceRequest)
					cast.Ticker.Price = -1
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(huo))

		ethTicker := "ETH"
		price, err := huo.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, InvalidResponseDataErr, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		expectedPrice := 4714.05000000
		huo := &huobi{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*huobiPriceRequest)
					cast.Ticker.Price = expectedPrice
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(huo))

		ethTicker := "ETH"
		price, err := huo.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Nil(t, err)
		require.Equal(t, expectedPrice, price)
		assert.Equal(t, "Huobi", huo.Name())
	})
}
