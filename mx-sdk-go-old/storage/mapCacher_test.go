package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeInBytesContained(t *testing.T) {
	cacher := NewMapCacher()
	cacher.Put([]byte("key"), "value", 0)

	assert.Equal(t, uint64(9), cacher.SizeInBytesContained())
}
