package mock

import "github.com/multiversx/mx-sdk-go/authentication"

// AuthTokenHandlerStub -
type AuthTokenHandlerStub struct {
	DecodeCalled                   func(accessToken string) (authentication.AuthToken, error)
	EncodeCalled                   func(authToken authentication.AuthToken) (string, error)
	GetUnsignedTokenCalled         func(authToken authentication.AuthToken) []byte
	GetSignableMessageCalled       func(address, unsignedToken []byte) []byte
	GetSignableMessageLegacyCalled func(address, unsignedToken []byte) []byte
}

// Decode -
func (stub *AuthTokenHandlerStub) Decode(accessToken string) (authentication.AuthToken, error) {
	if stub.DecodeCalled != nil {
		return stub.DecodeCalled(accessToken)
	}
	return nil, nil
}

// Encode -
func (stub *AuthTokenHandlerStub) Encode(authToken authentication.AuthToken) (string, error) {
	if stub.EncodeCalled != nil {
		return stub.EncodeCalled(authToken)
	}
	return "", nil
}

// GetTokenBody -
func (stub *AuthTokenHandlerStub) GetUnsignedToken(authToken authentication.AuthToken) []byte {
	if stub.GetUnsignedTokenCalled != nil {
		return stub.GetUnsignedTokenCalled(authToken)
	}
	return nil
}

// GetSignableMessage -
func (stub *AuthTokenHandlerStub) GetSignableMessage(address, unsignedToken []byte) []byte {
	if stub.GetSignableMessageCalled != nil {
		return stub.GetSignableMessageCalled(address, unsignedToken)
	}
	return nil
}

// GetSignableMessageLegacy -
func (stub *AuthTokenHandlerStub) GetSignableMessageLegacy(address, unsignedToken []byte) []byte {
	if stub.GetSignableMessageLegacyCalled != nil {
		return stub.GetSignableMessageLegacyCalled(address, unsignedToken)
	}
	return nil
}

// IsInterfaceNil -
func (stub *AuthTokenHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
