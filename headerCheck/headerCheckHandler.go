package headerCheck

import (
	"context"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/process/headerCheck"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/disabled"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerCheck/factory"
)

// NewHeaderCheckHandler will create all components needed for header
// verification and returns the header verifier component. It behaves like a
// main factory for header verification components
func NewHeaderCheckHandler(proxy Proxy) (HeaderVerifier, error) {
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

	enableEpochsConfig, err := proxy.GetEnableEpochsConfig(context.Background())
	if err != nil {
		return nil, err
	}

	coreComp, err := factory.CreateCoreComponents(ratingsConfig, networkConfig)
	if err != nil {
		return nil, err
	}

	cryptoComp, err := factory.CreateCryptoComponents()
	if err != nil {
		return nil, err
	}

	nodesCoordinator, err := factory.CreateNodesCoordinator(
		coreComp.Hasher,
		coreComp.Marshaller,
		coreComp.Rater,
		networkConfig,
		enableEpochsConfig,
	)
	if err != nil {
		return nil, err
	}

	headerSigArgs := &headerCheck.ArgsHeaderSigVerifier{
		Marshalizer:             coreComp.Marshaller,
		Hasher:                  coreComp.Hasher,
		NodesCoordinator:        nodesCoordinator,
		MultiSigVerifier:        cryptoComp.MultiSig,
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
