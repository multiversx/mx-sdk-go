package disabled

import "github.com/multiversx/mx-chain-core-go/core"

// EnableEpochsHandler is a disabled implementation of EnableEpochsHandler interface
type EnableEpochsHandler struct {
}

// GetCurrentEpoch returns 0 as it is disabled
func (eeh *EnableEpochsHandler) GetCurrentEpoch() uint32 {
	return 0
}

// IsFlagDefined returns true as it is disabled
func (eeh *EnableEpochsHandler) IsFlagDefined(_ core.EnableEpochFlag) bool {
	return true
}

// IsFlagEnabled returns true as it is disabled
func (eeh *EnableEpochsHandler) IsFlagEnabled(_ core.EnableEpochFlag) bool {
	return true
}

// IsFlagEnabledInEpoch returns true as it is disabled
func (eeh *EnableEpochsHandler) IsFlagEnabledInEpoch(_ core.EnableEpochFlag, _ uint32) bool {
	return true
}

// GetActivationEpoch returns 0 as it is disabled
func (eeh *EnableEpochsHandler) GetActivationEpoch(_ core.EnableEpochFlag) uint32 {
	return 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (eeh *EnableEpochsHandler) IsInterfaceNil() bool {
	return eeh == nil
}
