package disabled

// NodesCoordinatorCache is a disabled implementation of Cacher interface
type NodesCoordinatorCache struct {
}

// Clear does nothing
func (rm *NodesCoordinatorCache) Clear() {
}

// Put returns false
func (rm *NodesCoordinatorCache) Put(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
	return false
}

// Get returns false
func (rm *NodesCoordinatorCache) Get(key []byte) (value interface{}, ok bool) {
	return nil, false
}
