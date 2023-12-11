package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
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

		assert.Equal(t, httpAcceptType, req.Header.Get(httpAcceptTypeKey))
		assert.Equal(t, "", req.Header.Get(httpContentTypeKey)) // this is not set on a GET request
		assert.Equal(t, httpUserAgent, req.Header.Get(httpUserAgentKey))
	}))
	wrapper := NewHttpClientWrapper(nil, testHttpServer.URL)

	t.Run("context done", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		resp, code, err := wrapper.GetHTTP(ctx, "endpoint")
		assert.Nil(t, resp)
		require.NotNil(t, err)
		assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
		assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
		assert.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		resp, code, err := wrapper.GetHTTP(context.Background(), "endpoint")
		assert.Equal(t, response, resp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, code)
	})
}

func TestClientWrapper_PostHTTP(t *testing.T) {
	t.Parallel()

	response := []byte("response")
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// simulating that the operation takes a lot of time
		time.Sleep(time.Second * 2)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(response)

		assert.Equal(t, httpAcceptType, req.Header.Get(httpAcceptTypeKey))
		assert.Equal(t, httpContentType, req.Header.Get(httpContentTypeKey))
		assert.Equal(t, httpUserAgent, req.Header.Get(httpUserAgentKey))
	}))
	wrapper := NewHttpClientWrapper(nil, testHttpServer.URL)

	t.Run("context done", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		resp, code, err := wrapper.PostHTTP(ctx, "endpoint", nil)
		assert.Nil(t, resp)
		require.NotNil(t, err)
		assert.Equal(t, "*url.Error", fmt.Sprintf("%T", err))
		assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
		assert.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		resp, code, err := wrapper.PostHTTP(context.Background(), "endpoint", nil)
		assert.Equal(t, response, resp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, code)
	})
}
