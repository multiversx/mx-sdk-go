package txcheck

import (
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// VerifyTransactionSignature handles the signature verification for a given transaction
func VerifyTransactionSignature(
	tx *data.Transaction,
	pk []byte,
	signature []byte,
	verifier builders.TxSigVerifier,
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

	return verifier.Verify(pk, unsignedMessage, signature)
}

func checkParams(
	tx *data.Transaction,
	pk []byte,
	signature []byte,
	verifier builders.TxSigVerifier,
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
