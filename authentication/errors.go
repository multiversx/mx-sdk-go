package authentication

import "errors"

// ErrNilTxSigner signals that a nil transaction signer was provided
var ErrNilTxSigner = errors.New("nil transaction signer")

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilPrivateKey signals that a nil private key was provided
var ErrNilPrivateKey = errors.New("nil private key")
