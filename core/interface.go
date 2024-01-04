package core

import crypto "github.com/multiversx/mx-chain-crypto-go"

// AddressHandler will handle different implementations of an address
type AddressHandler interface {
	AddressAsBech32String() (string, error)
	AddressBytes() []byte
	AddressSlice() [32]byte
	IsValid() bool
	Pretty() string
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
