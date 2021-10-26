package serde

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/serde/testingMocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const (
	SRC_BASIC_TYPES       = "./testingMocks/basicDataTypes"
	SRC_NESTED_STRUCTURES = "./testingMocks/nestedStructures"
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
