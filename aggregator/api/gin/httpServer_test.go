package gin

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	apiErrors "github.com/multiversx/mx-sdk-go/aggregator/api/errors"
	testsServer "github.com/multiversx/mx-sdk-go/testsCommon/server"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpServer(t *testing.T) {
	t.Parallel()

	t.Run("nil server should error", func(t *testing.T) {
		t.Parallel()

		hs, err := NewHttpServer(nil)
		assert.Equal(t, apiErrors.ErrNilHttpServer, err)
		assert.True(t, check.IfNil(hs))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hs, err := NewHttpServer(&testsServer.ServerStub{})
		assert.Nil(t, err)
		assert.False(t, check.IfNil(hs))
	})
}

func TestNewHttpServer_Start(t *testing.T) {
	t.Parallel()

	t.Run("ListenAndServe returns closed server", func(t *testing.T) {
		t.Parallel()

		s := &testsServer.ServerStub{
			ListenAndServeCalled: func() error {
				return http.ErrServerClosed
			},
		}

		hs, _ := NewHttpServer(s)
		assert.False(t, check.IfNil(hs))

		hs.Start()
	})
	t.Run("ListenAndServe returns other error", func(t *testing.T) {
		t.Parallel()

		s := &testsServer.ServerStub{
			ListenAndServeCalled: func() error {
				return http.ErrContentLength
			},
		}

		hs, _ := NewHttpServer(s)
		assert.False(t, check.IfNil(hs))

		hs.Start()
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected err")
		s := &testsServer.ServerStub{
			ShutdownCalled: func(ctx context.Context) error {
				return expectedErr
			},
		}
		hs, _ := NewHttpServer(s)
		assert.False(t, check.IfNil(hs))

		hs.Start()

		err := hs.Close()
		assert.Equal(t, expectedErr, err)
	})
}
