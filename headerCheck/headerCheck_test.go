package headerCheck_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerCheck"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/headerCheck/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgHeaderVerifier() headerCheck.ArgsHeaderVerifier {
	return headerCheck.ArgsHeaderVerifier{
		HeaderHandler:     &mock.RawHeaderHandlerStub{},
		HeaderSigVerifier: &mock.HeaderSigVerifierStub{},
		NodesCoordinator:  &mock.NodesCoordinatorStub{},
	}
}

func TestNewHeaderVerifier(t *testing.T) {
	t.Parallel()

	t.Run("nil raw header handler", func(t *testing.T) {
		t.Parallel()

		args := createMockArgHeaderVerifier()
		args.HeaderHandler = nil
		hv, err := headerCheck.NewHeaderVerifier(args)

		assert.True(t, check.IfNil(hv))
		assert.True(t, errors.Is(err, headerCheck.ErrNilRawHeaderHandler))
	})
	t.Run("nil header sig verifier", func(t *testing.T) {
		t.Parallel()

		args := createMockArgHeaderVerifier()
		args.HeaderSigVerifier = nil
		hv, err := headerCheck.NewHeaderVerifier(args)

		assert.True(t, check.IfNil(hv))
		assert.True(t, errors.Is(err, headerCheck.ErrNilHeaderSigVerifier))
	})
	t.Run("nil nodes coordinator", func(t *testing.T) {
		t.Parallel()

		args := createMockArgHeaderVerifier()
		args.NodesCoordinator = nil
		hv, err := headerCheck.NewHeaderVerifier(args)

		assert.True(t, check.IfNil(hv))
		assert.True(t, errors.Is(err, headerCheck.ErrNilNodesCoordinator))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgHeaderVerifier()
		hv, err := headerCheck.NewHeaderVerifier(args)

		assert.False(t, check.IfNil(hv))
		assert.Nil(t, err)
	})
}

func TestNewHeaderVerifier_VerifyHeaderByHash_ShouldFail(t *testing.T) {
	t.Parallel()

	rawHeaderHandler := &mock.RawHeaderHandlerStub{
		GetShardBlockByHashCalled: func(shardID uint32, hash string) (data.HeaderHandler, error) {
			return &block.Header{Epoch: 1}, nil
		},
	}

	expectedErr := errors.New("signature verifier error")
	headerSigVerifier := &mock.HeaderSigVerifierStub{
		VerifySignatureCalled: func(_ data.HeaderHandler) error {
			return expectedErr
		},
	}
	args := createMockArgHeaderVerifier()
	args.HeaderSigVerifier = headerSigVerifier
	args.HeaderHandler = rawHeaderHandler
	hv, err := headerCheck.NewHeaderVerifier(args)
	require.Nil(t, err)

	status, err := hv.VerifyHeaderByHash(context.Background(), 0, "aaaa")
	assert.False(t, status)
	assert.True(t, errors.Is(expectedErr, err))
}

func TestNewHeaderVerifier_VerifyHeaderByHash_ShouldWork(t *testing.T) {
	t.Parallel()

	rawHeaderHandler := &mock.RawHeaderHandlerStub{
		GetShardBlockByHashCalled: func(shardID uint32, hash string) (data.HeaderHandler, error) {
			return &block.Header{Epoch: 1}, nil
		},
	}

	headerSigVerifier := &mock.HeaderSigVerifierStub{
		VerifySignatureCalled: func(_ data.HeaderHandler) error {
			return nil
		},
	}

	args := createMockArgHeaderVerifier()
	args.HeaderHandler = rawHeaderHandler
	args.HeaderSigVerifier = headerSigVerifier
	hv, err := headerCheck.NewHeaderVerifier(args)
	require.Nil(t, err)

	status, err := hv.VerifyHeaderByHash(context.Background(), 0, "aaaa")
	assert.Nil(t, err)
	assert.True(t, status)
}

func TestNewHeaderVerifier_FetchHeaderByHashAndShard_ShouldFail(t *testing.T) {
	t.Parallel()

	t.Run("fail to get meta block", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("fail to fetch meta block")
		rawHeaderHandler := &mock.RawHeaderHandlerStub{
			GetMetaBlockByHashCalled: func(_ string) (data.MetaHeaderHandler, error) {
				return nil, expectedErr
			},
		}
		args := createMockArgHeaderVerifier()
		args.HeaderHandler = rawHeaderHandler
		hv, _ := headerCheck.NewHeaderVerifier(args)

		header, err := hv.FetchHeaderByHashAndShard(context.Background(), core.MetachainShardId, "aaaa")
		assert.Nil(t, header)
		assert.True(t, errors.Is(expectedErr, err))
	})
	t.Run("fail to get shard block", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("fail to fetch shard block")
		rawHeaderHandler := &mock.RawHeaderHandlerStub{
			GetShardBlockByHashCalled: func(_ uint32, _ string) (data.HeaderHandler, error) {
				return nil, expectedErr
			},
		}
		args := createMockArgHeaderVerifier()
		args.HeaderHandler = rawHeaderHandler
		hv, _ := headerCheck.NewHeaderVerifier(args)

		header, err := hv.FetchHeaderByHashAndShard(context.Background(), 0, "aaaa")
		assert.Nil(t, header)
		assert.True(t, errors.Is(expectedErr, err))
	})
}

func TestNewHeaderVerifier_FetchHeaderByHashAndShard_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedMetaBlock := &block.MetaBlock{
		Nonce: 1,
		Epoch: 1,
	}

	expectedShardBlock := &block.Header{
		Nonce: 2,
		Epoch: 2,
	}

	rawHeaderHandler := &mock.RawHeaderHandlerStub{
		GetMetaBlockByHashCalled: func(_ string) (data.MetaHeaderHandler, error) {
			return expectedMetaBlock, nil
		},
		GetShardBlockByHashCalled: func(_ uint32, _ string) (data.HeaderHandler, error) {
			return expectedShardBlock, nil
		},
	}
	args := createMockArgHeaderVerifier()
	args.HeaderHandler = rawHeaderHandler

	hv, err := headerCheck.NewHeaderVerifier(args)
	require.Nil(t, err)

	shardBlock, err := hv.FetchHeaderByHashAndShard(context.Background(), 0, "aaaa")
	assert.Nil(t, err)
	assert.Equal(t, expectedShardBlock, shardBlock)

	metaBlock, err := hv.FetchHeaderByHashAndShard(context.Background(), core.MetachainShardId, "aaaa")
	assert.Nil(t, err)
	assert.Equal(t, expectedMetaBlock, metaBlock)
}
