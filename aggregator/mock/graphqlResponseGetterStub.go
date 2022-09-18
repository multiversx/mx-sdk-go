package mock

import "context"

// GraphqlResponseGetterStub -
type GraphqlResponseGetterStub struct {
	GetCalled func(ctx context.Context, url string, query string, variables string) ([]byte, error)
}

// Query -
func (stub *GraphqlResponseGetterStub) Query(ctx context.Context, url string, query string, variables string) ([]byte, error) {
	if stub.GetCalled != nil {
		return stub.GetCalled(ctx, url, query, variables)
	}
	return nil, nil
}
