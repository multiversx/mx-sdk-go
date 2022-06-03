package fetchers

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

const (
	maiarPriceUrl = "https://api.elrond.com/mex-pairs/%s/%s"
)

type maiarPriceRequest struct {
	BasePrice  float64 `json:"basePrice"`
	QuotePrice float64 `json:"quotePrice"`
}

type maiar struct {
	aggregator.ResponseGetter
	baseFetcher
	maiarTokensMap map[string]MaiarTokensPair
}

// FetchPrice will fetch the price using the http client
func (m *maiar) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !m.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	maiarTokensPair, ok := m.fetchMaiarTokensPair(base, quote)
	if !ok {
		return 0, errInvalidPair
	}

	var mpr maiarPriceRequest
	err := m.ResponseGetter.Get(ctx, fmt.Sprintf(maiarPriceUrl, maiarTokensPair.Base, maiarTokensPair.Quote), &mpr)
	if err != nil {
		return 0, err
	}
	if mpr.BasePrice <= 0 {
		return 0, errInvalidResponseData
	}
	if mpr.QuotePrice <= 0 {
		return 0, errInvalidResponseData
	}
	price := mpr.BasePrice / mpr.QuotePrice
	return price, nil
}

func (m *maiar) fetchMaiarTokensPair(base, quote string) (MaiarTokensPair, bool) {
	pair := fmt.Sprintf("%s-%s", base, quote)
	mtp, ok := m.maiarTokensMap[pair]
	return mtp, ok
}

// Name returns the name
func (m *maiar) Name() string {
	return MaiarName
}

// IsInterfaceNil returns true if there is no value under the interface
func (m *maiar) IsInterfaceNil() bool {
	return m == nil
}
