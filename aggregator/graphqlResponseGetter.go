package aggregator

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"golang.org/x/oauth2"
)

// graphqlResponseGetter wraps over the default http client
type graphqlResponseGetter struct {
	authClient authentication.AuthClient
}

type graphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// NewGraphqlResponseGetter returns a new graphql response getter instance
func NewGraphqlResponseGetter(authClient authentication.AuthClient) (*graphqlResponseGetter, error) {
	if check.IfNil(authClient) {
		return nil, ErrNilAuthClient
	}
	return &graphqlResponseGetter{
		authClient: authClient,
	}, nil
}

// Query does a get operation on the specified url and tries to cast the response bytes over the response object through
// the json serializer
func (getter *graphqlResponseGetter) Query(ctx context.Context, url string, query string, variables string) ([]byte, error) {

	accessToken, err := getter.authClient.GetAccessToken()
	if err != nil {
		return nil, err
	}

	client := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		),
	)

	request := graphQLRequest{
		Query:     query,
		Variables: variables,
	}
	gqlMarshalled, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(url, "application/json", strings.NewReader(string(gqlMarshalled)))
	if err != nil {
		return nil, err
	}
	responseBytes, _ := io.ReadAll(resp.Body)

	return responseBytes, nil
}
