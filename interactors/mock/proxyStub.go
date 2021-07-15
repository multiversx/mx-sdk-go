package mock

import (
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// ProxyStub -
type ProxyStub struct {
	GetNetworkConfigCalled func() (*data.NetworkConfig, error)
	GetAccountCalled       func(address core.AddressHandler) (*data.Account, error)
	SendTransactionCalled  func(tx *data.Transaction) (string, error)
	SendTransactionsCalled func(txs []*data.Transaction) ([]string, error)
}

// GetNetworkConfig -
func (ps *ProxyStub) GetNetworkConfig() (*data.NetworkConfig, error) {
	return ps.GetNetworkConfigCalled()
}

// GetAccount -
func (ps *ProxyStub) GetAccount(address core.AddressHandler) (*data.Account, error) {
	return ps.GetAccountCalled(address)
}

// SendTransaction -
func (ps *ProxyStub) SendTransaction(tx *data.Transaction) (string, error) {
	return ps.SendTransactionCalled(tx)
}

// SendTransactions -
func (ps *ProxyStub) SendTransactions(txs []*data.Transaction) ([]string, error) {
	return ps.SendTransactionsCalled(txs)
}

// IsInterfaceNil -
func (ps *ProxyStub) IsInterfaceNil() bool {
	return ps == nil
}
