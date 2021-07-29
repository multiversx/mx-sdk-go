package disabled

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Accounts is a disabled implementation of the AccountAdapter interface
type Accounts struct {
}

// GetCode returns nil
func (a *Accounts) GetCode(_ []byte) []byte {
	return nil
}

// LoadAccount returns a nil account and nil error
func (a *Accounts) LoadAccount(_ []byte) (vmcommon.AccountHandler, error) {
	return nil, nil
}

// SaveAccount returns nil
func (a *Accounts) SaveAccount(_ vmcommon.AccountHandler) error {
	return nil
}

// Commit returns nil byte slice and nil
func (a *Accounts) Commit() ([]byte, error) {
	return nil, nil
}

// GetExistingAccount returns nil  account handler and nil error
func (a *Accounts) GetExistingAccount(_ []byte) (vmcommon.AccountHandler, error) {
	return nil, nil
}

// JournalLen returns 0
func (a *Accounts) JournalLen() int {
	return 0
}

// RemoveAccount returns nil
func (a *Accounts) RemoveAccount(_ []byte) error {
	return nil
}

// RevertToSnapshot returns nil
func (a *Accounts) RevertToSnapshot(_ int) error {
	return nil
}

// RootHash returns nil byte slice and nil error
func (a *Accounts) RootHash() ([]byte, error) {
	return nil, nil
}

// RecreateTrie returns nil
func (a *Accounts) RecreateTrie(_ []byte) error {
	return nil
}

// GetNumCheckpoints returns 0
func (a *Accounts) GetNumCheckpoints() uint32 {
	return 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Accounts) IsInterfaceNil() bool {
	return a == nil
}
