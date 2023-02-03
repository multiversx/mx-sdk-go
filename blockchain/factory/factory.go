package factory

import (
	"fmt"

	"github.com/multiversx/mx-sdk-go/blockchain/endpointProviders"
	"github.com/multiversx/mx-sdk-go/blockchain/finalityProvider"
	"github.com/multiversx/mx-sdk-go/core"
)

// CreateEndpointProvider creates a new instance of EndpointProvider
func CreateEndpointProvider(entityType core.RestAPIEntityType) (EndpointProvider, error) {
	switch entityType {
	case core.ObserverNode:
		return endpointProviders.NewNodeEndpointProvider(), nil
	case core.Proxy:
		return endpointProviders.NewProxyEndpointProvider(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownRestAPIEntityType, entityType)
	}
}

// CreateFinalityProvider creates a new instance of FinalityProvider
func CreateFinalityProvider(proxy proxy, checkFinality bool) (FinalityProvider, error) {
	if !checkFinality {
		return finalityProvider.NewDisabledFinalityProvider(), nil
	}

	switch proxy.GetRestAPIEntityType() {
	case core.ObserverNode:
		return finalityProvider.NewNodeFinalityProvider(proxy)
	case core.Proxy:
		return finalityProvider.NewProxyFinalityProvider(proxy)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownRestAPIEntityType, proxy.GetRestAPIEntityType())
	}
}
