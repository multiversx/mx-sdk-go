package testsCommon

import (
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/core"
)

// CryptoComponentsHolderStub -
type CryptoComponentsHolderStub struct {
	GetPublicKeyCalled      func() crypto.PublicKey
	GetPrivateKeyCalled     func() crypto.PrivateKey
	GetBech32Called         func() string
	GetAddressHandlerCalled func() core.AddressHandler
}

// GetPublicKey -
func (stub *CryptoComponentsHolderStub) GetPublicKey() crypto.PublicKey {
	if stub.GetPublicKeyCalled != nil {
		return stub.GetPublicKeyCalled()
	}
	return nil
}

// GetPrivateKey -
func (stub *CryptoComponentsHolderStub) GetPrivateKey() crypto.PrivateKey {
	if stub.GetPrivateKeyCalled != nil {
		return stub.GetPrivateKeyCalled()
	}
	return nil
}

// GetBech32 -
func (stub *CryptoComponentsHolderStub) GetBech32() string {
	if stub.GetBech32Called != nil {
		return stub.GetBech32Called()
	}
	return ""
}

// GetAddressHandler -
func (stub *CryptoComponentsHolderStub) GetAddressHandler() core.AddressHandler {
	if stub.GetAddressHandlerCalled != nil {
		return stub.GetAddressHandlerCalled()
	}
	return nil
}

// IsInterfaceNil -
func (stub *CryptoComponentsHolderStub) IsInterfaceNil() bool {
	return stub == nil
}
