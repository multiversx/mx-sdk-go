package factory

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/endProcess"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	commonFactory "github.com/ElrondNetwork/elrond-go/common/factory"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/disabled"
)

type initialNode struct {
	assignedShard uint32
	eligible      bool
	pubKey        []byte
	address       []byte
	initialRating uint32
}

func NewInitialNode(pubKeyConverter core.PubkeyConverter, pubKey string, address string, rating uint32) (*initialNode, error) {
	pubKeyBytes, err := pubKeyConverter.Decode(pubKey)
	if err != nil {
		return nil, err
	}

	return &initialNode{
		pubKey:        pubKeyBytes,
		address:       []byte(address),
		initialRating: rating,
	}, nil
}

func (in *initialNode) AssignedShard() uint32    { return in.assignedShard }
func (in *initialNode) AddressBytes() []byte     { return in.address }
func (in *initialNode) PubKeyBytes() []byte      { return in.pubKey }
func (in *initialNode) GetInitialRating() uint32 { return in.initialRating }
func (in *initialNode) IsInterfaceNil() bool     { return in == nil }

type validator = nodesCoordinator.Validator

func convertGenesisNodesConfigToValidators(
	genesisNodesConfig *data.GenesisNodesConfig,
) (map[uint32][]nodesCoordinator.GenesisNodeInfoHandler,
	map[uint32][]nodesCoordinator.GenesisNodeInfoHandler) {

	pubkeyConfig := &config.PubkeyConfig{
		Length:          96,
		Type:            "hex",
		SignatureLength: 48,
	}
	validatorPubkeyConverter, _ := commonFactory.NewPubkeyConverter(*pubkeyConfig)

	el := genesisNodesConfig.Eligible
	wt := genesisNodesConfig.Waiting

	eligible := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)
	waiting := make(map[uint32][]nodesCoordinator.GenesisNodeInfoHandler)

	for shardID, nodes := range el {
		for _, node := range nodes {
			nd, err := NewInitialNode(
				validatorPubkeyConverter,
				node.PubKey,
				node.Address,
				node.InitialRating,
			)
			if err != nil {
				continue
			}
			eligible[shardID] = append(eligible[shardID], nd)
		}
	}

	for shardID, nodes := range wt {
		for _, node := range nodes {
			nd, err := NewInitialNode(
				validatorPubkeyConverter,
				node.PubKey,
				node.Address,
				node.InitialRating,
			)
			if err != nil {
				continue
			}
			waiting[shardID] = append(waiting[shardID], nd)
		}
	}

	return eligible, waiting
}

// CreateNodesCoordinator creates nodes coordinator which will be used for header verification
func CreateNodesCoordinator(
	hasher hashing.Hasher,
	marshaller marshal.Marshalizer,
	rater nodesCoordinator.ChanceComputer,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
	publicKey crypto.PublicKey,
	genesisNodesConfig *data.GenesisNodesConfig,
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
