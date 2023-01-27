package fetchers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/multiversx/mx-sdk-go/aggregator"
)

const (
	// TODO EN-13146: extract this urls constants in a file
	dataApiUrl = "https://tools.multiversx.com/data-api/graphql"
	query      = "query MaiarPriceUrl($base: String!, $quote: String!) { trading { pair(first_token: $base, second_token: $quote) { price { last time } } } }"
)

type variables struct {
	BasePrice  string `json:"base"`
	QuotePrice string `json:"quote"`
}

type priceResponse struct {
	Last float64   `json:"last"`
	Time time.Time `json:"time"`
}

type graphqlResponse struct {
	Data struct {
		Trading struct {
			Pair struct {
				Price []priceResponse `json:"price"`
			} `json:"pair"`
		} `json:"trading"`
	} `json:"data"`
}

type xExchange struct {
	aggregator.GraphqlGetter
	baseFetcher
	xExchangeTokensMap map[string]XExchangeTokensPair
}

// FetchPrice will fetch the price using the http client
func (x *xExchange) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !x.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	xExchangeTokensPair, ok := x.fetchXExchangeTokensPair(base, quote)
	if !ok {
		return 0, errInvalidPair
	}

	vars, err := json.Marshal(variables{
		BasePrice:  xExchangeTokensPair.Base,
		QuotePrice: xExchangeTokensPair.Quote,
	})
	if err != nil {
		return 0, err
	}

	resp, err := x.GraphqlGetter.Query(ctx, dataApiUrl, query, string(vars))
	if err != nil {
		return 0, err
	}

	var graphqlResp graphqlResponse
	err = json.Unmarshal(resp, &graphqlResp)
	if err != nil {
		return 0, errInvalidGraphqlResponse
	}

	price := graphqlResp.Data.Trading.Pair.Price[0].Last

	if price <= 0 {
		return 0, errInvalidResponseData
	}
	return price, nil
}

func (x *xExchange) fetchXExchangeTokensPair(base, quote string) (XExchangeTokensPair, bool) {
	pair := fmt.Sprintf("%s-%s", base, quote)
	mtp, ok := x.xExchangeTokensMap[pair]
	return mtp, ok
}

// Name returns the name
func (x *xExchange) Name() string {
	return XExchangeName
}

// IsInterfaceNil returns true if there is no value under the interface
func (x *xExchange) IsInterfaceNil() bool {
	return x == nil
}
