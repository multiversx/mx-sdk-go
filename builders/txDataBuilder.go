package elrond

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

const dataSeparator = "@"

// generalTxDataBuilder can be used to easy construct a transaction's data field for a smart contract call
// can also be used to construct a VmValueRequest instance ready to be used on a VM query
type generalTxDataBuilder struct {
	address    string
	function   string
	callerAddr string
	args       []string
	log        logger.Logger
	err        error
}

// NewGeneralTxDataBuilder creates a new general transaction data builder
func NewGeneralTxDataBuilder(log logger.Logger) *generalTxDataBuilder {
	builder := &generalTxDataBuilder{
		log: log,
	}
	if check.IfNil(log) {
		builder.err = ErrNilLogger
	}

	return builder
}

// Function sets the function to be called
func (builder *generalTxDataBuilder) Function(function string) TxDataBuilder {
	builder.function = function

	return builder
}

// CallerAddress sets the caller address
func (builder *generalTxDataBuilder) CallerAddress(address core.AddressHandler) TxDataBuilder {
	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return builder
	}

	builder.callerAddr = address.AddressAsBech32String()

	return builder
}

// Address sets the destination address
func (builder *generalTxDataBuilder) Address(address core.AddressHandler) TxDataBuilder {
	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return builder
	}

	builder.address = address.AddressAsBech32String()

	return builder
}

// ArgHexString adds the provided hex string to the arguments list
func (builder *generalTxDataBuilder) ArgHexString(hexed string) TxDataBuilder {
	_, err := hex.DecodeString(hexed)
	if err != nil {
		builder.err = fmt.Errorf("%w in builder.ArgHexString for string %s", err, hexed)
		return builder
	}

	builder.args = append(builder.args, hexed)

	return builder
}

// ArgAddress adds the provided address to the arguments list
func (builder *generalTxDataBuilder) ArgAddress(address core.AddressHandler) TxDataBuilder {
	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return builder
	}

	return builder.addBytes(address.AddressBytes())
}

func (builder *generalTxDataBuilder) checkAddress(address core.AddressHandler) error {
	if check.IfNil(address) {
		return fmt.Errorf("%w in builder.checkAddress", ErrNilAddress)
	}
	if len(address.AddressBytes()) == 0 {
		return fmt.Errorf("%w in builder.checkAddress", ErrInvalidAddress)
	}

	return nil
}

// ArgBigInt adds the provided value to the arguments list
func (builder *generalTxDataBuilder) ArgBigInt(value *big.Int) TxDataBuilder {
	if value == nil {
		builder.err = fmt.Errorf("%w in builder.ArgBigInt", ErrNilValue)
		return builder
	}

	return builder.addBytes(value.Bytes())
}

// ArgInt64 adds the provided value to the arguments list
func (builder *generalTxDataBuilder) ArgInt64(value int64) TxDataBuilder {
	b := big.NewInt(value)

	return builder.addBytes(b.Bytes())
}

// ArgBytes adds the provided bytes to the arguments list. The parameter should contain at least one byte
func (builder *generalTxDataBuilder) ArgBytes(bytes []byte) TxDataBuilder {
	if len(bytes) == 0 {
		builder.err = fmt.Errorf("%w in builder.ArgBytes", ErrInvalidValue)
	}

	builder.args = append(builder.args, hex.EncodeToString(bytes))

	return builder
}

func (builder *generalTxDataBuilder) addBytes(bytes []byte) TxDataBuilder {
	if len(bytes) == 0 {
		bytes = []byte{0}
	}

	builder.args = append(builder.args, hex.EncodeToString(bytes))

	return builder
}

// ToDataString returns the formatted data string ready to be used in a transaction call
func (builder *generalTxDataBuilder) ToDataString() (string, error) {
	if builder.err != nil {
		return "", builder.err
	}

	parts := append([]string{builder.function}, builder.args...)

	return strings.Join(parts, dataSeparator), nil
}

// ToDataBytes returns the formatted data string ready to be used in a transaction call as bytes
func (builder *generalTxDataBuilder) ToDataBytes() ([]byte, error) {
	dataField, err := builder.ToDataString()
	if err != nil {
		return nil, err
	}

	return []byte(dataField), err
}

// ToVmValueRequest returns the VmValueRequest structure to be used in a VM call
func (builder *generalTxDataBuilder) ToVmValueRequest() (*data.VmValueRequest, error) {
	if builder.err != nil {
		return nil, builder.err
	}

	return &data.VmValueRequest{
		Address:    builder.address,
		FuncName:   builder.function,
		CallerAddr: builder.callerAddr,
		Args:       builder.args,
	}, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (builder *generalTxDataBuilder) IsInterfaceNil() bool {
	return builder == nil
}
