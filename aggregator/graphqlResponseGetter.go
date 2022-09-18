package aggregator

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"golang.org/x/oauth2"
)

// GraphqlResponseGetter wraps over the default http client
type GraphqlResponseGetter struct {
	AuthClient authentication.AuthClient
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// Query does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *GraphqlResponseGetter) Query(ctx context.Context, url string, query string, variables string) (interface{}, error) {

	accessToken, err := getter.AuthClient.GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		),
	)

	var request graphQLRequest
	request.Query = query
	request.Variables = variables
	gqlMarshalled, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return client.Post(url, "application/json", strings.NewReader(string(gqlMarshalled)))
}
