package testsCommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplySignatureCalled func(cryptoHolder erdgoCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

// ApplySignature -
func (stub *TxBuilderStub) ApplySignature(cryptoHolder erdgoCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if stub.ApplySignatureCalled != nil {
		return stub.ApplySignatureCalled(cryptoHolder, tx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
