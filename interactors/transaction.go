package interactors

import (
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var log = logger.GetOrCreate("interactors")

// ArgCreateTransaction will hold the transaction fields
type ArgCreateTransaction struct {
	Nonce     uint64
	Value     string
	RcvAddr   string
	SndAddr   string
	GasPrice  uint64
	GasLimit  uint64
	Data      []byte
	Signature string
	ChainID   string
	Version   uint32
	Options   uint32
}

const defaultTimeBetweenBunches = time.Second

type transactionInteractor struct {
	Proxy
	TxSigner
	mutTxAccumulator      sync.RWMutex
	mutTimeBetweenBunches sync.RWMutex
	timeBetweenBunches    time.Duration
	txAccumulator         []*data.Transaction
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
		Proxy:              proxy,
		TxSigner:           txSigner,
		timeBetweenBunches: defaultTimeBetweenBunches,
	}, nil
}

func (ti *transactionInteractor) SetTimeBetweenBunches(timeBetweenBunches time.Duration) {
	ti.mutTimeBetweenBunches.Lock()
	ti.timeBetweenBunches = timeBetweenBunches
	ti.mutTimeBetweenBunches.Unlock()
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

// CreateTransaction assembles a transaction from the provided arguments
func (ti *transactionInteractor) CreateTransaction(arg ArgCreateTransaction) *data.Transaction {
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

// ApplySignatureAndSender will apply the corresponding sender and compute the signature field upon the provided argument
func (ti *transactionInteractor) ApplySignatureAndSender(skBytes []byte, arg ArgCreateTransaction) (ArgCreateTransaction, error) {
	pkBytes, err := ti.TxSigner.GeneratePkBytes(skBytes)
	if err != nil {
		return ArgCreateTransaction{}, err
	}

	copyArg := arg
	copyArg.Signature = ""
	copyArg.SndAddr = core.AddressPublicKeyConverter.Encode(pkBytes)

	unsignedMessage, err := ti.createUnsignedMessage(copyArg)
	if err != nil {
		return ArgCreateTransaction{}, err
	}

	signature, err := ti.TxSigner.SignMessage(unsignedMessage, skBytes)
	if err != nil {
		return ArgCreateTransaction{}, err
	}

	copyArg.Signature = hex.EncodeToString(signature)

	return copyArg, nil
}

func (ti *transactionInteractor) createUnsignedMessage(arg ArgCreateTransaction) ([]byte, error) {
	arg.Signature = ""
	tx := ti.CreateTransaction(arg)

	return json.Marshal(tx)
}

func (ti *transactionInteractor) SendTransactionsAsBunch(bunchSize int) ([]string, error) {
	if bunchSize <= 0 {
		return nil, ErrInvalidValue
	}

	ti.mutTimeBetweenBunches.RLock()
	timeBetweenBunches := ti.timeBetweenBunches
	ti.mutTimeBetweenBunches.RUnlock()

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

		time.Sleep(timeBetweenBunches)
	}

	return allHashes, nil
}
