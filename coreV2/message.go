package coreV2

import (
	"encoding/hex"
	"github.com/multiversx/mx-chain-core-go/hashing/keccak"
	"strconv"
)

type Message struct {
	Data      []byte
	Signature []byte
}

type messageComputer struct{}

type MessageComputer interface {
	ComputeBytesForSigning(message Message) []byte
}

func NewMessageComputer() MessageComputer {
	return &messageComputer{}
}

func (mc *messageComputer) ComputeBytesForSigning(message Message) []byte {
	prefix, _ := hex.DecodeString("17456c726f6e64205369676e6564204d6573736167653a0a")
	msgSize := strconv.FormatInt(int64(len(message.Data)), 10)
	msg := append([]byte(msgSize), message.Data...)
	msg = append(prefix, msg...)

	return keccak.NewKeccak().Compute(string(msg))
}
