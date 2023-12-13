package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"math/big"
)

var (
	blake2bHasher          = blake2b.NewBlake2b()
	nodeInternalMarshaller = &marshal.GogoProtoMarshalizer{}
)

type Amount = big.Int

type Transaction struct {
	Sender   string `json:"sender,omitempty"`
	Receiver string `json:"receiver,omitempty"`
	GasLimit uint64 `json:"gas_limit,omitempty"`
	ChainID  string `json:"chain_id,omitempty"`

	Nonce            uint64  `json:"nonce,omitempty"`
	Value            *Amount `json:"value"`
	SenderUsername   string  `json:"sender_username,omitempty"`
	ReceiverUsername string  `json:"receiver_username,omitempty"`
	GasPrice         uint64  `json:"gas_price,omitempty"`

	Data     []byte `json:"data,omitempty"`
	Version  uint32 `json:"version,omitempty"`
	Options  uint32 `json:"options,omitempty"`
	Guardian string `json:"guardian,omitempty"`

	Signature         []byte `json:"signature,omitempty"`
	GuardianSignature []byte `json:"guardian_signature,omitempty"`
}

type TransactionComputer interface {
	ComputeTransactionFee(transaction Transaction, networkConfig NetworkConfig) (*Amount, error)
	ComputeBytesForSigning(transaction Transaction) ([]byte, error)
	ComputeTransactionHash(transaction Transaction) ([]byte, error)
}

type transactionComputer struct {
}

func NewTransactionComputer() TransactionComputer {
	return &transactionComputer{}
}

func (tc *transactionComputer) ComputeTransactionFee(transaction Transaction, networkConfig NetworkConfig) (*Amount, error) {
	moveBalanceGas := networkConfig.MinGasLimit + len(transaction.Data)*networkConfig.GasPerDataByte

	if moveBalanceGas > int(transaction.GasLimit) {
		return nil, fmt.Errorf("not enough gas provided: %q", transaction.GasLimit)
	}

	feeForMove := moveBalanceGas * int(transaction.GasPrice)
	if moveBalanceGas == int(transaction.GasLimit) {

		return big.NewInt(int64(feeForMove)), nil
	}

	diff := int(transaction.GasLimit) - moveBalanceGas
	modifiedGasPrice := float32(transaction.GasPrice) * networkConfig.GasPriceModifier
	processingFee := float32(diff) * modifiedGasPrice

	return big.NewInt(int64(feeForMove) + int64(processingFee)), nil
}

func (tc *transactionComputer) ComputeBytesForSigning(transaction Transaction) ([]byte, error) {
	bytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	return bytes, nil
}

func (tc *transactionComputer) ComputeTransactionHash(transaction Transaction) ([]byte, error) {
	if len(transaction.Signature) == 0 {
		return nil, errors.New("transaction is missing signature")
	}

	nodeTx, err := transactionToNodeTransaction(transaction)
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

func transactionToNodeTransaction(tx Transaction) (*transaction.Transaction, error) {
	receiverBytes, err := AddressPublicKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver address: %v", err)
	}

	senderBytes, err := AddressPublicKeyConverter.Decode(tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sender address: %v", err)
	}

	return &transaction.Transaction{
		Nonce:     tx.Nonce,
		Value:     tx.Value,
		RcvAddr:   receiverBytes,
		SndAddr:   senderBytes,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      tx.Data,
		ChainID:   []byte(tx.ChainID),
		Version:   tx.Version,
		Signature: tx.Signature,
		Options:   tx.Options,
	}, nil
}
