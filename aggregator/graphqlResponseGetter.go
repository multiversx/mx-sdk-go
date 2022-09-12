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
	authentication.AuthClient
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// Query does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *GraphqlResponseGetter) Query(ctx context.Context, url string, query string, variables string, response interface{}) error {

	accessToken, err := getter.GetAccessToken()
	if err != nil {
		return err
	}

	client := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		),
	)

	gqlMarshalled, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})

	response, err = client.Post(url, "application/json", strings.NewReader(string(gqlMarshalled)))
	return err
}
