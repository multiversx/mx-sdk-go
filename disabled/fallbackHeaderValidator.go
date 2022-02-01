package disabled

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
)

// FallBackHeaderValidator is a disabled implementation of FallBackHeaderValidator interface
type FallBackHeaderValidator struct {
}

// ShouldApplyFallbackValidation returns false
func (fhvs *FallBackHeaderValidator) ShouldApplyFallbackValidation(_ data.HeaderHandler) bool {
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (fhvs *FallBackHeaderValidator) IsInterfaceNil() bool {
	return fhvs == nil
}
