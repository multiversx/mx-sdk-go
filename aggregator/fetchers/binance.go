package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	binancePriceUrl = "https://api.binance.com/api/v3/ticker/price?symbol=%s%s"
)

type binancePriceRequest struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type binance struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (b *binance) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !b.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = b.normalizeQuoteName(quote, BinanceName)

	var bpr binancePriceRequest
	err := b.ResponseGetter.Get(ctx, fmt.Sprintf(binancePriceUrl, base, quote), &bpr)
	if err != nil {
		return 0, err
	}
	if bpr.Price == "" {
		return 0, errInvalidResponseData
	}

	return StrToPositiveFloat64(bpr.Price)
}

// Name returns the name
func (b *binance) Name() string {
	return BinanceName
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *binance) IsInterfaceNil() bool {
	return b == nil
}
