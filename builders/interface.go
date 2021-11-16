package builders

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxDataBuilder defines the behavior of a transaction data builder
type TxDataBuilder interface {
	Function(function string) TxDataBuilder
	CallerAddress(address core.AddressHandler) TxDataBuilder
	Address(address core.AddressHandler) TxDataBuilder

	ArgHexString(hexed string) TxDataBuilder
	ArgAddress(address core.AddressHandler) TxDataBuilder
	ArgBigInt(value *big.Int) TxDataBuilder
	ArgInt64(value int64) TxDataBuilder
	ArgBytes(bytes []byte) TxDataBuilder

	ToDataString() (string, error)
	ToDataBytes() ([]byte, error)
	ToVmValueRequest() (*data.VmValueRequest, error)
}
