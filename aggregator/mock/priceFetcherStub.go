package mock

import "context"

// PriceFetcherStub -
type PriceFetcherStub struct {
	NameCalled       func() string
	FetchPriceCalled func(ctx context.Context, base string, quote string) (float64, error)
	AddPairCalled    func(base, quote string)
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

// AddPair -
func (stub *PriceFetcherStub) AddPair(base, quote string) {
	if stub.AddPairCalled != nil {
		stub.AddPairCalled(base, quote)
	}
}

// IsInterfaceNil -
func (stub *PriceFetcherStub) IsInterfaceNil() bool {
	return stub == nil
}
