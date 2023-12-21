package authentication

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrNilTokenHandler signals that a nil token handler has been provided
var ErrNilTokenHandler = errors.New("nil token handler")

// ErrNilSigner signals that a nil signer has been provided
var ErrNilSigner = errors.New("nil signer")

// ErrNilSignature signals that the token has a nil signature
var ErrNilSignature = errors.New("nil token signature")

// ErrNilAddress signals that the token has a nil address
var ErrNilAddress = errors.New("nil token address")

// ErrNilBody signals that the token has a nil body
var ErrNilBody = errors.New("nil token body")

// ErrTokenExpired signals that the provided token is expired
var ErrTokenExpired = errors.New("token expired")

// ErrNilCryptoComponentsHolder signals that a nil cryptoComponentsHolder has been provided
var ErrNilCryptoComponentsHolder = errors.New("nil cryptoComponentsHolder")

// ErrNilHttpClientWrapper signals that a nil http client wrapper was provided
var ErrNilHttpClientWrapper = errors.New("nil http client wrapper")

// ErrHTTPStatusCodeIsNotOK signals that the returned HTTP status code is not OK
var ErrHTTPStatusCodeIsNotOK = errors.New("HTTP status code is not OK")

// ErrNilCacher signals that a nil cacher has been provided
var ErrNilCacher = errors.New("nil cacher")

// ErrInvalidValue signals that an invalid value has been provided
var ErrInvalidValue = errors.New("invalid value")

// CreateHTTPStatusError creates an error with the provided http status code and error
func CreateHTTPStatusError(httpStatusCode int, err error) error {
	if err == nil {
		err = ErrHTTPStatusCodeIsNotOK
	}

	return fmt.Errorf("%w, returned http status: %d, %s",
		err, httpStatusCode, http.StatusText(httpStatusCode))
}
