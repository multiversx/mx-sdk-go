package disabled

import vmcommon "github.com/multiversx/mx-chain-vm-common-go"

// BlockChainHookCounter is a disabled implementation of BlockChainHookCounter interface
type BlockChainHookCounter struct {
}

// ProcessCrtNumberOfTrieReadsCounter returns nil
func (bhc *BlockChainHookCounter) ProcessCrtNumberOfTrieReadsCounter() error {
	return nil
}

// ProcessMaxBuiltInCounters returns nil
func (bhc *BlockChainHookCounter) ProcessMaxBuiltInCounters(_ *vmcommon.ContractCallInput) error {
	return nil
}

// ResetCounters does nothing
func (bhc *BlockChainHookCounter) ResetCounters() {
}

// SetMaximumValues does nothing
func (bhc *BlockChainHookCounter) SetMaximumValues(_ map[string]uint64) {
}

// GetCounterValues returns an empty map
func (bhc *BlockChainHookCounter) GetCounterValues() map[string]uint64 {
	return make(map[string]uint64)
}

// IsInterfaceNil returns true if there is no value under the interface
func (bhc *BlockChainHookCounter) IsInterfaceNil() bool {
	return bhc == nil
}
