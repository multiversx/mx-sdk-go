package headerCheck

import (
	"context"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/headerCheck")

// ArgsHeaderVerifier holds all dependencies required by headerVerifier in
// order to create a new instance
type ArgsHeaderVerifier struct {
	HeaderHandler     RawHeaderHandler
	HeaderSigVerifier HeaderSigVerifierHandler
	NodesCoordinator  nodesCoordinator.EpochsConfigUpdateHandler
}

type headerVerifier struct {
	rawHeaderHandler  RawHeaderHandler
	headerSigVerifier HeaderSigVerifierHandler
	nodesCoordinator  nodesCoordinator.EpochsConfigUpdateHandler
}

// NewHeaderVerifier creates new instance of headerVerifier
func NewHeaderVerifier(args ArgsHeaderVerifier) (*headerVerifier, error) {
	err := checkArguments(args)
	if err != nil {
		return nil, err
	}

	return &headerVerifier{
		rawHeaderHandler:  args.HeaderHandler,
		headerSigVerifier: args.HeaderSigVerifier,
		nodesCoordinator:  args.NodesCoordinator,
	}, nil
}

func checkArguments(arguments ArgsHeaderVerifier) error {
	if check.IfNil(arguments.HeaderHandler) {
		return ErrNilRawHeaderHandler
	}
	if check.IfNil(arguments.NodesCoordinator) {
		return ErrNilNodesCoordinator
	}
	if check.IfNil(arguments.HeaderSigVerifier) {
		return ErrNilHeaderSigVerifier
	}

	return nil
}

// VerifyHeaderSignatureByHash verifies wether a header signature matches by providing
// the hash and shard where the header belongs to
func (hch *headerVerifier) VerifyHeaderSignatureByHash(ctx context.Context, shardId uint32, hash string) (bool, error) {
	header, err := hch.fetchHeaderByHashAndShard(ctx, shardId, hash)
	if err != nil {
		return false, err
	}

	headerEpoch := header.GetEpoch()
	log.Debug("fetched header in", "epoch", headerEpoch)

	if !hch.nodesCoordinator.IsEpochInConfig(headerEpoch) {
		log.Info("nodes config is set for epoch", "epoch", headerEpoch)
		err := hch.updateNodesConfigPerEpoch(ctx, headerEpoch)
		if err != nil {
			return false, err
		}
	}

	err = hch.headerSigVerifier.VerifySignature(header)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (hch *headerVerifier) fetchHeaderByHashAndShard(ctx context.Context, shardId uint32, hash string) (coreData.HeaderHandler, error) {
	var err error
	var header coreData.HeaderHandler

	if shardId == core.MetachainShardId {
		header, err = hch.rawHeaderHandler.GetMetaBlockByHash(ctx, hash)
		if err != nil {
			return nil, err
		}
	} else {
		header, err = hch.rawHeaderHandler.GetShardBlockByHash(ctx, shardId, hash)
		if err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (hch *headerVerifier) updateNodesConfigPerEpoch(ctx context.Context, epoch uint32) error {
	log.Debug("epoch", epoch, "not in cache")

	validatorInfo, randomness, err := hch.rawHeaderHandler.GetValidatorsInfoPerEpoch(ctx, epoch)
	if err != nil {
		return err
	}

	err = hch.nodesCoordinator.SetNodesConfigFromValidatorsInfo(epoch, randomness, validatorInfo)
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (hch *headerVerifier) IsInterfaceNil() bool {
	return hch == nil
}
