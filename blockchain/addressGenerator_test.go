package blockchain

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/process/smartContract/hooks"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_GetShard(t *testing.T) {
	t.Parallel()

	pubkey := make([]byte, 32)
	_, _ = rand.Read(pubkey)

	numShardsWithoutMeta := uint32(2)
	shardCoordinator, _ := NewShardCoordinator(numShardsWithoutMeta, 0)

	pubkey[31] &= 0xFE
	addr0 := data.NewAddressFromBytes(pubkey)

	pubkey[31] |= 0x01
	addr1 := data.NewAddressFromBytes(pubkey)

	sh0, err := shardCoordinator.ComputeShardId(addr0)
	assert.Nil(t, err)

	sh1, err := shardCoordinator.ComputeShardId(addr1)
	assert.Nil(t, err)

	assert.Equal(t, sh0, uint32(0))
	assert.Equal(t, sh1, uint32(1))
}

func TestGenerateSameDNSAddress(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	args := ArgsAddressGenerator{
		Coordinator: coord,
		PubkeyConv:  core.AddressPublicKeyConverter,
	}
	ag, err := NewAddressGenerator(args)
	require.Nil(t, err)

	newDNS, err := ag.CompatibleDNSAddressFromUsername("laura.elrond")
	require.Nil(t, err)

	fmt.Printf("Compatibile DNS address is %s\n", newDNS.AddressAsBech32String())
	assert.Equal(t, "erd1qqqqqqqqqqqqqpgqvrsdh798pvd4x09x0argyscxc9h7lzfhqz4sttlatg", newDNS.AddressAsBech32String())
}

func TestAddressGenerator_ComputeArwenScAddress(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	args := ArgsAddressGenerator{
		Coordinator: coord,
		PubkeyConv:  core.AddressPublicKeyConverter,
	}
	ag, err := NewAddressGenerator(args)
	require.Nil(t, err)
	owner, err := data.NewAddressFromBech32String("erd1dglncxk6sl9a3xumj78n6z2xux4ghp5c92cstv5zsn56tjgtdwpsk46qrs")
	require.Nil(t, err)

	scAddress, err := ag.ComputeArwenScAddress(owner, 10)
	require.Nil(t, err)

	assert.Equal(t, "erd1qqqqqqqqqqqqqpgqxcy5fma93yhw44xcmt3zwrl0tlhaqmxrdwpsr2vh8p", scAddress.AddressAsBech32String())
}

func TestBlockChainHookImpl_NewAddressLengthNoGood(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	args := ArgsAddressGenerator{
		Coordinator: coord,
		PubkeyConv:  core.AddressPublicKeyConverter,
	}
	ag, err := NewAddressGenerator(args)
	require.Nil(t, err)

	address := []byte("test")
	nonce := uint64(10)

	scAddress, err := ag.NewAddress(address, nonce, []byte("00"))
	assert.Equal(t, hooks.ErrAddressLengthNotCorrect, err)
	assert.Nil(t, scAddress)

	address = []byte("1234567890123456789012345678901234567890")
	scAddress, err = ag.NewAddress(address, nonce, []byte("00"))
	assert.Equal(t, hooks.ErrAddressLengthNotCorrect, err)
	assert.Nil(t, scAddress)
}

func TestBlockChainHookImpl_NewAddressVMTypeTooLong(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	args := ArgsAddressGenerator{
		Coordinator: coord,
		PubkeyConv:  core.AddressPublicKeyConverter,
	}
	ag, err := NewAddressGenerator(args)
	require.Nil(t, err)

	address := []byte("01234567890123456789012345678900")
	nonce := uint64(10)

	vmType := []byte("010")
	scAddress, err := ag.NewAddress(address, nonce, vmType)
	assert.Equal(t, hooks.ErrVMTypeLengthIsNotCorrect, err)
	assert.Nil(t, scAddress)
}

func TestBlockChainHookImpl_NewAddress(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	args := ArgsAddressGenerator{
		Coordinator: coord,
		PubkeyConv:  core.AddressPublicKeyConverter,
	}
	ag, err := NewAddressGenerator(args)
	require.Nil(t, err)

	address := []byte("01234567890123456789012345678900")
	nonce := uint64(10)

	vmType := []byte("11")
	scAddress1, err := ag.NewAddress(address, nonce, vmType)
	assert.Nil(t, err)

	for i := 0; i < 8; i++ {
		assert.Equal(t, scAddress1[i], uint8(0))
	}
	assert.True(t, bytes.Equal(vmType, scAddress1[8:10]))

	nonce++
	scAddress2, err := ag.NewAddress(address, nonce, []byte("00"))
	assert.Nil(t, err)

	assert.False(t, bytes.Equal(scAddress1, scAddress2))

	fmt.Printf("%s \n%s \n", hex.EncodeToString(scAddress1), hex.EncodeToString(scAddress2))
}
