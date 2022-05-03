package endpointProviders

import "fmt"

const (
	proxyGetNodeStatusEndpoint      = "network/status/%d"
	proxyRawBlockByHashEndpoint     = "internal/%d/raw/block/by-hash/%s"
	proxyRawBlockByNonceEndpoint    = "internal/%d/raw/block/by-nonce/%d"
	proxyRawMiniBlockByHashEndpoint = "internal/%d/raw/miniblock/by-hash/%s/epoch/%d"
)

// proxyEndpointProvider is suitable to work with an Elrond Proxy
type proxyEndpointProvider struct {
	*baseEndpointProvider
}

// NewProxyEndpointProvider returns a new instance of a proxyEndpointProvider
func NewProxyEndpointProvider() *proxyEndpointProvider {
	return &proxyEndpointProvider{}
}

// GetNodeStatusEndpoint returns the node status endpoint
func (proxy *proxyEndpointProvider) GetNodeStatusEndpoint(shardID uint32) string {
	return fmt.Sprintf(proxyGetNodeStatusEndpoint, shardID)
}

// GetRawBlockByHashEndpoint returns the raw block by hash endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByHashEndpoint(shardID uint32, hexHash string) string {
	return fmt.Sprintf(proxyRawBlockByHashEndpoint, shardID, hexHash)
}

// GetRawBlockByNonceEndpoint returns the raw block by nonce endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByNonceEndpoint(shardID uint32, nonce uint64) string {
	return fmt.Sprintf(proxyRawBlockByNonceEndpoint, shardID, nonce)
}

// GetRawMiniBlockByHashEndpoint returns the raw miniblock by hash endpoint
func (proxy *proxyEndpointProvider) GetRawMiniBlockByHashEndpoint(shardID uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(proxyRawMiniBlockByHashEndpoint, shardID, hexHash, epoch)
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *proxyEndpointProvider) IsInterfaceNil() bool {
	return proxy == nil
}
