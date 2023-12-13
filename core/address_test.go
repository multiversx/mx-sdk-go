package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAddressFromBech32(t *testing.T) {
	value := "erd104ct3knz45kt2s7ap93haqz3nrahzgf35zjpgq28fl99qrlepd6sqx6nf4"

	a, err := NewAddressFromBech32(value)
	require.NoError(t, err)

	a2, err := NewAddressFromHex(a.ToHex(), "erd")
	require.NoError(t, err)

	require.Equal(t, a, a2)
}

func TestAddressWithCustomHRP(t *testing.T) {
	a, err := NewAddressFromHex("0139472eff6886771a982f3083da5d421f24c29181e63888228dc81ca60d69e1", "test")
	require.NoError(t, err, "failed to create address from hex")

	bech32, err := a.ToBech32()
	require.NoError(t, err, "failed to retrieve bech32 from address")
	a1, err := NewAddressFromBech32(bech32)
	require.NoError(t, err, "failed to create address from bech32")

	a2, err := NewAddressFromHex(a.ToHex(), "test")
	require.NoError(t, err)

	require.Equal(t, a, a1, a2)
}

func TestAddress_IsSmartContract(t *testing.T) {

}
