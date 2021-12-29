package aggregator

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

const epsilon = 0.0001

// ArgsPriceNotifier is the argument DTO for the price notifier
type ArgsPriceNotifier struct {
	Pairs   []*ArgsPair
	Fetcher PriceFetcher
	Notifee PriceNotifee
}

// ArgsPair is the argument DTO for a pair
type ArgsPair struct {
	Base                      string
	Quote                     string
	PercentDifferenceToNotify uint32
	TrimPrecision             float64
}

type notifyArgs struct {
	*ArgsPair
	newPrice          float64
	lastNotifiedPrice float64
	index             int
}

type priceNotifier struct {
	priceFetcher          PriceFetcher
	pairs                 []*ArgsPair
	mutLastNotifiedPrices sync.RWMutex
	lastNotifiedPrices    []float64
	notifee               PriceNotifee
}

// NewPriceNotifier will create a new priceNotifier instance
func NewPriceNotifier(args ArgsPriceNotifier) (*priceNotifier, error) {
	err := checkArgsPriceNotifier(args)
	if err != nil {
		return nil, err
	}

	return &priceNotifier{
		priceFetcher:       args.Fetcher,
		pairs:              args.Pairs,
		lastNotifiedPrices: make([]float64, len(args.Pairs)),
		notifee:            args.Notifee,
	}, nil
}

func checkArgsPriceNotifier(args ArgsPriceNotifier) error {
	if len(args.Pairs) < 1 {
		return errEmptyArgsPairsSlice
	}

	for idx, argsPair := range args.Pairs {
		if argsPair == nil {
			return fmt.Errorf("%w, index %d", errNilArgsPair, idx)
		}
		if argsPair.PercentDifferenceToNotify == 0 {
			return fmt.Errorf("%w, got %d for pair %s-%s", errInvalidPercentDifference,
				argsPair.PercentDifferenceToNotify, argsPair.Base, argsPair.Quote)
		}
		if argsPair.TrimPrecision < epsilon {
			return fmt.Errorf("%w, got %d for pair %s-%s", errInvalidTrimPrecision,
				argsPair.PercentDifferenceToNotify, argsPair.Base, argsPair.Quote)
		}
	}

	if check.IfNil(args.Notifee) {
		return errNilPriceNotifee
	}
	if check.IfNil(args.Fetcher) {
		return errNilPriceFetcher
	}

	return nil
}

// Execute will trigger the price fetching and notification if the new price exceeded provided percentage change
func (pn *priceNotifier) Execute(ctx context.Context) error {
	fetchedPrices, err := pn.getAllPrices(ctx)
	if err != nil {
		return err
	}

	notifyArgsSlice := pn.computeNotifyArgsSlice(fetchedPrices)

	return pn.notify(ctx, notifyArgsSlice)
}

func (pn *priceNotifier) getAllPrices(ctx context.Context) ([]float64, error) {
	fetchedPrices := make([]float64, len(pn.pairs))
	for idx, pair := range pn.pairs {
		price, err := pn.priceFetcher.FetchPrice(ctx, pair.Base, pair.Quote)
		if err != nil {
			return nil, fmt.Errorf("%w while querying the pair %s-%s", err, pair.Base, pair.Quote)
		}

		fetchedPrices[idx] = trim(price, pair.TrimPrecision)
	}

	return fetchedPrices, nil
}

func (pn *priceNotifier) computeNotifyArgsSlice(fetchedPrices []float64) []*notifyArgs {
	pn.mutLastNotifiedPrices.RLock()
	defer pn.mutLastNotifiedPrices.RUnlock()

	result := make([]*notifyArgs, 0, len(pn.pairs))
	for idx, pair := range pn.pairs {
		notifyArgsValue := &notifyArgs{
			ArgsPair:          pair,
			newPrice:          fetchedPrices[idx],
			lastNotifiedPrice: pn.lastNotifiedPrices[idx],
			index:             idx,
		}

		if shouldNotify(notifyArgsValue) {
			result = append(result, notifyArgsValue)
		}
	}

	return result
}

func shouldNotify(notifyArgsValue *notifyArgs) bool {
	if notifyArgsValue.lastNotifiedPrice < epsilon {
		return true
	}

	absoluteChange := math.Abs(notifyArgsValue.lastNotifiedPrice - notifyArgsValue.newPrice)
	percentageChange := absoluteChange * 100 / notifyArgsValue.lastNotifiedPrice

	return percentageChange >= float64(notifyArgsValue.PercentDifferenceToNotify)
}

func (pn *priceNotifier) notify(ctx context.Context, notifyArgsSlice []*notifyArgs) error {
	var lastErr error
	for _, notify := range notifyArgsSlice {
		err := pn.notifee.PriceChanged(ctx, notify.Base, notify.Quote, notify.newPrice)
		if err != nil {
			log.Error("error notifying", "base", notify.Base, "quote", notify.Quote,
				"new price", notify.newPrice, "error", err)
			lastErr = err
			continue
		}

		pn.mutLastNotifiedPrices.Lock()
		pn.lastNotifiedPrices[notify.index] = notify.newPrice
		pn.mutLastNotifiedPrices.Unlock()
	}

	return lastErr
}

// IsInterfaceNil returns true if there is no value under the interface
func (pn *priceNotifier) IsInterfaceNil() bool {
	return pn == nil
}
