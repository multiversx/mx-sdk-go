package txcheck

import "errors"

var errNilPubKey = errors.New("nil public key")
var errNilSignature = errors.New("nil signature")
var errNilSignatureVerifier = errors.New("err nil signature verifier")
var errNilMarshaller = errors.New("err nil marshaller")
var errNilHasher = errors.New("err nil hasher")
