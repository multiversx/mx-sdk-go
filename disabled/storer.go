package disabled

// Storer is a disabled implementation of Storer interface
type Storer struct {
}

// Put returns nil
func (s *Storer) Put(key, data []byte) error {
	return nil
}

// PutInEpoch returns nil
func (s *Storer) PutInEpoch(key, data []byte, _ uint32) error {
	return nil
}

// Get returns nil
func (s *Storer) Get(key []byte) ([]byte, error) {
	return nil, nil
}

// GetFromEpoch returns nil
func (s *Storer) GetFromEpoch(key []byte, _ uint32) ([]byte, error) {
	return nil, nil
}

// GetBulkFromEpoch returns nil
func (s *Storer) GetBulkFromEpoch(keys [][]byte, _ uint32) (map[string][]byte, error) {
	return nil, nil
}

// Has returns nil
func (s *Storer) Has(_ []byte) error {
	return nil
}

// SearchFirst returns nil
func (s *Storer) SearchFirst(key []byte) ([]byte, error) {
	return nil, nil
}

// RemoveFromCurrentEpoch returns nil
func (s *Storer) RemoveFromCurrentEpoch(_ []byte) error {
	return nil
}

// Remove return nil
func (s *Storer) Remove(_ []byte) error {
	return nil
}

// ClearCache does nothing
func (s *Storer) ClearCache() {
}

// DestroyUnit returns nil
func (s *Storer) DestroyUnit() error {
	return nil
}

// GetOldestEpoch return nil
func (s *Storer) GetOldestEpoch() (uint32, error) {
	return 0, nil
}

// Close returns nil
func (s *Storer) Close() error {
	return nil
}

// RangeKeys does nothing
func (s *Storer) RangeKeys(_ func(key []byte, val []byte) bool) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *Storer) IsInterfaceNil() bool {
	return s == nil
}
