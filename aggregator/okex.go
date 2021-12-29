package aggregator

import (
	"context"
	"fmt"
	"strings"
)

const (
	okexPriceUrl = "https://www.okex.com/api/v5/market/ticker?instId=%s-%s"
)

type okexPriceRequest struct {
	Data []okexTicker
}

type okexTicker struct {
	Price string `json:"last"`
}

type okex struct {
	ResponseGetter
}

// FetchPrice will fetch the price using the http client
func (o *okex) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	// TODO: (next PR) extract this to a base fetcher
	if strings.Contains(strings.ToUpper(quote), QuoteUSDFiat) {
		quote = QuoteUSDT
	}

	var opr okexPriceRequest
	err := o.ResponseGetter.Get(ctx, fmt.Sprintf(okexPriceUrl, base, quote), &opr)
	if err != nil {
		return 0, err
	}
	if len(opr.Data) == 0 {
		return 0, ErrInvalidResponseData
	}
	if opr.Data[0].Price == "" {
		return 0, ErrInvalidResponseData
	}

	return StrToPositiveFloat64(opr.Data[0].Price)
}

// Name returns the name
func (o *okex) Name() string {
	return "Okex"
}

// IsInterfaceNil returns true if there is no value under the interface
func (o *okex) IsInterfaceNil() bool {
	return o == nil
}
