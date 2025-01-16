package disabled

import (
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// SimpleESDTNFTStorageHandler is a disabled implementation of SimpleESDTNFTStorageHandler interface
type SimpleESDTNFTStorageHandler struct {
}

// SaveNFTMetaData returns nil
func (sns *SimpleESDTNFTStorageHandler) SaveNFTMetaData(_ data.TransactionHandler) error {
	return nil
}

// GetESDTNFTTokenOnDestination returns nil
func (sns *SimpleESDTNFTStorageHandler) GetESDTNFTTokenOnDestination(_ vmcommon.UserAccountHandler, _ []byte, _ uint64) (*esdt.ESDigitalToken, bool, error) {
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

// SaveNFTMetaData returns nil
func (sns *SimpleESDTNFTStorageHandler) SaveNFTMetaData(tx data.TransactionHandler) error {
	return nil
}
