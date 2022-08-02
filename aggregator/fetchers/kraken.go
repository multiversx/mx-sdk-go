package fetchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	krakenPriceUrl = "https://api.kraken.com/0/public/Ticker?pair=%s%s"
)

type krakenPriceRequest struct {
	Result map[string]krakenPricePair `json:"result"`
}

type krakenPricePair struct {
	Price []string `json:"c"`
}

type kraken struct {
	aggregator.ResponseGetter
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (k *kraken) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !k.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	quote = k.normalizeQuoteName(quote, KrakenName)

	var hpr krakenPriceRequest
	err := k.ResponseGetter.Get(ctx, fmt.Sprintf(krakenPriceUrl, base, quote), &hpr)
	if err != nil {
		return 0, err
	}
	if len(hpr.Result) == 0 {
		return 0, errInvalidResponseData
	}
	for k, v := range hpr.Result {
		if k == "" || v.Price[0] == "" {
			return 0, errInvalidResponseData
		}

		if strings.Contains(k, base) || strings.Contains(k, quote) {
			return StrToPositiveFloat64(v.Price[0])
		}
	}

	return 0, errInvalidResponseData
}

// Name returns the name
func (k *kraken) Name() string {
	return KrakenName
}

// IsInterfaceNil returns true if there is no value under the interface
func (k *kraken) IsInterfaceNil() bool {
	return k == nil
}
