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
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
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
		pk := hasher.Compute(fmt.Sprintf("pkeligible_%d", j))
		list = append(list, mock.NewValidatorMock(pk, 1, nodesCoordinator.DefaultSelectionChances))
	}

	return list
}

func createDummyNodesMap(nodesPerShard uint32, nbShards uint32) map[uint32][]validator {
	nodesMap := make(map[uint32][]validator)

	var shard uint32

	for i := uint32(0); i <= nbShards; i++ {
		shard = i
		if i == nbShards {
			shard = core.MetachainShardId
		}
		list := createDummyNodesList(nodesPerShard, "_i")
		//list := make([]validator, nodesPerShard)
		nodesMap[shard] = list
	}

	return nodesMap
}

func createArgsNodesShuffler(
	eec *data.EnableEpochsConfig,
	networkConfig *data.NetworkConfig,
) *nodesCoordinator.NodesShufflerArgs {
	maxNodesChangeConfigs := make([]config.MaxNodesChangeConfig, 0)
	for _, conf := range eec.MaxNodesChangeEnableEpoch {
		maxNodesChangeConfig := config.MaxNodesChangeConfig{
			EpochEnable:            conf.EpochEnable,
			MaxNumNodes:            conf.MaxNumNodes,
			NodesToShufflePerShard: conf.NodesToShufflePerShard,
		}

		maxNodesChangeConfigs = append(maxNodesChangeConfigs, maxNodesChangeConfig)
	}

	argsNodesShuffler := &nodesCoordinator.NodesShufflerArgs{
		NodesShard:                     networkConfig.NumNodesInShard,
		NodesMeta:                      networkConfig.NumMetachainNodes,
		Hysteresis:                     networkConfig.GetHysteresis(),
		Adaptivity:                     networkConfig.GetAdaptivity(),
		ShuffleBetweenShards:           true,
		MaxNodesEnableConfig:           maxNodesChangeConfigs,
		BalanceWaitingListsEnableEpoch: eec.BalanceWaitingListsEnableEpoch,
		WaitingListFixEnableEpoch:      eec.WaitingListFixEnableEpoch,
	}

	return argsNodesShuffler
}

func CreateNodesCoordinatorLite(
	hasher hashing.Hasher,
	rater sharding.ChanceComputer,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
) (*NodesCoordinatorLiteWithRater, error) {

	waitingMap := make(map[uint32][]validator)
	eligibleMap := createDummyNodesMap(networkConfig.MetaConsensusGroup, networkConfig.NumShardsWithoutMeta)

	argsNodesShuffler := createArgsNodesShuffler(enableEpochsConfig, networkConfig)
	nodeShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)
	initialEpoch := uint32(0)

	arguments := nodesCoordinator.ArgNodesCoordinatorLite{
		Epoch:                      initialEpoch,
		ShardConsensusGroupSize:    int(networkConfig.ShardConsensusGroupSize),
		MetaConsensusGroupSize:     int(networkConfig.MetaConsensusGroup),
		Hasher:                     hasher,
		NbShards:                   networkConfig.NumShardsWithoutMeta,
		EligibleNodes:              eligibleMap,
		WaitingNodes:               waitingMap,
		SelfPublicKey:              []byte("dummy"),
		ConsensusGroupCache:        &mock.NodesCoordinatorCacheMock{},
		WaitingListFixEnabledEpoch: enableEpochsConfig.WaitingListFixEnableEpoch,
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
