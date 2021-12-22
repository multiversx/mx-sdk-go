package blockchain

import (
	"context"
	"sync"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type elrondProxyWithCache struct {
	mut   sync.RWMutex
	proxy elrondProxy
}

func (cacher *elrondProxyWithCache) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	return proxy.GetNetworkConfig(ctx)
}

func (cacher *elrondProxyWithCache) GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
	// TODO implement me
	panic("implement me")
}

func (cacher *elrondProxyWithCache) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (cacher *elrondProxyWithCache) SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (cacher *elrondProxyWithCache) IsInterfaceNil() bool {
	// TODO implement me
	panic("implement me")
}
