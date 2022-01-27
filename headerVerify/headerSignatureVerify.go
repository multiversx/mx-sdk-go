package headerVerify

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/process/headerCheck"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerVerify/factory"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/headerVerify")

const maxEpochs = 3

type HeaderSignatureVerifier struct {
	HeaderVerifier *headerCheck.HeaderSigVerifier
	NdLite         *NodesCoordinatorLiteWithRater
	ConfigCache    EpochsConfigCacheHandler
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

	cache := NewEpochsConfigCache(maxEpochs)

	hsv := &HeaderSignatureVerifier{
		HeaderVerifier: headerVerifier,
		NdLite:         ndLite,
		ConfigCache:    cache,
	}

	return hsv, nil
}

func (hsv *HeaderSignatureVerifier) SetNodesConfigPerEpoch(
	validatorsInfo []*state.ShardValidatorInfo,
	epoch uint32,
	randomness []byte,
) {
	hsv.NdLite.SetNodesConfigOnNewEpochStart(epoch, randomness, validatorsInfo)
	return
}

func (hsv *HeaderSignatureVerifier) VerifyHeader(header coreData.HeaderHandler) bool {
	err := hsv.HeaderVerifier.VerifySignature(header)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}
