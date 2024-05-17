package txcheck

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/builders"
)

// VerifyTransactionSignature handles the signature verification for a given transaction
func VerifyTransactionSignature(
	tx *transaction.FrontendTransaction,
	pk crypto.PublicKey,
	signature []byte,
	verifier builders.Signer,
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

	shouldVerifyOnTxHash := unsignedTx.Version >= 2 && unsignedTx.Options&transaction.MaskSignedWithHash > 0
	if shouldVerifyOnTxHash {
		unsignedMessage = hasher.Compute(string(unsignedMessage))
	}

	return verifier.VerifyByteSlice(unsignedMessage, pk, signature)
}

func checkParams(
	tx *transaction.FrontendTransaction,
	pk crypto.PublicKey,
	signature []byte,
	verifier builders.Signer,
	marshaller coreData.Marshaller,
	hasher coreData.Hasher,
) error {
	if tx == nil {
		return ErrNilTransaction
	}
	if len(signature) == 0 {
		return ErrNilSignature
	}
	if check.IfNil(pk) {
		return ErrNilPubKey
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
