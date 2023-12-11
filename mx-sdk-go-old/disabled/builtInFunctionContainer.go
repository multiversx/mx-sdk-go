package disabled

import (
	"errors"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

var errNotImplemented = errors.New("not implemented")

// BuiltInFunctionContainer is a disabled implementation of the vmcommon.BuiltInFunctionContainer interface
type BuiltInFunctionContainer struct {
}

// Get returns nil and error
func (bifc *BuiltInFunctionContainer) Get(_ string) (vmcommon.BuiltinFunction, error) {
	return nil, errNotImplemented
}

// Add does nothing and returns nil error
func (bifc *BuiltInFunctionContainer) Add(_ string, _ vmcommon.BuiltinFunction) error {
	return nil
}

// Replace does nothing and returns nil error
func (bifc *BuiltInFunctionContainer) Replace(_ string, _ vmcommon.BuiltinFunction) error {
	return nil
}

// Remove does nothing
func (bifc *BuiltInFunctionContainer) Remove(_ string) {
}

// Len returns 0
func (bifc *BuiltInFunctionContainer) Len() int {
	return 0
}

// Keys returns an empty map
func (bifc *BuiltInFunctionContainer) Keys() map[string]struct{} {
	return make(map[string]struct{})
}

// IsInterfaceNil returns true if there is no value under the interface
func (bifc *BuiltInFunctionContainer) IsInterfaceNil() bool {
	return bifc == nil
}
