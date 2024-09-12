package disabled

import (
	"context"

	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/state"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// Accounts is a disabled implementation of the AccountAdapter interface
type Accounts struct {
}

// SetSyncer returns nil
func (a *Accounts) SetSyncer(_ state.AccountsDBSyncer) error {
	return nil
}

// StartSnapshotIfNeeded returns nil
func (a *Accounts) StartSnapshotIfNeeded() error {
	return nil
}

// RecreateTrieFromEpoch returns nil
func (a *Accounts) RecreateTrieFromEpoch(_ common.RootHashHolder) error {
	return nil
}

// GetStackDebugFirstEntry return nil
func (a *Accounts) GetStackDebugFirstEntry() []byte {
	return nil
}

// SnapshotState does nothing
func (a *Accounts) SnapshotState(_ []byte, _ uint32) {
}

// GetTrie returns nil trie and nil error
func (a *Accounts) GetTrie(_ []byte) (common.Trie, error) {
	return nil, nil
}

// Close does nothing and returns nil
func (a *Accounts) Close() error {
	return nil
}

// GetCode returns nil
func (a *Accounts) GetCode(_ []byte) []byte {
	return nil
}

// RecreateAllTries return a nil map and nil error
func (a *Accounts) RecreateAllTries(_ []byte) (map[string]common.Trie, error) {
	return nil, nil
}

// LoadAccount returns a nil account and nil error
func (a *Accounts) LoadAccount(_ []byte) (vmcommon.AccountHandler, error) {
	return nil, nil
}

// SaveAccount returns nil
func (a *Accounts) SaveAccount(_ vmcommon.AccountHandler) error {
	return nil
}

// GetAllLeaves returns nil
func (a *Accounts) GetAllLeaves(_ *common.TrieIteratorChannels, _ context.Context, _ []byte, _ common.TrieLeafParser) error {
	return nil
}

// Commit returns nil byte slice and nil
func (a *Accounts) Commit() ([]byte, error) {
	return nil, nil
}

// CommitInEpoch returns nil byte slice and nil
func (a *Accounts) CommitInEpoch(uint32, uint32) ([]byte, error) {
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
func (a *Accounts) RecreateTrie(_ common.RootHashHolder) error {
	return nil
}

// PruneTrie does nothing
func (a *Accounts) PruneTrie(_ []byte, _ state.TriePruningIdentifier, _ state.PruningHandler) {
}

// CancelPrune does nothing
func (a *Accounts) CancelPrune(_ []byte, _ state.TriePruningIdentifier) {
}

// SetStateCheckpoint does nothing
func (a *Accounts) SetStateCheckpoint(_ []byte) {
}

// IsPruningEnabled returns false
func (a *Accounts) IsPruningEnabled() bool {
	return false
}

// GetNumCheckpoints returns 0
func (a *Accounts) GetNumCheckpoints() uint32 {
	return 0
}

// GetAccountFromBytes returns a nil account and nil error
func (a *Accounts) GetAccountFromBytes(_ []byte, _ []byte) (vmcommon.AccountHandler, error) {
	return nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Accounts) IsInterfaceNil() bool {
	return a == nil
}
