package mock

import (
	"context"

	"github.com/ElrondNetwork/elrond-go-core/data/block"
	state "github.com/ElrondNetwork/elrond-go/state"
)

// RawHeaderHandlerStub -
type RawHeaderHandlerStub struct {
	GetMetaBlockByHashCalled        func(hash string) (*block.MetaBlock, error)
	GetShardBlockByHashCalled       func(shardId uint32, hash string) (*block.Header, error)
	GetValidatorsInfoPerEpochCalled func(epoch uint32) ([]*state.ShardValidatorInfo, []byte, error)
}

// GetMetaBlockByHash -
func (rh *RawHeaderHandlerStub) GetMetaBlockByHash(_ context.Context, hash string) (*block.MetaBlock, error) {
	if rh.GetMetaBlockByHashCalled != nil {
		return rh.GetMetaBlockByHashCalled(hash)
	}
	return nil, nil
}

// GetShardBlockByHash -
func (rh *RawHeaderHandlerStub) GetShardBlockByHash(_ context.Context, shardId uint32, hash string) (*block.Header, error) {
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
