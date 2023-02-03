package endpointProviders

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	erdgoCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/stretchr/testify/assert"
)

func TestNewProxyEndpointProvider(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.False(t, check.IfNil(provider))
}

func TestProxyEndpointProvider_GetNodeStatus(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.Equal(t, "network/status/0", provider.GetNodeStatus(0))
	assert.Equal(t, "network/status/4294967295", provider.GetNodeStatus(core.MetachainShardId))
}

func TestProxyEndpointProvider_Getters(t *testing.T) {
	t.Parallel()

	provider := NewProxyEndpointProvider()
	assert.Equal(t, "internal/2/raw/block/by-hash/hex", provider.GetRawBlockByHash(2, "hex"))
	assert.Equal(t, "internal/2/raw/block/by-nonce/3", provider.GetRawBlockByNonce(2, 3))
	assert.Equal(t, "internal/2/raw/miniblock/by-hash/hex/epoch/4", provider.GetRawMiniBlockByHash(2, "hex", 4))
	assert.Equal(t, erdgoCore.Proxy, provider.GetRestAPIEntityType())
	assert.False(t, provider.ShouldCheckShardIDForNodeStatus())
}
