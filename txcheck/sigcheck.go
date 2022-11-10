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

func VerifyTransactionSignature(
	tx data.Transaction,
	pk []byte,
	signature []byte,
	verifier SigVerifier,
	marshaller coreData.Marshaller,
	hasher coreData.Hasher,
) error {
	err := checkParams(pk, signature, verifier, marshaller, hasher)
	if err != nil {
		return err
	}

	unsignedTx := builders.TransactionToUnsignedTx(&tx)
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
	pk []byte,
	signature []byte,
	verifier SigVerifier,
	marshaller coreData.Marshaller,
	hasher coreData.Hasher, ) error {
	if len(pk) == 0 {
		return errNilPubKey
	}
	if len(signature) == 0 {
		return errNilSignature
	}
	if check.IfNil(verifier) {
		return errNilSignatureVerifier
	}
	if check.IfNil(marshaller) {
		return errNilMarshaller
	}
	if check.IfNil(hasher) {
		return errNilHasher
	}
	return nil
}
