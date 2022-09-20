package nonceHandlerV1

import "github.com/ElrondNetwork/elrond-sdk-erdgo/core"

func (nth *nonceTransactionsHandlerV1) GetOrCreateAddressNonceHandler(address core.AddressHandler) *addressNonceHandler {
	return nth.getOrCreateAddressNonceHandler(address)
}
