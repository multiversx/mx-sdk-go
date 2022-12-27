package builders

import (
	"math/big"

	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxDataBuilder defines the behavior of a transaction data builder
type TxDataBuilder interface {
	Function(function string) TxDataBuilder

	ArgHexString(hexed string) TxDataBuilder
	ArgAddress(address core.AddressHandler) TxDataBuilder
	ArgBigInt(value *big.Int) TxDataBuilder
	ArgInt64(value int64) TxDataBuilder
	ArgBytes(bytes []byte) TxDataBuilder

	ToDataString() (string, error)
	ToDataBytes() ([]byte, error)

	IsInterfaceNil() bool
}

// VMQueryBuilder defines the behavior of a vm query builder
type VMQueryBuilder interface {
	Function(function string) VMQueryBuilder
	CallerAddress(address core.AddressHandler) VMQueryBuilder
	Address(address core.AddressHandler) VMQueryBuilder

	ArgHexString(hexed string) VMQueryBuilder
	ArgAddress(address core.AddressHandler) VMQueryBuilder
	ArgBigInt(value *big.Int) VMQueryBuilder
	ArgInt64(value int64) VMQueryBuilder
	ArgBytes(bytes []byte) VMQueryBuilder

	ToVmValueRequest() (*data.VmValueRequest, error)

	IsInterfaceNil() bool
}

// Signer defines the method used by a struct used to create valid signatures
type Signer interface {
	SignMessage(msg []byte, privateKey crypto.PrivateKey) ([]byte, error)
	VerifyMessage(msg []byte, publicKey crypto.PublicKey, sig []byte) error
	SignTransaction(tx *data.Transaction, privateKey crypto.PrivateKey) ([]byte, error)
	SignByteSlice(msg []byte, privateKey crypto.PrivateKey) ([]byte, error)
	VerifyByteSlice(msg []byte, publicKey crypto.PublicKey, sig []byte) error
	IsInterfaceNil() bool
}
