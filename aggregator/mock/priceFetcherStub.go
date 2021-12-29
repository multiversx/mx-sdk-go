package mock

import "context"

// PriceFetcherStub -
type PriceFetcherStub struct {
	NameCalled       func() string
	FetchPriceCalled func(ctx context.Context, base string, quote string) (float64, error)
}

// Name -
func (stub *PriceFetcherStub) Name() string {
	if stub.NameCalled != nil {
		return stub.NameCalled()
	}

	return ""
}

// FetchPrice -
func (stub *PriceFetcherStub) FetchPrice(ctx context.Context, base string, quote string) (float64, error) {
	if stub.FetchPriceCalled != nil {
		return stub.FetchPriceCalled(ctx, base, quote)
	}

	return 1, nil
}

// IsInterfaceNil -
func (stub *PriceFetcherStub) IsInterfaceNil() bool {
	return stub == nil
}
