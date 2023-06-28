package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueryResponseError(t *testing.T) {
	t.Parallel()

	err := NewQueryResponseError("code", "message", "function", "address", "arg1", "arg2")
	expectedErrorString := "got response code 'code' and message 'message' while querying function 'function' with arguments [arg1 arg2] and address address"

	assert.Equal(t, expectedErrorString, err.Error())
}
