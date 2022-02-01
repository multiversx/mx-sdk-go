package disabled

// ShuffledOutHandler is a disabled implementation of ShuffledOutHandler interface
type ShuffledOutHandler struct {
}

// Process returns nil
func (s *ShuffledOutHandler) Process(newShardID uint32) error {
	return nil
}

// RegisterHandler does nothing
func (s *ShuffledOutHandler) RegisterHandler(handler func(newShardID uint32)) {
}

// CurrentShardID return zero
func (s *ShuffledOutHandler) CurrentShardID() uint32 {
	return 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *ShuffledOutHandler) IsInterfaceNil() bool {
	return s == nil
}
