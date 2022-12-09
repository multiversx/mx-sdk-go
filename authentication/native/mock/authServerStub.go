package mock

// AuthServerStub -
type AuthServerStub struct {
	ValidateCalled func(accessToken string) (string, error)
}

// Validate -
func (stub *AuthServerStub) Validate(accessToken string) (string, error) {
	if stub.ValidateCalled != nil {
		return stub.ValidateCalled(accessToken)
	}
	return "", nil
}

// IsInterfaceNil -
func (stub *AuthServerStub) IsInterfaceNil() bool {
	return stub == nil
}
