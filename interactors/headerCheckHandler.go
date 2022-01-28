package interactors

import (
	"context"
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerVerify"
)

type headerCheckHandler struct {
	proxy          Proxy
	headerVerifier HeaderVerifierHandler
	marshaller     marshal.Marshalizer
}

func NewHeaderCheckHandler(proxy Proxy) (*headerCheckHandler, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}

	networkConfig, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ratingsConfig, err := proxy.GetRatingsConfig(context.Background())
	if err != nil {
		return nil, err
	}

	enableEpochsConfig, err := proxy.GetEnableEpochsConfig(context.Background())
	if err != nil {
		return nil, err
	}

	headerVerifyArgs := headerVerify.ArgHeaderVerifier{
		RatingsConfig:      ratingsConfig,
		NetworkConfig:      networkConfig,
		EnableEpochsConfig: enableEpochsConfig,
	}

	headerVerifier, err := headerVerify.NewHeaderVerifier(headerVerifyArgs)
	if err != nil {
		return nil, err
	}

	return &headerCheckHandler{
		proxy:          proxy,
		headerVerifier: headerVerifier,
		marshaller:     headerVerifier.Marshaller(),
	}, nil
}

func (hch *headerCheckHandler) VerifyHeaderByHash(ctx context.Context, shardId uint32, hash string) (bool, error) {
	var err error

	var header coreData.HeaderHandler
	if shardId == core.MetachainShardId {
		header, err = hch.getMetaBlockByHash(ctx, hash)
		if err != nil {
			return false, err
		}
	} else {
		header, err = hch.getShardBlockByHash(ctx, shardId, hash)
		if err != nil {
			return false, err
		}
	}
	headerEpoch := header.GetEpoch()
	log.Info("fetched header in", "epoch", headerEpoch)

	if !hch.headerVerifier.IsInCache(headerEpoch) {
		log.Info("epoch", headerEpoch, "not in cache")

		validatorInfo, randomness, err := hch.getValidatorsInfoPerEpoch(ctx, headerEpoch)
		if err != nil {
			return false, err
		}

		err = hch.headerVerifier.SetNodesConfigPerEpoch(validatorInfo, headerEpoch, randomness)
		if err != nil {
			return false, err
		}
	}

	return hch.headerVerifier.VerifyHeader(header), nil
}

func (hch *headerCheckHandler) getLastStartOfEpochMetaBlock(ctx context.Context) (*block.MetaBlock, error) {
	nonce, err := hch.proxy.GetNonceAtEpochStart(ctx, core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	metaBlockBytes, err := hch.proxy.GetRawBlockByNonce(ctx, core.MetachainShardId, uint64(nonce))
	if err != nil {
		return nil, err
	}

	blockHeader := &block.MetaBlock{}
	err = hch.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

func (hch *headerCheckHandler) getValidatorsInfoPerEpoch(ctx context.Context, epoch uint32) ([]*state.ShardValidatorInfo, []byte, error) {
	metaBlock, err := hch.getLastStartOfEpochMetaBlock(ctx)
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
		metaBlock, err = hch.getMetaBlockByHash(ctx, newHash)
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

	validatorsInfoPerEpoch, err := hch.getValidatorsInfo(ctx, metaBlock)
	if err != nil {
		return nil, nil, err
	}

	return validatorsInfoPerEpoch, randomness, nil
}

func (hch *headerCheckHandler) getValidatorsInfo(ctx context.Context, metaBlock *block.MetaBlock) ([]*state.ShardValidatorInfo, error) {
	allValidatorInfo := make([]*state.ShardValidatorInfo, 0)
	for _, miniBlockHeader := range metaBlock.MiniBlockHeaders {
		hash := hex.EncodeToString(miniBlockHeader.Hash)

		miniBlockBytes, err := hch.proxy.GetRawMiniBlockByHash(ctx, core.MetachainShardId, hash)
		if err != nil {
			return nil, err
		}

		miniBlock := &block.MiniBlock{}
		err = hch.marshaller.Unmarshal(miniBlock, miniBlockBytes)
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
			err := hch.marshaller.Unmarshal(vid, txHash)
			if err != nil {
				return nil, err
			}

			allValidatorInfo = append(allValidatorInfo, vid)
		}
	}

	return allValidatorInfo, nil
}

func (hch *headerCheckHandler) parseMiniBlocks(miniblocks []*block.MiniBlock) ([]*state.ShardValidatorInfo, error) {
	allValidatorInfo := make([]*state.ShardValidatorInfo, 0)
	for _, peerMiniBlock := range miniblocks {
		if peerMiniBlock.Type != block.PeerBlock {
			continue
		}

		for _, txHash := range peerMiniBlock.TxHashes {
			vid := &state.ShardValidatorInfo{}
			err := hch.marshaller.Unmarshal(vid, txHash)
			if err != nil {
				return nil, err
			}

			allValidatorInfo = append(allValidatorInfo, vid)
		}
	}

	return allValidatorInfo, nil
}

func (hch *headerCheckHandler) getMetaBlockByHash(ctx context.Context, hash string) (*block.MetaBlock, error) {
	metaBlockBytes, err := hch.proxy.GetRawBlockByHash(ctx, core.MetachainShardId, hash)
	if err != nil {
		return nil, err
	}

	blockHeader := &block.MetaBlock{}
	err = hch.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

func (hch *headerCheckHandler) getShardBlockByHash(ctx context.Context, shardId uint32, hash string) (*block.Header, error) {
	metaBlockBytes, err := hch.proxy.GetRawBlockByHash(ctx, shardId, hash)
	if err != nil {
		return nil, err
	}

	blockHeader := &block.Header{}
	err = hch.marshaller.Unmarshal(blockHeader, metaBlockBytes)
	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}
