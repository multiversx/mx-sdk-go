package finalityProvider

import (
	"context"
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
)

type nonces struct {
	current  uint64
	highest  uint64
	probable uint64
}

type nodeFinalityProvider struct {
	proxy proxy
}

// NewNodeFinalityProvider creates a new instance of type nodeFinalityProvider
func NewNodeFinalityProvider(proxy proxy) (*nodeFinalityProvider, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}

	return &nodeFinalityProvider{
		proxy: proxy,
	}, nil
}

// CheckShardFinalization will query the proxy and check if the target shard ID has a current nonce close to the highest nonce
// nonce <= highest_nonce + maxNoncesDelta
// it also checks the probable nonce to determine (with high degree of precision) if the node is syncing:
// nonce + maxNoncesDelta < probable_nonce
func (provider *nodeFinalityProvider) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
	if maxNoncesDelta < core.MinAllowedDeltaToFinal {
		return fmt.Errorf("%w, provided: %d, minimum: %d", ErrInvalidAllowedDeltaToFinal, maxNoncesDelta, core.MinAllowedDeltaToFinal)
	}

	result, err := provider.getNonces(ctx, targetShardID)
	if err != nil {
		return err
	}

	if result.current+maxNoncesDelta < result.probable {
		return fmt.Errorf("shardID %d is syncing, probable nonce is %d, current nonce is %d, max delta: %d",
			targetShardID, result.probable, result.current, maxNoncesDelta)
	}
	if result.current <= result.highest+maxNoncesDelta {
		log.Trace("nodeFinalityProvider.CheckShardFinalization - shard is in sync",
			"shardID", targetShardID, "highest nonce", result.highest, "probable nonce", result.probable,
			"current nonce", result.current, "max delta", maxNoncesDelta)
		return nil
	}

	return fmt.Errorf("shardID %d is stuck, highest nonce is %d, current nonce is %d, max delta: %d",
		targetShardID, result.highest, result.current, maxNoncesDelta)
}

func (provider *nodeFinalityProvider) getNonces(ctx context.Context, targetShardID uint32) (nonces, error) {
	networkStatusShard, err := provider.proxy.GetNetworkStatus(ctx, targetShardID)
	if err != nil {
		return nonces{}, err
	}

	result := nonces{
		current:  networkStatusShard.Nonce,
		highest:  networkStatusShard.HighestNonce,
		probable: networkStatusShard.ProbableHighestNonce,
	}

	isEmpty := result.current == 0 && result.highest == 0 && result.probable == 0
	if isEmpty {
		return nonces{}, ErrNodeNotStarted
	}

	return result, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *nodeFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}
