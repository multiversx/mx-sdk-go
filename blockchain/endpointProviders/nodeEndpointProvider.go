package endpointProviders

import "fmt"

const (
	nodeGetNodeStatusEndpoint      = "node/status"
	nodeRawBlockByHashEndpoint     = "internal/raw/block/by-hash/%s"
	nodeRawBlockByNonceEndpoint    = "internal/raw/block/by-nonce/%d"
	nodeRawMiniBlockByHashEndpoint = "internal/raw/miniblock/by-hash/%s/epoch/%d"
)

// proxyEndpointProvider is suitable to work with an Elrond node
type nodeEndpointProvider struct {
	*baseEndpointProvider
}

// NewNodeEndpointProvider returns a new instance of a nodeEndpointProvider
func NewNodeEndpointProvider() *nodeEndpointProvider {
	return &nodeEndpointProvider{}
}

// GetNodeStatusEndpoint returns the node status endpoint
func (node *nodeEndpointProvider) GetNodeStatusEndpoint(_ uint32) string {
	return nodeGetNodeStatusEndpoint
}

// GetRawBlockByHashEndpoint returns the raw block by hash endpoint
func (node *nodeEndpointProvider) GetRawBlockByHashEndpoint(_ uint32, hexHash string) string {
	return fmt.Sprintf(nodeRawBlockByHashEndpoint, hexHash)
}

// GetRawBlockByNonceEndpoint returns the raw block by nonce endpoint
func (node *nodeEndpointProvider) GetRawBlockByNonceEndpoint(_ uint32, nonce uint64) string {
	return fmt.Sprintf(nodeRawBlockByNonceEndpoint, nonce)
}

// GetRawMiniBlockByHashEndpoint returns the raw miniblock by hash endpoint
func (node *nodeEndpointProvider) GetRawMiniBlockByHashEndpoint(_ uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(nodeRawMiniBlockByHashEndpoint, hexHash, epoch)
}

// IsInterfaceNil returns true if there is no value under the interface
func (node *nodeEndpointProvider) IsInterfaceNil() bool {
	return node == nil
}
