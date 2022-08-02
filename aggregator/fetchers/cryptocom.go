package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	cryptocomPriceUrl = "https://api.crypto.com/v2/public/get-ticker?instrument_name=%s_%s"
)

type cryptocomPriceRequest struct {
	Result cryptocomData `json:"result"`
}

type cryptocomData struct {
	Data cryptocomPair `json:"data"`
}

type cryptocomPair struct {
	Price float64 `json:"a"`
}

type cryptocom struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (c *cryptocom) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if !c.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = c.normalizeQuoteName(quote, CryptocomName)

	var cpr cryptocomPriceRequest
	err := c.ResponseGetter.Get(ctx, fmt.Sprintf(cryptocomPriceUrl, base, quote), &cpr)
	if err != nil {
		return 0, err
	}
	if cpr.Result.Data.Price <= 0 {
		return 0, errInvalidResponseData
	}
	return cpr.Result.Data.Price, nil
}

// Name returns the name
func (c *cryptocom) Name() string {
	return CryptocomName
}

// IsInterfaceNil returns true if there is no value under the interface
func (c *cryptocom) IsInterfaceNil() bool {
	return c == nil
}
