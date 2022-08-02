package aggregator

import "errors"

var (
	// ErrNotEnoughResponses signals that not enough responses were fetched
	ErrNotEnoughResponses = errors.New("not enough responses to compute a valid price")
	// ErrInvalidMinNumberOfResults signals that an invalid minimum number of results value was provided
	ErrInvalidMinNumberOfResults = errors.New("invalid minimum number of results")
	// ErrInvalidNumberOfPriceFetchers signals that an invalid number of price fetchers were provided
	ErrInvalidNumberOfPriceFetchers = errors.New("invalid number of price fetchers")
	// ErrNilPriceFetcher signals that a nil price fetcher was provided
	ErrNilPriceFetcher = errors.New("nil price fetcher")
	// ErrNilPriceAggregator signals that a nil price aggregator was provided
	ErrNilPriceAggregator = errors.New("nil price aggregator")
	// ErrEmptyArgsPairsSlice signals that an empty arguments pair slice was provided
	ErrEmptyArgsPairsSlice = errors.New("empty pair arguments slice")
	// ErrNilArgsPair signals that a nil argument pair was found
	ErrNilArgsPair = errors.New("nil pair argument")
	// ErrNilPriceNotifee signals that a nil price notifee was provided
	ErrNilPriceNotifee = errors.New("nil price notifee")
	// ErrInvalidNumOfElementsToComputeMedian signals that an invalid number of elements to compute the median was provided
	ErrInvalidNumOfElementsToComputeMedian = errors.New("invalid number of elements to compute the median")
	// ErrInvalidDecimals signals that an invalid number of decimals was provided
	ErrInvalidDecimals = errors.New("invalid decimals")
	// ErrInvalidAutoSendInterval signals that an invalid auto send interval value was provided
	ErrInvalidAutoSendInterval = errors.New("invalid auto send interval")
	// ErrPairNotSupported signals that the pair is not supported by the fetcher
	ErrPairNotSupported = errors.New("pair not supported")
)
