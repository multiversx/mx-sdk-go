package testsCommon

import "github.com/multiversx/mx-chain-crypto-go"

// PrivateKeyStub -
type PrivateKeyStub struct {
	ToByteArrayCalled    func() ([]byte, error)
	SuiteCalled          func() crypto.Suite
	GeneratePublicCalled func() crypto.PublicKey
	ScalarCalled         func() crypto.Scalar
}

// ToByteArray -
func (stub *PrivateKeyStub) ToByteArray() ([]byte, error) {
	if stub.ToByteArrayCalled != nil {
		return stub.ToByteArrayCalled()
	}

	return make([]byte, 0), nil
}

// Suite -
func (stub *PrivateKeyStub) Suite() crypto.Suite {
	if stub.SuiteCalled != nil {
		return stub.SuiteCalled()
	}

	return nil
}

// GeneratePublic -
func (stub *PrivateKeyStub) GeneratePublic() crypto.PublicKey {
	if stub.GeneratePublicCalled != nil {
		return stub.GeneratePublicCalled()
	}

	return &PublicKeyStub{}
}

// Scalar -
func (stub *PrivateKeyStub) Scalar() crypto.Scalar {
	if stub.ScalarCalled != nil {
		return stub.ScalarCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *PrivateKeyStub) IsInterfaceNil() bool {
	return stub == nil
}
