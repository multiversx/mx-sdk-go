package endpointProviders

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewProxyEndpointProvider(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.False(t, check.IfNil(provider))
}

func TestProxyEndpointProvider_GetNodeStatusEndpoint(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.Equal(t, "network/status/0", provider.GetNodeStatusEndpoint(0))
	assert.Equal(t, "network/status/4294967295", provider.GetNodeStatusEndpoint(core.MetachainShardId))
}

func TestProxyEndpointProvider_Getters(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.Equal(t, "internal/2/raw/block/by-hash/hex", provider.GetRawBlockByHashEndpoint(2, "hex"))
	assert.Equal(t, "internal/2/raw/block/by-nonce/3", provider.GetRawBlockByNonceEndpoint(2, 3))
	assert.Equal(t, "internal/2/raw/miniblock/by-hash/hex/epoch/4", provider.GetRawMiniBlockByHashEndpoint(2, "hex", 4))
}
