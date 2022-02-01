package disabled

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

// EpochNotifier is a disabled implementation of EpochNotifier interface
type EpochNotifier struct {
}

// RegisterNotifyHandler does nothing
func (en *EpochNotifier) RegisterNotifyHandler(handler vmcommon.EpochSubscriberHandler) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (en *EpochNotifier) IsInterfaceNil() bool {
	return en == nil
}
