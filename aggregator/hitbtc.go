package aggregator

import (
	"context"
	"fmt"
	"strings"
)

const (
	hitbtcPriceUrl = "https://api.hitbtc.com/api/3/public/ticker/%s%s"
)

type hitbtc struct {
	ResponseGetter
}

type hitbtcPriceRequest struct {
	Price string `json:"last"`
}

// FetchPrice will fetch the price using the http client
func (b *hitbtc) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if strings.Contains(quote, QuoteUSDFiat) {
		quote = QuoteUSDT
	}

	var hpr hitbtcPriceRequest
	err := b.ResponseGetter.Get(ctx, fmt.Sprintf(hitbtcPriceUrl, base, quote), &hpr)
	if err != nil {
		return 0, err
	}
	if hpr.Price == "" {
		return 0, ErrInvalidResponseData
	}
	price, err := StrToFloat64(hpr.Price)
	if err != nil {
		return 0, err
	}
	if price <= 0 {
		return 0, ErrInvalidResponseData
	}
	return price, nil
}

// Name returns the name
func (b *hitbtc) Name() string {
	return "HitBTC"
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *hitbtc) IsInterfaceNil() bool {
	return b == nil
}
