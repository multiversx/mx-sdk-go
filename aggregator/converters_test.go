package aggregator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 1.0, trim(1, 1))
	assert.Equal(t, 2.0, trim(1.83892672, 1))
	assert.Equal(t, 1.8, trim(1.83892672, 0.1))
	assert.Equal(t, 1.84, trim(1.83892672, 0.01))
	assert.Equal(t, 0.0, trim(1.83892672, 10))
	assert.Equal(t, 10.0, trim(11.83892672, 10))
	assert.Equal(t, 12.0, trim(11.83892672, 1))
}
