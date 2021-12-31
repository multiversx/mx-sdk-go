package fetchers

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPriceFetcher(t *testing.T) {
	t.Parallel()

	t.Run("invalid fetcher name should error", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher("invalid name", &mock.HttpResponseGetterStub{})
		assert.Nil(t, pf)
		assert.Equal(t, errInvalidFetcherName, err)
	})
	t.Run("nil responseGetter should error", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(binanceName, nil)
		assert.Nil(t, pf)
		assert.Equal(t, errNilResponseGetter, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		pf, err := NewPriceFetcher(binanceName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.binance", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(bitfinexName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.bitfinex", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(cryptocomName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.cryptocom", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(geminiName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.gemini", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(hitbtcName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.hitbtc", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(huobiName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.huobi", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(krakenName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.kraken", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
		pf, err = NewPriceFetcher(okexName, &mock.HttpResponseGetterStub{})
		assert.Equal(t, "*fetchers.okex", fmt.Sprintf("%T", pf))
		assert.Nil(t, err)
	})
}
