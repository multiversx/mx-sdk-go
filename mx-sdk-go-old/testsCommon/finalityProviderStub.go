package testsCommon

import "context"

// FinalityProviderStub -
type FinalityProviderStub struct {
	CheckShardFinalizationCalled func(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error
}

// CheckShardFinalization -
func (stub *FinalityProviderStub) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
	if stub.CheckShardFinalizationCalled != nil {
		return stub.CheckShardFinalizationCalled(ctx, targetShardID, maxNoncesDelta)
	}

	return nil
}

// IsInterfaceNil -
func (stub *FinalityProviderStub) IsInterfaceNil() bool {
	return stub == nil
}
