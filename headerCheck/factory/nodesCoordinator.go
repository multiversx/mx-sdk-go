package factory

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/disabled"
)

const defaultSelectionChances = uint32(1)

type initialNode struct {
	assignedShard uint32
	eligible      bool
	pubKey        []byte
	address       []byte
	initialRating uint32
}

func NewInitialNode(pubKey []byte) *initialNode {
	return &initialNode{
		pubKey:  pubKey,
		address: []byte{},
	}
}

func (in *initialNode) AssignedShard() uint32    { return in.assignedShard }
func (in *initialNode) AddressBytes() []byte     { return in.address }
func (in *initialNode) PubKeyBytes() []byte      { return in.pubKey }
func (in *initialNode) GetInitialRating() uint32 { return in.initialRating }
func (in *initialNode) IsInterfaceNil() bool     { return in == nil }

type validator = nodesCoordinator.Validator

func convertGenesisNodesConfigToValidators(
	genesisNodesConfig *data.GenesisNodes,
) (map[uint32][]nodesCoordinator.GenesisNodeInfoHandler,
	map[uint32][]nodesCoordinator.GenesisNodeInfoHandler) {

	el := genesisNodesConfig.Eligible
	wt := genesisNodesConfig.Waiting

	eligible := generateGenesisNodes(el)
	waiting := generateGenesisNodes(wt)

	return eligible, waiting
}

func generateGenesisNodes(nodesConfig map[uint32][][]byte) map[uint32][]nodesCoordinator.GenesisNodeInfoHandler {
	nodes := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)

	for shardID, nodesPubKeys := range nodesConfig {
		for _, pubKey := range nodesPubKeys {
			nd := NewInitialNode(pubKey)
			nodes[shardID] = append(nodes[shardID], nd)
		}
	}

	return nodes
}

func generateGenesisNodesV2(nodesConfig map[uint32][][]byte) (map[uint32][]validator, error) {
	validatorsMap := make(map[uint32][]validator)

	for shardID, nodesPubKeys := range nodesConfig {
		validators := make([]validator, 0, len(nodesPubKeys))
		for i, pubKey := range nodesPubKeys {
			validatorObj, err := nodesCoordinator.NewValidator(pubKey, defaultSelectionChances, uint32(i))
			if err != nil {
				return nil, err
			}

			validators = append(validators, validatorObj)
		}
		validatorsMap[shardID] = validators
	}

	return validatorsMap, nil
}

// CreateNodesCoordinator creates nodes coordinator which will be used for header verification
func CreateNodesCoordinator(
	hasher hashing.Hasher,
	marshaller marshal.Marshalizer,
	rater nodesCoordinator.ChanceComputer,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
	publicKey crypto.PublicKey,
	genesisNodesConfig *data.GenesisNodes,
) (nodesCoordinator.EpochsConfigUpdateHandler, error) {

	fmt.Println(genesisNodesConfig)

	// TODO: manage epoch 0 from real nodes config
	// waitingMap := make(map[uint32][]validator)
	// eligibleMap := createDummyNodesMap(networkConfig.MetaConsensusGroup, networkConfig.NumShardsWithoutMeta, hasher)

	eligible, waiting := convertGenesisNodesConfigToValidators(genesisNodesConfig)

	eligibleValidators, err := nodesCoordinator.NodesInfoToValidators(eligible)
	if err != nil {
		return nil, err
	}

	waitingValidators, err := nodesCoordinator.NodesInfoToValidators(waiting)
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
		Epoch:                      initialEpoch,
		ShardConsensusGroupSize:    int(networkConfig.ShardConsensusGroupSize),
		MetaConsensusGroupSize:     int(networkConfig.MetaConsensusGroup),
		Marshalizer:                marshaller,
		EpochStartNotifier:         &disabled.EpochStartNotifier{},
		BootStorer:                 &disabled.Storer{},
		Hasher:                     hasher,
		NbShards:                   networkConfig.NumShardsWithoutMeta,
		EligibleNodes:              eligibleValidators,
		WaitingNodes:               waitingValidators,
		SelfPublicKey:              publicKeyBytes,
		ConsensusGroupCache:        &disabled.NodesCoordinatorCache{},
		WaitingListFixEnabledEpoch: enableEpochsConfig.WaitingListFixEnableEpoch,
		ChanStopNode:               make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:           &disabled.NodeTypeProvider{},
		Shuffler:                   nodeShuffler,
		ShuffledOutHandler:         &disabled.ShuffledOutHandler{},
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

func createDummyNodesList(nbNodes uint32, suffix string, hasher hashing.Hasher) []validator {
	list := make([]validator, 0)

	for j := uint32(0); j < nbNodes; j++ {
		pk := hasher.Compute(fmt.Sprintf("pkeligible_%d", j))
		val, _ := nodesCoordinator.NewValidator(pk, 1, 1)
		list = append(list, val)
	}

	return list
}

func createDummyNodesMap(nodesPerShard uint32, nbShards uint32, hasher hashing.Hasher) map[uint32][]validator {
	nodesMap := make(map[uint32][]validator)

	var shard uint32

	for i := uint32(0); i <= nbShards; i++ {
		shard = i
		if i == nbShards {
			shard = core.MetachainShardId
		}
		list := createDummyNodesList(nodesPerShard, "_i", hasher)
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
		Hysteresis:                     networkConfig.Hysteresys,
		Adaptivity:                     networkConfig.Adaptivity,
		ShuffleBetweenShards:           true,
		MaxNodesEnableConfig:           maxNodesChangeConfigs,
		BalanceWaitingListsEnableEpoch: eec.BalanceWaitingListsEnableEpoch,
		WaitingListFixEnableEpoch:      eec.WaitingListFixEnableEpoch,
	}

	return argsNodesShuffler
}
