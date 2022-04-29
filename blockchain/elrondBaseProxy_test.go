package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const goodResponseExample = "0: 9169897, 1: 9166353, 2: 9170524, "

func createMockArgsElrondBaseProxy() argsElrondBaseProxy {
	return argsElrondBaseProxy{
		httpClientWrapper: &testsCommon.HTTPClientWrapperStub{},
		expirationTime:    time.Second,
	}
}

func TestNewElrondBaseProxy(t *testing.T) {
	t.Parallel()

	t.Run("nil http client wrapper", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = nil
		baseProxy, err := newElrondBaseProxy(args)

		assert.True(t, check.IfNil(baseProxy))
		assert.True(t, errors.Is(err, ErrNilHTTPClientWrapper))
	})
	t.Run("invalid caching duration", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		args.expirationTime = time.Second - time.Nanosecond
		baseProxy, err := newElrondBaseProxy(args)

		assert.True(t, check.IfNil(baseProxy))
		assert.True(t, errors.Is(err, ErrInvalidCacherDuration))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		baseProxy, err := newElrondBaseProxy(args)

		assert.False(t, check.IfNil(baseProxy))
		assert.Nil(t, err)
	})
}

func TestElrondBaseProxy_GetNetworkConfig(t *testing.T) {
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
	response := &data.NetworkConfigResponse{
		Data: struct {
			Config *data.NetworkConfig `json:"config"`
		}{
			Config: expectedReturnedNetworkConfig,
		},
		Error: "",
		Code:  "",
	}
	networkConfigBytes, _ := json.Marshal(response)

	t.Run("cache time expired", func(t *testing.T) {
		t.Parallel()

		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return networkConfigBytes, nil
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		args.expirationTime = minimumCachingInterval * 2
		baseProxy, _ := newElrondBaseProxy(args)
		baseProxy.sinceTimeHandler = func(t time.Time) time.Duration {
			return minimumCachingInterval
		}

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, err)
		require.True(t, wasCalled)
		assert.Equal(t, expectedReturnedNetworkConfig, configs)
	})
	t.Run("fetchedConfigs is nil", func(t *testing.T) {
		t.Parallel()

		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return networkConfigBytes, nil
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		args.expirationTime = minimumCachingInterval * 2
		baseProxy, _ := newElrondBaseProxy(args)
		baseProxy.sinceTimeHandler = func(t time.Time) time.Duration {
			return minimumCachingInterval*2 + time.Millisecond
		}

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, err)
		require.True(t, wasCalled)
		assert.Equal(t, expectedReturnedNetworkConfig, configs)
	})
	t.Run("Proxy.GetNetworkConfig returns error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return nil, expectedError
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		baseProxy, _ := newElrondBaseProxy(args)

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, configs)
		require.True(t, wasCalled)
		assert.Equal(t, expectedError, err)
	})
	t.Run("and Proxy.GetNetworkConfig returns malformed data", func(t *testing.T) {
		t.Parallel()

		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return []byte("malformed data"), nil
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		baseProxy, _ := newElrondBaseProxy(args)

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, configs)
		require.True(t, wasCalled)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "invalid character"))
	})
	t.Run("and Proxy.GetNetworkConfig returns a response error", func(t *testing.T) {
		t.Parallel()

		errMessage := "error message"
		erroredResponse := &data.NetworkConfigResponse{
			Data: struct {
				Config *data.NetworkConfig `json:"config"`
			}{},
			Error: errMessage,
			Code:  "",
		}
		erroredNetworkConfigBytes, _ := json.Marshal(erroredResponse)

		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return erroredNetworkConfigBytes, nil
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		baseProxy, _ := newElrondBaseProxy(args)

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, configs)
		require.True(t, wasCalled)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), errMessage))
	})
	t.Run("getCachedConfigs returns valid fetchedConfigs", func(t *testing.T) {
		t.Parallel()

		mockWrapper := &testsCommon.HTTPClientWrapperStub{}
		wasCalled := false
		mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
			wasCalled = true
			return nil, nil
		}

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = mockWrapper
		args.expirationTime = minimumCachingInterval * 2
		baseProxy, _ := newElrondBaseProxy(args)
		baseProxy.fetchedConfigs = expectedReturnedNetworkConfig
		baseProxy.sinceTimeHandler = func(t time.Time) time.Duration {
			return minimumCachingInterval
		}

		configs, err := baseProxy.GetNetworkConfig(context.Background())

		require.Nil(t, err)
		assert.False(t, wasCalled)
		assert.Equal(t, expectedReturnedNetworkConfig, configs)
	})
}

