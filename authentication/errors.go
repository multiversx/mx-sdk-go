package authentication

import (
	"errors"
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

// ErrNilBlockhashHandler signals that a nil blockhash handler was provided
var ErrNilBlockhashHandler = errors.New("nil blockhash handler")
