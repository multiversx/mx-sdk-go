package headerCheck

import (
	"context"
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/prometheus/common/log"
)

type rawHeaderHandler struct {
	proxy      Proxy
	marshaller marshal.Marshalizer
}

func NewRawHeaderHandler(proxy Proxy, marshaller marshal.Marshalizer) (RawHeaderHandler, error) {
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

func (rh *rawHeaderHandler) GetMetaBlockByHash(ctx context.Context, hash string) (*block.MetaBlock, error) {
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

func (rh *rawHeaderHandler) GetShardBlockByHash(ctx context.Context, shardId uint32, hash string) (*block.Header, error) {
	metaBlockBytes, err := rh.proxy.GetRawBlockByHash(ctx, shardId, hash)
	if err != nil {
		return nil, err
	}

	blockHeader := &block.Header{}
	err = rh.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

func (rh *rawHeaderHandler) GetValidatorsInfoPerEpoch(ctx context.Context, epoch uint32) ([]*state.ShardValidatorInfo, []byte, error) {
	metaBlock, err := rh.getLastStartOfEpochMetaBlock(ctx)
	if err != nil {
		return nil, nil, err
	}
	randomness := metaBlock.GetPrevRandSeed()

	currEpoch := metaBlock.GetEpoch()
	for epoch <= currEpoch {
		if epoch == 0 {
			break
		}

		if epoch == currEpoch {
			break
		}

		newHash := hex.EncodeToString(metaBlock.EpochStart.Economics.PrevEpochStartHash)
		metaBlock, err = rh.GetMetaBlockByHash(ctx, newHash)
		if err != nil {
			return nil, nil, err
		}
		if metaBlock == nil {
			break
		}
		log.Info("fetched previous epoch")
		randomness = metaBlock.GetPrevRandSeed()

		currEpoch = metaBlock.GetEpoch()
	}

	validatorsInfoPerEpoch, err := rh.getValidatorsInfo(ctx, metaBlock)
	if err != nil {
		return nil, nil, err
	}

	return validatorsInfoPerEpoch, randomness, nil
}

func (rh *rawHeaderHandler) getLastStartOfEpochMetaBlock(ctx context.Context) (*block.MetaBlock, error) {
	nonce, err := rh.proxy.GetNonceAtEpochStart(ctx, core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	metaBlockBytes, err := rh.proxy.GetRawBlockByNonce(ctx, core.MetachainShardId, uint64(nonce))
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

func (rh *rawHeaderHandler) getValidatorsInfo(ctx context.Context, metaBlock *block.MetaBlock) ([]*state.ShardValidatorInfo, error) {
	allValidatorInfo := make([]*state.ShardValidatorInfo, 0)
	for _, miniBlockHeader := range metaBlock.MiniBlockHeaders {
		hash := hex.EncodeToString(miniBlockHeader.Hash)

		miniBlockBytes, err := rh.proxy.GetRawMiniBlockByHash(ctx, core.MetachainShardId, hash)
		if err != nil {
			return nil, err
		}

		miniBlock := &block.MiniBlock{}
		err = rh.marshaller.Unmarshal(miniBlock, miniBlockBytes)
		if err != nil {
			return nil, err
		}

		if miniBlock.Type != block.PeerBlock {
			continue
		}

		for _, txHash := range miniBlock.TxHashes {
			log.Info("miniblock",
				"epoch", metaBlock.GetEpoch(),
				"txHash", txHash)

			vid := &state.ShardValidatorInfo{}
			err := rh.marshaller.Unmarshal(vid, txHash)
			if err != nil {
				return nil, err
			}

			allValidatorInfo = append(allValidatorInfo, vid)
		}
	}

	return allValidatorInfo, nil
}

func (rh *rawHeaderHandler) parseMiniBlocks(miniblocks []*block.MiniBlock) ([]*state.ShardValidatorInfo, error) {
	allValidatorInfo := make([]*state.ShardValidatorInfo, 0)
	for _, peerMiniBlock := range miniblocks {
		if peerMiniBlock.Type != block.PeerBlock {
			continue
		}

		for _, txHash := range peerMiniBlock.TxHashes {
			vid := &state.ShardValidatorInfo{}
			err := rh.marshaller.Unmarshal(vid, txHash)
			if err != nil {
				return nil, err
			}

			allValidatorInfo = append(allValidatorInfo, vid)
		}
	}

	return allValidatorInfo, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rh *rawHeaderHandler) IsInterfaceNil() bool {
	return rh == nil
}
