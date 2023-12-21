package testsCommon

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	crypto "github.com/multiversx/mx-chain-crypto-go"
)

// SignerStub -
type SignerStub struct {
	SignTransactionCalled       func(tx *transaction.FrontendTransaction, privateKey crypto.PrivateKey) ([]byte, error)
	SignTransactionBufferCalled func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error)
	SignMessageCalled           func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error)
	VerifyMessageCalled         func(msg []byte, publicKey crypto.PublicKey, sig []byte) error
	SignByteSliceCalled         func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error)
	VerifyByteSliceCalled       func(msg []byte, publicKey crypto.PublicKey, sig []byte) error
}

// SignTransaction -
func (stub *SignerStub) SignTransaction(tx *transaction.FrontendTransaction, privateKey crypto.PrivateKey) ([]byte, error) {
	if stub.SignTransactionCalled != nil {
		return stub.SignTransactionCalled(tx, privateKey)
	}

	return make([]byte, 0), nil
}

// SignMessage -
func (stub *SignerStub) SignMessage(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	if stub.SignMessageCalled != nil {
		return stub.SignMessageCalled(msg, privateKey)
	}

	return make([]byte, 0), nil
}

// VerifyMessage -
func (stub *SignerStub) VerifyMessage(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
	if stub.VerifyMessageCalled != nil {
		return stub.VerifyMessageCalled(msg, publicKey, sig)
	}

	return nil
}

// VerifyByteSlice -
func (stub *SignerStub) VerifyByteSlice(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
	if stub.VerifyByteSliceCalled != nil {
		return stub.VerifyByteSliceCalled(msg, publicKey, sig)
	}

	return nil
}

// SignByteSlice -
func (stub *SignerStub) SignByteSlice(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	if stub.SignByteSliceCalled != nil {
		return stub.SignByteSliceCalled(msg, privateKey)
	}

	return make([]byte, 0), nil
}

// IsInterfaceNil -
func (stub *SignerStub) IsInterfaceNil() bool {
	return stub == nil
}
