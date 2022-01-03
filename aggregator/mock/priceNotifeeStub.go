package mock

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
)

// PriceNotifeeStub -
type PriceNotifeeStub struct {
	PricesChangedCalled func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error
}

// PricesChanged -
func (stub *PriceNotifeeStub) PricesChanged(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
	if stub.PricesChangedCalled != nil {
		return stub.PricesChangedCalled(ctx, args)
	}

	return nil
}

// IsInterfaceNil -
func (stub *PriceNotifeeStub) IsInterfaceNil() bool {
	return stub == nil
}
