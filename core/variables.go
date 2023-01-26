package core

import (
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/core")

// AddressPublicKeyConverter represents the default address public key converter
var AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, log)
