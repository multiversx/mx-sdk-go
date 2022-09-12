package aggregator

import (
	"context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
)

// HttpResponseGetter wraps over the default http client
type GraphqlResponseGetter struct {
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// Get does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *GraphqlResponseGetter) Get(ctx context.Context, url string, response interface{}) error {
	client := &http.Client{}

	client := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		))

	req, err := http.NewRequestWithContext(ctx, httpGetVerb, url, nil)
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
