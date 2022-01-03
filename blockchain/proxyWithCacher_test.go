package blockchain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewElrondProxyWithCache(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy", func(t *testing.T) {
		t.Parallel()

		nth, err := NewElrondProxyWithCache(nil, time.Minute)
		assert.True(t, nth.IsInterfaceNil())
		assert.Equal(t, ErrNilProxy, err)
	})
	t.Run("invalid caching duration", func(t *testing.T) {
		t.Parallel()

		mockProxy := &testsCommon.ProxyStub{}
		nth, err := NewElrondProxyWithCache(mockProxy, time.Millisecond)
		require.Nil(t, nth)
		assert.Equal(t, ErrInvalidCacherDuration, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		mockProxy := &testsCommon.ProxyStub{}
		nth, err := NewElrondProxyWithCache(mockProxy, minimumCachingInterval*2)

		require.NotNil(t, nth)
		require.Nil(t, err)
	})
}

func TestGetNetworkConfig(t *testing.T) {
	t.Parallel()
	expectedReturnedNetworkConfig := &data.NetworkConfig{
		ChainID:                  "test",
		Denomination:             1,
		GasPerDataByte:           2,
		LatestTagSoftwareVersion: "test",
		MetaConsensusGroup:       3,
		MinGasLimit:              4,
		MinGasPrice:              5,
		MinTransactionVersion:    6,
		NumMetachainNodes:        7,
		NumNodesInShard:          8,
		NumShardsWithoutMeta:     9,
		RoundDuration:            10,
		ShardConsensusGroupSize:  11,
		StartTime:                12,
	}
	t.Run("getCachedConfigs returns nil", func(t *testing.T) {
		t.Parallel()
		t.Run("cache time expired", func(t *testing.T) {
			t.Parallel()
			mockProxy := &testsCommon.ProxyStub{}
			wasCalled := false
			mockProxy.GetNetworkConfigCalled = func() (*data.NetworkConfig, error) {
				wasCalled = true
				return expectedReturnedNetworkConfig, nil
			}

			nth := NewElrondProxyWithCacheWithHandlers(
				mockProxy,
				minimumCachingInterval*2,
				func(t time.Time) time.Duration {
					return minimumCachingInterval
				})

			configs, err := nth.GetNetworkConfig(context.Background())

			require.Nil(t, err)
			require.True(t, wasCalled)
			assert.Equal(t, expectedReturnedNetworkConfig, configs)
		})
		t.Run("fetchedConfigs is nil", func(t *testing.T) {
			t.Parallel()
			mockProxy := &testsCommon.ProxyStub{}
			wasCalled := false
			mockProxy.GetNetworkConfigCalled = func() (*data.NetworkConfig, error) {
				wasCalled = true
				return expectedReturnedNetworkConfig, nil
			}

			nth := NewElrondProxyWithCacheWithHandlers(
				mockProxy,
				minimumCachingInterval*2,
				func(t time.Time) time.Duration {
					return minimumCachingInterval*2 + time.Millisecond
				})
			nth.fetchedConfigs = nil

			configs, err := nth.GetNetworkConfig(context.Background())

			require.Nil(t, err)
			require.True(t, wasCalled)
			assert.Equal(t, expectedReturnedNetworkConfig, configs)
		})
		t.Run("and Proxy.GetNetworkConfig returns error", func(t *testing.T) {
			t.Parallel()
			expectedError := errors.New("expected error")
			mockProxy := &testsCommon.ProxyStub{}
			wasCalled := false
			mockProxy.GetNetworkConfigCalled = func() (*data.NetworkConfig, error) {
				wasCalled = true
				return nil, expectedError
			}

			nth := NewElrondProxyWithCacheWithHandlers(
				mockProxy,
				minimumCachingInterval*2,
				func(t time.Time) time.Duration {
					return minimumCachingInterval*2 + time.Millisecond
				})
			nth.fetchedConfigs = nil

			configs, err := nth.GetNetworkConfig(context.Background())

			require.Nil(t, configs)
			require.True(t, wasCalled)
			assert.Equal(t, expectedError, err)
		})
	})
	t.Run("getCachedConfigs returns valid fetchedConfigs", func(t *testing.T) {
		t.Parallel()
		mockProxy := &testsCommon.ProxyStub{}
		wasCalled := false
		mockProxy.GetNetworkConfigCalled = func() (*data.NetworkConfig, error) {
			wasCalled = true
			return nil, nil
		}

		nth := NewElrondProxyWithCacheWithHandlers(
			mockProxy,
			minimumCachingInterval*2,
			func(t time.Time) time.Duration {
				return minimumCachingInterval
			})
		nth.fetchedConfigs = expectedReturnedNetworkConfig

		configs, err := nth.GetNetworkConfig(context.Background())

		require.Nil(t, err)
		assert.False(t, wasCalled)
		assert.Equal(t, expectedReturnedNetworkConfig, configs)
	})
}
