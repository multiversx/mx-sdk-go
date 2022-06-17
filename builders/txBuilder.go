package builders

import (
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
	skBytes []byte,
	arg data.ArgCreateTransaction,
) (*data.Transaction, error) {

	pkBytes, err := builder.txSigner.GeneratePkBytes(skBytes)
	if err != nil {
		return nil, err
	}

	arg.SndAddr = core.AddressPublicKeyConverter.Encode(pkBytes)
	unsignedMessage, err := builder.createUnsignedMessage(arg)
	if err != nil {
		return nil, err
	}

	shouldSignOnTxHash := arg.Version >= 2 && arg.Options&1 > 0
	if shouldSignOnTxHash {
		log.Debug("signing the transaction using the hash of the message")
		unsignedMessage = txHasher.Compute(string(unsignedMessage))
	}

	signature, err := builder.txSigner.SignMessage(unsignedMessage, skBytes)
	if err != nil {
		return nil, err
	}

	arg.Signature = hex.EncodeToString(signature)

	return builder.createTransaction(arg), nil
}

// ApplySignatureAndGenerateTxHash will sign the transaction and return it's hash
func (builder *txBuilder) ApplySignatureAndGenerateTxHash(
	skBytes []byte,
	arg data.ArgCreateTransaction,
) ([]byte, *data.Transaction, error) {
	signedTx, err := builder.ApplySignatureAndGenerateTx(skBytes, arg)
	if err != nil {
		return nil, nil, err
	}

	nodeTx, err := transactionToNodeTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}

	txBytes, err := nodeInternalMarshaller.Marshal(nodeTx)
	if err != nil {
		return nil, nil, err
	}

	txHash := blake2bHasher.Compute(string(txBytes))
	return txHash, signedTx, nil
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

func (builder *txBuilder) createUnsignedMessage(arg data.ArgCreateTransaction) ([]byte, error) {
	arg.Signature = ""
	tx := builder.createTransaction(arg)

	return json.Marshal(tx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txBuilder) IsInterfaceNil() bool {
	return builder == nil
}
