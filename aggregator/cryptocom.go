package aggregator

import (
	"context"
	"fmt"
	"strings"
)

const (
	cryptocomPriceUrl = "https://api.crypto.com/v2/public/get-ticker?instrument_name=%s_%s"
)

type cryptocom struct {
	ResponseGetter
}

type cryptocomPriceRequest struct {
	Result cryptocomData `json:"result"`
}

type cryptocomData struct {
	Data cryptocomPair `json:"data"`
}

type cryptocomPair struct {
	Price float64 `json:"a"`
}

// FetchPrice will fetch the price using the http client
func (c *cryptocom) FetchPrice(ctx context.Context, base, quote string) (float64, error) {
	if strings.Contains(quote, QuoteUSDFiat) {
		quote = QuoteUSDT
	}

	var cpr cryptocomPriceRequest
	err := c.ResponseGetter.Get(ctx, fmt.Sprintf(cryptocomPriceUrl, base, quote), &cpr)
	if err != nil {
		return 0, err
	}
	if cpr.Result.Data.Price <= 0 {
		return 0, InvalidResponseDataErr
	}
	return cpr.Result.Data.Price, nil
}

// Name returns the name
func (c *cryptocom) Name() string {
	return "Crypto.com"
}

// IsInterfaceNil returns true if there is no value under the interface
func (c *cryptocom) IsInterfaceNil() bool {
	return c == nil
}
