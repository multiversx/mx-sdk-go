package disabled

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// SimpleESDTNFTStorageHandler is a disabled implementation of SimpleESDTNFTStorageHandler interface
type SimpleESDTNFTStorageHandler struct {
}

// GetESDTNFTTokenOnDestination returns nil
func (sns *SimpleESDTNFTStorageHandler) GetESDTNFTTokenOnDestination(_ vmcommon.UserAccountHandler, _ []byte, nonce uint64) (*esdt.ESDigitalToken, bool, error) {
	return nil, false, nil
}

// SaveNFTMetaDataToSystemAccount returns nil
func (sns *SimpleESDTNFTStorageHandler) SaveNFTMetaDataToSystemAccount(_ data.TransactionHandler) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sns *SimpleESDTNFTStorageHandler) IsInterfaceNil() bool {
	return sns == nil
}
