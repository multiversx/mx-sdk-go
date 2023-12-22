package mock

// AuthTokenStub -
type AuthTokenStub struct {
	GetTtlCalled         func() int64
	GetAddressCalled     func() []byte
	GetHostCalled        func() []byte
	GetSignatureCalled   func() []byte
	GetBlockHashCalled   func() string
	GetExtraInfoCalled   func() []byte
	IsInterfaceNilCalled func() bool
}

// GetTtl -
func (stub *AuthTokenStub) GetTtl() int64 {
	if stub.GetTtlCalled != nil {
		return stub.GetTtlCalled()
	}
	return 0
}

// GetAddress -
func (stub *AuthTokenStub) GetAddress() []byte {
	if stub.GetAddressCalled != nil {
		return stub.GetAddressCalled()
	}
	return []byte("")
}

// GetHost -
func (stub *AuthTokenStub) GetHost() []byte {
	if stub.GetHostCalled != nil {
		return stub.GetHostCalled()
	}
	return []byte("")
}

// GetSignature -
func (stub *AuthTokenStub) GetSignature() []byte {
	if stub.GetSignatureCalled != nil {
		return stub.GetSignatureCalled()
	}
	return []byte("")
}

// GetBlockHash -
func (stub *AuthTokenStub) GetBlockHash() string {
	if stub.GetBlockHashCalled != nil {
		return stub.GetBlockHashCalled()
	}
	return ""
}

// GetExtraInfo -
func (stub *AuthTokenStub) GetExtraInfo() []byte {
	if stub.GetExtraInfoCalled != nil {
		return stub.GetExtraInfoCalled()
	}
	return []byte("")
}

// IsInterfaceNil -
func (stub *AuthTokenStub) IsInterfaceNil() bool {
	return stub == nil
}
