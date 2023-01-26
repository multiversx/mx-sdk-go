package disabled

// MultiversXShardCoordinator is the disabled shard coordinator instance that satisfies the mx-chain-go.sharding.Coordinator interface
type MultiversXShardCoordinator struct {
}

// NumberOfShards returns 0
func (msc *MultiversXShardCoordinator) NumberOfShards() uint32 {
	return 0
}

// ComputeId returns 0
func (msc *MultiversXShardCoordinator) ComputeId(_ []byte) uint32 {
	return 0
}

// SelfId returns 0
func (msc *MultiversXShardCoordinator) SelfId() uint32 {
	return 0
}

// SameShard returns false
func (msc *MultiversXShardCoordinator) SameShard(_, _ []byte) bool {
	return false
}

// CommunicationIdentifier returns empty string
func (msc *MultiversXShardCoordinator) CommunicationIdentifier(_ uint32) string {
	return ""
}

// IsInterfaceNil returns true if there is no value under the interface
func (msc *MultiversXShardCoordinator) IsInterfaceNil() bool {
	return msc == nil
}
