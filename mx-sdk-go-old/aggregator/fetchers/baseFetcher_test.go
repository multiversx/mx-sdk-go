package fetchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_normalizeQuoteName(t *testing.T) {
	t.Parallel()

	t.Run("updating to usdt", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		quote := "AAA USD AAA"
		quote = base.normalizeQuoteName(quote, BinanceName)
		assert.Equal(t, quoteUSDT, quote)
	})
	t.Run("updating to usd", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		quote := "AAA USD AAA"
		quote = base.normalizeQuoteName(quote, "other fetcher name")
		assert.Equal(t, quoteUSDFiat, quote)
	})
	t.Run("update not needed", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		providedQuote := "custom quote"
		quote := providedQuote
		quote = base.normalizeQuoteName(quote, "other fetcher name")
		assert.Equal(t, providedQuote, quote)
	})
}

func TestBaseFetcher_knownPairs(t *testing.T) {
	t.Parallel()

	b := baseFetcher{
		knownPairs: make(map[string]struct{}),
	}
	base := "base"
	quote := "quote"
	assert.False(t, b.hasPair(base, quote))
	b.AddPair(base, quote)
	assert.True(t, b.hasPair(base, quote))
}
