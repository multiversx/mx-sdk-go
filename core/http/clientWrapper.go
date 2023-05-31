package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
)

const (
	httpUserAgentKey = "User-Agent"
	httpUserAgent    = "MultiversX/1.0.1"

	httpAcceptTypeKey = "Accept"
	httpAcceptType    = "application/json"

	httpContentTypeKey = "Content-Type"
	httpContentType    = "application/json"
)

type clientWrapper struct {
	url    string
	client Client
}

// NewHttpClientWrapper will create a new instance of type httpClientWrapper
func NewHttpClientWrapper(client Client, url string) *clientWrapper {
	providedClient := client
	if check.IfNilReflect(providedClient) {
		providedClient = http.DefaultClient
	}

	return &clientWrapper{
		url:    url,
		client: providedClient,
	}
}

// GetHTTP does a GET method operation on the specified endpoint
func (wrapper *clientWrapper) GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	applyGetHeaderParams(request)

	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, nil
}

// PostHTTP does a POST method operation on the specified endpoint with the provided raw data bytes
func (wrapper *clientWrapper) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	applyPostHeaderParams(request)

	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	buff, err := ioutil.ReadAll(response.Body)

	return buff, response.StatusCode, err
}

// IsInterfaceNil returns true if there is no value under the interface
func (wrapper *clientWrapper) IsInterfaceNil() bool {
	return wrapper == nil
}

func applyGetHeaderParams(request *http.Request) {
	request.Header.Set(httpAcceptTypeKey, httpAcceptType)
	request.Header.Set(httpUserAgentKey, httpUserAgent)
}

func applyPostHeaderParams(request *http.Request) {
	applyGetHeaderParams(request)
	request.Header.Set(httpContentTypeKey, httpContentType)
}
