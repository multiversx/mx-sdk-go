package mock

import "context"

// ExecutorStub -
type ExecutorStub struct {
	ExecuteCalled func(ctx context.Context) error
}

// Execute -
func (es *ExecutorStub) Execute(ctx context.Context) error {
	if es.ExecuteCalled != nil {
		return es.ExecuteCalled(ctx)
	}

	return nil
}

// IsInterfaceNil -
func (es *ExecutorStub) IsInterfaceNil() bool {
	return es == nil
}
