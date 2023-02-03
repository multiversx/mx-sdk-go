package testsCommon

import (
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplySignatureAndGenerateTxCalled func(cryptoHolder erdgoCore.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error)
}

// ApplySignatureAndGenerateTx -
func (stub *TxBuilderStub) ApplyUserSignatureAndGenerateTx(cryptoHolder erdgoCore.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error) {
	if stub.ApplySignatureAndGenerateTxCalled != nil {
		return stub.ApplySignatureAndGenerateTxCalled(cryptoHolder, arg)
	}

	return &data.Transaction{}, nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
