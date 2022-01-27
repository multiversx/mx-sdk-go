package headerVerify

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/hashing/sha256"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/sharding/mock"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon/nodeTypeProviderMock"
)

type validator = nodesCoordinator.Validator

type NodesCoordinatorLiteWithRater struct {
	*nodesCoordinator.IndexHashedNodesCoordinatorLite
	chanceComputer sharding.ChanceComputer
}

func NewNodesCoordinatorLiteWithRater(
	nodesCoordinatorLite *nodesCoordinator.IndexHashedNodesCoordinatorLite,
	rater sharding.ChanceComputer,
) (*NodesCoordinatorLiteWithRater, error) {

	ndL := &NodesCoordinatorLiteWithRater{
		IndexHashedNodesCoordinatorLite: nodesCoordinatorLite,
		chanceComputer:                  rater,
	}

	ndL.SetNodesCoordinatorHelper(ndL)

	return ndL, nil
}

// GetChance returns the chance from an actual rating
func (ndwr *NodesCoordinatorLiteWithRater) GetChance(rating uint32) uint32 {
	return ndwr.chanceComputer.GetChance(rating)
}

// ValidatorsWeights returns the weights/chances for each given validator
func (ndwr *NodesCoordinatorLiteWithRater) ValidatorsWeights(validators []validator) ([]uint32, error) {
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

func (ndwr *NodesCoordinatorLiteWithRater) SetNodesConfigPerEpoch(
	validatorsInfo []*state.ShardValidatorInfo,
	epoch uint32,
	randomness []byte,
) error {
	err := ndwr.SetNodesConfigFromValidatorsInfo(epoch, randomness, validatorsInfo)
	return err
}

func createDummyNodesList(nbNodes uint32, suffix string) []validator {
	list := make([]validator, 0)
	hasher := sha256.NewSha256()

	for j := uint32(0); j < nbNodes; j++ {
		pk := hasher.Compute(fmt.Sprintf("pk%s_%d", suffix, j))
		list = append(list, mock.NewValidatorMock(pk, 1, nodesCoordinator.DefaultSelectionChances))
	}

	return list
}

func createDummyNodesMap(nodesPerShard uint32, nbShards uint32, suffix string) map[uint32][]validator {
	nodesMap := make(map[uint32][]validator)

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

func createArgsNodesShuffler() *nodesCoordinator.NodesShufflerArgs {
	maxNodesChangeConfigs := make([]config.MaxNodesChangeConfig, 0)
	maxNodesChangeConfig1 := config.MaxNodesChangeConfig{
		EpochEnable:            0,
		MaxNumNodes:            36,
		NodesToShufflePerShard: 4,
	}
	maxNodesChangeConfigs = append(maxNodesChangeConfigs, maxNodesChangeConfig1)
	maxNodesChangeConfig2 := config.MaxNodesChangeConfig{
		EpochEnable:            1,
		MaxNumNodes:            56,
		NodesToShufflePerShard: 2,
	}
	maxNodesChangeConfigs = append(maxNodesChangeConfigs, maxNodesChangeConfig2)

	argsNodesShuffler := &nodesCoordinator.NodesShufflerArgs{
		NodesShard:                     3,
		NodesMeta:                      3,
		Hysteresis:                     0,
		Adaptivity:                     false,
		ShuffleBetweenShards:           true,
		MaxNodesEnableConfig:           maxNodesChangeConfigs,
		BalanceWaitingListsEnableEpoch: 1,
		WaitingListFixEnableEpoch:      1000000,
	}

	return argsNodesShuffler
}

func CreateNodesCoordinatorLite(
	hasher hashing.Hasher,
	rater sharding.ChanceComputer,
	shardConsensusGroupSize uint64,
	metaConsensusGroupSize uint64,
	nShards uint32,
) (*NodesCoordinatorLiteWithRater, error) {

	waitingMap := make(map[uint32][]validator)
	eligibleMap := createDummyNodesMap(uint32(metaConsensusGroupSize), nShards, "eligible")

	argsNodesShuffler := createArgsNodesShuffler()
	nodeShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)

	arguments := nodesCoordinator.ArgNodesCoordinatorLite{
		Epoch:                      uint32(0),
		ShardConsensusGroupSize:    int(shardConsensusGroupSize),
		MetaConsensusGroupSize:     int(metaConsensusGroupSize),
		Hasher:                     hasher,
		NbShards:                   nShards,
		EligibleNodes:              eligibleMap,
		WaitingNodes:               waitingMap,
		SelfPublicKey:              []byte("key"),
		ConsensusGroupCache:        &mock.NodesCoordinatorCacheMock{},
		WaitingListFixEnabledEpoch: 1000000,
		ChanStopNode:               make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:           &nodeTypeProviderMock.NodeTypeProviderStub{},
		Shuffler:                   nodeShuffler,
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
