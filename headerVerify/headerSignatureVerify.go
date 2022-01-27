package headerVerify

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/process/headerCheck"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerVerify/factory"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/headerVerify")

type HeaderSignatureVerifier struct {
	headerVerifier *headerCheck.HeaderSigVerifier
	ndLite         *NodesCoordinatorLiteWithRater
	marshaller     marshal.Marshalizer
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
		networkConfig.ShardConsensusGroupSize,
		networkConfig.MetaConsensusGroup,
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
		marshaller:     coreComp.Marshalizer,
	}

	return hsv, nil
}

func (hsv *HeaderSignatureVerifier) IsInCache(epoch uint32) bool {
	status := hsv.ndLite.IsEpochInConfig(epoch)
	return status
}

func (hsv *HeaderSignatureVerifier) SetNodesConfigPerEpoch(
	validatorsInfo []*state.ShardValidatorInfo,
	epoch uint32,
	randomness []byte,
) error {
	err := hsv.ndLite.SetNodesConfigFromValidatorsInfo(epoch, randomness, validatorsInfo)
	return err
}

func (hsv *HeaderSignatureVerifier) VerifyHeader(header coreData.HeaderHandler) bool {
	err := hsv.headerVerifier.VerifySignature(header)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}

func (hsv *HeaderSignatureVerifier) Marshaller() marshal.Marshalizer {
	return hsv.marshaller
}
