package endpointProviders

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewNodeEndpointProvider(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.False(t, check.IfNil(provider))
}

func TestNodeEndpointProvider_GetNodeStatusEndpoint(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.Equal(t, nodeGetNodeStatusEndpoint, provider.GetNodeStatusEndpoint(2))
}

func TestNodeEndpointProvider_Getters(t *testing.T) {
	t.Parallel()

	provider := NewNodeEndpointProvider()
	assert.Equal(t, "internal/raw/block/by-hash/hex", provider.GetRawBlockByHashEndpoint(2, "hex"))
	assert.Equal(t, "internal/raw/block/by-nonce/3", provider.GetRawBlockByNonceEndpoint(2, 3))
	assert.Equal(t, "internal/raw/miniblock/by-hash/hex/epoch/4", provider.GetRawMiniBlockByHashEndpoint(2, "hex", 4))
}
