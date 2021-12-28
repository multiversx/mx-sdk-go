package aggregator

import (
	"fmt"
	"strings"
)

const (
	binancePriceUrl = "https://api.binance.com/api/v3/ticker/price?symbol=%s%s"
)

type binancePriceRequest struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type binance struct {
	ResponseGetter
}

// FetchPrice will fetch the price using the http client
func (b *binance) FetchPrice(base, quote string) (float64, error) {
	if strings.Contains(strings.ToUpper(quote), QuoteUSDFiat) {
		quote = QuoteUSDT
	}

	var bpr binancePriceRequest
	err := b.ResponseGetter.Get(fmt.Sprintf(binancePriceUrl, base, quote), &bpr)
	if err != nil {
		return 0, err
	}
	if bpr.Price == "" {
		return 0, InvalidResponseDataErr
	}

	return StrToFloat64(bpr.Price)
}

// Name returns the name
func (b *binance) Name() string {
	return "Binance"
}
