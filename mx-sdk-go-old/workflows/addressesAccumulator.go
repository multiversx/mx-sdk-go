package workflows

import "sync"

// addressesAccumulator is a mutex protected map that has 2 concurrent safe operations: pop and push
type addressesAccumulator struct {
	mut          sync.Mutex
	addressesMap map[string]struct{}
}

// newAddressesAccumulator creates a new instance of the address accumulator
func newAddressesAccumulator() *addressesAccumulator {
	return &addressesAccumulator{
		addressesMap: make(map[string]struct{}),
	}
}

func (aa *addressesAccumulator) push(address string) {
	aa.mut.Lock()
	aa.addressesMap[address] = struct{}{}
	aa.mut.Unlock()
}

func (aa *addressesAccumulator) pop() []string {
	aa.mut.Lock()
	addresses := make([]string, 0, len(aa.addressesMap))
	for address := range aa.addressesMap {
		addresses = append(addresses, address)
	}
	aa.addressesMap = make(map[string]struct{})
	aa.mut.Unlock()

	return addresses
}
