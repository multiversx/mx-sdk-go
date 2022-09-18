package mock

// NativeStub -
type NativeStub struct {
	GetAccessTokenCalled func() (string, error)
}

// GetAccessToken -
func (stub *NativeStub) GetAccessToken() (string, error) {
	if stub.GetAccessTokenCalled != nil {
		return stub.GetAccessTokenCalled()
	}
	return "", nil
}
