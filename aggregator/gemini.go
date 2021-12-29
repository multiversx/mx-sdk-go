package aggregator

import (
	"context"
	"fmt"
	"strings"
)

const (
	geminiPriceUrl = "https://api.gemini.com/v2/ticker/%s%s"
)

type geminiPriceRequest struct {
	Price string `json:"close"`
}

type gemini struct {
	ResponseGetter
}

// FetchPrice will fetch the price using the http client
func (g *gemini) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if strings.Contains(strings.ToUpper(base), QuoteUSDFiat) {
		quote = QuoteUSDFiat
	}

	var gpr geminiPriceRequest
	err := g.ResponseGetter.Get(ctx, fmt.Sprintf(geminiPriceUrl, base, quote), &gpr)
	if err != nil {
		return 0, err
	}
	if gpr.Price == "" {
		return 0, ErrInvalidResponseData
	}

	return StrToPositiveFloat64(gpr.Price)
}

// Name returns the name
func (g *gemini) Name() string {
	return "Gemini"
}

// IsInterfaceNil returns true if there is no value under the interface
func (g *gemini) IsInterfaceNil() bool {
	return g == nil
}
