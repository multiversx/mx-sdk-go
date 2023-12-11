package testsCommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	sdkCore "github.com/multiversx/mx-sdk-go/mx-sdk-go-old/core"
)

// TxBuilderStub -
type TxBuilderStub struct {
	ApplySignatureCalled func(cryptoHolder sdkCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

// ApplySignature -
func (stub *TxBuilderStub) ApplySignature(cryptoHolder sdkCore.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if stub.ApplySignatureCalled != nil {
		return stub.ApplySignatureCalled(cryptoHolder, tx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxBuilderStub) IsInterfaceNil() bool {
	return stub == nil
}
