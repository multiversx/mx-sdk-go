package headerCheck

import (
	"context"

	coreData "github.com/multiversx/mx-chain-core-go/data"
)

func (hch *headerVerifier) FetchHeaderByHashAndShard(ctx context.Context, shardId uint32, hash string) (coreData.HeaderHandler, error) {
	return hch.fetchHeaderByHashAndShard(ctx, shardId, hash)
}

func (hch *headerVerifier) UpdateNodesConfigPerEpoch(ctx context.Context, epoch uint32) error {
	return hch.updateNodesConfigPerEpoch(ctx, epoch)
}
