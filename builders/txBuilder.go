package builders

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/builders")
var txHasher = keccak.NewKeccak()

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

func (builder *txBuilder) createUnsignedMessage(arg data.ArgCreateTransaction) ([]byte, error) {
	arg.Signature = ""
	tx := builder.createTransaction(arg)

	return json.Marshal(tx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *txBuilder) IsInterfaceNil() bool {
	return builder == nil
}
