package aggregator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const minResultsNum = 1

var log = logger.GetOrCreate("elrond-sdk-erdgo/aggregator")

// ArgsPriceAggregator is the DTO used in the NewPriceAggregator function
type ArgsPriceAggregator struct {
	PriceFetchers []PriceFetcher
	MinResultsNum int
}

type priceAggregator struct {
	priceFetchers []PriceFetcher
	minResultsNum int
}

// NewPriceAggregator creates a new priceAggregator instance
func NewPriceAggregator(args ArgsPriceAggregator) (*priceAggregator, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &priceAggregator{
		priceFetchers: args.PriceFetchers,
		minResultsNum: args.MinResultsNum,
	}, nil
}

func checkArgs(args ArgsPriceAggregator) error {
	if args.MinResultsNum < minResultsNum {
		return fmt.Errorf("%w, provided: %d, minimum accepted: %d", ErrInvalidMinNumberOfResults, args.MinResultsNum, minResultsNum)
	}
	if len(args.PriceFetchers) < args.MinResultsNum {
		return fmt.Errorf("%w, len(args.PriceFetchers): %d, MinResultsNum: %d", ErrInvalidNumberOfPriceFetchers,
			len(args.PriceFetchers), args.MinResultsNum)
	}
	for idx, pf := range args.PriceFetchers {
		if check.IfNil(pf) {
			return fmt.Errorf("%w, index: %d", ErrNilPriceFetcher, idx)
		}
	}

	return nil
}

// FetchPrice will try to fetch the price based on the provided array of price fetchers
func (pa *priceAggregator) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	var wg sync.WaitGroup
	var mut sync.Mutex
	var prices []float64

	baseUpper := strings.ToUpper(base)
	quoteUpper := strings.ToUpper(quote)

	wg.Add(len(pa.priceFetchers))
	for _, pf := range pa.priceFetchers {
		go func(priceFetcher PriceFetcher) {
			defer wg.Done()
			price, err := priceFetcher.FetchPrice(ctx, baseUpper, quoteUpper)

			if err == ErrPairNotSupported {
				log.Trace("pair not supported",
					"price fetcher", priceFetcher.Name(),
					"base", baseUpper,
					"quote", quoteUpper,
				)
				return
			}

			if err != nil {
				log.Debug("failed to fetch price",
					"price fetcher", priceFetcher.Name(),
					"base", baseUpper,
					"quote", quoteUpper,
					"err", err.Error(),
				)
				return
			}

			mut.Lock()
			prices = append(prices, price)
			mut.Unlock()
		}(pf)
	}
	wg.Wait()

	if len(prices) < pa.minResultsNum {
		return 0, ErrNotEnoughResponses
	}

	return computeMedian(prices)
}

// Name returns the name
func (pa *priceAggregator) Name() string {
	return "price aggregator"
}

// IsInterfaceNil returns true if there is no value under the interface
func (pa *priceAggregator) IsInterfaceNil() bool {
	return pa == nil
}
