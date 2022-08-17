package disabled

type GlobalSettingsHandler struct {
}

func (gsh *GlobalSettingsHandler) IsPaused(_ []byte) bool {
	return false
}

func (gsh *GlobalSettingsHandler) IsLimitedTransfer(esdtTokenKey []byte) bool {
	return false
}

func (gsh *GlobalSettingsHandler) IsInterfaceNil() bool {
	return gsh == nil
}
