package authentication

import "errors"

// ErrNilTxSigner signals that a nil transaction signer was provided
var ErrNilTxSigner = errors.New("nil transaction signer")

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilCryptoComponentsHolder signals that a nil cryptoComponentsHolder has been provided
var ErrNilCryptoComponentsHolder = errors.New("nil cryptoComponentsHolder")
