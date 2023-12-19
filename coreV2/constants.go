package coreV2

import (
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
)

const (
	// AddressBytesLen represents the number of bytes of an address
	AddressBytesLen = 32

	// MinAllowedDeltaToFinal is the minimum value between nonces allowed when checking finality on a shard
	MinAllowedDeltaToFinal = 1

	// WebServerOffString represents the constant used to switch off the web server
	WebServerOffString = "off"

	// DefaultAddressPrefix is the default hrp of MultiversX/Elrond
	DefaultAddressPrefix = "erd"

	ScHexPubKeyPrefix = "0000000000000000"
)

var (
	// AddressPublicKeyConverter represents the default address public key converter
	AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, DefaultAddressPrefix)
)

func GetVMTypeWASMVM() []byte {
	return []byte{0x05, 0x00}
}
