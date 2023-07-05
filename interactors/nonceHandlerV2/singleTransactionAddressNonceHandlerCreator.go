package nonceHandlerV2

import (
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// SingleTransactionAddressNonceHandlerCreator is used to create singleTransactionAddressNonceHandler instances
type SingleTransactionAddressNonceHandlerCreator struct{}

// Create will create
func (anhc *SingleTransactionAddressNonceHandlerCreator) Create(proxy interactors.Proxy, address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	return NewSingleTransactionAddressNonceHandler(proxy, address)
}

// IsInterfaceNil returns true if there is no value under the interface
func (anhc *SingleTransactionAddressNonceHandlerCreator) IsInterfaceNil() bool {
	return anhc == nil
}
