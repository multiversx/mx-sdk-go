package aggregator

import (
	"context"
	"fmt"
	"strings"
)

const (
	bitfinexPriceUrl = "https://api.bitfinex.com/v1/pubticker/%s%s"
)

type bitfinex struct {
	ResponseGetter
}

type bitfinexPriceRequest struct {
	Price string `json:"last_price"`
}

// FetchPrice will fetch the price using the http client
func (b *bitfinex) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if strings.Contains(quote, QuoteUSDFiat) {
		quote = QuoteUSDFiat
	}

	var bit bitfinexPriceRequest
	err := b.ResponseGetter.Get(ctx, fmt.Sprintf(bitfinexPriceUrl, "t"+base, quote), &bit)
	if err != nil {
		return 0, err
	}
	if err != nil {
		return 0, nil
	}
	if bit.Price == "" {
		return 0, InvalidResponseDataErr
	}
	return StrToFloat64(bit.Price)
}

// Name returns the name
func (b *bitfinex) Name() string {
	return "Bitfinex"
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *bitfinex) IsInterfaceNil() bool {
	return b == nil
}
