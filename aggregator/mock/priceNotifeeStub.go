package mock

import "context"

// PriceNotifeeStub -
type PriceNotifeeStub struct {
	PriceChangedCalled func(ctx context.Context, base string, quote string, price float64) error
}

// PriceChanged -
func (stub *PriceNotifeeStub) PriceChanged(ctx context.Context, base string, quote string, price float64) error {
	if stub.PriceChangedCalled != nil {
		return stub.PriceChangedCalled(ctx, base, quote, price)
	}

	return nil
}

// IsInterfaceNil -
func (stub *PriceNotifeeStub) IsInterfaceNil() bool {
	return stub == nil
}
