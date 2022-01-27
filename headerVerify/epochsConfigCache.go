package headerVerify

import (
	"math"
	"sync"

	"github.com/ElrondNetwork/elrond-go/state"
)

type epochsConfigCache struct {
	cache       map[uint32][]*state.ShardValidatorInfo
	cacheMutex  sync.RWMutex
	oldestEpoch uint32
	cacheSize   uint64
}

func NewEpochsConfigCache(maxEpochs uint64) *epochsConfigCache {
	return &epochsConfigCache{
		cache:       make(map[uint32][]*state.ShardValidatorInfo),
		cacheMutex:  sync.RWMutex{},
		oldestEpoch: math.MaxUint32,
		cacheSize:   maxEpochs,
	}
}

func (ec *epochsConfigCache) Add(epoch uint32, validatorsInfo []*state.ShardValidatorInfo) error {
	ec.cacheMutex.Lock()
	defer ec.cacheMutex.Unlock()

	if ec.isCacheFull(epoch) {
		ec.updateOldestEpoch()
	}

	if epoch < ec.oldestEpoch {
		ec.oldestEpoch = epoch
	}

	ec.cache[epoch] = validatorsInfo

	return nil
}

func (ec *epochsConfigCache) IsInCache(epoch uint32) bool {
	_, isEpochInCache := ec.cache[epoch]
	return isEpochInCache
}

func (ec *epochsConfigCache) isCacheFull(epoch uint32) bool {
	isEpochInCache := ec.IsInCache(epoch)
	return len(ec.cache) >= int(ec.cacheSize) && !isEpochInCache
}

func (ec *epochsConfigCache) updateOldestEpoch() {
	delete(ec.cache, ec.oldestEpoch)

	min := uint32(math.MaxUint32)
	for epoch := range ec.cache {
		if epoch < min {
			min = epoch
		}
	}

	ec.oldestEpoch = min
}

// IsInterfaceNil checks if the underlying pointer is nil
func (rhc *epochsConfigCache) IsInterfaceNil() bool {
	return rhc == nil
}
