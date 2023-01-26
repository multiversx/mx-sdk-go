package headerCheck

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/factory/crypto"
	"github.com/multiversx/mx-chain-go/process/headerCheck"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/disabled"
	"github.com/multiversx/mx-sdk-go/headerCheck/factory"
)

// NewHeaderCheckHandler will create all components needed for header
// verification and returns the header verifier component. It behaves like a
// main factory for header verification components
func NewHeaderCheckHandler(
	proxy Proxy,
	enableEpochsConfig *data.EnableEpochsConfig,
) (HeaderVerifier, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}

	networkConfig, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ratingsConfig, err := proxy.GetRatingsConfig(context.Background())
	if err != nil {
		return nil, err
	}

	coreComp, err := factory.CreateCoreComponents(ratingsConfig, networkConfig, enableEpochsConfig)
	if err != nil {
		return nil, err
	}

	cryptoComp, err := factory.CreateCryptoComponents()
	if err != nil {
		return nil, err
	}

	args := crypto.MultiSigArgs{
		MultiSigHasherType:   "blake2b",
		BlSignKeyGen:         cryptoComp.KeyGen,
		ConsensusType:        "bls",
		ImportModeNoSigCheck: false,
	}

	multiSignerContainer, err := crypto.NewMultiSignerContainer(args, enableEpochsConfig.EnableEpochs.BLSMultiSignerEnableEpoch)
	if err != nil {
		return nil, err
	}

	genesisNodesConfig, err := proxy.GetGenesisNodesPubKeys(context.Background())
	if err != nil {
		return nil, err
	}

	nodesCoordinator, err := factory.CreateNodesCoordinator(
		coreComp,
		networkConfig,
		enableEpochsConfig,
		cryptoComp.PublicKey,
		genesisNodesConfig,
	)
	if err != nil {
		return nil, err
	}

	headerSigArgs := &headerCheck.ArgsHeaderSigVerifier{
		Marshalizer:             coreComp.Marshaller,
		Hasher:                  coreComp.Hasher,
		NodesCoordinator:        nodesCoordinator,
		MultiSigContainer:       multiSignerContainer,
		SingleSigVerifier:       cryptoComp.SingleSig,
		KeyGen:                  cryptoComp.KeyGen,
		FallbackHeaderValidator: &disabled.FallBackHeaderValidator{},
	}
	headerSigVerifier, err := headerCheck.NewHeaderSigVerifier(headerSigArgs)
	if err != nil {
		return nil, err
	}

	rawHeaderHandler, err := NewRawHeaderHandler(proxy, coreComp.Marshaller)
	if err != nil {
		return nil, err
	}

	headerVerifierArgs := ArgsHeaderVerifier{
		HeaderHandler:     rawHeaderHandler,
		HeaderSigVerifier: headerSigVerifier,
		NodesCoordinator:  nodesCoordinator,
	}
	headerVerifier, err := NewHeaderVerifier(headerVerifierArgs)
	if err != nil {
		return nil, err
	}

	return headerVerifier, nil
}
