package headerVerify

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/process/headerCheck"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerVerify/factory"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/headerVerify")

type HeaderSignatureVerifier struct {
	headerVerifier *headerCheck.HeaderSigVerifier
	ndLite         *NodesCoordinatorLiteWithRater
}

func NewHeaderSignatureVerifier(
	ratingsConfig *data.RatingsConfig,
	networkConfig *data.NetworkConfig,
) (*HeaderSignatureVerifier, error) {
	coreComp, err := factory.CreateCoreComponents(ratingsConfig, networkConfig)
	if err != nil {
		return nil, err
	}

	cryptoComp, err := factory.CreateCryptoComponents()
	if err != nil {
		return nil, err
	}

	ndLite, err := CreateNodesCoordinatorLite(
		coreComp.Hasher,
		coreComp.Rater,
		int(networkConfig.ShardConsensusGroupSize),
		int(networkConfig.MetaConsensusGroup),
		networkConfig.NumShardsWithoutMeta,
	)
	if err != nil {
		return nil, err
	}

	args := &headerCheck.ArgsHeaderSigVerifier{
		Marshalizer:             coreComp.Marshalizer,
		Hasher:                  coreComp.Hasher,
		NodesCoordinator:        ndLite,
		MultiSigVerifier:        cryptoComp.MultiSigVerifier,
		SingleSigVerifier:       cryptoComp.SingleSigVerifier,
		KeyGen:                  cryptoComp.KeyGen,
		FallbackHeaderValidator: &testscommon.FallBackHeaderValidatorStub{},
	}

	headerVerifier, err := headerCheck.NewHeaderSigVerifier(args)
	if err != nil {
		return nil, err
	}

	hsv := &HeaderSignatureVerifier{
		headerVerifier: headerVerifier,
		ndLite:         ndLite,
	}

	return hsv, nil
}

func (hsv *HeaderSignatureVerifier) SetNodesConfigPerEpoch(
	validatorsInfoPerEpoch map[uint32][]*state.ShardValidatorInfo,
) error {

	previousEpochConfig := &nodesCoordinator.EpochNodesConfig{}

	for epoch, validatorsInfo := range validatorsInfoPerEpoch {
		epochNodesConfig, err := hsv.ndLite.ComputeNodesConfigFromList(previousEpochConfig, validatorsInfo)
		if err != nil {
			return err
		}

		// for shard, list := range epochNodesConfig.EligibleMap {
		// 	for j, validators := range list {
		// 		fmt.Println("shard:", shard, "validator", j, ": ", hex.EncodeToString(validators.PubKey()))
		// 		fmt.Println("shard:", shard, "validator", j, ": ", validators.PubKey())
		// 	}
		// }

		hsv.ndLite.SetEpochNodesConfig(epoch, epochNodesConfig)
	}

	return nil
}

func (hsv *HeaderSignatureVerifier) VerifyHeader(header coreData.HeaderHandler) bool {
	err := hsv.headerVerifier.VerifySignature(header)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}
