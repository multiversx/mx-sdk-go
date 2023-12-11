package disabled

// ShardCoordinator is the disabled shard coordinator instance that satisfies the mx-chain-go.sharding.Coordinator interface
type ShardCoordinator struct {
}

// NumberOfShards returns 0
func (sc *ShardCoordinator) NumberOfShards() uint32 {
	return 0
}

// ComputeId returns 0
func (sc *ShardCoordinator) ComputeId(_ []byte) uint32 {
	return 0
}

// SelfId returns 0
func (sc *ShardCoordinator) SelfId() uint32 {
	return 0
}

// SameShard returns false
func (sc *ShardCoordinator) SameShard(_, _ []byte) bool {
	return false
}

// CommunicationIdentifier returns empty string
func (sc *ShardCoordinator) CommunicationIdentifier(_ uint32) string {
	return ""
}

// IsInterfaceNil returns true if there is no value under the interface
func (sc *ShardCoordinator) IsInterfaceNil() bool {
	return sc == nil
}
