package nonceHandlerV2

import (
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
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
		transactions: make(map[uint64]*data.Transaction),
	}, nil
}
