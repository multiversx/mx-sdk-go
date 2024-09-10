package blockchain

// DisabledBlockDataCache is a no-op implementation of the BlockDataCache interface.
// It does nothing and always returns empty results.
type DisabledBlockDataCache struct{}

func (n *DisabledBlockDataCache) Get(key []byte) (interface{}, bool) {
	return nil, false
}

func (n *DisabledBlockDataCache) Put(key []byte, value interface{}, sizeInBytes int) bool {
	return false
}

func (n *DisabledBlockDataCache) IsInterfaceNil() bool {
	return n == nil
}
