package txcheck

import "errors"

// ErrNilTransaction signals a nil transaction was provided
var ErrNilTransaction = errors.New("nil transaction")

// ErrNilPubKey signals a nil public key was provided
var ErrNilPubKey = errors.New("nil public key")

// ErrNilSignature signals a nil signature was provided
var ErrNilSignature = errors.New("nil signature")

// ErrNilSignatureVerifier signals a nil signature verifier was provided
var ErrNilSignatureVerifier = errors.New("err nil signature verifier")

// ErrNilMarshaller signals a nil marshaller was provided
var ErrNilMarshaller = errors.New("err nil marshaller")

// ErrNilHasher signals a nil hasher was provided
var ErrNilHasher = errors.New("err nil hasher")
