package disabled

// GlobalSettingsHandler implements a disabled vmcommon.ESDTGlobalSettingsHandler
type GlobalSettingsHandler struct {
}

// IsPaused returns false
func (handler *GlobalSettingsHandler) IsPaused(_ []byte) bool {
	return false
}

// IsLimitedTransfer returns false
func (handler *GlobalSettingsHandler) IsLimitedTransfer(_ []byte) bool {
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *GlobalSettingsHandler) IsInterfaceNil() bool {
	return handler == nil
}
