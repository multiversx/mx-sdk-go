package finalityProvider

import (
	"context"

	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/data"
)

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	IsInterfaceNil() bool
}
