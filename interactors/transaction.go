package interactors

import (
	"encoding/hex"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/hashing/keccak"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("interactors")

const defaultMillisecondsBetweenBunches = 1000 //1second
var txHasher = &keccak.Keccak{}

type transactionInteractor struct {
	Proxy
	TxSigner
	mutTxAccumulator           sync.RWMutex
	millisecondsBetweenBunches uint32
	txAccumulator              []*data.Transaction
}

// NewTransactionInteractor will create an interactor that extends the proxy functionality with some transaction-oriented functionality
func NewTransactionInteractor(proxy Proxy, txSigner TxSigner) (*transactionInteractor, error) {
	if check.IfNil(proxy) {
		return nil, ErrNilProxy
	}
	if check.IfNil(txSigner) {
		return nil, ErrNilTxSigner
	}

	return &transactionInteractor{
		Proxy:                      proxy,
		TxSigner:                   txSigner,
		millisecondsBetweenBunches: defaultMillisecondsBetweenBunches,
	}, nil
}

func (ti *transactionInteractor) SetMillisecondsBetweenBunches(millisecondsBetweenBunches uint32) {
	atomic.StoreUint32(&ti.millisecondsBetweenBunches, millisecondsBetweenBunches)
}

// AddTransaction will add the provided transaction in the transaction accumulator
func (ti *transactionInteractor) AddTransaction(tx *data.Transaction) {
	if tx == nil {
		return
	}

	ti.mutTxAccumulator.Lock()
	ti.txAccumulator = append(ti.txAccumulator, tx)
	ti.mutTxAccumulator.Unlock()
}

// PopAccumulatedTransactions will return the whole accumulated contents emptying the accumulator
func (ti *transactionInteractor) PopAccumulatedTransactions() []*data.Transaction {
	ti.mutTxAccumulator.Lock()
	result := make([]*data.Transaction, len(ti.txAccumulator))
	copy(result, ti.txAccumulator)
	ti.txAccumulator = make([]*data.Transaction, 0)
	ti.mutTxAccumulator.Unlock()

	return result
}

// createTransaction assembles a transaction from the provided arguments
func (ti *transactionInteractor) createTransaction(arg data.ArgCreateTransaction) *data.Transaction {
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

// ApplySignatureAndGenerateTransaction will apply the corresponding sender and compute the signature field and
// generate the transaction instance
func (ti *transactionInteractor) ApplySignatureAndGenerateTransaction(
	skBytes []byte,
	arg data.ArgCreateTransaction,
) (*data.Transaction, error) {

	pkBytes, err := ti.TxSigner.GeneratePkBytes(skBytes)
	if err != nil {
		return nil, err
	}

	arg.Signature = ""
	arg.SndAddr = core.AddressPublicKeyConverter.Encode(pkBytes)

	unsignedMessaged, err := ti.createUnsignedMessage(arg)
	if err != nil {
		return nil, err
	}

	shouldSignOnTxHash := arg.Version >= 2 && arg.Options&1 > 0
	if shouldSignOnTxHash {
		log.Debug("signing the transaction using the hash of the message")
		unsignedMessaged = txHasher.Compute(string(unsignedMessaged))
	}

	signature, err := ti.TxSigner.SignMessage(unsignedMessaged, skBytes)
	if err != nil {
		return nil, err
	}

	arg.Signature = hex.EncodeToString(signature)

	return ti.createTransaction(arg), nil
}

func (ti *transactionInteractor) createUnsignedMessage(arg data.ArgCreateTransaction) ([]byte, error) {
	tx := ti.createTransaction(arg)

	return json.Marshal(tx)
}

func (ti *transactionInteractor) SendTransactionsAsBunch(bunchSize int) ([]string, error) {
	if bunchSize <= 0 {
		return nil, ErrInvalidValue
	}

	millisecondsBetweenBunches := atomic.LoadUint32(&ti.millisecondsBetweenBunches)

	transactions := ti.PopAccumulatedTransactions()
	allHashes := make([]string, 0)
	for bunchIndex := 0; len(transactions) > 0; bunchIndex++ {
		var bunch []*data.Transaction

		log.Debug("sending bunch", "index", bunchIndex)

		if len(transactions) > bunchSize {
			bunch = transactions[0:bunchSize]
			transactions = transactions[bunchSize:]
		} else {
			bunch = transactions
			transactions = make([]*data.Transaction, 0)
		}

		hashes, err := ti.Proxy.SendTransactions(bunch)
		if err != nil {
			return nil, err
		}

		allHashes = append(allHashes, hashes...)

		time.Sleep(time.Duration(millisecondsBetweenBunches) * time.Millisecond)
	}

	return allHashes, nil
}
