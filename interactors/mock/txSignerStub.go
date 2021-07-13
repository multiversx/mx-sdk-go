package mock

// TxSignerStub -
type TxSignerStub struct {
	SignMessageCalled     func(msg []byte, skBytes []byte) ([]byte, error)
	GeneratePkBytesCalled func(skBytes []byte) ([]byte, error)
}

// SignMessage -
func (ts *TxSignerStub) SignMessage(msg []byte, skBytes []byte) ([]byte, error) {
	return ts.SignMessageCalled(msg, skBytes)
}

// GeneratePkBytes -
func (ts *TxSignerStub) GeneratePkBytes(skBytes []byte) ([]byte, error) {
	return ts.GeneratePkBytesCalled(skBytes)
}

// IsInterfaceNil -
func (ts *TxSignerStub) IsInterfaceNil() bool {
	return ts == nil
}
