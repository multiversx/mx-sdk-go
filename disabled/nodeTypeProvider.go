package disabled

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
)

// NodeTypeProvider is a disabled implementation of NodeTypeProviderHandler interface
type NodeTypeProvider struct {
}

// SetType does nothing
func (n *NodeTypeProvider) SetType(_ core.NodeType) {
}

// GetType returns empty string
func (n *NodeTypeProvider) GetType() core.NodeType {
	return ""
}

// IsInterfaceNil returns true if there is no value under the interface
func (n *NodeTypeProvider) IsInterfaceNil() bool {
	return n == nil
}
