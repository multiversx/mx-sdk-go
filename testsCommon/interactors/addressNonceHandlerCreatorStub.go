package interactors

import (
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// AddressNonceHandlerCreatorStub -
type AddressNonceHandlerCreatorStub struct {
	CreateCalled func(proxy interactors.Proxy, address core.AddressHandler) (interactors.AddressNonceHandler, error)
}

// Create -
func (stub *AddressNonceHandlerCreatorStub) Create(proxy interactors.Proxy, address core.AddressHandler) (interactors.AddressNonceHandler, error) {
	if stub.CreateCalled != nil {
		return stub.CreateCalled(proxy, address)
	}
	return nil, nil
}

// IsInterfaceNil -
func (stub *AddressNonceHandlerCreatorStub) IsInterfaceNil() bool {
	return stub == nil
}
