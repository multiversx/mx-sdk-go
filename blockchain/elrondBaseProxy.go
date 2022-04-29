package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	elrondCore "github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/blockchain")

const (
	minimumCachingInterval = time.Second
	minAllowedDeltaToFinal = 1
)

const (
	// endpoints
	networkConfigEndpoint    = "network/config"
	getNetworkStatusEndpoint = "network/status/%v"
)

type argsElrondBaseProxy struct {
	expirationTime    time.Duration
	httpClientWrapper httpClientWrapper
}

type elrondBaseProxy struct {
	httpClientWrapper
	mut                 sync.RWMutex
	fetchedConfigs      *data.NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
}

// newElrondBaseProxy will create a base elrond proxy with cache instance
func newElrondBaseProxy(args argsElrondBaseProxy) (*elrondBaseProxy, error) {
	err := checkArgsBaseProxy(args)
	if err != nil {
		return nil, err
	}

	return &elrondBaseProxy{
		httpClientWrapper:   args.httpClientWrapper,
		cacheExpiryDuration: args.expirationTime,
		sinceTimeHandler:    since,
	}, nil
}

func checkArgsBaseProxy(args argsElrondBaseProxy) error {
	if args.expirationTime < minimumCachingInterval {
		return fmt.Errorf("%w, provided: %v, minimum: %v", ErrInvalidCacherDuration, args.expirationTime, minimumCachingInterval)
	}
	if check.IfNil(args.httpClientWrapper) {
		return ErrNilHTTPClientWrapper
	}

	return nil
}

func since(t time.Time) time.Duration {
	return time.Since(t)
}

// GetNetworkConfig will return the cached network configs fetching new values and saving them if necessary
func (proxy *elrondBaseProxy) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.RLock()
	cachedConfigs := proxy.getCachedConfigs()
	proxy.mut.RUnlock()

	if cachedConfigs != nil {
		return cachedConfigs, nil
	}

	return proxy.cacheConfigs(ctx)
}

func (proxy *elrondBaseProxy) getCachedConfigs() *data.NetworkConfig {
	if proxy.sinceTimeHandler(proxy.lastFetchedTime) > proxy.cacheExpiryDuration {
		return nil
	}

	return proxy.fetchedConfigs
}

func (proxy *elrondBaseProxy) cacheConfigs(ctx context.Context) (*data.NetworkConfig, error) {
	proxy.mut.Lock()
	defer proxy.mut.Unlock()

	// maybe another parallel running go routine already did the fetching
	cachedConfig := proxy.getCachedConfigs()
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	log.Debug("Network config not cached. caching...")
	configs, err := proxy.getNetworkConfigFromSource(ctx)
	if err != nil {
		return nil, err
	}

	proxy.lastFetchedTime = time.Now()
	proxy.fetchedConfigs = configs

	return configs, nil
}

// getNetworkConfigFromSource retrieves the network configuration from the proxy
func (proxy *elrondBaseProxy) getNetworkConfigFromSource(ctx context.Context) (*data.NetworkConfig, error) {
	buff, err := proxy.GetHTTP(ctx, networkConfigEndpoint)
	if err != nil {
		return nil, err
	}

	response := &data.NetworkConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// CheckShardFinalization will query the proxy and check if the target shard ID has a current nonce close to the cross
// check nonce from the metachain
// nonce(target shard ID) <= nonce(target shard ID notarized by meta) + maxNoncesDelta
func (proxy *elrondBaseProxy) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta uint64) error {
	if maxNoncesDelta < minAllowedDeltaToFinal {
		return fmt.Errorf("%w, provided: %d, minimum: %d", ErrInvalidAllowedDeltaToFinal, maxNoncesDelta, minAllowedDeltaToFinal)
	}
	if targetShardID == elrondCore.MetachainShardId {
		// we consider this final since the minAllowedDeltaToFinal is 1
		return nil
	}

	nonceFromMeta, nonceFromShard, err := proxy.getNoncesFromMetaAndShard(ctx, targetShardID)
	if err != nil {
		return err
	}

	if nonceFromShard < nonceFromMeta {
		return fmt.Errorf("shardID %d is syncing, meta cross check nonce is %d, current nonce is %d, max delta: %d",
			targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
	}
	if nonceFromShard <= nonceFromMeta+maxNoncesDelta {
		return nil
	}

	return fmt.Errorf("shardID %d is stuck, meta cross check nonce is %d, current nonce is %d, max delta: %d",
		targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
}

func (proxy *elrondBaseProxy) getNoncesFromMetaAndShard(ctx context.Context, targetShardID uint32) (uint64, uint64, error) {
	networkStatusMeta, err := proxy.GetNetworkStatus(ctx, elrondCore.MetachainShardId)
	if err != nil {
		return 0, 0, err
	}

	crossCheckValue := networkStatusMeta.CrossCheckBlockHeight
	nonceFromMeta, err := extractNonceOfShardID(crossCheckValue, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	networkStatusShard, err := proxy.GetNetworkStatus(ctx, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	nonceFromShard := networkStatusShard.Nonce

	return nonceFromMeta, nonceFromShard, nil
}

// GetShardOfAddress returns the shard ID of a provided address by using a shardCoordinator object and querying the
// network config route
func (proxy *elrondBaseProxy) GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error) {
	addr, err := data.NewAddressFromBech32String(bech32Address)
	if err != nil {
		return 0, err
	}

	networkConfigs, err := proxy.GetNetworkConfig(ctx)
	if err != nil {
		return 0, err
	}

	shardCoordinatorInstance, err := NewShardCoordinator(networkConfigs.NumShardsWithoutMeta, 0)
	if err != nil {
		return 0, err
	}

	return shardCoordinatorInstance.ComputeShardId(addr)
}

// GetNetworkStatus will return the network status of a provided shard
func (proxy *elrondBaseProxy) GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
	endpoint := fmt.Sprintf(getNetworkStatusEndpoint, shardID)
	buff, err := proxy.GetHTTP(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	response := &data.NetworkStatusResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Status, nil
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
func (proxy *elrondBaseProxy) IsInterfaceNil() bool {
	return proxy == nil
}
