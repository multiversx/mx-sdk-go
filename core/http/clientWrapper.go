package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
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
func (wrapper *clientWrapper) GetHTTP(ctx context.Context, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// PostHTTP does a POST method operation on the specified endpoint with the provided raw data bytes
func (wrapper *clientWrapper) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "")
	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	return ioutil.ReadAll(response.Body)
}

// IsInterfaceNil returns true if there is no value under the interface
func (wrapper *clientWrapper) IsInterfaceNil() bool {
	return wrapper == nil
}
