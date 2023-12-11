package builders

import (
	"encoding/hex"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/core"
)

var (
	blake2bHasher          = blake2b.NewBlake2b()
	nodeInternalMarshaller = &marshal.GogoProtoMarshalizer{}
)

type txBuilder struct {
	signer Signer
}

// NewTxBuilder will create a new transaction builder able to build and correctly sign a transaction
func NewTxBuilder(signer Signer) (*txBuilder, error) {
	if check.IfNil(signer) {
		return nil, ErrNilSigner
	}

	return &txBuilder{
		signer: signer,
	}, nil
}

// ApplySignature will apply the corresponding sender and compute and set the signature field
func (builder *txBuilder) ApplySignature(
	cryptoHolder core.CryptoComponentsHolder,
	tx *transaction.FrontendTransaction,
) error {
	tx.Sender = cryptoHolder.GetBech32()
	unsignedMessage := builder.createUnsignedTx(tx)

	signature, err := builder.signer.SignTransaction(unsignedMessage, cryptoHolder.GetPrivateKey())
	if err != nil {
		return err
	}

	tx.Signature = hex.EncodeToString(signature)

	return nil
}

// ComputeTxHash will return the hash of the provided transaction. It assumes that the transaction is already signed,
// otherwise it will return an error.
func (builder *txBuilder) ComputeTxHash(tx *transaction.FrontendTransaction) ([]byte, error) {
	if len(tx.Signature) == 0 {
		return nil, ErrMissingSignature
	}

	nodeTx, err := transactionToNodeTransaction(tx)
	if err != nil {
		return nil, err
	}

	txBytes, err := nodeInternalMarshaller.Marshal(nodeTx)
	if err != nil {
		return nil, err
	}

	txHash := blake2bHasher.Compute(string(txBytes))
	return txHash, nil
}

func transactionToNodeTransaction(tx *transaction.FrontendTransaction) (*transaction.Transaction, error) {
	receiverBytes, err := core.AddressPublicKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return nil, err
	}

	senderBytes, err := core.AddressPublicKeyConverter.Decode(tx.Sender)
	if err != nil {
		return nil, err
	}

	signaturesBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return nil, err
	}

	valueBI, ok := big.NewInt(0).SetString(tx.Value, 10)
	if !ok {
		return nil, ErrInvalidValue
	}

	return &transaction.Transaction{
		Nonce:     tx.Nonce,
		Value:     valueBI,
		RcvAddr:   receiverBytes,
		SndAddr:   senderBytes,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      tx.Data,
		ChainID:   []byte(tx.ChainID),
		Version:   tx.Version,
		Signature: signaturesBytes,
		Options:   tx.Options,
	}, nil
}

func (builder *txBuilder) createUnsignedTx(tx *transaction.FrontendTransaction) *transaction.FrontendTransaction {
	copiedTransaction := *tx
	copiedTransaction.Signature = ""

	return &copiedTransaction
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txBuilder) IsInterfaceNil() bool {
	return builder == nil
}
