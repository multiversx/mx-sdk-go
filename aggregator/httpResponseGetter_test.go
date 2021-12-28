package aggregator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	IntVal    int
	StringVal string
}

func TestHttpResponseGetter_InvalidURLShouldError(t *testing.T) {
	t.Parallel()

	responseGetter := &HttpResponseGetter{}
	responseStruct := &testStruct{}

	err := responseGetter.Get("invalid URL", responseStruct)
	require.NotNil(t, err)
	require.IsType(t, err, &url.Error{})
}

func TestHttpResponseGetter_NilResponseObjectShouldError(t *testing.T) {
	t.Parallel()

	httpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(nil)
	}))
	defer httpServer.Close()

	responseGetter := &HttpResponseGetter{}

	err := responseGetter.Get(httpServer.URL, nil)
	require.NotNil(t, err)
	require.IsType(t, err, &json.SyntaxError{})
}

func TestHttpResponseGetter_InvalidResponseShouldError(t *testing.T) {
	t.Parallel()

	httpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("invalid bytes"))
	}))
	defer httpServer.Close()

	responseGetter := &HttpResponseGetter{}
	err := responseGetter.Get(httpServer.URL, responseGetter)
	require.NotNil(t, err)
	require.IsType(t, err, &json.SyntaxError{})
}

func TestHttpResponseGetter_GetShouldWork(t *testing.T) {
	t.Parallel()

	expectedStruct := &testStruct{
		IntVal:    1232,
		StringVal: "string value",
	}

	httpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		buff, err := json.Marshal(expectedStruct)
		require.Nil(t, err)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(buff)
	}))
	defer httpServer.Close()

	responseGetter := &HttpResponseGetter{}
	responseStruct := &testStruct{}

	err := responseGetter.Get(httpServer.URL, responseStruct)
	require.Nil(t, err)
	require.Equal(t, expectedStruct, responseStruct)
}
