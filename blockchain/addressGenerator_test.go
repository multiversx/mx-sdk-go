package blockchain

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_GetShard(t *testing.T) {
	t.Parallel()

	pubkey := make([]byte, 32)
	_, _ = rand.Read(pubkey)

	numShardsWithoutMeta := uint32(2)
	shardCoordinatorInstance, _ := NewShardCoordinator(numShardsWithoutMeta, 0)

	pubkey[31] &= 0xFE
	addr0 := data.NewAddressFromBytes(pubkey)

	pubkey[31] |= 0x01
	addr1 := data.NewAddressFromBytes(pubkey)

	sh0, err := shardCoordinatorInstance.ComputeShardId(addr0)
	assert.Nil(t, err)

	sh1, err := shardCoordinatorInstance.ComputeShardId(addr1)
	assert.Nil(t, err)

	assert.Equal(t, sh0, uint32(0))
	assert.Equal(t, sh1, uint32(1))
}

func TestGenerateSameDNSAddress(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	ag, err := NewAddressGenerator(coord)
	require.Nil(t, err)

	newDNS, err := ag.CompatibleDNSAddressFromUsername("laura.elrond")
	require.Nil(t, err)

	newDnsAsBech32, err := newDNS.AddressAsBech32String()
	require.Nil(t, err)

	fmt.Printf("Compatibile DNS address is %s\n", newDnsAsBech32)
	assert.Equal(t, "erd1qqqqqqqqqqqqqpgqvrsdh798pvd4x09x0argyscxc9h7lzfhqz4sttlatg", newDnsAsBech32)
}

func TestAddressGenerator_ComputeArwenScAddress(t *testing.T) {
	t.Parallel()

	coord, err := NewShardCoordinator(3, 0)
	require.Nil(t, err)

	ag, err := NewAddressGenerator(coord)
	require.Nil(t, err)
	owner, err := data.NewAddressFromBech32String("erd1dglncxk6sl9a3xumj78n6z2xux4ghp5c92cstv5zsn56tjgtdwpsk46qrs")
	require.Nil(t, err)

	scAddress, err := ag.ComputeArwenScAddress(owner, 10)
	require.Nil(t, err)

	scAddressAsBech32, err := scAddress.AddressAsBech32String()
	require.Nil(t, err)

	assert.Equal(t, "erd1qqqqqqqqqqqqqpgqxcy5fma93yhw44xcmt3zwrl0tlhaqmxrdwpsr2vh8p", scAddressAsBech32)
}
