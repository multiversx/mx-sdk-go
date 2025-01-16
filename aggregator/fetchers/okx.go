package fetchers

import (
	"context"
	"fmt"

	"github.com/multiversx/mx-sdk-go/aggregator"
)

const (
	okxPriceUrl = "https://www.okx.com/api/v5/market/ticker?instId=%s-%s"
)

type okxPriceRequest struct {
	Data []okxTicker
}

type okxTicker struct {
	Price string `json:"last"`
}

type okx struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (o *okx) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !o.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = o.normalizeQuoteName(quote, OkxName)

	var opr okxPriceRequest
	err := o.ResponseGetter.Get(ctx, fmt.Sprintf(okxPriceUrl, base, quote), &opr)
	if err != nil {
		return 0, err
	}
	if len(opr.Data) == 0 {
		return 0, errInvalidResponseData
	}
	if opr.Data[0].Price == "" {
		return 0, errInvalidResponseData
	}

	return StrToPositiveFloat64(opr.Data[0].Price)
}

// Name returns the name
func (o *okx) Name() string {
	return OkxName
}

// IsInterfaceNil returns true if there is no value under the interface
func (o *okx) IsInterfaceNil() bool {
	return o == nil
}
