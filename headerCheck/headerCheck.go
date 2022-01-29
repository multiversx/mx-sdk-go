package headerCheck

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/process/headerCheck"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerCheck/factory"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/headerCheck")

type ArgHeaderVerifier struct {
	RatingsConfig      *data.RatingsConfig
	NetworkConfig      *data.NetworkConfig
	EnableEpochsConfig *data.EnableEpochsConfig
}

type headerVerifier struct {
	headerSigVerifier *headerCheck.HeaderSigVerifier
	ndLite            nodesCoordinator.EpochsConfigUpdateHandler
	marshaller        marshal.Marshalizer
}

func NewHeaderVerifier(args ArgHeaderVerifier) (*headerVerifier, error) {
	err := checkArguments(args)
	if err != nil {
		return nil, err
	}

	coreComp, err := factory.CreateCoreComponents(args.RatingsConfig, args.NetworkConfig)
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
		args.NetworkConfig,
		args.EnableEpochsConfig,
	)
	if err != nil {
		return nil, err
	}

	headerSigArgs := &headerCheck.ArgsHeaderSigVerifier{
		Marshalizer:             coreComp.Marshalizer,
		Hasher:                  coreComp.Hasher,
		NodesCoordinator:        ndLite,
		MultiSigVerifier:        cryptoComp.MultiSigVerifier,
		SingleSigVerifier:       cryptoComp.SingleSigVerifier,
		KeyGen:                  cryptoComp.KeyGen,
		FallbackHeaderValidator: &testscommon.FallBackHeaderValidatorStub{},
	}

	headerSigVerifier, err := headerCheck.NewHeaderSigVerifier(headerSigArgs)
	if err != nil {
		return nil, err
	}

	hsv := &headerVerifier{
		headerSigVerifier: headerSigVerifier,
		ndLite:            ndLite,
		marshaller:        coreComp.Marshalizer,
	}

	return hsv, nil
}

func checkArguments(args ArgHeaderVerifier) error {
	if args.NetworkConfig == nil {
		return ErrNilNetworkConfig
	}
	if args.RatingsConfig == nil {
		return ErrNilRatingsConfig
	}
	if args.EnableEpochsConfig == nil {
		return ErrNilEnableEpochsConfig
	}

	return nil
}

func (hsv *headerVerifier) IsInCache(epoch uint32) bool {
	status := hsv.ndLite.IsEpochInConfig(epoch)
	return status
}

func (hsv *headerVerifier) SetNodesConfigPerEpoch(
	validatorsInfo []*state.ShardValidatorInfo,
	epoch uint32,
	randomness []byte,
) error {
	err := hsv.ndLite.SetNodesConfigFromValidatorsInfo(epoch, randomness, validatorsInfo)
	return err
}

func (hsv *headerVerifier) VerifyHeader(header coreData.HeaderHandler) bool {
	err := hsv.headerSigVerifier.VerifySignature(header)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}

func (hsv *headerVerifier) Marshaller() marshal.Marshalizer {
	return hsv.marshaller
}
