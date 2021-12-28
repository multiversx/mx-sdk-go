package mock

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

func (ps *ProxyStub) ExecuteVMQuery(_ context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error) {
	return ps.ExecuteVMQueryCalled(vmRequest)
}

// GetNetworkConfig -
func (ps *ProxyStub) GetNetworkConfig(_ context.Context) (*data.NetworkConfig, error) {
	return ps.GetNetworkConfigCalled()
}

// GetAccount -
func (ps *ProxyStub) GetAccount(_ context.Context, address core.AddressHandler) (*data.Account, error) {
	return ps.GetAccountCalled(address)
}

// SendTransaction -
func (ps *ProxyStub) SendTransaction(_ context.Context, tx *data.Transaction) (string, error) {
	return ps.SendTransactionCalled(tx)
}

// SendTransactions -
func (ps *ProxyStub) SendTransactions(_ context.Context, txs []*data.Transaction) ([]string, error) {
	return ps.SendTransactionsCalled(txs)
}

// IsInterfaceNil -
func (ps *ProxyStub) IsInterfaceNil() bool {
	return ps == nil
}
