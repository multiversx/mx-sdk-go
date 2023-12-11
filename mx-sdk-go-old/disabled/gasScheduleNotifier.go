package disabled

import "github.com/multiversx/mx-chain-core-go/core"

// GasScheduleNotifier -
type GasScheduleNotifier struct {
}

// RegisterNotifyHandler does nothing
func (gsn *GasScheduleNotifier) RegisterNotifyHandler(_ core.GasScheduleSubscribeHandler) {
}

// LatestGasSchedule returns an empty map
func (gsn *GasScheduleNotifier) LatestGasSchedule() map[string]map[string]uint64 {
	return make(map[string]map[string]uint64)
}

// UnRegisterAll does nothing
func (gsn *GasScheduleNotifier) UnRegisterAll() {
}

// ChangeGasSchedule does nothing
func (gsn *GasScheduleNotifier) ChangeGasSchedule(_ map[string]map[string]uint64) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (gsn *GasScheduleNotifier) IsInterfaceNil() bool {
	return gsn == nil
}
