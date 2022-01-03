package aggregator

import "context"

// ResponseGetter is the component able to execute a get operation on the provided URL
type ResponseGetter interface {
	Get(ctx context.Context, url string, response interface{}) error
}

// PriceFetcher defines the behavior of a component able to query the price of a provided pair
type PriceFetcher interface {
	Name() string
	FetchPrice(ctx context.Context, base string, quote string) (float64, error)
	IsInterfaceNil() bool
}

// ArgsPriceChanged is the argument used when notifying the notifee instance
type ArgsPriceChanged struct {
	Base               string
	Quote              string
	DenominatedPrice   uint64
	DenominationFactor uint64
}

// PriceNotifee defines the behavior of a component able to be notified over a price change
type PriceNotifee interface {
	PriceChanged(ctx context.Context, priceChanges []*ArgsPriceChanged) error
	IsInterfaceNil() bool
}
