package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockhashHandler(t *testing.T) {
	t.Run("invalid polling interval", func(t *testing.T) {
		args := createMockArgsBlockhashHandler()
		args.pollingInterval = time.Millisecond
		handler, err := NewBlockhashHandler(args)

		assert.True(t, check.IfNil(handler))
		assert.Equal(t, err, fmt.Errorf("%w in checkArgs for value PollingInterval", ErrInvalidValue))
	})
	t.Run("invalid block ttl", func(t *testing.T) {
		args := createMockArgsBlockhashHandler()
		args.blockTtl = time.Millisecond
		handler, err := NewBlockhashHandler(args)

		assert.True(t, check.IfNil(handler))
		assert.Equal(t, err, fmt.Errorf("%w in checkArgs for value BlockTtl", ErrInvalidValue))
	})
	t.Run("nil http client wrapper", func(t *testing.T) {
		args := createMockArgsBlockhashHandler()
		args.httpClientWrapper = nil
		handler, err := NewBlockhashHandler(args)

		assert.True(t, check.IfNil(handler))
		assert.Equal(t, err, ErrNilHTTPClientWrapper)
	})
	t.Run("should work", func(t *testing.T) {
		args := createMockArgsBlockhashHandler()
		handler, err := NewBlockhashHandler(args)

		assert.False(t, check.IfNil(handler))
		assert.Nil(t, err)
	})
	t.Run("process loop should start", func(t *testing.T) {
		args := createMockArgsBlockhashHandler()
		handler, _ := NewBlockhashHandler(args)
		time.Sleep(2 * time.Second) // wait for the process loop to start

		assert.True(t, handler.loopStatus.IsSet())
	})
}

func TestGetBlockTimestampByHash_Success(t *testing.T) {
	expectedTimestamp := int(time.Now().Unix())

	args := createMockArgsBlockhashHandler()
	args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
		GetHTTPCalled: func(ctx context.Context, url string) ([]byte, int, error) {
			block := data.Block{Timestamp: expectedTimestamp}
			blockBytes, _ := json.Marshal(&block)
			return blockBytes, http.StatusOK, nil
		},
	}
	bh, _ := NewBlockhashHandler(args)
	actualTimestamp, err := bh.GetBlockTimestampByHash(context.Background(), "blockHash")

	assert.Nil(t, err)
	assert.NotNil(t, actualTimestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}

func TestGetBlockTimestampByHash_HTTPError(t *testing.T) {
	args := createMockArgsBlockhashHandler()
	args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
		GetHTTPCalled: func(ctx context.Context, url string) ([]byte, int, error) {
			return nil, http.StatusBadRequest, fmt.Errorf("Bad Request")
		},
	}
	bh, _ := NewBlockhashHandler(args)

	actualTimestamp, err := bh.GetBlockTimestampByHash(context.Background(), "blockHash")

	assert.NotNil(t, err)
	assert.Equal(t, actualTimestamp, 0)
}

func TestGetBlockTimestampByHash_InvalidJSON(t *testing.T) {
	args := createMockArgsBlockhashHandler()
	args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
		GetHTTPCalled: func(ctx context.Context, url string) ([]byte, int, error) {
			return []byte(`{"timestamp": "invalid"}`), http.StatusOK, nil
		},
	}

	bh, _ := NewBlockhashHandler(args)

	actualTimestamp, err := bh.GetBlockTimestampByHash(context.Background(), "blockHash")

	assert.NotNil(t, err)
	assert.Equal(t, actualTimestamp, 0)
}

func TestGetBlockTimestampByHash_ExpiredBlock(t *testing.T) {
	args := createMockArgsBlockhashHandler()
	args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
		GetHTTPCalled: func(ctx context.Context, url string) ([]byte, int, error) {
			block := data.Block{Timestamp: int(time.Now().Add(-2 * time.Second).Unix())}
			blockBytes, _ := json.Marshal(&block)
			return blockBytes, http.StatusOK, nil
		},
	}

	bh, _ := NewBlockhashHandler(args)
	currentTime := time.Now()
	bh.getTimeHandler = func() time.Time {
		return currentTime.Add(time.Second * 1120)
	}

	timestamp, err := bh.GetBlockTimestampByHash(context.Background(), "blockHash")

	assert.Nil(t, err)
	assert.NotNil(t, timestamp)
	// check bh.blockhashes does not containt the blockHash
	_, ok := bh.blockhashes["blockHash"]
	assert.False(t, ok)
}

func TestGetBlockTimestampByHash_CachedBlock(t *testing.T) {
	blockTimestamp := int(time.Now().Unix() - 10) // block timestamp 10 seconds ago
	args := createMockArgsBlockhashHandler()
	args.httpClientWrapper = &testsCommon.HTTPClientWrapperStub{
		GetHTTPCalled: func(ctx context.Context, url string) ([]byte, int, error) {
			panic("should not call this function")
		},
	}
	bh, _ := NewBlockhashHandler(args)
	bh.blockhashes["blockHash"] = blockTimestamp

	expectedTimestamp := blockTimestamp
	actualTimestamp, err := bh.GetBlockTimestampByHash(context.Background(), "blockHash")

	assert.Nil(t, err)
	assert.NotNil(t, actualTimestamp)
	assert.Equal(t, expectedTimestamp, actualTimestamp)
}

func TestProcessLoop_BlockExpiration(t *testing.T) {
	now := time.Now()
	args := createMockArgsBlockhashHandler()
	bh, _ := NewBlockhashHandler(args)

	bh.getTimeHandler = func() time.Time {
		return now
	}
	bh.blockhashes["block1"] = int(now.Add(-args.blockTtl / 2).Unix())
	bh.blockhashes["block2"] = int(now.Add(-args.blockTtl * 2).Unix())
	bh.blockhashes["block3"] = int(now.Add(-args.blockTtl).Unix())

	// Wait for the loop to run for a few seconds
	time.Sleep(2 * time.Second)

	// Check that expired blocks were deleted from cache
	_, ok1 := bh.blockhashes["block1"]
	_, ok2 := bh.blockhashes["block2"]
	_, ok3 := bh.blockhashes["block3"]
	assert.True(t, ok1)
	assert.False(t, ok2)
	assert.False(t, ok3)
}

func TestProcessLoop_ContextCancellation(t *testing.T) {
	args := createMockArgsBlockhashHandler()
	bh, _ := NewBlockhashHandler(args)

	bh.cancel()
	time.Sleep(time.Second)
	assert.False(t, bh.loopStatus.IsSet())
}

func createMockArgsBlockhashHandler() argsBlockhashHandler {
	return argsBlockhashHandler{
		blockTtl:          time.Minute,
		pollingInterval:   time.Second,
		httpClientWrapper: &testsCommon.HTTPClientWrapperStub{},
	}
}
