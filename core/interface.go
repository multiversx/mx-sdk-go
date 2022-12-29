package core

import crypto "github.com/ElrondNetwork/elrond-go-crypto"

// AddressHandler will handle different implementations of an address
type AddressHandler interface {
	AddressAsBech32String() string
	AddressBytes() []byte
	AddressSlice() [32]byte
	IsValid() bool
	IsInterfaceNil() bool
}

// CryptoComponentsHolder is able to holder and provide all the crypto components
type CryptoComponentsHolder interface {
	GetPublicKey() crypto.PublicKey
	GetPrivateKey() crypto.PrivateKey
	GetBech32() string
	GetAddressHandler() AddressHandler
	IsInterfaceNil() bool
}
