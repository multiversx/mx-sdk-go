package mock

import "context"

// HttpResponseGetterStub -
type HttpResponseGetterStub struct {
	GetCalled func(ctx context.Context, url string, response interface{}) error
}

// Get -
func (stub *HttpResponseGetterStub) Get(ctx context.Context, url string, response interface{}) error {
	if stub.GetCalled != nil {
		return stub.GetCalled(ctx, url, response)
	}

	return nil
}
