package core

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
)

// AddressPublicKeyConverter represents the default address public key converter
var AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, core.DefaultAddressPrefix)
