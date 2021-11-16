package elrond

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = logger.GetOrCreate("test")

func TestNewTxDataBuilder(t *testing.T) {
	t.Parallel()

	t.Run("without logger should contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(nil)
		assert.False(t, check.IfNil(builder))
		_, err := builder.ToDataString()
		assert.Equal(t, ErrNilLogger, err)
	})
	t.Run("with logger should not contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		assert.False(t, check.IfNil(builder))
		_, err := builder.ToDataString()
		assert.Nil(t, err)
	})
}

func TestTxDataBuilder_Address(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.Address(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("invalid address should contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.Address(data.NewAddressFromBytes(make([]byte, 0)))
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrInvalidAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("should work", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.Address(address)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)
		assert.Equal(t, address.AddressAsBech32String(), valueRequest.Address)
	})
}

func TestTxDataBuilder_CallerAddress(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.CallerAddress(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("invalid address should contain error", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.CallerAddress(data.NewAddressFromBytes(make([]byte, 0)))
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrInvalidAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("should work", func(t *testing.T) {
		builder := NewTxDataBuilder(log)
		builder.CallerAddress(address)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)
		assert.Equal(t, address.AddressAsBech32String(), valueRequest.CallerAddr)
	})
}

func TestTxDataBuilder_Function(t *testing.T) {
	t.Parallel()

	function := "sc call function"
	builder := NewTxDataBuilder(log)
	builder.Function(function)

	valueRequest, err := builder.ToVmValueRequest()
	assert.Nil(t, err)
	assert.Equal(t, function, valueRequest.FuncName)

	txData, err := builder.ToDataString()
	assert.Nil(t, err)
	assert.Equal(t, function, txData)
}

func TestTxDataBuilder_AllGoodArguments(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	builder := NewTxDataBuilder(log).
		Function("function").
		ArgBigInt(big.NewInt(15)).
		ArgInt64(14).
		ArgAddress(address).
		ArgHexString("eeff00").
		ArgBytes([]byte("aa")).
		ArgBigInt(big.NewInt(0))

	expectedTxData := "function@" + hex.EncodeToString([]byte{15}) +
		"@" + hex.EncodeToString([]byte{14}) + "@" +
		hex.EncodeToString(address.AddressBytes()) + "@eeff00@" +
		hex.EncodeToString([]byte("aa")) + "@00"

	txData, err := builder.ToDataString()
	require.Nil(t, err)
	require.Equal(t, expectedTxData, txData)

	txDataBytes, err := builder.ToDataBytes()
	require.Nil(t, err)

	require.Equal(t, []byte(expectedTxData), txDataBytes)
}

func TestTxDataBuilder_InvalidArguments(t *testing.T) {
	t.Parallel()

	builder := NewTxDataBuilder(log).
		Function("function").
		ArgInt64(4)

	t.Run("invalid hex string argument", func(t *testing.T) {
		builder.ArgHexString("invalid hex string")

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, strings.Contains(errString.Error(), "builder.ArgHexString for string"))
	})
	t.Run("nil address argument", func(t *testing.T) {
		builder.ArgAddress(nil)

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrNilAddress))
	})
	t.Run("nil big int argument", func(t *testing.T) {
		builder.ArgBigInt(nil)

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrNilValue))
	})
	t.Run("empty bytes argument", func(t *testing.T) {
		builder.ArgBytes(make([]byte, 0))

		txData, errString := builder.ToDataString()
		txDataBytes, errBytes := builder.ToDataBytes()
		assert.Equal(t, errString, errBytes)
		assert.Equal(t, "", txData)
		assert.Nil(t, txDataBytes)
		assert.NotNil(t, errString)
		assert.True(t, errors.Is(errString, ErrInvalidValue))
	})

}
