package testsCommon

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain/cryptoProvider"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplySignatureAndGenerateTxCalled func(cryptoHolder cryptoProvider.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error)
}

// ApplySignatureAndGenerateTx -
func (stub *TxBuilderStub) ApplySignatureAndGenerateTx(cryptoHolder cryptoProvider.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error) {
	if stub.ApplySignatureAndGenerateTxCalled != nil {
		return stub.ApplySignatureAndGenerateTxCalled(cryptoHolder, arg)
	}

	return &data.Transaction{}, nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
