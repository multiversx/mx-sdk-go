package builders

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var (
	log                    = logger.GetOrCreate("elrond-sdk-erdgo/builders")
	txHasher               = keccak.NewKeccak()
	blake2bHasher          = blake2b.NewBlake2b()
	nodeInternalMarshaller = &marshal.GogoProtoMarshalizer{}
)

type txBuilder struct {
	txSigner TxSigner
}

// NewTxBuilder will create a new transaction builder able to build and correctly sign a transaction
func NewTxBuilder(txSigner TxSigner) (*txBuilder, error) {
	if check.IfNil(txSigner) {
		return nil, ErrNilTxSigner
	}

	return &txBuilder{
		txSigner: txSigner,
	}, nil
}

// createTransaction assembles a transaction from the provided arguments
func (builder *txBuilder) createTransaction(arg data.ArgCreateTransaction) *data.Transaction {
	return &data.Transaction{
		Nonce:             arg.Nonce,
		Value:             arg.Value,
		RcvAddr:           arg.RcvAddr,
		SndAddr:           arg.SndAddr,
		GasPrice:          arg.GasPrice,
		GasLimit:          arg.GasLimit,
		Data:              arg.Data,
		Signature:         arg.Signature,
		ChainID:           arg.ChainID,
		Version:           arg.Version,
		Options:           arg.Options,
		GuardianAddr:      arg.GuardianAddr,
		GuardianSignature: arg.GuardianSignature,
	}
}

// ApplyUserSignatureAndGenerateTx will apply the corresponding sender and compute the signature field and
// generate the transaction instance
func (builder *txBuilder) ApplyUserSignatureAndGenerateTx(
	skBytes []byte,
	arg data.ArgCreateTransaction,
) (*data.Transaction, error) {
	newArg := arg
	pkBytes, err := builder.txSigner.GeneratePkBytes(skBytes)
	if err != nil {
		return nil, err
	}

	newArg.SndAddr = core.AddressPublicKeyConverter.Encode(pkBytes)
	unsignedTx, err := builder.CreateUnsignedTransaction(newArg)
	if err != nil {
		return nil, err
	}

	signature, err := builder.signTx(unsignedTx, skBytes)
	if err != nil {
		return nil, err
	}

	newArg.Signature = hex.EncodeToString(signature)

	return builder.createTransaction(newArg), nil
}

func (builder *txBuilder) signTx(unsignedTx *data.Transaction, skBytes []byte) ([]byte, error) {
	unsignedMessage, err := json.Marshal(unsignedTx)
	if err != nil {
		return nil, err
	}

	shouldSignOnTxHash := unsignedTx.Version >= 2 && unsignedTx.Options&1 > 0
	if shouldSignOnTxHash {
		log.Debug("signing the transaction using the hash of the message")
		unsignedMessage = txHasher.Compute(string(unsignedMessage))
	}

	return builder.txSigner.SignMessage(unsignedMessage, skBytes)
}

// ApplyGuardianSignature applies the guardian signature over the transaction.
// Does a basic check for the transaction options and guardian address.
func (builder *txBuilder) ApplyGuardianSignature(
	skGuardianBytes []byte,
	tx *data.Transaction,
) error {
	nodeTx, err:= transactionToNodeTransaction(tx)
	if err!= nil{
		return err
	}

	if !nodeTx.HasOptionGuardianSet() {
		return ErrMissingGuardianOption
	}

	pkGuardianBytes, err := builder.txSigner.GeneratePkBytes(skGuardianBytes)
	if err != nil {
		return err
	}

	txGuardianAddrBytes, err := core.AddressPublicKeyConverter.Decode(tx.GuardianAddr)
	if err != nil {
		return err
	}

	if !bytes.Equal(txGuardianAddrBytes, pkGuardianBytes) {
		return ErrGuardianDoesNotMatch
	}

	unsignedTx := TransactionToUnsignedTx(tx)
	guardianSignature, err := builder.signTx(unsignedTx, skGuardianBytes)
	if err != nil {
		return err
	}

	tx.GuardianSignature = hex.EncodeToString(guardianSignature)

	return err
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

	var guardianAddrBytes, guardianSigBytes []byte
	if len(tx.GuardianAddr) > 0 {
		guardianAddrBytes, err = core.AddressPublicKeyConverter.Decode(tx.GuardianAddr)
		if err != nil {
			return nil, err
		}

		guardianSigBytes, err = hex.DecodeString(tx.GuardianSignature)
		if err != nil {
			return nil, err
		}
	}

	return &transaction.Transaction{
		Nonce:             tx.Nonce,
		Value:             valueBI,
		RcvAddr:           receiverBytes,
		SndAddr:           senderBytes,
		GasPrice:          tx.GasPrice,
		GasLimit:          tx.GasLimit,
		Data:              tx.Data,
		ChainID:           []byte(tx.ChainID),
		Version:           tx.Version,
		Signature:         signaturesBytes,
		Options:           tx.Options,
		GuardianAddr:      guardianAddrBytes,
		GuardianSignature: guardianSigBytes,
	}, nil
}

func TransactionToUnsignedTx(tx *data.Transaction) *data.Transaction {
	unsignedTx := *tx
	unsignedTx.Signature = ""
	unsignedTx.GuardianSignature = ""
	return &unsignedTx
}

func (builder *txBuilder) CreateUnsignedTransaction(arg data.ArgCreateTransaction) (*data.Transaction, error) {
	tx := builder.createTransaction(arg)

	return TransactionToUnsignedTx(tx), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txBuilder) IsInterfaceNil() bool {
	return builder == nil
}
