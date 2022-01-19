package factory

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	hasherFactory "github.com/ElrondNetwork/elrond-go-core/hashing/factory"
	"github.com/ElrondNetwork/elrond-go-core/hashing/sha256"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/sharding/mock"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon/nodeTypeProviderMock"
)

type validatorList []nodesCoordinator.Validator

// Len will return the length of the validatorList
func (v validatorList) Len() int { return len(v) }

// Swap will interchange the objects on input indexes
func (v validatorList) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// Less will return true if object on index i should appear before object in index j
// Sorting of validators should be by index and public key
func (v validatorList) Less(i, j int) bool {
	if v[i].Index() == v[j].Index() {
		return bytes.Compare(v[i].PubKey(), v[j].PubKey()) < 0
	}
	return v[i].Index() < v[j].Index()
}

type nodesCoordinatorLiteWithRater struct {
	*nodesCoordinator.IndexHashedNodesCoordinatorLite
	chanceComputer sharding.ChanceComputer
}

func NewNodesCoordinatorLiteWithRater(
	nodesCoordinatorLite *nodesCoordinator.IndexHashedNodesCoordinatorLite,
	rater sharding.ChanceComputer,
) (*nodesCoordinatorLiteWithRater, error) {

	ndL := &nodesCoordinatorLiteWithRater{
		IndexHashedNodesCoordinatorLite: nodesCoordinatorLite,
		chanceComputer:                  rater,
	}

	ndL.SetNodesCoordinatorHelper(ndL)

	// currentEpoch := ndL.GetCurrentEpoch()

	// nodesConfig := nodesConfigPerShard[currentEpoch]
	// nodesConfig.Selectors, _ = nodesCoordinatorLite.CreateSelectors(nodesConfig)

	// ndL.SetNodesConfigPerEpoch(currentEpoch, nodesConfig)

	return ndL, nil
}

func (ndwr *nodesCoordinatorLiteWithRater) SetEpochNodesConfig(epoch uint32, epochNodesConfig *nodesCoordinator.EpochNodesConfig) {
	epochNodesConfig.Selectors, _ = ndwr.CreateSelectors(epochNodesConfig)

	ndwr.SetNodesConfigPerEpoch(epoch, epochNodesConfig)
}

// GetChance returns the chance from an actual rating
func (ndwr *nodesCoordinatorLiteWithRater) GetChance(rating uint32) uint32 {
	return ndwr.chanceComputer.GetChance(rating)
}

// ValidatorsWeights returns the weights/chances for each given validator
func (ndwr *nodesCoordinatorLiteWithRater) ValidatorsWeights(validators []Validator) ([]uint32, error) {
	minChance := ndwr.GetChance(0)
	weights := make([]uint32, len(validators))

	for i, validatorInShard := range validators {
		weights[i] = validatorInShard.Chances()
		if weights[i] < minChance {
			//default weight if all validators need to be selected
			weights[i] = minChance
		}
	}

	return weights, nil
}

