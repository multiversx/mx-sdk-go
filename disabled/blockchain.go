package disabled

import (
	"github.com/multiversx/mx-chain-core-go/data"
)

// Blockchain is a disabled implementation of the ChainHandler interface
type Blockchain struct {
}

// GetFinalBlockInfo return 0 and empty slices
func (b *Blockchain) GetFinalBlockInfo() (uint64, []byte, []byte) {
	return 0, make([]byte, 0), make([]byte, 0)
}

// SetFinalBlockInfo does nothing
func (b *Blockchain) SetFinalBlockInfo(_ uint64, _ []byte, _ []byte) {
}

// GetGenesisHeader returns nil
func (b *Blockchain) GetGenesisHeader() data.HeaderHandler {
	return nil
}

// SetGenesisHeader returns nil
func (b *Blockchain) SetGenesisHeader(_ data.HeaderHandler) error {
	return nil
}

// GetGenesisHeaderHash returns nil
func (b *Blockchain) GetGenesisHeaderHash() []byte {
	return nil
}

// SetGenesisHeaderHash does nothing
func (b *Blockchain) SetGenesisHeaderHash(_ []byte) {
}

// GetCurrentBlockHeader returns nil
func (b *Blockchain) GetCurrentBlockHeader() data.HeaderHandler {
	return nil
}

// SetCurrentBlockHeader returns nil
func (b *Blockchain) SetCurrentBlockHeader(_ data.HeaderHandler) error {
	return nil
}

// GetCurrentBlockHeaderHash returns nil
func (b *Blockchain) GetCurrentBlockHeaderHash() []byte {
	return nil
}

// SetCurrentBlockHeaderHash does nothing
func (b *Blockchain) SetCurrentBlockHeaderHash(_ []byte) {
}

// CreateNewHeader returns nil
func (b *Blockchain) CreateNewHeader() data.HeaderHandler {
	return nil
}

// GetCurrentBlockRootHash returns nil
func (b *Blockchain) GetCurrentBlockRootHash() []byte {
	return nil
}

// SetCurrentBlockHeaderAndRootHash return nil
func (b *Blockchain) SetCurrentBlockHeaderAndRootHash(_ data.HeaderHandler, _ []byte) error {
	return nil
}

// GetCurrentBlockHeaderAndHash returns nil
func (b *Blockchain) GetCurrentBlockHeaderAndHash() (data.HeaderHandler, []byte) {
	return nil, nil
}

// GetLastExecutedBlockInfo returns nil
func (b *Blockchain) GetLastExecutedBlockInfo() (uint64, []byte, []byte) {
	return 0, nil, nil
}

// GetLastExecutedBlockHeader returns nil
func (b *Blockchain) GetLastExecutedBlockHeader() data.HeaderHandler {
	return nil
}

// SetLastExecutedBlockHeaderAndRootHash does nothing
func (b *Blockchain) SetLastExecutedBlockHeaderAndRootHash(
	header data.HeaderHandler,
	headerHash []byte,
	rootHash []byte,
) {
}

// GetLastExecutionResult returns nil
func (b *Blockchain) GetLastExecutionResult() data.BaseExecutionResultHandler {
	return nil
}

// SetCurrentBlockHeaderAndHash returns nil
func (b *Blockchain) SetCurrentBlockHeaderAndHash(
	headerHash []byte,
	header data.HeaderHandler,
) error {
	return nil
}

// SetLastExecutionInfo does nothing
func (b *Blockchain) SetLastExecutionInfo(
	header data.HeaderHandler,
	result data.BaseExecutionResultHandler,
) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (b *Blockchain) IsInterfaceNil() bool {
	return b == nil
}
