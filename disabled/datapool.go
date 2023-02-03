package disabled

import (
	"github.com/multiversx/mx-chain-go/dataRetriever"
	"github.com/multiversx/mx-chain-go/storage"
)

// DataPool is the disabled implementation of a PoolsHolder interface
type DataPool struct {
}

// PeerAuthentications returns nil
func (dp *DataPool) PeerAuthentications() storage.Cacher {
	return nil
}

// Heartbeats returns nil
func (dp *DataPool) Heartbeats() storage.Cacher {
	return nil
}

// Close returns nil
func (dp *DataPool) Close() error {
	return nil
}

// TrieNodesChunks returns nil
func (dp *DataPool) TrieNodesChunks() storage.Cacher {
	return nil
}

// Transactions returns nil
func (dp *DataPool) Transactions() dataRetriever.ShardedDataCacherNotifier {
	return nil
}

// UnsignedTransactions returns nil
func (dp *DataPool) UnsignedTransactions() dataRetriever.ShardedDataCacherNotifier {
	return nil
}

// RewardTransactions returns nil
func (dp *DataPool) RewardTransactions() dataRetriever.ShardedDataCacherNotifier {
	return nil
}

// Headers returns nil
func (dp *DataPool) Headers() dataRetriever.HeadersPool {
	return nil
}

// MiniBlocks returns nil
func (dp *DataPool) MiniBlocks() storage.Cacher {
	return nil
}

//PeerChangesBlocks returns nil
func (dp *DataPool) PeerChangesBlocks() storage.Cacher {
	return nil
}

// TrieNodes returns nil
func (dp *DataPool) TrieNodes() storage.Cacher {
	return nil
}

// SmartContracts returns nil
func (dp *DataPool) SmartContracts() storage.Cacher {
	return nil
}

// CurrentBlockTxs returns nil
func (dp *DataPool) CurrentBlockTxs() dataRetriever.TransactionCacher {
	return nil
}

// CurrentEpochValidatorInfo returns nil
func (dp *DataPool) CurrentEpochValidatorInfo() dataRetriever.ValidatorInfoCacher {
	return nil
}

// ValidatorsInfo returns nil
func (dp *DataPool) ValidatorsInfo() dataRetriever.ShardedDataCacherNotifier {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (dp *DataPool) IsInterfaceNil() bool {
	return dp == nil
}
