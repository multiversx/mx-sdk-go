package authentication

import "errors"

// ErrNilAcceptedHosts signals that a nil or empty accepted hosts map has been provided
var ErrNilAcceptedHosts = errors.New("nil token Address")

// ErrNilTokenHandler signals that a nil token handler has been provided
var ErrNilTokenHandler = errors.New("nil token handler")

// ErrNilSigner signals that a nil signer has been provided
var ErrNilSigner = errors.New("nil signer")

// ErrNilSignature signals that the token has a nil Signature
var ErrNilSignature = errors.New("nil token signature")

// ErrNilAddress signals that the token has a nil Address
var ErrNilAddress = errors.New("nil token address")

// ErrNilBody signals that the token has a nil body
var ErrNilBody = errors.New("nil token body")

// ErrHostNotAccepted signals that the given Host is not accepted
var ErrHostNotAccepted = errors.New("host not accepted")

// ErrTokenExpired signals that the provided token is expired
var ErrTokenExpired = errors.New("token expired")

// ErrCannotConvertToken signals that the given interface{} cannot
// be converted to an access token struct
var ErrCannotConvertToken = errors.New("cannot convert token")
