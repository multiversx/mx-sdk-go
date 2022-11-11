package txcheck

import (
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// SigVerifier defines the methods available for a signature verifier
type SigVerifier interface {
	Verify(pk []byte, msg []byte, sigBytes []byte) error
	IsInterfaceNil() bool
}

// VerifyTransactionSignature handles the signature verification for a given transaction
func VerifyTransactionSignature(
	tx *data.Transaction,
	pk []byte,
	signature []byte,
	verifier SigVerifier,
	marshaller coreData.Marshaller,
	hasher coreData.Hasher,
) error {
	err := checkParams(tx, pk, signature, verifier, marshaller, hasher)
	if err != nil {
		return err
	}

	unsignedTx := builders.TransactionToUnsignedTx(tx)
	unsignedMessage, err := marshaller.Marshal(unsignedTx)
	if err != nil {
		return err
	}

	shouldSignOnTxHash := unsignedTx.Version >= 2 && unsignedTx.Options&1 > 0
	if shouldSignOnTxHash {
		unsignedMessage = hasher.Compute(string(unsignedMessage))
	}

	return verifier.Verify(pk, unsignedMessage, signature)
}

func checkParams(
	tx *data.Transaction,
	pk []byte,
	signature []byte,
	verifier SigVerifier,
	marshaller coreData.Marshaller,
	hasher coreData.Hasher,
) error {
	if tx == nil {
		return ErrNilTransaction
	}
	if len(pk) == 0 {
		return ErrNilPubKey
	}
	if len(signature) == 0 {
		return ErrNilSignature
	}
	if check.IfNil(verifier) {
		return ErrNilSignatureVerifier
	}
	if check.IfNil(marshaller) {
		return ErrNilMarshaller
	}
	if check.IfNil(hasher) {
		return ErrNilHasher
	}
	return nil
}
