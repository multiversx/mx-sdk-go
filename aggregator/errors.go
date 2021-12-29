package aggregator

import "errors"

var (
	errNotEnoughResponses                  = errors.New("not enough responses to compute a valid price")
	errInvalidMinNumberOfResults           = errors.New("invalid minimum number of results")
	errInvalidNumberOfPriceFetchers        = errors.New("invalid number of price fetchers")
	errNilPriceFetcher                     = errors.New("nil price fetcher")
	errInvalidNumOfElementsToComputeMedian = errors.New("invalid number of elements to compute the median")
)
