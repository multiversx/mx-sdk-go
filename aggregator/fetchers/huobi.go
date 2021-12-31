package fetchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	huobiPriceUrl = "https://api.huobi.pro/market/detail/merged?symbol=%s%s"
)

type huobiPriceRequest struct {
	Ticker huobiPriceTicker `json:"tick"`
}

type huobiPriceTicker struct {
	Price float64 `json:"close"`
}

type huobi struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (h *huobi) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	quote = h.normalizeQuoteName(quote, huobiName)

	var hpr huobiPriceRequest
	err := h.ResponseGetter.Get(ctx, fmt.Sprintf(huobiPriceUrl, strings.ToLower(base), strings.ToLower(quote)), &hpr)
	if err != nil {
		return 0, err
	}
	if hpr.Ticker.Price <= 0 {
		return 0, errInvalidResponseData
	}

	return hpr.Ticker.Price, nil
}

// Name returns the name
func (h *huobi) Name() string {
	return huobiName
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *huobi) IsInterfaceNil() bool {
	return h == nil
}
