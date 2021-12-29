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

func TestKraken_FunctionalTesting(t *testing.T) {
	t.Skip("this test should be run only when doing debugging work on the component")

	t.Parallel()

	kra := &kraken{
		ResponseGetter: &HttpResponseGetter{},
	}
	ethTicker := "ETH"
	price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
	require.Nil(t, err)
	fmt.Printf("price between %s and %s is: %v\n", ethTicker, QuoteUSDFiat, price)
	require.True(t, price > 0)
}

func TestKraken_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	t.Run("response getter errors should error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		kra := &kraken{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					return expectedError
				},
			},
		}
		assert.False(t, check.IfNil(kra))

		ethTicker := "ETH"
		price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, expectedError, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("empty string for price should error", func(t *testing.T) {
		t.Parallel()

		ethTicker := "ETH"
		pair := ethTicker + QuoteUSDFiat
		kra := &kraken{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*krakenPriceRequest)
					cast.Result = map[string]krakenPricePair{
						pair: {[]string{"", ""}},
					}
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(kra))

		price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, ErrInvalidResponseData, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("empty string for key should error", func(t *testing.T) {
		t.Parallel()

		kra := &kraken{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*krakenPriceRequest)
					cast.Result = map[string]krakenPricePair{
						"": {[]string{"4714.05000000", ""}},
					}
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(kra))

		ethTicker := "ETH"
		price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, ErrInvalidResponseData, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("invalid string for price should error", func(t *testing.T) {
		t.Parallel()

		ethTicker := "ETH"
		pair := ethTicker + QuoteUSDFiat
		kra := &kraken{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*krakenPriceRequest)
					cast.Result = map[string]krakenPricePair{
						pair: {[]string{"not a number", ""}},
					}
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(kra))

		price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.NotNil(t, err)
		require.Equal(t, float64(0), price)
		require.IsType(t, err, &strconv.NumError{})
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ethTicker := "ETH"
		pair := ethTicker + QuoteUSDFiat
		kra := &kraken{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*krakenPriceRequest)
					cast.Result = map[string]krakenPricePair{
						pair: {[]string{"4714.05000000"}},
					}
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(kra))

		price, err := kra.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Nil(t, err)
		require.Equal(t, 4714.05, price)
		assert.Equal(t, "Kraken", kra.Name())
	})
}
