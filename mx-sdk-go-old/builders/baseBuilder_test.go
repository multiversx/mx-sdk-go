package builders

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseBuilder_AnErrorWillPreventProcessingOfTheArguments(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	b := &baseBuilder{
		err: expectedErr,
	}

	b.addArgBytes(nil)
	b.addArgAddress(nil)
	b.addArgBigInt(nil)
	b.addArgHexString("not hexed")
	b.addArgInt64(0)

	assert.Equal(t, expectedErr, b.err)
	assert.Equal(t, 0, len(b.args))
}
