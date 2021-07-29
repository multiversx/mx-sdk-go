package core

import (
	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("core")

// AddressPublicKeyConverter represents the default address public key converter
var AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressLen, log)
