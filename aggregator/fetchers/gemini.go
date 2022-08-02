package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	geminiPriceUrl = "https://api.gemini.com/v2/ticker/%s%s"
)

type geminiPriceRequest struct {
	Price string `json:"close"`
}

type gemini struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (g *gemini) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if !g.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = g.normalizeQuoteName(quote, GeminiName)

	var gpr geminiPriceRequest
	err := g.ResponseGetter.Get(ctx, fmt.Sprintf(geminiPriceUrl, base, quote), &gpr)
	if err != nil {
		return 0, err
	}
	if gpr.Price == "" {
		return 0, errInvalidResponseData
	}

	return StrToPositiveFloat64(gpr.Price)
}

// Name returns the name
func (g *gemini) Name() string {
	return GeminiName
}

// IsInterfaceNil returns true if there is no value under the interface
func (g *gemini) IsInterfaceNil() bool {
	return g == nil
}
