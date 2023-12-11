package finalityProvider

import "context"

type disabledFinalityProvider struct {
}

// NewDisabledFinalityProvider returns a new instance of type disabledFinalityProvider
func NewDisabledFinalityProvider() *disabledFinalityProvider {
	return &disabledFinalityProvider{}
}

// CheckShardFinalization will always return nil
func (provider *disabledFinalityProvider) CheckShardFinalization(_ context.Context, _ uint32, _ uint64) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *disabledFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}
