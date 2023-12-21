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
	owner, err := data.NewAddressFromBech32String("erd1uzk2g5rhvg8prk9y50d0q7qsxg7tm7f320q0q4qlpmfu395wjmdqqy0n9q")
	require.Nil(t, err)

	for i := 0; i < 10; i++ {
		scAddress, err := ag.ComputeArwenScAddress(owner, uint64(i))
		require.Nil(t, err)

		scAddressAsBech32, err := scAddress.AddressAsBech32String()
		require.Nil(t, err)

		fmt.Println(scAddressAsBech32)
	}
}
