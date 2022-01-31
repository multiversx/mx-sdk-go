package headerCheck

import "errors"

// ErrNilProxy signals that a nil proxy was provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilNetworkConfig signals that a nil network config was provided
var ErrNilNetworkConfig = errors.New("nil network config")

// ErrNilRatingsConfig signals that a nil ratings config was provided
var ErrNilRatingsConfig = errors.New("nil ratings config")

// ErrNilEnableEpochsConfig signals that a nil enable epochs config was provided
var ErrNilEnableEpochsConfig = errors.New("nil enable epochs config")

// ErrNilNodesCoordinator signals that a nil nodes coordinator was provided
var ErrNilNodesCoordinator = errors.New("nil nodes coordinator")

// ErrNilHeaderSigVerifier signals that a nil header sig verifier was provided
var ErrNilHeaderSigVerifier = errors.New("nil header signature verifier")

// ErrNilRawHeaderHandler signals that a nil raw header handler was provided
var ErrNilRawHeaderHandler = errors.New("nil raw header handler")

// ErrNilMarshaller signals that a nil marshaller was provided
var ErrNilMarshaller = errors.New("nil marshaller")

// ErrNilHasher signals that a nil hasher was provided
var ErrNilHasher = errors.New("nil hasher")

// ErrNilMultiSig signals that a nil multisig was provided
var ErrNilMultiSig = errors.New("nil multisig")

// ErrNilSingleSigner signals that a nil single signer was provided
var ErrNilSingleSigner = errors.New("nil singlesig")

// ErrNilKeyGen signals that a nil key generator was provided
var ErrNilKeyGen = errors.New("nil keygen")
