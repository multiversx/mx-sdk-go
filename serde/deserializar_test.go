package serde

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/serde/testingMocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

const (
	SRC_BASIC_TYPES       = "./testingMocks/basicDataTypes"
	SRC_NESTED_STRUCTURES = "./testingMocks/nestedStructures"
	SRC_PRIMITIVE         = "./testingMocks/basicTypeBigInt"
)

func TestDeserializer_CreateStruct_BasicTypes(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(SRC_BASIC_TYPES)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bt := &testingMocks.DataBasics{}

	_, err = ds.CreateStruct(bt, data)

	assert.Nil(t, err)
}

func TestDeserializer_CreateStruct_NestedStructures(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(SRC_NESTED_STRUCTURES)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bt := &testingMocks.NestedStructure{}

	_, err = ds.CreateStruct(bt, data)

	assert.Nil(t, err)
}

func Test_deserializer_CreatePrimitiveDataType(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(SRC_PRIMITIVE)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bigInt := &big.Int{}

	err = ds.CreatePrimitiveDataType(bigInt, data)

	assert.Nil(t, err)
}
