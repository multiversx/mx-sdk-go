package disabled

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go/common"
)

type StorageManager struct {
}

// Database returns nil
func (s *StorageManager) Database() common.DBWriteCacher {
	return nil
}

// TakeSnapshot does nothing
func (s *StorageManager) TakeSnapshot(_ []byte, _ bool, _ chan core.KeyValueHolder) {
}

// SetCheckpoint does nothing
func (s *StorageManager) SetCheckpoint(_ []byte, _ chan core.KeyValueHolder) {
}

// GetSnapshotThatContainsHash returns nil
func (s *StorageManager) GetSnapshotThatContainsHash(_ []byte) common.SnapshotDbHandler {
	return nil
}

// IsPruningEnabled returns false
func (s *StorageManager) IsPruningEnabled() bool {
	return false
}

// IsPruningBlocked returns false
func (s *StorageManager) IsPruningBlocked() bool {
	return false
}

// EnterPruningBufferingMode does nothing
func (s *StorageManager) EnterPruningBufferingMode() {
}

// ExitPruningBufferingMode does nothing
func (s *StorageManager) ExitPruningBufferingMode() {
}

// GetSnapshotDbBatchDelay returns 0
func (s *StorageManager) GetSnapshotDbBatchDelay() int {
	return 0
}

// AddDirtyCheckpointHashes returns false
func (s *StorageManager) AddDirtyCheckpointHashes(_ []byte, _ common.ModifiedHashes) bool {
	return false
}

// Remove does nothing
func (s *StorageManager) Remove(_ []byte) error {
	return nil
}

// Close does nothing
func (s *StorageManager) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value
func (s *StorageManager) IsInterfaceNil() bool {
	return s == nil
}
