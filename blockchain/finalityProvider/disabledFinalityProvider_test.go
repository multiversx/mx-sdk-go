package finalityProvider

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewDisabledFinalityProvider(t *testing.T) {
	t.Parallel()

	provider := NewDisabledFinalityProvider()
	assert.False(t, check.IfNil(provider))
}

func TestDisabledFinalityProvider_CheckShardFinalization(t *testing.T) {
	t.Parallel()

	provider := NewDisabledFinalityProvider()
	err := provider.CheckShardFinalization(context.Background(), 0, 0)
	assert.Nil(t, err)
}
