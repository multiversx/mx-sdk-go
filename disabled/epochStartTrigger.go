package disabled

import (
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
)

// EpochStartTrigger is a disabled epoch start trigger
type EpochStartTrigger struct {
}

// LastCommitedEpochStartHdr returns an empty handler and nil as it is disabled
func (trigger *EpochStartTrigger) LastCommitedEpochStartHdr() (data.HeaderHandler, error) {
	return &block.HeaderV2{}, nil
}

// GetEpochStartHdrFromStorage returns an empty handler and nil as it is disabled
func (trigger *EpochStartTrigger) GetEpochStartHdrFromStorage(_ uint32) (data.HeaderHandler, error) {
	return &block.HeaderV2{}, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (trigger *EpochStartTrigger) IsInterfaceNil() bool {
	return trigger == nil
}
