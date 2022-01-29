package headerCheck

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/hashing/sha256"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/sharding/mock"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/testscommon/nodeTypeProviderMock"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type validator = nodesCoordinator.Validator

func CreateNodesCoordinatorLite(
	hasher hashing.Hasher,
	rater nodesCoordinator.ChanceComputer,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
) (nodesCoordinator.EpochsConfigUpdateHandler, error) {

	waitingMap := make(map[uint32][]validator)
	eligibleMap := createDummyNodesMap(networkConfig.MetaConsensusGroup, networkConfig.NumShardsWithoutMeta)

	argsNodesShuffler := createArgsNodesShuffler(enableEpochsConfig, networkConfig)
	nodeShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)
	if err != nil {
		return nil, err
	}

	initialEpoch := uint32(0)
	dummySelfPublicKey := []byte("dummy")
	arguments := nodesCoordinator.ArgNodesCoordinator{
		Epoch:                      initialEpoch,
		ShardConsensusGroupSize:    int(networkConfig.ShardConsensusGroupSize),
		MetaConsensusGroupSize:     int(networkConfig.MetaConsensusGroup),
		Marshalizer:                &mock.MarshalizerMock{},
		EpochStartNotifier:         &mock.EpochStartNotifierStub{},
		BootStorer:                 mock.NewStorerMock(),
		Hasher:                     hasher,
		NbShards:                   networkConfig.NumShardsWithoutMeta,
		EligibleNodes:              eligibleMap,
		WaitingNodes:               waitingMap,
		SelfPublicKey:              dummySelfPublicKey,
		ConsensusGroupCache:        &mock.NodesCoordinatorCacheMock{},
		WaitingListFixEnabledEpoch: enableEpochsConfig.WaitingListFixEnableEpoch,
		ChanStopNode:               make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:           &nodeTypeProviderMock.NodeTypeProviderStub{},
		Shuffler:                   nodeShuffler,
		ShuffledOutHandler:         &mock.ShuffledOutHandlerStub{},
	}

	baseNodesCoordinator, err := nodesCoordinator.NewIndexHashedNodesCoordinator(arguments)
	if err != nil {
		return nil, err
	}

	nd, err := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseNodesCoordinator, rater)
	if err != nil {
		return nil, err
	}

	return nd, nil
}

func createDummyNodesList(nbNodes uint32, suffix string) []validator {
	list := make([]validator, 0)
	hasher := sha256.NewSha256()

	for j := uint32(0); j < nbNodes; j++ {
		pk := hasher.Compute(fmt.Sprintf("pkeligible_%d", j))
		list = append(list, mock.NewValidatorMock(pk, 1, 1))
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
