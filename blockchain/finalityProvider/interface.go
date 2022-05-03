package finalityProvider

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	IsInterfaceNil() bool
}
