package mock

import "context"

// GraphqlResponseGetterStub -
type GraphqlResponseGetterStub struct {
	GetCalled func(ctx context.Context, url string, query string, variables string) (interface{}, error)
}

// Query -
func (stub *GraphqlResponseGetterStub) Query(ctx context.Context, url string, query string, variables string) (interface{}, error) {
	if stub.GetCalled != nil {
		return stub.GetCalled(ctx, url, query, variables)
	}
	return nil, nil
}
