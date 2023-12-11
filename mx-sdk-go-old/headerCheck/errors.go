package headerCheck

import "errors"

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilNodesCoordinator signals that a nil nodes coordinator was provided
var ErrNilNodesCoordinator = errors.New("nil nodes coordinator")

// ErrNilHeaderSigVerifier signals that a nil header sig verifier was provided
var ErrNilHeaderSigVerifier = errors.New("nil header signature verifier")

// ErrNilRawHeaderHandler signals that a nil raw header handler was provided
var ErrNilRawHeaderHandler = errors.New("nil raw header handler")

// ErrNilMarshaller signals that a nil marshaller was provided
var ErrNilMarshaller = errors.New("nil marshaller")
