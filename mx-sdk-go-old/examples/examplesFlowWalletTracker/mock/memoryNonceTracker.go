package mock

import "sync"

// MemoryNonceTracker is a test implementation of the LastProcessedNonceHandler that will keep
// the current nonce in memory
type MemoryNonceTracker struct {
	sync.RWMutex
	nonce uint64
}

// ProcessedNonce -
func (nt *MemoryNonceTracker) ProcessedNonce(nonce uint64) {
	nt.Lock()
	nt.nonce = nonce
	nt.Unlock()
}

// GetLastProcessedNonce -
func (nt *MemoryNonceTracker) GetLastProcessedNonce() uint64 {
	nt.RLock()
	defer nt.RUnlock()

	return nt.nonce
}

// IsInterfaceNil -
func (nt *MemoryNonceTracker) IsInterfaceNil() bool {
	return nt == nil
}
