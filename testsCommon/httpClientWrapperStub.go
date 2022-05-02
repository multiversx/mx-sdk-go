package testsCommon

import "context"

// HTTPClientWrapperStub -
type HTTPClientWrapperStub struct {
	GetHTTPCalled  func(ctx context.Context, endpoint string) ([]byte, error)
	PostHTTPCalled func(ctx context.Context, endpoint string, data []byte) ([]byte, error)
}

// GetHTTP -
func (stub *HTTPClientWrapperStub) GetHTTP(ctx context.Context, endpoint string) ([]byte, error) {
	if stub.GetHTTPCalled != nil {
		return stub.GetHTTPCalled(ctx, endpoint)
	}

	return make([]byte, 0), nil
}

// PostHTTP -
func (stub *HTTPClientWrapperStub) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, error) {
	if stub.PostHTTPCalled != nil {
		return stub.PostHTTPCalled(ctx, endpoint, data)
	}

	return make([]byte, 0), nil
}

// IsInterfaceNil -
func (stub *HTTPClientWrapperStub) IsInterfaceNil() bool {
	return stub == nil
}
