package builders

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVMQueryBuilder(t *testing.T) {
	t.Parallel()

	builder := NewVMQueryBuilder()
	assert.False(t, check.IfNil(builder))
	_, err := builder.ToVmValueRequest()
	assert.Nil(t, err)
}

func TestTxDataBuilder_Address(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.Address(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("invalid address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.Address(data.NewAddressFromBytes(make([]byte, 0)))
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrInvalidAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("should work", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.Address(address)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)

		addressAsBech32String, err := address.AddressAsBech32String()
		assert.Nil(t, err)
		assert.Equal(t, addressAsBech32String, valueRequest.Address)
	})
}

func TestTxDataBuilder_CallerAddress(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	t.Run("nil address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.CallerAddress(nil)
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrNilAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("invalid address should contain error", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.CallerAddress(data.NewAddressFromBytes(make([]byte, 0)))
		valueRequest, err := builder.ToVmValueRequest()
		assert.True(t, errors.Is(err, ErrInvalidAddress))
		assert.Nil(t, valueRequest)
	})
	t.Run("should work", func(t *testing.T) {
		builder := NewVMQueryBuilder()
		builder.CallerAddress(address)
		valueRequest, err := builder.ToVmValueRequest()
		assert.Nil(t, err)
		addressAsBech32String, err := address.AddressAsBech32String()
		assert.Nil(t, err)
		assert.Equal(t, addressAsBech32String, valueRequest.CallerAddr)
	})
}

func TestVmQueryBuilder_AllGoodArguments(t *testing.T) {
	t.Parallel()

	address, errBech32 := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, errBech32)

	builder := NewVMQueryBuilder().
		Function("function").
		ArgBigInt(big.NewInt(15)).
		ArgInt64(14).
		ArgAddress(address).
		ArgHexString("eeff00").
		ArgBytes([]byte("aa")).
		ArgBigInt(big.NewInt(0))

	valueRequest, err := builder.ToVmValueRequest()
	assert.Nil(t, err)
	assert.Equal(t, "function", valueRequest.FuncName)

	expectedArgs := []string{
		hex.EncodeToString([]byte{15}),
		hex.EncodeToString([]byte{14}),
		hex.EncodeToString(address.AddressBytes()),
		"eeff00",
		hex.EncodeToString([]byte("aa")),
		"00",
	}

	require.Equal(t, expectedArgs, valueRequest.Args)
}
