package fetchers

import "errors"

var (
	errInvalidResponseData = errors.New("invalid response data")
	errInvalidFetcherName  = errors.New("invalid fetcher name")
	errNilResponseGetter   = errors.New("nil response getter")
)
