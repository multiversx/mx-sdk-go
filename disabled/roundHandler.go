package disabled

import "time"

// RoundHandler is a disabled round handler
type RoundHandler struct {
}

// TimeDuration returns time.Second as it is disabled
func (handler *RoundHandler) TimeDuration() time.Duration {
	return time.Second
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *RoundHandler) IsInterfaceNil() bool {
	return handler == nil
}
