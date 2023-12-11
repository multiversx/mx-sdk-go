package testsCommon

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-go/state"
)

// RawHeaderHandlerStub -
type RawHeaderHandlerStub struct {
	GetMetaBlockByHashCalled        func(hash string) (data.MetaHeaderHandler, error)
	GetShardBlockByHashCalled       func(shardId uint32, hash string) (data.HeaderHandler, error)
	GetValidatorsInfoPerEpochCalled func(epoch uint32) ([]*state.ShardValidatorInfo, []byte, error)
}

// GetMetaBlockByHash -
func (rh *RawHeaderHandlerStub) GetMetaBlockByHash(_ context.Context, hash string) (data.MetaHeaderHandler, error) {
	if rh.GetMetaBlockByHashCalled != nil {
		return rh.GetMetaBlockByHashCalled(hash)
	}
	return nil, nil
}

// GetShardBlockByHash -
func (rh *RawHeaderHandlerStub) GetShardBlockByHash(_ context.Context, shardId uint32, hash string) (data.HeaderHandler, error) {
	if rh.GetShardBlockByHashCalled != nil {
		return rh.GetShardBlockByHashCalled(shardId, hash)
	}
	return nil, nil
}

// GetValidatorsInfoPerEpoch -
func (rh *RawHeaderHandlerStub) GetValidatorsInfoPerEpoch(_ context.Context, epoch uint32) ([]*state.ShardValidatorInfo, []byte, error) {
	if rh.GetMetaBlockByHashCalled != nil {
		return rh.GetValidatorsInfoPerEpochCalled(epoch)
	}
	return nil, nil, nil
}

// IsInterfaceNil -
func (rh *RawHeaderHandlerStub) IsInterfaceNil() bool {
	return rh == nil
}