func TestElrondBaseProxy_GetNetworkStatus(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("get errors", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, error) {
				return nil, expectedErr
			},
		}
		baseProxy, _ := newElrondBaseProxy(args)

		result, err := baseProxy.GetNetworkStatus(context.Background(), 0)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("malformed response", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, error) {
				return []byte("malformed response"), nil
			},
		}
		baseProxy, _ := newElrondBaseProxy(args)

		result, err := baseProxy.GetNetworkStatus(context.Background(), 0)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "invalid character 'm'"))
	})
	t.Run("response error", func(t *testing.T) {
		t.Parallel()

		resp := &data.NetworkStatusResponse{
			Data: struct {
				Status *data.NetworkStatus `json:"status"`
			}{},
			Error: expectedErr.Error(),
			Code:  "",
		}
		respBytes, _ := json.Marshal(resp)

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, error) {
				return respBytes, nil
			},
		}
		baseProxy, _ := newElrondBaseProxy(args)

		result, err := baseProxy.GetNetworkStatus(context.Background(), 0)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), expectedErr.Error()))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		providedNetworkStatus := &data.NetworkStatus{
			CurrentRound:               1,
			EpochNumber:                2,
			Nonce:                      3,
			NonceAtEpochStart:          4,
			NoncesPassedInCurrentEpoch: 5,
			RoundAtEpochStart:          6,
			RoundsPassedInCurrentEpoch: 7,
			RoundsPerEpoch:             8,
			CrossCheckBlockHeight:      "aaa",
		}

		resp := &data.NetworkStatusResponse{
			Data: struct {
				Status *data.NetworkStatus `json:"status"`
			}{Status: providedNetworkStatus},
		}
		respBytes, _ := json.Marshal(resp)

		args := createMockArgsElrondBaseProxy()
		args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, error) {
				return respBytes, nil
			},
		}
		baseProxy, _ := newElrondBaseProxy(args)

		result, err := baseProxy.GetNetworkStatus(context.Background(), 0)
		assert.Nil(t, err)
		assert.Equal(t, providedNetworkStatus, result)
	})
}

func TestElrondBaseProxy_GetShardOfAddress(t *testing.T) {
	t.Parallel()

	t.Run("invalid address", func(t *testing.T) {
		t.Parallel()

		baseProxy := createBaseProxyForGetShardOfAddress(3, nil)

		addrShard1 := "invalid"
		shardID, err := baseProxy.GetShardOfAddress(context.Background(), addrShard1)

		assert.Zero(t, shardID)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "invalid bech32 string length 7"))
	})
	t.Run("get network config errors", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		baseProxy := createBaseProxyForGetShardOfAddress(3, expectedErr)

		addrShard1 := "erd1qqqqqqqqqqqqqpgqzyuaqg3dl7rqlkudrsnm5ek0j3a97qevd8sszj0glf"
		shardID, err := baseProxy.GetShardOfAddress(context.Background(), addrShard1)

		assert.Zero(t, shardID)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("num shards without meta is 0", func(t *testing.T) {
		t.Parallel()

		baseProxy := createBaseProxyForGetShardOfAddress(0, nil)

		addrShard1 := "erd1qqqqqqqqqqqqqpgqzyuaqg3dl7rqlkudrsnm5ek0j3a97qevd8sszj0glf"
		shardID, err := baseProxy.GetShardOfAddress(context.Background(), addrShard1)

		assert.Zero(t, shardID)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "the number of shards must be greater than zero"))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		baseProxy := createBaseProxyForGetShardOfAddress(3, nil)

		addrShard1 := "erd1qqqqqqqqqqqqqpgqzyuaqg3dl7rqlkudrsnm5ek0j3a97qevd8sszj0glf"
		shardID, err := baseProxy.GetShardOfAddress(context.Background(), addrShard1)

		assert.Nil(t, err)
		assert.Equal(t, uint32(1), shardID)

		addrShardMeta := "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"
		shardID, err = baseProxy.GetShardOfAddress(context.Background(), addrShardMeta)

		assert.Nil(t, err)
		assert.Equal(t, core.MetachainShardId, shardID)
	})
}

func createBaseProxyForGetShardOfAddress(numShards uint32, errGet error) *elrondBaseProxy {
	expectedReturnedNetworkConfig := &data.NetworkConfig{
		NumShardsWithoutMeta: numShards,
	}
	response := &data.NetworkConfigResponse{
		Data: struct {
			Config *data.NetworkConfig `json:"config"`
		}{
			Config: expectedReturnedNetworkConfig,
		},
	}
	networkConfigBytes, _ := json.Marshal(response)

	mockWrapper := &testsCommon.HTTPClientWrapperStub{}
	mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
		if errGet != nil {
			return nil, errGet
		}

		return networkConfigBytes, nil
	}

	args := createMockArgsElrondBaseProxy()
	args.httpClientWrapper = mockWrapper
	baseProxy, _ := newElrondBaseProxy(args)

	return baseProxy
}

