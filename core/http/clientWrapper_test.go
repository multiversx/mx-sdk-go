package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientWrapper(t *testing.T) {
	t.Parallel()

	httpClientWrapperInstance := NewHttpClientWrapper(nil, "")
	assert.False(t, check.IfNil(httpClientWrapperInstance))
}

func TestClientWrapper_GetHTTP(t *testing.T) {
	t.Parallel()

	response := []byte("response")
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// simulating that the operation takes a lot of time
		time.Sleep(time.Second * 2)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(response)
	}))
	wrapper := NewHttpClientWrapper(nil, testHttpServer.URL)

	t.Run("context done", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		resp, err := wrapper.GetHTTP(ctx, "endpoint")
		assert.Nil(t, resp)
		require.NotNil(t, err)
		assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
		assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		resp, err := wrapper.GetHTTP(context.Background(), "endpoint")
		assert.Equal(t, response, resp)
		assert.Nil(t, err)
	})
}

func TestElrondBaseProxy_PostHTTP(t *testing.T) {
	t.Parallel()

	response := []byte("response")
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// simulating that the operation takes a lot of time
		time.Sleep(time.Second * 2)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(response)
	}))
	wrapper := NewHttpClientWrapper(nil, testHttpServer.URL)

	t.Run("context done", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		resp, err := wrapper.PostHTTP(ctx, "endpoint", nil)
		assert.Nil(t, resp)
		require.NotNil(t, err)
		assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
		assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		resp, err := wrapper.PostHTTP(context.Background(), "endpoint", nil)
		assert.Equal(t, response, resp)
		assert.Nil(t, err)
	})
}
