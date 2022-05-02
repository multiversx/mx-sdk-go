package testsCommon

import (
	"context"
	"net/http"
)

// HTTPClientWrapperStub -
type HTTPClientWrapperStub struct {
	GetHTTPCalled  func(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTPCalled func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
}

// GetHTTP -
func (stub *HTTPClientWrapperStub) GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error) {
	if stub.GetHTTPCalled != nil {
		return stub.GetHTTPCalled(ctx, endpoint)
	}

	return make([]byte, 0), http.StatusOK, nil
}

// PostHTTP -
func (stub *HTTPClientWrapperStub) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
	if stub.PostHTTPCalled != nil {
		return stub.PostHTTPCalled(ctx, endpoint, data)
	}

	return make([]byte, 0), http.StatusOK, nil
}

// IsInterfaceNil -
func (stub *HTTPClientWrapperStub) IsInterfaceNil() bool {
	return stub == nil
}
