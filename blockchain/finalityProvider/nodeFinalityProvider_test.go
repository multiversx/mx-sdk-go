package finalityProvider

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
)

func TestNewNodeFinalityProvider(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		provider, err := NewNodeFinalityProvider(nil)
		assert.True(t, check.IfNil(provider))
		assert.Equal(t, ErrNilProxy, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		provider, err := NewNodeFinalityProvider(&testsCommon.ProxyStub{})
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
	})
}

func TestNodeFinalityProvider_CheckShardFinalization(t *testing.T) {
	t.Parallel()

	t.Run("invalid maxNoncesDelta", func(t *testing.T) {
		t.Parallel()

		provider, _ := NewNodeFinalityProvider(&testsCommon.ProxyStub{})

		err := provider.CheckShardFinalization(context.Background(), 0, 0)
		assert.True(t, errors.Is(err, ErrInvalidAllowedDeltaToFinal))
	})
	t.Run("get status fails", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		provider := createNodeProviderForCheckShardFinalization(nonces{}, expectedErr)

		err := provider.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, expectedErr))
	})
	t.Run("node not started", func(t *testing.T) {
		provider := createNodeProviderForCheckShardFinalization(nonces{}, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, ErrNodeNotStarted))
	})
	t.Run("node is syncing", func(t *testing.T) {
		n := nonces{
			current:  1,
			highest:  0,
			probable: 9,
		}

		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "shardID 1 is syncing, probable nonce is 9, current nonce is 1"))
	})
	t.Run("shard is stuck", func(t *testing.T) {
		n := nonces{
			current:  10,
			highest:  2,
			probable: 13,
		}

		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "shardID 1 is stuck, highest nonce is 2, current nonce is 10"))
	})
	t.Run("should work -> probable under current", func(t *testing.T) {
		n := nonces{
			current:  10,
			highest:  8,
			probable: 9,
		}
		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.Nil(t, err)
	})
	t.Run("should work -> probable == current", func(t *testing.T) {
		n := nonces{
			current:  10,
			highest:  8,
			probable: 10,
		}
		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.Nil(t, err)
	})
	t.Run("should work -> probable == current + maxDelta", func(t *testing.T) {
		n := nonces{
			current:  10,
			highest:  8,
			probable: 17,
		}
		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.Nil(t, err)
	})
	t.Run("should work -> probable == current + maxDelta && highest + maxDelta == current", func(t *testing.T) {
		n := nonces{
			current:  10,
			highest:  3,
			probable: 17,
		}
		provider := createNodeProviderForCheckShardFinalization(n, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 7)
		assert.Nil(t, err)
	})
}

func createNodeProviderForCheckShardFinalization(
	nonces nonces,
	err error,
) *nodeFinalityProvider {

	stub := &testsCommon.ProxyStub{
		GetNetworkStatusCalled: func(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
			if err != nil {
				return nil, err
			}

			return &data.NetworkStatus{
				Nonce:                nonces.current,
				HighestNonce:         nonces.highest,
				ProbableHighestNonce: nonces.probable,
			}, nil
		},
	}

	provider, _ := NewNodeFinalityProvider(stub)

	return provider
}
