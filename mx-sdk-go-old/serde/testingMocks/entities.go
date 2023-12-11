package testingMocks

import "math/big"

// DataBasics defines a data structure that contains the implemented primitives
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

// OtherStruct another structure containing a small sub set of primitives
type OtherStruct struct {
	String string
	Bool   bool
}

// NestingStructure contains a nested structure alongside other primitive types
type NestingStructure struct {
	String        string
	Ticker        string
	Bool          bool
	Int64         int64
	BigInt        big.Int
	OtherStruct   OtherStruct
	AnotherBigInt big.Int
}
