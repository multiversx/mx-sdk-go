package finalityProvider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	mxChainCore "github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/core"
)

var log = logger.GetOrCreate("mx-sdk-go/blockchain/finalityprovider")

type proxyFinalityProvider struct {
	proxy proxy
}

// NewProxyFinalityProvider creates a new instance of type proxyFinalityProvider
func NewProxyFinalityProvider(proxy proxy) (*proxyFinalityProvider, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}

	return &proxyFinalityProvider{
		proxy: proxy,
	}, nil
}

// CheckShardFinalization will query the proxy and check if the target shard ID has a current nonce close to the cross
// check nonce from the metachain
// nonce(target shard ID) <= nonce(target shard ID notarized by meta) + maxNoncesDelta
func (provider *proxyFinalityProvider) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
	if maxNoncesDelta < core.MinAllowedDeltaToFinal {
		return fmt.Errorf("%w, provided: %d, minimum: %d", ErrInvalidAllowedDeltaToFinal, maxNoncesDelta, core.MinAllowedDeltaToFinal)
	}
	if targetShardID == mxChainCore.MetachainShardId {
		// we consider this final since the minAllowedDeltaToFinal is 1
		return nil
	}

	nonceFromMeta, nonceFromShard, err := provider.getNoncesFromMetaAndShard(ctx, targetShardID)
	if err != nil {
		return err
	}

	if nonceFromShard < nonceFromMeta {
		return fmt.Errorf("shardID %d is syncing, meta cross check nonce is %d, current nonce is %d, max delta: %d",
			targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
	}
	if nonceFromShard <= nonceFromMeta+maxNoncesDelta {
		log.Trace("proxyFinalityProvider.CheckShardFinalization - shard is in sync",
			"shardID", targetShardID, "meta cross check nonce", nonceFromMeta,
			"current nonce", nonceFromShard, "max delta", maxNoncesDelta)
		return nil
	}

	return fmt.Errorf("shardID %d is stuck, meta cross check nonce is %d, current nonce is %d, max delta: %d",
		targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
}

func (provider *proxyFinalityProvider) getNoncesFromMetaAndShard(ctx context.Context, targetShardID uint32) (uint64, uint64, error) {
	networkStatusMeta, err := provider.proxy.GetNetworkStatus(ctx, mxChainCore.MetachainShardId)
	if err != nil {
		return 0, 0, err
	}

	crossCheckValue := networkStatusMeta.CrossCheckBlockHeight
	nonceFromMeta, err := extractNonceOfShardID(crossCheckValue, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	networkStatusShard, err := provider.proxy.GetNetworkStatus(ctx, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	nonceFromShard := networkStatusShard.Nonce

	return nonceFromMeta, nonceFromShard, nil
}

func extractNonceOfShardID(crossCheckValue string, shardID uint32) (uint64, error) {
	// the value will come in this format: "0: 9169897, 1: 9166353, 2: 9170524, "
	if len(crossCheckValue) == 0 {
		return 0, fmt.Errorf("%w: empty value, maybe bad observer version", ErrInvalidNonceCrossCheckValueFormat)
	}
	shardsData := strings.Split(crossCheckValue, ",")
	shardIdAsString := fmt.Sprintf("%d", shardID)

	for _, shardData := range shardsData {
		shardNonce := strings.Split(shardData, ":")
		if len(shardNonce) != 2 {
			continue
		}

		shardNonce[0] = strings.TrimSpace(shardNonce[0])
		shardNonce[1] = strings.TrimSpace(shardNonce[1])
		if shardNonce[0] != shardIdAsString {
			continue
		}

		val, ok := big.NewInt(0).SetString(shardNonce[1], 10)
		if !ok {
			return 0, fmt.Errorf("%w: %s is not a valid number as found in this response: %s",
				ErrInvalidNonceCrossCheckValueFormat, shardNonce[1], crossCheckValue)
		}

		return val.Uint64(), nil
	}

	return 0, fmt.Errorf("%w: value not found for shard %d from this response: %s",
		ErrInvalidNonceCrossCheckValueFormat, shardID, crossCheckValue)
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *proxyFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}
