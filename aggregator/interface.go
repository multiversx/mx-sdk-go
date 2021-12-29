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

// PriceNotifee defines the behavior of a component able to be notified over a price change
type PriceNotifee interface {
	PriceChanged(ctx context.Context, base string, quote string, price float64) error
	IsInterfaceNil() bool
}
