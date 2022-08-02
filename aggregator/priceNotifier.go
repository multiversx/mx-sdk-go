package aggregator

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
)

const epsilon = 0.0001
const minAutoSendInterval = time.Second

// ArgsPriceNotifier is the argument DTO for the price notifier
type ArgsPriceNotifier struct {
	Pairs            []*ArgsPair
	Aggregator       PriceAggregator
	Notifee          PriceNotifee
	AutoSendInterval time.Duration
}

// ArgsPair is the argument DTO for a pair
type ArgsPair struct {
	Base                      string
	Quote                     string
	PercentDifferenceToNotify uint32
	TrimPrecision             float64
	DenominationFactor        uint64
	Exchanges                 map[string]struct{}
}

type priceInfo struct {
	price     float64
	timestamp int64
}

type notifyArgs struct {
	*ArgsPair
	newPrice          priceInfo
	lastNotifiedPrice float64
	index             int
}

type priceNotifier struct {
	mut                sync.Mutex
	priceAggregator    PriceAggregator
	pairs              []*ArgsPair
	lastNotifiedPrices []float64
	notifee            PriceNotifee
	autoSendInterval   time.Duration
	lastTimeAutoSent   time.Time
	timeSinceHandler   func(t time.Time) time.Duration
}

// NewPriceNotifier will create a new priceNotifier instance
func NewPriceNotifier(args ArgsPriceNotifier) (*priceNotifier, error) {
	err := checkArgsPriceNotifier(args)
	if err != nil {
		return nil, err
	}

	return &priceNotifier{
		priceAggregator:    args.Aggregator,
		pairs:              args.Pairs,
		lastNotifiedPrices: make([]float64, len(args.Pairs)),
		notifee:            args.Notifee,
		autoSendInterval:   args.AutoSendInterval,
		lastTimeAutoSent:   time.Now(),
		timeSinceHandler:   time.Since,
	}, nil
}

func checkArgsPriceNotifier(args ArgsPriceNotifier) error {
	if len(args.Pairs) < 1 {
		return ErrEmptyArgsPairsSlice
	}

	for idx, argsPair := range args.Pairs {
		if argsPair == nil {
			return fmt.Errorf("%w, index %d", ErrNilArgsPair, idx)
		}
		if argsPair.TrimPrecision < epsilon {
			return fmt.Errorf("%w, got %f for pair %s-%s", ErrInvalidTrimPrecision,
				argsPair.TrimPrecision, argsPair.Base, argsPair.Quote)
		}
		if argsPair.DenominationFactor == 0 {
			return fmt.Errorf("%w, got %d for pair %s-%s", ErrInvalidDenominationFactor,
				argsPair.DenominationFactor, argsPair.Base, argsPair.Quote)
		}
	}
	if args.AutoSendInterval < minAutoSendInterval {
		return fmt.Errorf("%w, minimum %v, got %v", ErrInvalidAutoSendInterval, minAutoSendInterval, args.AutoSendInterval)
	}
	if check.IfNil(args.Notifee) {
		return ErrNilPriceNotifee
	}
	if check.IfNil(args.Aggregator) {
		return ErrNilPriceAggregator
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

func (pn *priceNotifier) getAllPrices(ctx context.Context) ([]priceInfo, error) {
	fetchedPrices := make([]priceInfo, len(pn.pairs))
	for idx, pair := range pn.pairs {
		price, err := pn.priceAggregator.FetchPrice(ctx, pair.Base, pair.Quote)
		if err != nil {
			return nil, fmt.Errorf("%w while querying the pair %s-%s", err, pair.Base, pair.Quote)
		}

		fetchedPrice := priceInfo{
			price:     trim(price, pair.TrimPrecision),
			timestamp: time.Now().Unix(),
		}
		fetchedPrices[idx] = fetchedPrice
	}

	return fetchedPrices, nil
}

func (pn *priceNotifier) computeNotifyArgsSlice(fetchedPrices []priceInfo) []*notifyArgs {
	pn.mut.Lock()
	defer pn.mut.Unlock()

	shouldNotifyAll := pn.timeSinceHandler(pn.lastTimeAutoSent) > pn.autoSendInterval

	result := make([]*notifyArgs, 0, len(pn.pairs))
	for idx, pair := range pn.pairs {
		notifyArgsValue := &notifyArgs{
			ArgsPair:          pair,
			newPrice:          fetchedPrices[idx],
			lastNotifiedPrice: pn.lastNotifiedPrices[idx],
			index:             idx,
		}

		if shouldNotifyAll || shouldNotify(notifyArgsValue) {
			result = append(result, notifyArgsValue)
		}
	}

	if shouldNotifyAll {
		pn.lastTimeAutoSent = time.Now()
	}

	return result
}

func shouldNotify(notifyArgsValue *notifyArgs) bool {
	percentValue := float64(notifyArgsValue.PercentDifferenceToNotify) / 100
	shouldBypassPercentCheck := notifyArgsValue.lastNotifiedPrice < epsilon || percentValue < epsilon
	if shouldBypassPercentCheck {
		return true
	}

	absoluteChange := math.Abs(notifyArgsValue.lastNotifiedPrice - notifyArgsValue.newPrice.price)
	percentageChange := absoluteChange * 100 / notifyArgsValue.lastNotifiedPrice

	return percentageChange >= float64(notifyArgsValue.PercentDifferenceToNotify)
}

func (pn *priceNotifier) notify(ctx context.Context, notifyArgsSlice []*notifyArgs) error {
	if len(notifyArgsSlice) == 0 {
		return nil
	}

	args := make([]*ArgsPriceChanged, 0, len(notifyArgsSlice))
	for _, notify := range notifyArgsSlice {
		priceTrimmed := trim(notify.newPrice.price, notify.TrimPrecision)
		denominatedPrice := uint64(priceTrimmed * float64(notify.DenominationFactor))

		argPriceChanged := &ArgsPriceChanged{
			Base:               notify.Base,
			Quote:              notify.Quote,
			DenominatedPrice:   denominatedPrice,
			DenominationFactor: notify.DenominationFactor,
			Timestamp:          notify.newPrice.timestamp,
		}

		args = append(args, argPriceChanged)

		pn.mut.Lock()
		pn.lastNotifiedPrices[notify.index] = priceTrimmed
		pn.mut.Unlock()
	}

	return pn.notifee.PriceChanged(ctx, args)
}

// IsInterfaceNil returns true if there is no value under the interface
func (pn *priceNotifier) IsInterfaceNil() bool {
	return pn == nil
}
