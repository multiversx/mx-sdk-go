package interactors

import "errors"

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilTxSigner signals that a nil transaction signer was provided
var ErrNilTxSigner = errors.New("nil transaction signer")

// ErrInvalidValue signals that an invalid value was provided
var ErrInvalidValue = errors.New("invalid value")

// ErrWrongPassword signals that a wrong password was provided
var ErrWrongPassword = errors.New("wrong password")

// ErrDifferentAccountRecovered signals that a different account was recovered
var ErrDifferentAccountRecovered = errors.New("different account recovered")

// ErrInvalidPemFile signals that an invalid pem file was provided
var ErrInvalidPemFile = errors.New("invalid .PEM file")

// ErrNilMarshalizer signals that the given marshalizer is nil
var ErrNilMarshalizer = errors.New("the given marshalizer is nil")

// ErrNilHasher signals that the given hasher is nil
var ErrNilHasher = errors.New("the given hasher is nil")
