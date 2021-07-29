package blockchain

import (
	"bytes"
	elrondCore "github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

const accountStartNonce = uint64(0)

var initialDNSAddress = bytes.Repeat([]byte{1}, 32)

// addressGenerator is used to generate some addresses based on elrond-go logic
type addressGenerator struct {
	coordinator    *shardCoordinator
	hasher         hashing.Hasher
}

// NewAddressGenerator will create an address generator instance
func NewAddressGenerator(coordinator *shardCoordinator) (*addressGenerator, error) {
	if check.IfNil(coordinator) {
		return nil, ErrNilShardCoordinator
	}

	return &addressGenerator{
		coordinator:    coordinator,
		hasher:         keccak.NewKeccak(),
	}, nil
}

// CompatibleDNSAddress will return the compatible DNS address providing the shard ID
func (ag *addressGenerator) CompatibleDNSAddress(shardId byte) (core.AddressHandler, error) {
	addressLen := len(initialDNSAddress)
	shardInBytes := []byte{0, shardId}

	newDNSPk := string(initialDNSAddress[:(addressLen-elrondCore.ShardIdentiferLen)]) + string(shardInBytes)
	newDNSAddress, err :=  elrondCore.NewAddress([]byte(newDNSPk), core.AddressPublicKeyConverter.Len(), accountStartNonce, factory.ArwenVirtualMachine)
	if err != nil {
		return nil, err
	}

	return data.NewAddressFromBytes(newDNSAddress), err
}

// CompatibleDNSAddressFromUsername will return the compatible DNS address providing the username
func (ag *addressGenerator) CompatibleDNSAddressFromUsername(username string) (core.AddressHandler, error) {
	hash := ag.hasher.Compute(username)
	lastByte := hash[len(hash)-1]
	return ag.CompatibleDNSAddress(lastByte)
}
