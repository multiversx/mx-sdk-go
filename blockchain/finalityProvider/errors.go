package finalityProvider

import "errors"

// ErrInvalidNonceCrossCheckValueFormat signals that an invalid nonce cross-check value has been provided
var ErrInvalidNonceCrossCheckValueFormat = errors.New("invalid nonce cross check value format")

// ErrInvalidAllowedDeltaToFinal signals that an invalid allowed delta to final value has been provided
var ErrInvalidAllowedDeltaToFinal = errors.New("invalid allowed delta to final value")

// ErrNilProxy signals that a nil proxy has been provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNodeNotStarted signals that the node is not started
var ErrNodeNotStarted = errors.New("node not started")
