package testingMocks

import "math/big"

type DataBasics struct {
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64

	I8  int8
	I16 int16
	I32 int32
	I64 int64

	Bool bool

	BoxedBytes      string
	TokenIdentifier string
	BigInt          big.Int
	BigUint         big.Int
}
