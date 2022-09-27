package nonceHandlerV2

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

// singleTransactionAddressNonceHandlerCreator is used to create singleTransactionAddressNonceHandler instances
type singleTransactionAddressNonceHandlerCreator struct{}

// Create will create
func (anhc *singleTransactionAddressNonceHandlerCreator) Create(proxy interactors.Proxy, address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	return NewSingleTransactionAddressNonceHandler(proxy, address)
}

// IsInterfaceNil returns true if there is no value under the interface
func (anhc *singleTransactionAddressNonceHandlerCreator) IsInterfaceNil() bool {
	return anhc == nil
}
