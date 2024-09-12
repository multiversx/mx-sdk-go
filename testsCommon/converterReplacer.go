package testsCommon

import (
	mxChainCore "github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-sdk-go/core"
)

func ReplaceConverter() {
	core.AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(core.AddressBytesLen, mxChainCore.DefaultAddressPrefix)
}
