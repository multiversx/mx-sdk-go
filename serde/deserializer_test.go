package serde

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/serde/testingMocks"
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

	_, err = ds.CreateStruct(bt, data)

	assert.Nil(t, err)
}

func TestDeserializer_CreateStruct_NestedStructures(t *testing.T) {
	t.Parallel()

	data, err := ioutil.ReadFile(srcNestedStructures)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	ds := NewDeserializer()

	bt := &testingMocks.NestedStructure{}

	_, err = ds.CreateStruct(bt, data)

	assert.Nil(t, err)
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
}
