package testsCommon

import (
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-go/sharding/nodesCoordinator"
	"github.com/multiversx/mx-chain-go/state"
)

// NodesCoordinatorStub -
type NodesCoordinatorStub struct {
	ComputeValidatorsGroupCalled                      func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]nodesCoordinator.Validator, error)
	GetValidatorsPublicKeysCalled                     func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]string, error)
	GetValidatorsRewardsAddressesCalled               func(randomness []byte, round uint64, shardId uint32, epoch uint32) ([]string, error)
	GetValidatorWithPublicKeyCalled                   func(publicKey []byte) (validator nodesCoordinator.Validator, shardId uint32, err error)
	GetAllValidatorsPublicKeysCalled                  func() (map[uint32][][]byte, error)
	ConsensusGroupSizeCalled                          func(shardID uint32) int
	SetNodesConfigFromValidatorsInfoCalled            func(epoch uint32, randomness []byte, validatorsInfo []*state.ShardValidatorInfo) error
	IsEpochInConfigCalled                             func(epoch uint32) bool
	GetAllShuffledOutValidatorsPublicKeysCalled       func(epoch uint32) (map[uint32][][]byte, error)
	GetWaitingEpochsLeftForPublicKeyCalled            func(publicKey []byte) (uint32, error)
	GetShuffledOutToAuctionValidatorsPublicKeysCalled func(epoch uint32) (map[uint32][][]byte, error)
}

func (ncm *NodesCoordinatorStub) EpochStartPrepare(metaHdr data.HeaderHandler, body data.BodyHandler) {
	//TODO implement me
	panic("implement me")
}

func (ncm *NodesCoordinatorStub) NodesCoordinatorToRegistry(epoch uint32) nodesCoordinator.NodesCoordinatorRegistryHandler {
	//TODO implement me
	panic("implement me")
}

// GetChance -
func (ncm *NodesCoordinatorStub) GetChance(uint32) uint32 {
	return 1
}

// ValidatorsWeights -
func (ncm *NodesCoordinatorStub) ValidatorsWeights(_ []nodesCoordinator.Validator) ([]uint32, error) {
	return nil, nil
}

// GetAllLeavingValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllLeavingValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// SetConfig -
func (ncm *NodesCoordinatorStub) SetConfig(_ *nodesCoordinator.NodesCoordinatorRegistry) error {
	return nil
}

// ComputeAdditionalLeaving -
func (ncm *NodesCoordinatorStub) ComputeAdditionalLeaving(_ []*state.ShardValidatorInfo) (map[uint32][]nodesCoordinator.Validator, error) {
	return nil, nil
}

// GetAllEligibleValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllEligibleValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetAllWaitingValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllWaitingValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	return nil, nil
}

// GetNumTotalEligible -
func (ncm *NodesCoordinatorStub) GetNumTotalEligible() uint64 {
	return 1
}

// GetAllValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllValidatorsPublicKeys(_ uint32) (map[uint32][][]byte, error) {
	if ncm.GetAllValidatorsPublicKeysCalled != nil {
		return ncm.GetAllValidatorsPublicKeysCalled()
	}

	return nil, nil
}

// GetValidatorsIndexes -
func (ncm *NodesCoordinatorStub) GetValidatorsIndexes(_ []string, _ uint32) ([]uint64, error) {
	return nil, nil
}

// ComputeConsensusGroup -
func (ncm *NodesCoordinatorStub) ComputeConsensusGroup(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) (validatorsGroup []nodesCoordinator.Validator, err error) {

	if ncm.ComputeValidatorsGroupCalled != nil {
		return ncm.ComputeValidatorsGroupCalled(randomness, round, shardId, epoch)
	}

	var list []nodesCoordinator.Validator

	return list, nil
}

// ConsensusGroupSize -
func (ncm *NodesCoordinatorStub) ConsensusGroupSize(shardID uint32) int {
	if ncm.ConsensusGroupSizeCalled != nil {
		return ncm.ConsensusGroupSizeCalled(shardID)
	}
	return 1
}

// GetConsensusValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetConsensusValidatorsPublicKeys(
	randomness []byte,
	round uint64,
	shardId uint32,
	epoch uint32,
) ([]string, error) {
	if ncm.GetValidatorsPublicKeysCalled != nil {
		return ncm.GetValidatorsPublicKeysCalled(randomness, round, shardId, epoch)
	}

	return nil, nil
}

