package fetchers

import "errors"

var (
	errInvalidResponseData    = errors.New("invalid response data")
	errInvalidFetcherName     = errors.New("invalid fetcher name")
	errNilResponseGetter      = errors.New("nil response getter")
	errNilGraphqlGetter       = errors.New("nil graphql getter")
	errNilMaiarTokensMap      = errors.New("nil maiar tokens map")
	errInvalidPair            = errors.New("invalid pair")
	errInvalidGraphqlResponse = errors.New("invalid graphql reponse")
)
