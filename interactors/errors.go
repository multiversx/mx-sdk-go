package interactors

import "errors"

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilTxBuilder signals that a nil transaction builder was provided
var ErrNilTxBuilder = errors.New("nil tx builder")

// ErrInvalidValue signals that an invalid value was provided
var ErrInvalidValue = errors.New("invalid value")

// ErrWrongPassword signals that a wrong password was provided
var ErrWrongPassword = errors.New("wrong password")

// ErrDifferentAccountRecovered signals that a different account was recovered
var ErrDifferentAccountRecovered = errors.New("different account recovered")

// ErrInvalidPemFile signals that an invalid pem file was provided
var ErrInvalidPemFile = errors.New("invalid .PEM file")

// ErrNilAddress signals that the provided address is nil
var ErrNilAddress = errors.New("nil address")

// ErrNilArgCreateTransaction signals that the provided ArgCreateTransaction is nil
var ErrNilArgCreateTransaction = errors.New("nil ArgCreateTransaction")

// ErrNilTransaction signals that provided transaction is nil
var ErrNilTransaction = errors.New("nil transaction")

// ErrTxAlreadySent signals that a transaction was already sent
var ErrTxAlreadySent = errors.New("transaction already sent")

// ErrGapNonce signals that a gap nonce between the lowest nonce of the transactions from the cache and the blockchain nonce has been detected
var ErrGapNonce = errors.New("gap nonce detected")

// ErrNilAddressNonceHandlerConstructorFunc signals that a nil constructor for the address nonce handler was provided
var ErrNilAddressNonceHandlerConstructorFunc = errors.New("nil constructor for AddressNonceHandler")