func TestElrondBaseProxy_CheckShardFinalization(t *testing.T) {
	t.Parallel()

	t.Run("invalid maxNoncesDelta", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		baseProxy, _ := newElrondBaseProxy(args)

		err := baseProxy.CheckShardFinalization(context.Background(), 0, 0)
		assert.True(t, errors.Is(err, ErrInvalidAllowedDeltaToFinal))
	})
	t.Run("for metachain it will return nil", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsElrondBaseProxy()
		baseProxy, _ := newElrondBaseProxy(args)

		err := baseProxy.CheckShardFinalization(context.Background(), core.MetachainShardId, 1)
		assert.Nil(t, err)
	})
	t.Run("get cross check from meta fails", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		baseProxy := createBaseProxyForCheckShardFinalization(0, goodResponseExample, expectedErr, nil)

		err := baseProxy.CheckShardFinalization(context.Background(), 1, 1)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("invalid response from meta", func(t *testing.T) {
		baseProxy := createBaseProxyForCheckShardFinalization(0, "invalid", nil, nil)

		err := baseProxy.CheckShardFinalization(context.Background(), 1, 1)
		assert.True(t, errors.Is(err, ErrInvalidNonceCrossCheckValueFormat))
	})
	t.Run("get nonce from shard fails", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		baseProxy := createBaseProxyForCheckShardFinalization(0, goodResponseExample, nil, expectedErr)

		err := baseProxy.CheckShardFinalization(context.Background(), 1, 1)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("finalization checks", func(t *testing.T) {
		t.Parallel()

		t.Run("0 difference", func(t *testing.T) {
			nonce := uint64(9166353)
			baseProxy := createBaseProxyForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := baseProxy.CheckShardFinalization(context.Background(), 1, 1)
			assert.Nil(t, err)

			err = baseProxy.CheckShardFinalization(context.Background(), 1, 10)
			assert.Nil(t, err)
		})
		t.Run("10 difference", func(t *testing.T) {
			nonce := uint64(9166353 + 10)
			baseProxy := createBaseProxyForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := baseProxy.CheckShardFinalization(context.Background(), 1, 10)
			assert.Nil(t, err)

			err = baseProxy.CheckShardFinalization(context.Background(), 1, 11)
			assert.Nil(t, err)

			err = baseProxy.CheckShardFinalization(context.Background(), 1, 9)
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "shardID 1 is stuck"))
		})
		t.Run("shard is syncing", func(t *testing.T) {
			nonce := uint64(9166353 - 1)
			baseProxy := createBaseProxyForCheckShardFinalization(nonce, goodResponseExample, nil, nil)

			err := baseProxy.CheckShardFinalization(context.Background(), 1, 9)
			assert.NotNil(t, err)
			assert.True(t, strings.Contains(err.Error(), "shardID 1 is syncing"))
		})
	})
}

func createBaseProxyForCheckShardFinalization(
	crtNonce uint64,
	crossCheck string,
	fetchFromMetaError error,
	fetchFromShardError error,
) *elrondBaseProxy {

	mockWrapper := &testsCommon.HTTPClientWrapperStub{}
	mockWrapper.GetHTTPCalled = func(ctx context.Context, endpoint string) ([]byte, error) {
		if strings.Contains(endpoint, fmt.Sprintf("%d", core.MetachainShardId)) {
			if fetchFromMetaError != nil {
				return nil, fetchFromMetaError
			}

			return createNetworkStatusBytes(0, crossCheck), nil
		}

		if fetchFromShardError != nil {
			return nil, fetchFromShardError
		}

		return createNetworkStatusBytes(crtNonce, ""), nil
	}

	args := createMockArgsElrondBaseProxy()
	args.httpClientWrapper = mockWrapper
	baseProxy, _ := newElrondBaseProxy(args)

	return baseProxy
}

func createNetworkStatusBytes(crtNonce uint64, crossCheck string) []byte {
	response := &data.NetworkStatusResponse{
		Data: struct {
			Status *data.NetworkStatus `json:"status"`
		}{
			Status: &data.NetworkStatus{
				Nonce:                 crtNonce,
				CrossCheckBlockHeight: crossCheck,
			},
		},
	}

	networkStatusBytes, _ := json.Marshal(response)

	return networkStatusBytes
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
