package serde

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/multiversx/mx-sdk-go/serde/testingMocks"
	"github.com/stretchr/testify/assert"
)

const (
	srcBasicTypes       = "./testingMocks/basicDataTypes"
	srcNestedStructures = "./testingMocks/nestedStructures"
	srcPrimitive        = "./testingMocks/basicTypeBigInt"
)

func TestDeserializer_CreateStruct_BasicTypes(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(srcBasicTypes)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bt := &testingMocks.DataBasics{}

	usedBytes, err := ds.CreateStruct(bt, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(len(data)), usedBytes)

	bigInt := big.Int{}
	bigInt.SetBytes([]byte{23, 0, 0, 0})
	expectedStruct := &testingMocks.DataBasics{
		U8:              8,
		U16:             16,
		U32:             32,
		U64:             64,
		I8:              8,
		I16:             16,
		I32:             32,
		I64:             64,
		Bool:            true,
		BoxedBytes:      "BoxedBytes",
		TokenIdentifier: "ALC-6258d2",
		BigInt:          bigInt,
		BigUint:         bigInt,
	}

	assert.EqualValues(t, expectedStruct, bt)

}

func TestDeserializer_CreateStruct_NestedStructures(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(srcNestedStructures)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	nesting := &testingMocks.NestingStructure{}

	usedBytes, err := ds.CreateStruct(nesting, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(len(data)), usedBytes)

	bigInt := big.Int{}
	bigInt.SetBytes([]byte{23, 0, 0, 0})
	expectedStruct := &testingMocks.NestingStructure{
		String: "BoxedBytes",
		Ticker: "ALC-6258d2",
		Bool:   false,
		Int64:  385875968,
		BigInt: bigInt,
		OtherStruct: testingMocks.OtherStruct{
			String: "BoxedBytes",
			Bool:   true,
		},
		AnotherBigInt: bigInt,
	}

	assert.EqualValues(t, expectedStruct, nesting)
}

func Test_deserializer_CreatePrimitiveDataType(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(srcPrimitive)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bigInt := &big.Int{}

	err = ds.CreatePrimitiveDataType(bigInt, data)
	assert.Nil(t, err)

	expected := &big.Int{}
	expected.SetBytes(data)

	assert.EqualValues(t, expected, bigInt)
}
