package testsCommon

import "github.com/ElrondNetwork/elrond-go-crypto"

// PublicKeyStub -
type PublicKeyStub struct {
	ToByteArrayCalled func() ([]byte, error)
	SuiteCalled       func() crypto.Suite
	PointCalled       func() crypto.Point
}

// ToByteArray -
func (stub *PublicKeyStub) ToByteArray() ([]byte, error) {
	if stub.ToByteArrayCalled != nil {
		return stub.ToByteArrayCalled()
	}

	return make([]byte, 0), nil
}

// Suite -
func (stub *PublicKeyStub) Suite() crypto.Suite {
	if stub.SuiteCalled != nil {
		return stub.SuiteCalled()
	}

	return nil
}

// Point -
func (stub *PublicKeyStub) Point() crypto.Point {
	if stub.PointCalled != nil {
		return stub.PointCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *PublicKeyStub) IsInterfaceNil() bool {
	return stub == nil
}
