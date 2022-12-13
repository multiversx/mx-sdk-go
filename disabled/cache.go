package disabled

// Cache is a disabled implementation of Cacher interface
type Cache struct {
}

// Clear does nothing
func (c *Cache) Clear() {
}

// Put returns false
func (c *Cache) Put(_ []byte, _ interface{}, _ int) bool {
	return false
}

// Get returns false
func (c *Cache) Get(_ []byte) (interface{}, bool) {
	return nil, false
}
