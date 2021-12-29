package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	bitfinexPriceUrl = "https://api.bitfinex.com/v1/pubticker/%s%s"
)

type bitfinexPriceRequest struct {
	Price string `json:"last_price"`
}

type bitfinex struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (b *bitfinex) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	b.updateQuoteIfNeeded(&quote, bitfinexName)

	var bit bitfinexPriceRequest
	err := b.ResponseGetter.Get(ctx, fmt.Sprintf(bitfinexPriceUrl, base, quote), &bit)
	if err != nil {
		return 0, err
	}
	if bit.Price == "" {
		return 0, errInvalidResponseData
	}
	return StrToPositiveFloat64(bit.Price)
}

// Name returns the name
func (b *bitfinex) Name() string {
	return bitfinexName
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *bitfinex) IsInterfaceNil() bool {
	return b == nil
}