func (ndwr *nodesCoordinatorLiteWithRater) computeNodesConfigFromList(
	previousEpochConfig *nodesCoordinator.EpochNodesConfig,
	validatorInfos []*state.ShardValidatorInfo,
) (*nodesCoordinator.EpochNodesConfig, error) {
	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	leavingMap := make(map[uint32][]Validator)
	newNodesList := make([]Validator, 0)

	if len(validatorInfos) == 0 {
		log.Warn("computeNodesConfigFromList - validatorInfos len is 0")
	}

	for _, validatorInfo := range validatorInfos {
		chance := ndwr.GetChance(validatorInfo.TempRating)
		currentValidator, err := nodesCoordinator.NewValidator(validatorInfo.PublicKey, chance, validatorInfo.Index)
		if err != nil {
			return nil, err
		}

		switch validatorInfo.List {
		case string(common.WaitingList):
			waitingMap[validatorInfo.ShardId] = append(waitingMap[validatorInfo.ShardId], currentValidator)
		case string(common.EligibleList):
			eligibleMap[validatorInfo.ShardId] = append(eligibleMap[validatorInfo.ShardId], currentValidator)
		case string(common.LeavingList):
			log.Debug("leaving node validatorInfo", "pk", validatorInfo.PublicKey)
			leavingMap[validatorInfo.ShardId] = append(leavingMap[validatorInfo.ShardId], currentValidator)
			// ihgs.addValidatorToPreviousMap(
			// 	previousEpochConfig,
			// 	eligibleMap,
			// 	waitingMap,
			// 	currentValidator,
			// 	validatorInfo.ShardId)
		case string(common.NewList):
			log.Debug("new node registered", "pk", validatorInfo.PublicKey)
			newNodesList = append(newNodesList, currentValidator)
		case string(common.InactiveList):
			log.Debug("inactive validator", "pk", validatorInfo.PublicKey)
		case string(common.JailedList):
			log.Debug("jailed validator", "pk", validatorInfo.PublicKey)
		}
	}

	sort.Sort(validatorList(newNodesList))
	for _, eligibleList := range eligibleMap {
		sort.Sort(validatorList(eligibleList))
	}
	for _, waitingList := range waitingMap {
		sort.Sort(validatorList(waitingList))
	}
	for _, leavingList := range leavingMap {
		sort.Sort(validatorList(leavingList))
	}

	if len(eligibleMap) == 0 {
		return nil, fmt.Errorf("eligible map size is zero. No validators found")
	}

	nbShards := len(eligibleMap) - 1

	newNodesConfig := &nodesCoordinator.EpochNodesConfig{
		EligibleMap: eligibleMap,
		WaitingMap:  waitingMap,
		LeavingMap:  leavingMap,
		NewList:     newNodesList,
		NbShards:    uint32(nbShards),
	}

	return newNodesConfig, nil
}

type Validator = nodesCoordinator.Validator

func createDummyNodesList(nbNodes uint32, suffix string) []Validator {
	list := make([]Validator, 0)
	hasher := sha256.NewSha256()

	for j := uint32(0); j < nbNodes; j++ {
		pk := hasher.Compute(fmt.Sprintf("pk%s_%d", suffix, j))
		list = append(list, mock.NewValidatorMock(pk, 1, nodesCoordinator.DefaultSelectionChances))
	}

	return list
}

func createDummyNodesMap(nodesPerShard uint32, nbShards uint32, suffix string) map[uint32][]Validator {
	nodesMap := make(map[uint32][]Validator)

	var shard uint32

	for i := uint32(0); i <= nbShards; i++ {
		shard = i
		if i == nbShards {
			shard = core.MetachainShardId
		}
		list := createDummyNodesList(nodesPerShard, suffix+"_i")
		nodesMap[shard] = list
	}

	return nodesMap
}

func CreateNodesCoordinatorLite(
	hasher hashing.Hasher,
	rater sharding.ChanceComputer,
	shardConsensusGroupSize int,
	metaConsensusGroupSize int,
	nShards uint32,
) (nodesCoordinator.NodesCoordinatorLite, error) {

	cache := &mock.NodesCoordinatorCacheMock{
		GetCalled: func(key []byte) (value interface{}, ok bool) {
			return nil, false
		},
		PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
			return false
		},
	}

	waitingMap := make(map[uint32][]nodesCoordinator.Validator)
	eligibleMap := createDummyNodesMap(3, 2, "eligible")

	hasher, err := hasherFactory.NewHasher("blake2b")
	if err != nil {
		return nil, err
	}

	arguments := nodesCoordinator.ArgNodesCoordinatorLite{
		Epoch:                      uint32(0),
		ShardConsensusGroupSize:    shardConsensusGroupSize,
		MetaConsensusGroupSize:     metaConsensusGroupSize,
		Hasher:                     hasher,
		NbShards:                   nShards,
		EligibleNodes:              eligibleMap,
		WaitingNodes:               waitingMap,
		SelfPublicKey:              []byte("key"),
		ConsensusGroupCache:        cache,
		WaitingListFixEnabledEpoch: 1000000,
		ChanStopNode:               make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:           &nodeTypeProviderMock.NodeTypeProviderStub{},
	}

	nd, err := nodesCoordinator.NewIndexHashedNodesCoordinatorLite(arguments)
	if err != nil {
		return nil, err
	}

	ndWithRater, err := NewNodesCoordinatorLiteWithRater(nd, rater)
	if err != nil {
		return nil, err
	}

	return ndWithRater, nil
}
