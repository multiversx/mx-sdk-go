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

type maiar struct {
	aggregator.GraphqlGetter
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

	variables, err := json.Marshal(variables{
		BasePrice:  maiarTokensPair.Base,
		QuotePrice: maiarTokensPair.Quote,
	})
	if err != nil {
		return 0, err
	}

	resp, err := m.GraphqlGetter.Query(ctx, dataApiUrl, query, string(variables))
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
