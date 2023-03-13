package testsCommon

import (
	"context"

	"github.com/multiversx/mx-sdk-go/data"
)

// BlockhashHandlerStub -
type BlockhashHandlerStub struct {
	GetBlockByHashCalled func(ctx context.Context, hash string) (*data.Block, error)
}

// GetBlockByHash -
func (b *BlockhashHandlerStub) GetBlockByHash(ctx context.Context, hash string) (*data.Block, error) {
	if b.GetBlockByHashCalled != nil {
		return b.GetBlockByHashCalled(ctx, hash)
	}

	return nil, nil
}

// IsInterfaceNil -
func (b *BlockhashHandlerStub) IsInterfaceNil() bool {
	return b == nil
}
