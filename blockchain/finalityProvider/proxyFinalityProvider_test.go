package finalityProvider

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
)

const goodResponseExample = "0: 9169897, 1: 9166353, 2: 9170524, "

func TestNewProxyFinalityProvider(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		provider, err := NewProxyFinalityProvider(nil)
		assert.True(t, check.IfNil(provider))
		assert.Equal(t, ErrNilProxy, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		provider, err := NewProxyFinalityProvider(&testsCommon.ProxyStub{})
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
	})
}

func TestExtractNonceOfShardID(t *testing.T) {
	t.Parallel()

	t.Run("empty response should error", func(t *testing.T) {
		t.Parallel()

		nonce, err := extractNonceOfShardID("", 0)
		assert.True(t, errors.Is(err, ErrInvalidNonceCrossCheckValueFormat))
		assert.True(t, strings.Contains(err.Error(), "empty value"))
		assert.Equal(t, uint64(0), nonce)
	})
	t.Run("shard data contains a NaN", func(t *testing.T) {
		t.Parallel()

		nonce, err := extractNonceOfShardID("0: aaa", 0)
		assert.True(t, errors.Is(err, ErrInvalidNonceCrossCheckValueFormat))
		assert.True(t, strings.Contains(err.Error(), "is not a valid number"))
		assert.Equal(t, uint64(0), nonce)
	})
	t.Run("shard not found", func(t *testing.T) {
		t.Parallel()

		nonce, err := extractNonceOfShardID(goodResponseExample, 3)
		assert.True(t, errors.Is(err, ErrInvalidNonceCrossCheckValueFormat))
		assert.True(t, strings.Contains(err.Error(), "value not found for shard 3"))
		assert.Equal(t, uint64(0), nonce)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		nonce, err := extractNonceOfShardID(goodResponseExample, 0)
		assert.Nil(t, err)
		assert.Equal(t, uint64(9169897), nonce)

		nonce, err = extractNonceOfShardID(goodResponseExample, 1)
		assert.Nil(t, err)
		assert.Equal(t, uint64(9166353), nonce)

		nonce, err = extractNonceOfShardID(goodResponseExample, 2)
		assert.Nil(t, err)
		assert.Equal(t, uint64(9170524), nonce)
	})
	t.Run("should work even if it contains extra data", func(t *testing.T) {
		t.Parallel()

		nonce, err := extractNonceOfShardID("extra,"+goodResponseExample, 0)
		assert.Nil(t, err)
		assert.Equal(t, uint64(9169897), nonce)
	})
}

func TestElrondBaseProxy_CheckShardFinalization(t *testing.T) {
	t.Parallel()

	t.Run("invalid maxNoncesDelta", func(t *testing.T) {
		t.Parallel()

		provider, _ := NewProxyFinalityProvider(&testsCommon.ProxyStub{})

		err := provider.CheckShardFinalization(context.Background(), 0, 0)
		assert.True(t, errors.Is(err, ErrInvalidAllowedDeltaToFinal))
	})
	t.Run("for metachain it will return nil", func(t *testing.T) {
		t.Parallel()

		provider, _ := NewProxyFinalityProvider(&testsCommon.ProxyStub{})

		err := provider.CheckShardFinalization(context.Background(), core.MetachainShardId, 1)
		assert.Nil(t, err)
	})
	t.Run("get cross check from meta fails", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		provider := createProxyProviderForCheckShardFinalization(0, goodResponseExample, expectedErr, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, expectedErr))
	})
	t.Run("invalid response from meta", func(t *testing.T) {
		provider := createProxyProviderForCheckShardFinalization(0, "invalid", nil, nil)

		err := provider.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, ErrInvalidNonceCrossCheckValueFormat))
	})
	t.Run("get nonce from shard fails", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		provider := createProxyProviderForCheckShardFinalization(0, goodResponseExample, nil, expectedErr)

		err := provider.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, expectedErr))
	})
	t.Run("finalization checks", func(t *testing.T) {
		t.Parallel()

		t.Run("0 difference", func(t *testing.T) {
			nonce := uint64(9166353)
			provider := createProxyProviderForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := provider.CheckShardFinalization(context.Background(), 1, 1)
			assert.Nil(t, err)

			err = provider.CheckShardFinalization(context.Background(), 1, 10)
			assert.Nil(t, err)
		})
		t.Run("10 difference", func(t *testing.T) {
			nonce := uint64(9166353 + 10)
			provider := createProxyProviderForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := provider.CheckShardFinalization(context.Background(), 1, 10)
			assert.Nil(t, err)

			err = provider.CheckShardFinalization(context.Background(), 1, 11)
			assert.Nil(t, err)

			err = provider.CheckShardFinalization(context.Background(), 1, 9)
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "shardID 1 is stuck"))
		})
		t.Run("shard is syncing", func(t *testing.T) {
			nonce := uint64(9166353 - 1)
			provider := createProxyProviderForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := provider.CheckShardFinalization(context.Background(), 1, 9)
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "shardID 1 is syncing"))
		})
	})
}

func createProxyProviderForCheckShardFinalization(
	crtNonce uint64,
	crossCheck string,
	fetchFromMetaError error,
	fetchFromShardError error,
) *proxyFinalityProvider {

	stub := &testsCommon.ProxyStub{
		GetNetworkStatusCalled: func(ctx context.Context, shardID uint32) (*data.NetworkStatus, error) {
			if shardID == core.MetachainShardId {
				if fetchFromMetaError != nil {
					return nil, fetchFromMetaError
				}

				return &data.NetworkStatus{
					CrossCheckBlockHeight: crossCheck,
				}, nil
			}

			if fetchFromShardError != nil {
				return nil, fetchFromShardError
			}

			return &data.NetworkStatus{
				Nonce: crtNonce,
			}, nil
		},
	}

	provider, _ := NewProxyFinalityProvider(stub)

	return provider
}
