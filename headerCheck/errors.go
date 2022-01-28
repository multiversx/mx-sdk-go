package headerCheck

import "errors"

// ErrNilNetworkConfig signals that a nil network config was provided
var ErrNilNetworkConfig = errors.New("nil network config")

// ErrNilRatingsConfig signals that a nil ratings config was provided
var ErrNilRatingsConfig = errors.New("nil ratings config")

// ErrNilEnableEpochsConfig signals that a nil enable epochs config was provided
var ErrNilEnableEpochsConfig = errors.New("nil enable epochs config")
