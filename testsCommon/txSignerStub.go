package testsCommon

// TxSignerStub -
type TxSignerStub struct {
	SignMessageCalled     func(msg []byte, skBytes []byte) ([]byte, error)
	GeneratePkBytesCalled func(skBytes []byte) ([]byte, error)
}

// SignMessage -
func (stub *TxSignerStub) SignMessage(msg []byte, skBytes []byte) ([]byte, error) {
	if stub.SignMessageCalled != nil {
		return stub.SignMessageCalled(msg, skBytes)
	}

	return make([]byte, 0), nil
}

// GeneratePkBytes -
func (stub *TxSignerStub) GeneratePkBytes(skBytes []byte) ([]byte, error) {
	if stub.GeneratePkBytesCalled != nil {
		return stub.GeneratePkBytesCalled(skBytes)
	}

	return make([]byte, 0), nil
}

// IsInterfaceNil -
func (stub *TxSignerStub) IsInterfaceNil() bool {
	return stub == nil
}
