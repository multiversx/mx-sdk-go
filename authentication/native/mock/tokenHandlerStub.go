package mock

import "github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"

// AuthTokenHandlerStub -
type AuthTokenHandlerStub struct {
	DecodeCalled func(accessToken string) (authentication.AuthToken, error)
	EncodeCalled func(authToken authentication.AuthToken) (string, error)
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

// IsInterfaceNil -
func (stub *AuthTokenHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
