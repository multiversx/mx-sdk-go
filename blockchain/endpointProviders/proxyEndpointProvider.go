package endpointProviders

import "fmt"

const (
	proxyGetNodeStatus      = "network/status/%d"
	proxyRawBlockByHash     = "internal/%d/raw/block/by-hash/%s"
	proxyRawBlockByNonce    = "internal/%d/raw/block/by-nonce/%d"
	proxyRawMiniBlockByHash = "internal/%d/raw/miniblock/by-hash/%s/epoch/%d"
)

// proxyEndpointProvider is suitable to work with an Elrond Proxy
type proxyEndpointProvider struct {
	*baseEndpointProvider
}

// NewProxyEndpointProvider returns a new instance of a proxyEndpointProvider
func NewProxyEndpointProvider() *proxyEndpointProvider {
	return &proxyEndpointProvider{}
}

// GetNodeStatus returns the node status endpoint
func (proxy *proxyEndpointProvider) GetNodeStatus(shardID uint32) string {
	return fmt.Sprintf(proxyGetNodeStatus, shardID)
}

// GetRawBlockByHash returns the raw block by hash endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByHash(shardID uint32, hexHash string) string {
	return fmt.Sprintf(proxyRawBlockByHash, shardID, hexHash)
}

// GetRawBlockByNonce returns the raw block by nonce endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByNonce(shardID uint32, nonce uint64) string {
	return fmt.Sprintf(proxyRawBlockByNonce, shardID, nonce)
}

// GetRawMiniBlockByHash returns the raw miniblock by hash endpoint
func (proxy *proxyEndpointProvider) GetRawMiniBlockByHash(shardID uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(proxyRawMiniBlockByHash, shardID, hexHash, epoch)
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *proxyEndpointProvider) IsInterfaceNil() bool {
	return proxy == nil
}
