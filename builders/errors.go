package builders

import "errors"

// ErrInvalidValue signals that an invalid value was provided
var ErrInvalidValue = errors.New("invalid value")

// ErrNilValue signals that a nil value was provided
var ErrNilValue = errors.New("nil value")

// ErrNilAddress signals that a nil address was provided
var ErrNilAddress = errors.New("nil address handler")

// ErrInvalidAddress signals that an invalid address was provided
var ErrInvalidAddress = errors.New("invalid address handler")

// ErrNilTxSigner signals that a nil transaction signer was provided
var ErrNilTxSigner = errors.New("nil transaction signer")

// ErrMissingSignature signals that a transaction's signature is empty when trying to compute it's hash
var ErrMissingSignature = errors.New("missing signature when computing the transaction's hash")

// ErrMissingGuardianOption signals that the guardian flag is missing in the transaction option field
var ErrMissingGuardianOption = errors.New("guardian flag is missing in the option field")

// ErrGuardianDoesNotMatch signals a mismatch between the configured guardian in tx and the signing guardian address
var ErrGuardianDoesNotMatch = errors.New("configured guardian does not match signing guardian")
