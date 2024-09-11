package fetchers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-sdk-go/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func createMockArgsPriceFetcher() ArgsPriceFetcher {
	return ArgsPriceFetcher{
		FetcherName:        BinanceName,
		ResponseGetter:     &mock.HttpResponseGetterStub{},
		GraphqlGetter:      &mock.GraphqlResponseGetterStub{},
		XExchangeTokensMap: createMockMap(),
		EVMGasConfig:       EVMGasPriceFetcherConfig{},
	}
}

func TestNewPriceFetcher(t *testing.T) {
	t.Parallel()

	t.Run("invalid fetcher name should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceFetcher()
		args.FetcherName = "invalid name"
		pf, err := NewPriceFetcher(args)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errInvalidFetcherName))
		assert.True(t, strings.Contains(err.Error(), args.FetcherName))
	})
	t.Run("nil responseGetter should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceFetcher()
		args.ResponseGetter = nil
		pf, err := NewPriceFetcher(args)
		assert.Nil(t, pf)
		assert.Equal(t, errNilResponseGetter, err)
	})
	t.Run("nil graphqlGetter should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceFetcher()
		args.FetcherName = XExchangeName
		args.GraphqlGetter = nil
		pf, err := NewPriceFetcher(args)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errNilGraphqlGetter))
	})
	t.Run("nil map for xExchange should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceFetcher()
		args.FetcherName = XExchangeName
		args.XExchangeTokensMap = nil
		pf, err := NewPriceFetcher(args)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errNilXExchangeTokensMap))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		t.Run("Binance", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = BinanceName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.binance", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Bitfinex", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = BitfinexName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.bitfinex", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("CryptoCom", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = CryptocomName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.cryptocom", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Gemini", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = GeminiName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.gemini", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Hitbtc", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = HitbtcName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.hitbtc", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Huobi", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = HuobiName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.huobi", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Kraken", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = KrakenName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.kraken", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("Okx", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = OkxName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.okx", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("xExchange", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = XExchangeName
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.xExchange", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
		t.Run("EVM gas price", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsPriceFetcher()
			args.FetcherName = EVMGasPriceStation
			pf, err := NewPriceFetcher(args)
			assert.Equal(t, "*fetchers.evmGasPriceFetcher", fmt.Sprintf("%T", pf))
			assert.Nil(t, err)
		})
	})
}
