package testsCommon

// XSignerStub -
type XSignerStub struct {
	SignTransactionCalled func(tx []byte, skBytes []byte) ([]byte, error)
	SignMessageCalled     func(msg []byte, skBytes []byte) ([]byte, error)
	GeneratePkBytesCalled func(skBytes []byte) ([]byte, error)
}

// SignTransaction -
func (stub *XSignerStub) SignTransaction(tx []byte, skBytes []byte) ([]byte, error) {
	if stub.SignTransactionCalled != nil {
		return stub.SignTransactionCalled(tx, skBytes)
	}

	return make([]byte, 0), nil
}

// SignMessage -
func (stub *XSignerStub) SignMessage(msg []byte, skBytes []byte) ([]byte, error) {
	if stub.SignMessageCalled != nil {
		return stub.SignMessageCalled(msg, skBytes)
	}

	return make([]byte, 0), nil
}

// GeneratePkBytes -
func (stub *XSignerStub) GeneratePkBytes(skBytes []byte) ([]byte, error) {
	if stub.GeneratePkBytesCalled != nil {
		return stub.GeneratePkBytesCalled(skBytes)
	}

	return make([]byte, 0), nil
}

// IsInterfaceNil -
func (stub *XSignerStub) IsInterfaceNil() bool {
	return stub == nil
}
