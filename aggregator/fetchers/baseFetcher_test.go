package fetchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateQuoteIfNeeded(t *testing.T) {
	t.Parallel()

	t.Run("updating to usdt", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		quote := "AAA USD AAA"
		base.updateQuoteIfNeeded(&quote, binanceName)
		assert.Equal(t, quoteUSDT, quote)
	})
	t.Run("updating to usd", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		quote := "AAA USD AAA"
		base.updateQuoteIfNeeded(&quote, "other fetcher name")
		assert.Equal(t, quoteUSDFiat, quote)
	})
	t.Run("update not needed", func(t *testing.T) {
		t.Parallel()

		base := baseFetcher{}
		providedQuote := "custom quote"
		quote := providedQuote
		base.updateQuoteIfNeeded(&quote, "other fetcher name")
		assert.Equal(t, providedQuote, quote)
	})
}
