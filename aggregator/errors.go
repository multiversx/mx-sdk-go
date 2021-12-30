package aggregator

import "errors"

var (
	errNotEnoughResponses                  = errors.New("not enough responses to compute a valid price")
	errInvalidMinNumberOfResults           = errors.New("invalid minimum number of results")
	errInvalidNumberOfPriceFetchers        = errors.New("invalid number of price fetchers")
	errNilPriceFetcher                     = errors.New("nil price fetcher")
	errEmptyArgsPairsSlice                 = errors.New("empty pair arguments slice")
	errNilArgsPair                         = errors.New("nil pair argument")
	errInvalidPercentDifference            = errors.New("invalid percentage difference")
	errInvalidTrimPrecision                = errors.New("invalid trim precision")
	errNilPriceNotifee                     = errors.New("nil price notifee")
	errInvalidNumOfElementsToComputeMedian = errors.New("invalid number of elements to compute the median")

	// ErrInvalidResponseData signals that an invalid response has been provided
	ErrInvalidResponseData = errors.New("invalid response data")
)
