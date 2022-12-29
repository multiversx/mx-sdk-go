package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	hitbtcPriceUrl = "https://api.hitbtc.com/api/3/public/ticker/%s%s"
)

type hitbtcPriceRequest struct {
	Price string `json:"last"`
}

type hitbtc struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (h *hitbtc) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if !h.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = h.normalizeQuoteName(quote, HitbtcName)

	var hpr hitbtcPriceRequest
	err := h.ResponseGetter.Get(ctx, fmt.Sprintf(hitbtcPriceUrl, base, quote), &hpr)
	if err != nil {
		return 0, err
	}
	if hpr.Price == "" {
		return 0, errInvalidResponseData
	}
	return StrToPositiveFloat64(hpr.Price)
}

// Name returns the name
func (h *hitbtc) Name() string {
	return HitbtcName
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *hitbtc) IsInterfaceNil() bool {
	return h == nil
}
