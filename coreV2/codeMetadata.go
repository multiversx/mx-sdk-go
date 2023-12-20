package coreV2

import (
	"encoding/hex"
)

const (
	Upgradeable       byte = 1
	Reserved2         byte = 2
	Readable          byte = 4
	Reserved1         byte = 1
	Payable           byte = 2
	PayableByContract      = 4
)

type Option func(metadata *codeMetadata)

type codeMetadata struct {
	Upgradeable       bool
	Readable          bool
	Payable           bool
	PayableByContract bool
}

func NewCodeMetadata(opts ...Option) *codeMetadata {
	cm := codeMetadata{Upgradeable: true, Readable: true, Payable: false, PayableByContract: false}

	for _, opt := range opts {
		opt(&cm)
	}

	return &cm
}

func WithUpgradeable(upgradeable bool) func(metadata *codeMetadata) {
	return func(metadata *codeMetadata) {
		metadata.Upgradeable = upgradeable
	}
}

func WithReadable(readable bool) func(metadata *codeMetadata) {
	return func(metadata *codeMetadata) {
		metadata.Readable = readable
	}
}

func WithPayable(payable bool) func(metadata *codeMetadata) {
	return func(metadata *codeMetadata) {
		metadata.Payable = payable
	}
}

func WithPayableByContract(payableByContract bool) func(metadata *codeMetadata) {
	return func(metadata *codeMetadata) {
		metadata.PayableByContract = payableByContract
	}
}

func (cm *codeMetadata) serialize() []byte {
	data := []byte{0, 0}

	if cm.Upgradeable {
		data[0] |= Upgradeable
	}

	if cm.Readable {
		data[0] |= Readable
	}

	if cm.Payable {
		data[1] |= Payable
	}

	if cm.PayableByContract {
		data[1] |= PayableByContract
	}

	return data
}

func (cm *codeMetadata) String() string {
	return hex.EncodeToString(cm.serialize())
}
