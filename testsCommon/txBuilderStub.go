package testsCommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplyUserSignatureCalled func(cryptoHolder erdgoCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

// ApplyUserSignature -
func (stub *TxBuilderStub) ApplyUserSignature(cryptoHolder erdgoCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if stub.ApplyUserSignatureCalled != nil {
		return stub.ApplyUserSignatureCalled(cryptoHolder, tx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
