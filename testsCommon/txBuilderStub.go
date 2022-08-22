package testsCommon

import "github.com/ElrondNetwork/elrond-sdk-erdgo/data"

// TxBuilderStub -
type TxBuilderStub struct {
	ApplySignatureAndGenerateTxCalled func(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error)
}

// ApplySignatureAndGenerateTx -
func (stub *TxBuilderStub) ApplyUserSignatureAndGenerateTx(skBytes []byte, arg data.ArgCreateTransaction) (*data.Transaction, error) {
	if stub.ApplySignatureAndGenerateTxCalled != nil {
		return stub.ApplySignatureAndGenerateTxCalled(skBytes, arg)
	}

	return &data.Transaction{}, nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
