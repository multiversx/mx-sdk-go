package blockchain

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedErr        = fmt.Errorf("expected error")
	expectedAddress, _ = data.NewAddressFromBech32String("erd1qqqqqqqqqqqqqpgqxcy5fma93yhw44xcmt3zwrl0tlhaqmxrdwpsr2vh8p")
)

func TestNewAddressGenerator(t *testing.T) {
	t.Parallel()

	t.Run("nil pubkey converter", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsAddressGenerator()
		args.PubkeyConv = nil

		ag, err := NewAddressGenerator(args)

		assert.Nil(t, ag)
		assert.Equal(t, process.ErrNilPubkeyConverter, err)
	})
	t.Run("nil address generator core", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = nil

		ag, err := NewAddressGenerator(args)

		assert.Nil(t, ag)
		assert.Equal(t, ErrNilAddressGenerator, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createMockArgsAddressGenerator()
		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)
		require.NotNil(t, ag)
	})
}

func TestGenerateSameDNSAddress(t *testing.T) {
	t.Parallel()

	t.Run("New address errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = &mock.AddressGeneratorStub{
			NewAddressCalled: func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
				return nil, expectedErr
			},
		}
		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)

		newDNS, err := ag.CompatibleDNSAddressFromUsername("laura.elrond")
		require.Nil(t, newDNS)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = &mock.AddressGeneratorStub{
			NewAddressCalled: func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
				return expectedAddress.AddressBytes(), nil
			},
		}

		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)

		newDNS, err := ag.CompatibleDNSAddressFromUsername("laura.elrond")
		require.Nil(t, err)

		fmt.Printf("Compatibile DNS address is %s\n", newDNS.AddressAsBech32String())
		assert.Equal(t, expectedAddress, newDNS)
	})
}

func TestAddressGenerator_ComputeArwenScAddress(t *testing.T) {
	t.Parallel()

	t.Run("New address errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = &mock.AddressGeneratorStub{
			NewAddressCalled: func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
				return nil, expectedErr
			},
		}
		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)

		scAddress, err := ag.ComputeArwenScAddress(expectedAddress, 10)
		require.Nil(t, scAddress)

		assert.Equal(t, expectedErr, err)
	})
	t.Run("nil address should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = &mock.AddressGeneratorStub{
			NewAddressCalled: func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
				return nil, nil
			},
		}
		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)

		scAddress, err := ag.ComputeArwenScAddress(nil, 10)
		require.Nil(t, scAddress)

		assert.Equal(t, ErrNilAddress, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createMockArgsAddressGenerator()
		args.AddressGeneratorCore = &mock.AddressGeneratorStub{
			NewAddressCalled: func(address []byte, nonce uint64, vmType []byte) ([]byte, error) {
				return expectedAddress.AddressBytes(), nil
			},
		}
		ag, err := NewAddressGenerator(args)
		require.Nil(t, err)
		owner, err := data.NewAddressFromBech32String("erd1dglncxk6sl9a3xumj78n6z2xux4ghp5c92cstv5zsn56tjgtdwpsk46qrs")
		require.Nil(t, err)

		scAddress, err := ag.ComputeArwenScAddress(owner, 10)
		require.Nil(t, err)

		assert.Equal(t, expectedAddress, scAddress)
	})
}

func createMockArgsAddressGenerator() ArgsAddressGenerator {
	return ArgsAddressGenerator{
		PubkeyConv:           core.AddressPublicKeyConverter,
		AddressGeneratorCore: &mock.AddressGeneratorStub{},
	}
}
