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

func TestCryptocom_FunctionalTesting(t *testing.T) {
	t.Skip("this test should be run only when doing debugging work on the component")

	t.Parallel()

	ccom := &cryptocom{
		ResponseGetter: &HttpResponseGetter{},
	}
	ethTicker := "ETH"
	price, err := ccom.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
	require.Nil(t, err)
	fmt.Printf("price between %s and %s is: %v\n", ethTicker, QuoteUSDFiat, price)
	require.True(t, price > 0)
}

func TestCryptocom_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	t.Run("response getter errors should error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		ccom := &cryptocom{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					return expectedError
				},
			},
		}
		assert.False(t, check.IfNil(ccom))

		ethTicker := "ETH"
		price, err := ccom.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, expectedError, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("invalid value for price should error", func(t *testing.T) {
		t.Parallel()

		ccom := &cryptocom{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*cryptocomPriceRequest)
					cast.Result.Data.Price = -1
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(ccom))

		ethTicker := "ETH"
		price, err := ccom.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Equal(t, ErrInvalidResponseData, err)
		require.Equal(t, float64(0), price)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		expectedPrice := 4714.05000000
		ccom := &cryptocom{
			ResponseGetter: &mock.HttpResponseGetterStub{
				GetCalled: func(ctx context.Context, url string, response interface{}) error {
					cast, _ := response.(*cryptocomPriceRequest)
					cast.Result.Data.Price = expectedPrice
					return nil
				},
			},
		}
		assert.False(t, check.IfNil(ccom))

		ethTicker := "ETH"
		price, err := ccom.FetchPrice(context.Background(), ethTicker, QuoteUSDFiat)
		require.Nil(t, err)
		require.Equal(t, expectedPrice, price)
		assert.Equal(t, "Crypto.com", ccom.Name())
	})
}
