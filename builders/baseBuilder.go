package builders

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/core"
)

type baseBuilder struct {
	args []string
	err  error
}

func (builder *baseBuilder) addBytes(bytes []byte) {
	if len(bytes) == 0 {
		bytes = []byte{0}
	}

	builder.args = append(builder.args, hex.EncodeToString(bytes))
}

func (builder *baseBuilder) checkAddress(address core.AddressHandler) error {
	if check.IfNil(address) {
		return fmt.Errorf("%w in builder.checkAddress", ErrNilAddress)
	}
	if !address.IsValid() {
		return fmt.Errorf("%w in builder.checkAddress", ErrInvalidAddress)
	}

	return nil
}

func (builder *baseBuilder) addArgHexString(hexed string) {
	if builder.err != nil {
		return
	}

	_, err := hex.DecodeString(hexed)
	if err != nil {
		builder.err = fmt.Errorf("%w in builder.ArgHexString for string %s", err, hexed)
		return
	}

	builder.args = append(builder.args, hexed)
}

func (builder *baseBuilder) addArgAddress(address core.AddressHandler) {
	if builder.err != nil {
		return
	}

	err := builder.checkAddress(address)
	if err != nil {
		builder.err = err
		return
	}

	builder.addBytes(address.AddressBytes())
}

func (builder *baseBuilder) addArgBigInt(value *big.Int) {
	if builder.err != nil {
		return
	}

	if value == nil {
		builder.err = fmt.Errorf("%w in builder.ArgBigInt", ErrNilValue)
		return
	}

	builder.addBytes(value.Bytes())
}

func (builder *baseBuilder) addArgInt64(value int64) {
	if builder.err != nil {
		return
	}

	b := big.NewInt(value)

	builder.addBytes(b.Bytes())
}

func (builder *baseBuilder) addArgBytes(bytes []byte) {
	if builder.err != nil {
		return
	}

	if len(bytes) == 0 {
		builder.err = fmt.Errorf("%w in builder.ArgBytes", ErrInvalidValue)
	}

	builder.addBytes(bytes)
}
