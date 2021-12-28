package aggregator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	httpGetVerb = "GET"
)

// HttpResponseGetter wraps over the default http client
type HttpResponseGetter struct {
}

// Get does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *HttpResponseGetter) Get(url string, response interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest(httpGetVerb, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return json.Unmarshal(respBytes, response)
}
