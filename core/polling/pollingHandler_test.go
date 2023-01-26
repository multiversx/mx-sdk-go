package polling

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/core/polling/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgs() ArgsPollingHandler {
	return ArgsPollingHandler{
		Log:              logger.GetOrCreate("test"),
		Name:             "test",
		PollingInterval:  time.Millisecond,
		PollingWhenError: time.Millisecond,
		Executor:         &mock.ExecutorStub{},
	}
}

func TestNewPollingHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil logger should error", func(t *testing.T) {
		args := createMockArgs()
		args.Log = nil

		ph, err := NewPollingHandler(args)
		assert.True(t, check.IfNil(ph))
		assert.Equal(t, ErrNilLogger, err)
	})
	t.Run("invalid polling interval should error", func(t *testing.T) {
		args := createMockArgs()
		args.PollingInterval = minimumPollingInterval - time.Nanosecond

		ph, err := NewPollingHandler(args)
		assert.True(t, check.IfNil(ph))
		assert.True(t, errors.Is(err, ErrInvalidValue))
		assert.True(t, strings.Contains(err.Error(), "PollingInterval"))
	})
	t.Run("invalid polling interval when error should error", func(t *testing.T) {
		args := createMockArgs()
		args.PollingWhenError = minimumPollingInterval - time.Nanosecond

		ph, err := NewPollingHandler(args)
		assert.True(t, check.IfNil(ph))
		assert.True(t, errors.Is(err, ErrInvalidValue))
		assert.True(t, strings.Contains(err.Error(), "PollingWhenError"))
	})
	t.Run("nil executor should error", func(t *testing.T) {
		args := createMockArgs()
		args.Executor = nil

		ph, err := NewPollingHandler(args)
		assert.True(t, check.IfNil(ph))
		assert.Equal(t, ErrNilExecutor, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createMockArgs()

		ph, err := NewPollingHandler(args)
		assert.False(t, check.IfNil(ph))
		assert.Nil(t, err)
	})
}

func TestPollingHandler_NotStartedShouldNotCallExecutor(t *testing.T) {
	t.Parallel()

	args := createMockArgs()
	args.Executor = &mock.ExecutorStub{
		ExecuteCalled: func(ctx context.Context) error {
			require.Fail(t, "should have not called execute")

			return nil
		},
	}

	ph, _ := NewPollingHandler(args)
	assert.False(t, ph.IsRunning())

	time.Sleep(time.Second)

	err := ph.Close()
	assert.Nil(t, err)
	assert.False(t, ph.IsRunning())
}

func TestPollingHandler_StartedShouldCallExecuteMultipleTimes(t *testing.T) {
	t.Parallel()

	numCalls := uint32(0)
	args := createMockArgs()
	args.Executor = &mock.ExecutorStub{
		ExecuteCalled: func(ctx context.Context) error {
			atomic.AddUint32(&numCalls, 1)

			return nil
		},
	}

	ph, _ := NewPollingHandler(args)
	assert.False(t, ph.IsRunning())

	err := ph.StartProcessingLoop()
	assert.Nil(t, err)
	assert.True(t, ph.IsRunning())

	time.Sleep(time.Millisecond * 200)
	assert.True(t, atomic.LoadUint32(&numCalls) > 2)

	err = ph.Close()
	assert.Nil(t, err)
	time.Sleep(time.Second)

	assert.False(t, ph.IsRunning())
}

func TestPollingHandler_StartedShouldTerminateExecutorOnClose(t *testing.T) {
	t.Parallel()

	args := createMockArgs()
	args.Executor = &mock.ExecutorStub{
		ExecuteCalled: func(ctx context.Context) error {
			<-ctx.Done()

			return nil
		},
	}

	ph, _ := NewPollingHandler(args)
	assert.False(t, ph.IsRunning())

	err := ph.StartProcessingLoop()
	assert.Nil(t, err)
	assert.True(t, ph.IsRunning())

	time.Sleep(time.Millisecond * 200)

	err = ph.Close()
	assert.Nil(t, err)
	time.Sleep(time.Second)

	assert.False(t, ph.IsRunning())
}

func TestPollingHandler_StartedShouldUseDifferentPollingTimeOnError(t *testing.T) {
	t.Parallel()

	args := createMockArgs()
	args.PollingInterval = time.Hour
	numCalls := uint32(0)
	args.Executor = &mock.ExecutorStub{
		ExecuteCalled: func(ctx context.Context) error {
			atomic.AddUint32(&numCalls, 1)

			return errors.New("expected error")
		},
	}

	ph, _ := NewPollingHandler(args)
	assert.False(t, ph.IsRunning())

	err := ph.StartProcessingLoop()
	assert.Nil(t, err)
	assert.True(t, ph.IsRunning())

	time.Sleep(time.Millisecond * 200)
	assert.True(t, atomic.LoadUint32(&numCalls) > 2)

	err = ph.Close()
	assert.Nil(t, err)
	time.Sleep(time.Second)

	assert.False(t, ph.IsRunning())
}

func TestPollingHandler_StartedMultipleTimesShouldError(t *testing.T) {
	t.Parallel()

	args := createMockArgs()

	ph, _ := NewPollingHandler(args)
	assert.False(t, ph.IsRunning())
	err := ph.StartProcessingLoop()
	assert.Nil(t, err)
	assert.True(t, ph.IsRunning())

	err = ph.StartProcessingLoop()
	assert.Equal(t, ErrLoopAlreadyStarted, err)

	err = ph.Close()
	assert.Nil(t, err)
	time.Sleep(time.Second)

	assert.False(t, ph.IsRunning())
}
