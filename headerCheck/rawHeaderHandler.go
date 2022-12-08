package headerCheck

import (
	"context"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/state"
)

type rawHeaderHandler struct {
	proxy      Proxy
	marshaller marshal.Marshalizer
}

// NewRawHeaderHandler will create a new rawHeaderHandler instance
func NewRawHeaderHandler(proxy Proxy, marshaller marshal.Marshalizer) (*rawHeaderHandler, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	return &rawHeaderHandler{
		proxy:      proxy,
		marshaller: marshaller,
	}, nil
}

// GetMetaBlockByHash will return the MetaBlock based on the raw marshalized
// data from proxy
func (rh *rawHeaderHandler) GetMetaBlockByHash(ctx context.Context, hash string) (data.MetaHeaderHandler, error) {
	metaBlockBytes, err := rh.proxy.GetRawBlockByHash(ctx, core.MetachainShardId, hash)
	if err != nil {
		return nil, err
	}

	blockHeader := &block.MetaBlock{}
	err = rh.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

// GetShardBlockByHash will return the Header based on the raw marshalized data
// from proxy
func (rh *rawHeaderHandler) GetShardBlockByHash(ctx context.Context, shardId uint32, hash string) (data.HeaderHandler, error) {
	headerBytes, err := rh.proxy.GetRawBlockByHash(ctx, shardId, hash)
	if err != nil {
		return nil, err
	}

	blockHeader, err := process.UnmarshalShardHeader(rh.marshaller, headerBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

// GetStartOfEpochMetaBlock will return the start of epoch metablock based on
// the raw marshalized data from proxy
func (rh *rawHeaderHandler) GetStartOfEpochMetaBlock(ctx context.Context, epoch uint32) (data.MetaHeaderHandler, error) {
	metaBlockBytes, err := rh.proxy.GetRawStartOfEpochMetaBlock(ctx, epoch)
	if err != nil {
		return nil, err
	}

	blockHeader := &block.MetaBlock{}
	err = rh.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

// GetValidatorsInfoPerEpoch will return validators info based on start of
// epoch metablock for a specific epoch
func (rh *rawHeaderHandler) GetValidatorsInfoPerEpoch(ctx context.Context, epoch uint32) ([]*state.ShardValidatorInfo, []byte, error) {
	metaBlock, err := rh.GetStartOfEpochMetaBlock(ctx, epoch)
	if err != nil {
		log.Error("failed on getting start of epochs metaBlock")
		return nil, nil, err
	}
	randomness := metaBlock.GetPrevRandSeed()

	validatorsInfoPerEpoch, err := rh.proxy.GetValidatorsInfoByEpoch(ctx, epoch)
	if err != nil {
		log.Error("faileddddd on getting start of epochs metaBlock")
		return nil, nil, err
	}

	return validatorsInfoPerEpoch, randomness, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rh *rawHeaderHandler) IsInterfaceNil() bool {
	return rh == nil
}
