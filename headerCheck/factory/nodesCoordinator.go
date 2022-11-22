package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/common/enablers"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/disabled"
)

const (
	defaultSelectionChances = uint32(1)
)

// CreateNodesCoordinator creates nodes coordinator which will be used for header verification
func CreateNodesCoordinator(
	coreComp *coreComponents,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
	publicKey crypto.PublicKey,
	genesisNodesConfig *data.GenesisNodes,
) (nodesCoordinator.EpochsConfigUpdateHandler, error) {
	eligibleValidators, err := generateGenesisNodes(coreComp.PubKeyConverter, genesisNodesConfig.Eligible)
	if err != nil {
		return nil, err
	}

	waitingValidators, err := generateGenesisNodes(coreComp.PubKeyConverter, genesisNodesConfig.Waiting)
	if err != nil {
		return nil, err
	}

	argsNodesShuffler := createArgsNodesShuffler(enableEpochsConfig, networkConfig)
	nodeShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	initialEpoch := uint32(0)
	arguments := nodesCoordinator.ArgNodesCoordinator{
		Epoch:                   initialEpoch,
		ShardConsensusGroupSize: int(networkConfig.ShardConsensusGroupSize),
		MetaConsensusGroupSize:  int(networkConfig.MetaConsensusGroup),
		Marshalizer:             coreComp.Marshaller,
		EpochStartNotifier:      &disabled.EpochStartNotifier{},
		BootStorer:              &disabled.Storer{},
		Hasher:                  coreComp.Hasher,
		NbShards:                networkConfig.NumShardsWithoutMeta,
		EligibleNodes:           eligibleValidators,
		WaitingNodes:            waitingValidators,
		SelfPublicKey:           publicKeyBytes,
		ConsensusGroupCache:     &disabled.NodesCoordinatorCache{},
		ChanStopNode:            make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:        &disabled.NodeTypeProvider{},
		Shuffler:                nodeShuffler,
		ShuffledOutHandler:      &disabled.ShuffledOutHandler{},
	}

	baseNodesCoordinator, err := nodesCoordinator.NewIndexHashedNodesCoordinator(arguments)
	if err != nil {
		return nil, err
	}

	nd, err := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseNodesCoordinator, coreComp.Rater)
	if err != nil {
		return nil, err
	}

	return nd, nil
}

func generateGenesisNodes(converter core.PubkeyConverter, nodesConfig map[uint32][]string) (map[uint32][]nodesCoordinator.Validator, error) {
	validatorsMap := make(map[uint32][]nodesCoordinator.Validator)

	for shardID, nodesPubKeys := range nodesConfig {
		validators := make([]nodesCoordinator.Validator, 0, len(nodesPubKeys))
		for i, pubKey := range nodesPubKeys {
			pubKeyBytes, err := converter.Decode(pubKey)
			if err != nil {
				return nil, err
			}

			validator, err := nodesCoordinator.NewValidator(pubKeyBytes, defaultSelectionChances, uint32(i))
			if err != nil {
				return nil, err
			}

			validators = append(validators, validator)
		}
		validatorsMap[shardID] = validators
	}

	return validatorsMap, nil
}

func createArgsNodesShuffler(
	eec *data.EnableEpochsConfig,
	networkConfig *data.NetworkConfig,
) *nodesCoordinator.NodesShufflerArgs {
	maxNodesChangeConfigs := make([]config.MaxNodesChangeConfig, 0)
	for _, conf := range eec.EnableEpochs.MaxNodesChangeEnableEpoch {
		maxNodesChangeConfig := config.MaxNodesChangeConfig{
			EpochEnable:            conf.EpochEnable,
			MaxNumNodes:            conf.MaxNumNodes,
			NodesToShufflePerShard: conf.NodesToShufflePerShard,
		}

		maxNodesChangeConfigs = append(maxNodesChangeConfigs, maxNodesChangeConfig)
	}

	enableEpochsHandler, err := enablers.NewEnableEpochsHandler(eec.EnableEpochs, &disabled.EpochNotifier{})
	if err != nil {
		return nil
	}

	argsNodesShuffler := &nodesCoordinator.NodesShufflerArgs{
		NodesShard:           networkConfig.NumNodesInShard,
		NodesMeta:            networkConfig.NumMetachainNodes,
		Hysteresis:           networkConfig.Hysteresys,
		Adaptivity:           networkConfig.Adaptivity,
		ShuffleBetweenShards: true,
		MaxNodesEnableConfig: maxNodesChangeConfigs,
		EnableEpochsHandler:  enableEpochsHandler,
	}

	return argsNodesShuffler
}
