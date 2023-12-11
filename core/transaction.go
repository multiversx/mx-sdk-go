package core

import "math/big"

type Amount = big.Int

type Transaction struct {
	Sender   string
	Receiver string
	GasLimit string
	ChainID  string

	Nonce            uint64
	Value            Amount
	SenderUsername   string
	ReceiverUsername string
	GasPrice         uint32

	Data     []byte
	Version  uint32
	Options  uint32
	Guardian string

	Signature         []byte
	GuardianSignature []byte
}
