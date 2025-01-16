package fetchers

import (
	"context"
	"errors"
	"testing"

	"github.com/multiversx/mx-sdk-go/aggregator"
	"github.com/multiversx/mx-sdk-go/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func createMockEVMGasPriceFetcher() *evmGasPriceFetcher {
	return &evmGasPriceFetcher{
		ResponseGetter: &mock.HttpResponseGetterStub{},
		config: EVMGasPriceFetcherConfig{
			ApiURL:   "api-url",
			Selector: "SafeGasPrice",
		},
		baseFetcher: newBaseFetcher(),
	}
}

func TestEvmGasPriceFetcher_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *evmGasPriceFetcher
	assert.True(t, instance.IsInterfaceNil())

	instance = &evmGasPriceFetcher{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestEvmGasPriceFetcher_Name(t *testing.T) {
	t.Parallel()

	fetcher := &evmGasPriceFetcher{
		config: EVMGasPriceFetcherConfig{
			Selector: "test selector",
		},
	}

	assert.Equal(t, "EVM gas price station when using selector test selector", fetcher.Name())
}

func TestKraken_FetchPrice(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("un-added pair should error", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		value, err := fetcher.FetchPrice(context.Background(), "missing", "pair")
		assert.Zero(t, value)
		assert.Equal(t, aggregator.ErrPairNotSupported, err)
	})
	t.Run("HTTP getter fails, should error", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		fetcher.AddPair("test", "pair")
		fetcher.ResponseGetter = &mock.HttpResponseGetterStub{
			GetCalled: func(ctx context.Context, url string, response interface{}) error {
				assert.Equal(t, "api-url", url)
				return expectedErr
			},
		}
		value, err := fetcher.FetchPrice(context.Background(), "test", "pair")
		assert.Zero(t, value)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("invalid selector, should error", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		fetcher.config.Selector = "invalid-selector"
		fetcher.AddPair("test", "pair")
		fetcher.ResponseGetter = &mock.HttpResponseGetterStub{
			GetCalled: func(ctx context.Context, url string, response interface{}) error {
				assert.Equal(t, "api-url", url)
				return nil
			},
		}
		value, err := fetcher.FetchPrice(context.Background(), "test", "pair")
		assert.Zero(t, value)
		assert.ErrorIs(t, err, errInvalidGasPriceSelector)
		assert.Contains(t, err.Error(), "invalid-selector")
	})
	testHTTPResponseGetter := &mock.HttpResponseGetterStub{
		GetCalled: func(ctx context.Context, url string, response interface{}) error {
			assert.Equal(t, "api-url", url)
			gasStationResp := response.(*gasStationResponse)
			gasStationResp.Result.SafeGasPrice = "37"
			gasStationResp.Result.ProposeGasPrice = "38"
			gasStationResp.Result.FastGasPrice = "39"

			return nil
		},
	}

	t.Run("with SafeGasPrice should work", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		fetcher.config.Selector = evmSafeGasPrice
		fetcher.AddPair("test", "pair")
		fetcher.ResponseGetter = testHTTPResponseGetter
		value, err := fetcher.FetchPrice(context.Background(), "test", "pair")
		assert.Nil(t, err)
		assert.Equal(t, float64(37), value)
	})
	t.Run("with ProposeGasPrice should work", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		fetcher.config.Selector = evmProposeGasPrice
		fetcher.AddPair("test", "pair")
		fetcher.ResponseGetter = testHTTPResponseGetter
		value, err := fetcher.FetchPrice(context.Background(), "test", "pair")
		assert.Nil(t, err)
		assert.Equal(t, float64(38), value)
	})
	t.Run("with FastGasPrice should work", func(t *testing.T) {
		t.Parallel()

		fetcher := createMockEVMGasPriceFetcher()
		fetcher.config.Selector = evmFastGasPrice
		fetcher.AddPair("test", "pair")
		fetcher.ResponseGetter = testHTTPResponseGetter
		value, err := fetcher.FetchPrice(context.Background(), "test", "pair")
		assert.Nil(t, err)
		assert.Equal(t, float64(39), value)
	})
}
