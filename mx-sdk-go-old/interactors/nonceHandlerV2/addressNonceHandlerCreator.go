package nonceHandlerV2

import (
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/core"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/interactors"
)

// AddressNonceHandlerCreator is used to create addressNonceHandler instances
type AddressNonceHandlerCreator struct{}

// Create will create
func (anhc *AddressNonceHandlerCreator) Create(proxy interactors.Proxy, address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	return NewAddressNonceHandler(proxy, address)
}

// IsInterfaceNil returns true if there is no value under the interface
func (anhc *AddressNonceHandlerCreator) IsInterfaceNil() bool {
	return anhc == nil
}
