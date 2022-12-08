package interactors

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
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
