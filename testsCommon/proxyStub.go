package testsCommon

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// ProxyStub -
type ProxyStub struct {
	GetNetworkConfigCalled func() (*data.NetworkConfig, error)
	GetAccountCalled       func(address core.AddressHandler) (*data.Account, error)
	SendTransactionCalled  func(tx *data.Transaction) (string, error)
	SendTransactionsCalled func(txs []*data.Transaction) ([]string, error)
	ExecuteVMQueryCalled   func(vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
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

// GetAccount -
func (stub *ProxyStub) GetAccount(_ context.Context, address core.AddressHandler) (*data.Account, error) {
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

// IsInterfaceNil -
func (stub *ProxyStub) IsInterfaceNil() bool {
	return stub == nil
}
