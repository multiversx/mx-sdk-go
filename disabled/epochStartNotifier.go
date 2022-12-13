package disabled

import (
	"github.com/ElrondNetwork/elrond-go/epochStart"
)

// EpochStartNotifier is a disabled implementation of EpochStartEventNotifier interface
type EpochStartNotifier struct {
}

// RegisterHandler does nothing
func (esn *EpochStartNotifier) RegisterHandler(_ epochStart.ActionHandler) {
}

// UnregisterHandler does nothing
func (esn *EpochStartNotifier) UnregisterHandler(_ epochStart.ActionHandler) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (esn *EpochStartNotifier) IsInterfaceNil() bool {
	return esn == nil
}
