package aggregator

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	httpGetVerb = "GET"
)

// httpResponseGetter wraps over the default http client
type httpResponseGetter struct {
}

// NewHttpResponseGetter returns a new http response getter instance
func NewHttpResponseGetter() (*httpResponseGetter, error) {
	return &httpResponseGetter{}, nil
}

// Get does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *httpResponseGetter) Get(ctx context.Context, url string, response interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, httpGetVerb, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return json.Unmarshal(respBytes, response)
}
