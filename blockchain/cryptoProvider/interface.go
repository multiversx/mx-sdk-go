package cryptoProvider

import (
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
)

// CryptoComponentsHolder is able to holder and provide all the crypto components
type CryptoComponentsHolder interface {
	GetPublicKey() crypto.PublicKey
	GetPrivateKey() crypto.PrivateKey
	GetBech32() string
	GetAddressHandler() core.AddressHandler
	IsInterfaceNil() bool
}
