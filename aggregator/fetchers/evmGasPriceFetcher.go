package fetchers

import (
	"context"
	"fmt"

	"github.com/multiversx/mx-sdk-go/aggregator"
)

const (
	evmFastGasPrice    = "FastGasPrice"
	evmSafeGasPrice    = "SafeGasPrice"
	evmProposeGasPrice = "ProposeGasPrice"
)

// EVMGasPriceFetcherConfig represents the config DTO used for the gas price fetcher
type EVMGasPriceFetcherConfig struct {
	ApiURL   string
	Selector string
}

type gasStationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		LastBlock       string `json:"LastBlock"`
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
		SuggestBaseFee  string `json:"suggestBaseFee"`
		GasUsedRatio    string `json:"gasUsedRatio"`
	} `json:"result"`
}

type evmGasPriceFetcher struct {
	aggregator.ResponseGetter
	config EVMGasPriceFetcherConfig
	baseFetcher
}

// FetchPrice will fetch the price using the http client
func (fetcher *evmGasPriceFetcher) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if !fetcher.hasPair(base, quote) {
		return 0, aggregator.ErrPairNotSupported
	}

	response := &gasStationResponse{}
	err := fetcher.ResponseGetter.Get(ctx, fmt.Sprintf(fetcher.config.ApiURL), response)
	if err != nil {
		return 0, err
	}

	latestGasPrice := 0
	switch fetcher.config.Selector {
	case evmFastGasPrice:
		_, err = fmt.Sscanf(response.Result.FastGasPrice, "%d", &latestGasPrice)
	case evmProposeGasPrice:
		_, err = fmt.Sscanf(response.Result.ProposeGasPrice, "%d", &latestGasPrice)
	case evmSafeGasPrice:
		_, err = fmt.Sscanf(response.Result.SafeGasPrice, "%d", &latestGasPrice)
	default:
		err = fmt.Errorf("%w: %q", errInvalidGasPriceSelector, fetcher.config.Selector)
	}

	return float64(latestGasPrice), err
}

// Name returns the name
func (fetcher *evmGasPriceFetcher) Name() string {
	return fmt.Sprintf("%s when using selector %s", EVMGasPriceStation, fetcher.config.Selector)
}

// IsInterfaceNil returns true if there is no value under the interface
func (fetcher *evmGasPriceFetcher) IsInterfaceNil() bool {
	return fetcher == nil
}
