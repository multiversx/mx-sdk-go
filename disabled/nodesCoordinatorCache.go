package disabled

// Cache is a disabled implementation of Cacher interface
type Cache struct {
}

// Clear does nothing
func (rm *Cache) Clear() {
}

// Put returns false
func (rm *Cache) Put(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
	return false
}

// Get returns false
func (rm *Cache) Get(key []byte) (value interface{}, ok bool) {
	return nil, false
}
