package builders

import (
	"encoding/hex"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var (
	blake2bHasher          = blake2b.NewBlake2b()
	nodeInternalMarshaller = &marshal.GogoProtoMarshalizer{}
)

type txBuilder struct {
	xSigner Signer
}

// NewTxBuilder will create a new transaction builder able to build and correctly sign a transaction
func NewTxBuilder(xSigner Signer) (*txBuilder, error) {
	if check.IfNil(xSigner) {
		return nil, ErrNilSigner
	}

	return &txBuilder{
		xSigner: xSigner,
	}, nil
}

// createTransaction assembles a transaction from the provided arguments
func (builder *txBuilder) createTransaction(arg data.ArgCreateTransaction) *data.Transaction {
	return &data.Transaction{
		Nonce:     arg.Nonce,
		Value:     arg.Value,
		RcvAddr:   arg.RcvAddr,
		SndAddr:   arg.SndAddr,
		GasPrice:  arg.GasPrice,
		GasLimit:  arg.GasLimit,
		Data:      arg.Data,
		Signature: arg.Signature,
		ChainID:   arg.ChainID,
		Version:   arg.Version,
		Options:   arg.Options,
	}
}

// ApplySignatureAndGenerateTx will apply the corresponding sender and compute the signature field and
// generate the transaction instance
func (builder *txBuilder) ApplySignatureAndGenerateTx(
	cryptoHolder core.CryptoComponentsHolder,
	arg data.ArgCreateTransaction,
) (*data.Transaction, error) {
	arg.SndAddr = cryptoHolder.GetBech32()
	unsignedMessage, err := builder.createUnsignedTx(arg)
	if err != nil {
		return nil, err
	}

	signature, err := builder.xSigner.SignTransaction(unsignedMessage, cryptoHolder.GetPrivateKey())
	if err != nil {
		return nil, err
	}

	arg.Signature = hex.EncodeToString(signature)

	return builder.createTransaction(arg), nil
}

// ComputeTxHash will return the hash of the provided transaction. It assumes that the transaction is already signed,
// otherwise it will return an error.
// The input can be the result of the ApplySignatureAndGenerateTx function
func (builder *txBuilder) ComputeTxHash(tx *data.Transaction) ([]byte, error) {
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

func transactionToNodeTransaction(tx *data.Transaction) (*transaction.Transaction, error) {
	receiverBytes, err := core.AddressPublicKeyConverter.Decode(tx.RcvAddr)
	if err != nil {
		return nil, err
	}

	senderBytes, err := core.AddressPublicKeyConverter.Decode(tx.SndAddr)
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

func (builder *txBuilder) createUnsignedTx(arg data.ArgCreateTransaction) (*data.Transaction, error) {
	arg.Signature = ""
	tx := builder.createTransaction(arg)

	return tx, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txBuilder) IsInterfaceNil() bool {
	return builder == nil
}
