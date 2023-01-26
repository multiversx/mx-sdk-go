package mock

import (
	"context"

	"github.com/multiversx/mx-sdk-go/aggregator"
)

// PriceNotifeeStub -
type PriceNotifeeStub struct {
	PriceChangedCalled func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error
}

// PriceChanged -
func (stub *PriceNotifeeStub) PriceChanged(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
	if stub.PriceChangedCalled != nil {
		return stub.PriceChangedCalled(ctx, args)
	}

	return nil
}

// IsInterfaceNil -
func (stub *PriceNotifeeStub) IsInterfaceNil() bool {
	return stub == nil
}
