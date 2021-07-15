package mock

import "sync"

// TrackableAddressProviderMock is the mock implementation of the trackable address provider to be used
// in the example
type TrackableAddressProviderMock struct {
	sync.RWMutex
	trackedAddresses map[string][]byte
}

// NewTrackableAddressProviderMock -
func NewTrackableAddressProviderMock() *TrackableAddressProviderMock {
	return &TrackableAddressProviderMock{
		trackedAddresses: make(map[string][]byte),
	}
}

// AllTrackableAddresses -
func (tap *TrackableAddressProviderMock) AllTrackableAddresses() []string {
	tap.Lock()
	allAddresses := make([]string, 0, len(tap.trackedAddresses))
	for address := range tap.trackedAddresses {
		allAddresses = append(allAddresses, address)
	}
	tap.Unlock()

	return allAddresses
}

// AddTrackableAddress -
func (tap *TrackableAddressProviderMock) AddTrackableAddress(addressAsBech32 string, skBytes []byte) {
	tap.Lock()
	tap.trackedAddresses[addressAsBech32] = skBytes
	tap.Unlock()
}

// IsTrackableAddresses -
func (tap *TrackableAddressProviderMock) IsTrackableAddresses(addressAsBech32 string) bool {
	tap.RLock()
	defer tap.RUnlock()

	_, ok := tap.trackedAddresses[addressAsBech32]

	return ok
}

// PrivateKeyOfBech32Address -
func (tap *TrackableAddressProviderMock) PrivateKeyOfBech32Address(addressAsBech32 string) []byte {
	tap.RLock()
	defer tap.RUnlock()

	skBytes, _ := tap.trackedAddresses[addressAsBech32]

	return skBytes
}

// IsInterfaceNil -
func (tap *TrackableAddressProviderMock) IsInterfaceNil() bool {
	return tap == nil
}
