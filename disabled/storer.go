package disabled

// Storer is a disabled implementation of Storer interface
type Storer struct {
}

// Put returns nil
func (s *Storer) Put(_, _ []byte) error {
	return nil
}

// PutInEpoch returns nil
func (s *Storer) PutInEpoch(_, _ []byte, _ uint32) error {
	return nil
}

// Get returns nil
func (s *Storer) Get(_ []byte) ([]byte, error) {
	return nil, nil
}

// GetFromEpoch returns nil
func (s *Storer) GetFromEpoch(_ []byte, _ uint32) ([]byte, error) {
	return nil, nil
}

// GetBulkFromEpoch returns nil
func (s *Storer) GetBulkFromEpoch(_ [][]byte, _ uint32) (map[string][]byte, error) {
	return nil, nil
}

// Has returns nil
func (s *Storer) Has(_ []byte) error {
	return nil
}

// SearchFirst returns nil
func (s *Storer) SearchFirst(_ []byte) ([]byte, error) {
	return nil, nil
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

func (s *Storer) RemoveFromCurrentEpoch(_ []byte) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *Storer) IsInterfaceNil() bool {
	return s == nil
}
