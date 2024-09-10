package testsCommon

// CacherStub -
type CacherStub struct {
	PutCalled func(key []byte, value interface{}, sizeInBytes int) (evicted bool)
	GetCalled func(key []byte) (value interface{}, ok bool)
}

// Put -
func (rm *CacherStub) Put(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
	if rm.PutCalled != nil {
		return rm.PutCalled(key, value, sizeInBytes)
	}
	return false
}

// Get -
func (rm *CacherStub) Get(key []byte) (value interface{}, ok bool) {
	if rm.GetCalled != nil {
		return rm.GetCalled(key)
	}
	return nil, false
}

// IsInterfaceNil -
func (stub *CacherStub) IsInterfaceNil() bool {
	return stub == nil
}
