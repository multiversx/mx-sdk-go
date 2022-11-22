package disabled

import crypto "github.com/ElrondNetwork/elrond-go-crypto"

type MultiSignerContainer struct {
}

// NewMultiSignerContainer returns nil
func NewMultiSignerContainer() *MultiSignerContainer {
	return nil
}

// GetMultiSigner returns nil
func (dmsc *MultiSignerContainer) GetMultiSigner(_ uint32) (crypto.MultiSigner, error) {
	return nil, nil
}

// IsInterfaceNil returns true if the underlying object is nil
func (dmsc *MultiSignerContainer) IsInterfaceNil() bool {
	return dmsc == nil
}
