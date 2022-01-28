package interactors

import (
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHeaderCheckHandler(t *testing.T) {
	t.Parallel()

	hch, err := NewHeaderCheckHandler(nil)
	require.Nil(t, hch)
	assert.Equal(t, ErrNilProxy, err)

	hch, err = NewHeaderCheckHandler(&testsCommon.ProxyStub{})
	require.NotNil(t, hch)
	require.Nil(t, err)
}