// SetNodesPerShards -
func (ncm *NodesCoordinatorStub) SetNodesPerShards(_ map[uint32][]nodesCoordinator.Validator, _ map[uint32][]nodesCoordinator.Validator, _ []nodesCoordinator.Validator, _ uint32) error {
	return nil
}

// LoadState -
func (ncm *NodesCoordinatorStub) LoadState(_ []byte) error {
	return nil
}

// GetSavedStateKey -
func (ncm *NodesCoordinatorStub) GetSavedStateKey() []byte {
	return []byte("key")
}

// ShardIdForEpoch returns the nodesCoordinator configured ShardId for specified epoch if epoch configuration exists,
// otherwise error
func (ncm *NodesCoordinatorStub) ShardIdForEpoch(_ uint32) (uint32, error) {
	panic("not implemented")
}

// ShuffleOutForEpoch verifies if the shards changed in the new epoch and calls the shuffleOutHandler
func (ncm *NodesCoordinatorStub) ShuffleOutForEpoch(_ uint32) {
	panic("not implemented")
}

// GetConsensusWhitelistedNodes return the whitelisted nodes allowed to send consensus messages, for each of the shards
func (ncm *NodesCoordinatorStub) GetConsensusWhitelistedNodes(
	_ uint32,
) (map[string]struct{}, error) {
	panic("not implemented")
}

// GetSelectedPublicKeys -
func (ncm *NodesCoordinatorStub) GetSelectedPublicKeys(_ []byte, _ uint32, _ uint32) ([]string, error) {
	panic("implement me")
}

// GetValidatorWithPublicKey -
func (ncm *NodesCoordinatorStub) GetValidatorWithPublicKey(publicKey []byte) (nodesCoordinator.Validator, uint32, error) {
	if ncm.GetValidatorWithPublicKeyCalled != nil {
		return ncm.GetValidatorWithPublicKeyCalled(publicKey)
	}
	return nil, 0, nil
}

// GetOwnPublicKey -
func (ncm *NodesCoordinatorStub) GetOwnPublicKey() []byte {
	return []byte("key")
}

// SetNodesConfigFromValidatorsInfo -
func (ncm *NodesCoordinatorStub) SetNodesConfigFromValidatorsInfo(epoch uint32, randomness []byte, validatorsInfo []*state.ShardValidatorInfo) error {
	if ncm.SetNodesConfigFromValidatorsInfoCalled != nil {
		return ncm.SetNodesConfigFromValidatorsInfoCalled(epoch, randomness, validatorsInfo)
	}
	return nil
}

// IsEpochInConfig -
func (ncm *NodesCoordinatorStub) IsEpochInConfig(epoch uint32) bool {
	if ncm.IsEpochInConfigCalled != nil {
		return ncm.IsEpochInConfigCalled(epoch)
	}
	return false
}

// GetAllShuffledOutValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetAllShuffledOutValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	if ncm.GetAllShuffledOutValidatorsPublicKeysCalled != nil {
		return ncm.GetAllShuffledOutValidatorsPublicKeysCalled(epoch)
	}

	return nil, nil
}

// GetWaitingEpochsLeftForPublicKey -
func (ncm *NodesCoordinatorStub) GetWaitingEpochsLeftForPublicKey(publicKey []byte) (uint32, error) {
	if ncm.GetWaitingEpochsLeftForPublicKeyCalled != nil {
		return ncm.GetWaitingEpochsLeftForPublicKeyCalled(publicKey)
	}

	return 0, nil
}

// GetShuffledOutToAuctionValidatorsPublicKeys -
func (ncm *NodesCoordinatorStub) GetShuffledOutToAuctionValidatorsPublicKeys(epoch uint32) (map[uint32][][]byte, error) {
	if ncm.GetShuffledOutToAuctionValidatorsPublicKeysCalled != nil {
		return ncm.GetShuffledOutToAuctionValidatorsPublicKeysCalled(epoch)
	}

	return make(map[uint32][][]byte), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ncm *NodesCoordinatorStub) IsInterfaceNil() bool {
	return ncm == nil
}
