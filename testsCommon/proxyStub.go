package testsCommon

import (
	"context"

	"github.com/ElrondNetwork/elrond-go-core/core"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// ProxyStub -
type ProxyStub struct {
	GetNetworkConfigCalled            func() (*data.NetworkConfig, error)
	GetRatingsConfigCalled            func() (*data.RatingsConfig, error)
	GetEnableEpochsConfigCalled       func() (*data.EnableEpochsConfig, error)
	GetAccountCalled                  func(address erdgoCore.AddressHandler) (*data.Account, error)
	SendTransactionCalled             func(tx *data.Transaction) (string, error)
	SendTransactionsCalled            func(txs []*data.Transaction) ([]string, error)
	ExecuteVMQueryCalled              func(vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
	GetNonceAtEpochStartCalled        func(shardId uint32) (uint64, error)
	GetRawMiniBlockByHashCalled       func(shardId uint32, hash string, epoch uint32) ([]byte, error)
	GetRawBlockByNonceCalled          func(shardId uint32, nonce uint64) ([]byte, error)
	GetRawBlockByHashCalled           func(shardId uint32, hash string) ([]byte, error)
	GetRawStartOfEpochMetaBlockCalled func(epoch uint32) ([]byte, error)
	GetGenesisNodesPubKeysCalled      func() (*data.GenesisNodes, error)
	GetNetworkStatusCalled            func(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	GetShardOfAddressCalled           func(ctx context.Context, bech32Address string) (uint32, error)
	GetRestAPIEntityTypeCalled        func() erdgoCore.RestAPIEntityType
	GetGuardianDataCalled             func(ctx context.Context, address erdgoCore.AddressHandler) (*data.GuardianData, error)
}

// ExecuteVMQuery -
func (stub *ProxyStub) ExecuteVMQuery(_ context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
	if stub.ExecuteVMQueryCalled != nil {
		return stub.ExecuteVMQueryCalled(vmRequest)
	}

	return &data.VmValuesResponseData{}, nil
}

// GetNetworkConfig -
func (stub *ProxyStub) GetNetworkConfig(_ context.Context) (*data.NetworkConfig, error) {
	if stub.GetNetworkConfigCalled != nil {
		return stub.GetNetworkConfigCalled()
	}

	return &data.NetworkConfig{}, nil
}

// GetRatingsConfig -
func (stub *ProxyStub) GetRatingsConfig(_ context.Context) (*data.RatingsConfig, error) {
	if stub.GetRatingsConfigCalled != nil {
		return stub.GetRatingsConfigCalled()
	}

	return &data.RatingsConfig{}, nil
}

// GetEnableEpochsConfig -
func (stub *ProxyStub) GetEnableEpochsConfig(_ context.Context) (*data.EnableEpochsConfig, error) {
	if stub.GetEnableEpochsConfigCalled != nil {
		return stub.GetEnableEpochsConfigCalled()
	}

	return &data.EnableEpochsConfig{}, nil
}

// GetAccount -
func (stub *ProxyStub) GetAccount(_ context.Context, address erdgoCore.AddressHandler) (*data.Account, error) {
	if stub.GetAccountCalled != nil {
		return stub.GetAccountCalled(address)
	}

	return &data.Account{}, nil
}

// SendTransaction -
func (stub *ProxyStub) SendTransaction(_ context.Context, tx *data.Transaction) (string, error) {
	if stub.SendTransactionCalled != nil {
		return stub.SendTransactionCalled(tx)
	}

	return "", nil
}

// SendTransactions -
func (stub *ProxyStub) SendTransactions(_ context.Context, txs []*data.Transaction) ([]string, error) {
	if stub.SendTransactionsCalled != nil {
		return stub.SendTransactionsCalled(txs)
	}

	return make([]string, 0), nil
}

// GetNonceAtEpochStart -
func (stub *ProxyStub) GetNonceAtEpochStart(_ context.Context, shardId uint32) (uint64, error) {
	if stub.GetNonceAtEpochStartCalled != nil {
		return stub.GetNonceAtEpochStartCalled(shardId)
	}

	return 0, nil
}

// GetRawMiniBlockByHash -
func (stub *ProxyStub) GetRawMiniBlockByHash(_ context.Context, shardId uint32, hash string, epoch uint32) ([]byte, error) {
	if stub.GetRawMiniBlockByHashCalled != nil {
		return stub.GetRawMiniBlockByHashCalled(shardId, hash, epoch)
	}

	return []byte{}, nil
}

// GetRawBlockByNonce -
func (stub *ProxyStub) GetRawBlockByNonce(_ context.Context, shardId uint32, nonce uint64) ([]byte, error) {
	if stub.GetRawBlockByNonceCalled != nil {
		return stub.GetRawBlockByNonceCalled(shardId, nonce)
	}

	return []byte{}, nil
}

// GetRawBlockByHash -
func (stub *ProxyStub) GetRawBlockByHash(_ context.Context, shardId uint32, hash string) ([]byte, error) {
	if stub.GetRawBlockByHashCalled != nil {
		return stub.GetRawBlockByHashCalled(shardId, hash)
	}

	return []byte{}, nil
}

// GetRawStartOfEpochMetaBlock -
func (stub *ProxyStub) GetRawStartOfEpochMetaBlock(_ context.Context, epoch uint32) ([]byte, error) {
	if stub.GetRawStartOfEpochMetaBlockCalled != nil {
		return stub.GetRawStartOfEpochMetaBlockCalled(epoch)
	}

	return []byte{}, nil
}

// GetGenesisNodesPubKeys -
func (stub *ProxyStub) GetGenesisNodesPubKeys(_ context.Context) (*data.GenesisNodes, error) {
	if stub.GetGenesisNodesPubKeysCalled != nil {
		return stub.GetGenesisNodesPubKeysCalled()
	}
	return nil, nil
}

// GetNetworkStatus -
func (stub *ProxyStub) GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
	if stub.GetNetworkStatusCalled != nil {
		return stub.GetNetworkStatusCalled(ctx, shardID)
	}

	return &data.NetworkStatus{}, nil
}

// GetShardOfAddress -
func (stub *ProxyStub) GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error) {
	if stub.GetShardOfAddressCalled != nil {
		return stub.GetShardOfAddressCalled(ctx, bech32Address)
	}
	return core.AllShardId, nil
}

// GetRestAPIEntityType -
func (stub *ProxyStub) GetRestAPIEntityType() erdgoCore.RestAPIEntityType {
	if stub.GetRestAPIEntityTypeCalled != nil {
		return stub.GetRestAPIEntityTypeCalled()
	}

	return ""
}

// GetGuardianData -
func (stub *ProxyStub) GetGuardianData(ctx context.Context, address erdgoCore.AddressHandler) (*data.GuardianData, error) {
	if stub.GetGuardianDataCalled != nil {
		return stub.GetGuardianDataCalled(ctx, address)
	}

	return &data.GuardianData{}, nil
}

// IsInterfaceNil -
func (stub *ProxyStub) IsInterfaceNil() bool {
	return stub == nil
}
