package blockchain

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// Proxy holds the primitive functions that the elrond proxy engine supports & implements
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*data.Transaction) ([]string, error)
	IsInterfaceNil() bool
}
