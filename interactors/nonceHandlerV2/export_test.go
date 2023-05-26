package nonceHandlerV2

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// NewAddressNonceHandlerWithPrivateAccess -
func NewAddressNonceHandlerWithPrivateAccess(proxy interactors.Proxy, address erdgoCore.AddressHandler) (*addressNonceHandler, error) {
	if check.IfNil(proxy) {
		return nil, interactors.ErrNilProxy
	}
	if check.IfNil(address) {
		return nil, interactors.ErrNilAddress
	}
	return &addressNonceHandler{
		address:      address,
		proxy:        proxy,
		transactions: make(map[uint64]*transaction.FrontendTransaction),
	}, nil
}
