package mock

import "github.com/multiversx/mx-sdk-go/authentication"

// AuthServerStub -
type AuthServerStub struct {
	ValidateCalled func(accessToken authentication.AuthToken) error
}

// Validate -
func (stub *AuthServerStub) Validate(accessToken authentication.AuthToken) error {
	if stub.ValidateCalled != nil {
		return stub.ValidateCalled(accessToken)
	}
	return nil
}

// IsInterfaceNil -
func (stub *AuthServerStub) IsInterfaceNil() bool {
	return stub == nil
}
