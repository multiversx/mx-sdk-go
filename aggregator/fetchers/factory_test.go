package fetchers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-sdk-go/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPriceFetcher(t *testing.T) {
	t.Parallel()

	t.Run("invalid fetcher name should error", func(t *testing.T) {
		t.Parallel()

		name := "invalid name"
		pf, err := NewPriceFetcher(name, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, nil)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errInvalidFetcherName))
		assert.True(t, strings.Contains(err.Error(), name))
	})
	t.Run("nil responseGetter should error", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(BinanceName, nil, &mock.GraphqlResponseGetterStub{}, nil)
		assert.Nil(t, pf)
		assert.Equal(t, errNilResponseGetter, err)
	})
	t.Run("nil graphqlGetter should error", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(MaiarName, &mock.HttpResponseGetterStub{}, nil, nil)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errNilGraphqlGetter))
	})
	t.Run("nil map for maiar should error", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(MaiarName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, nil)
		assert.Nil(t, pf)
		assert.True(t, errors.Is(err, errNilMaiarTokensMap))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(BinanceName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.binance", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(BitfinexName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.bitfinex", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(CryptocomName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.cryptocom", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(GeminiName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.gemini", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(HitbtcName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.hitbtc", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(HuobiName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.huobi", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(KrakenName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.kraken", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(OkexName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.okex", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(MaiarName, &mock.HttpResponseGetterStub{}, &mock.GraphqlResponseGetterStub{}, createMockMap())
		assert.Equal(t, "*fetchers.maiar", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
	})
}
