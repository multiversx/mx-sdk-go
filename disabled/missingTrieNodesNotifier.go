package disabled

import "github.com/multiversx/mx-chain-go/common"

// MissingTrieNodesNotifier is a disabled implementation of MissingTrieNodesNotifier interface
type MissingTrieNodesNotifier struct {
}

// RegisterHandler returns nil
func (m *MissingTrieNodesNotifier) RegisterHandler(_ common.StateSyncNotifierSubscriber) error {
	return nil
}

// AsyncNotifyMissingTrieNode does nothing
func (m *MissingTrieNodesNotifier) AsyncNotifyMissingTrieNode(_ []byte) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (m *MissingTrieNodesNotifier) IsInterfaceNil() bool {
	return m == nil
}
