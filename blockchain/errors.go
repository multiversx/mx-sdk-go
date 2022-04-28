package blockchain

import "errors"

// ErrInvalidAddress signals that the provided address is invalid
var ErrInvalidAddress = errors.New("invalid address")

// ErrNilAddress signals that the provided address is nil
var ErrNilAddress = errors.New("nil address")

// ErrNilShardCoordinator signals that the provided shard coordinator is nil
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrNilNetworkConfigs signals that the provided network configs is nil
var ErrNilNetworkConfigs = errors.New("nil network configs")

// ErrInvalidCacherDuration signals that the provided caching duration is invalid
var ErrInvalidCacherDuration = errors.New("invalid caching duration")

// ErrInvalidAllowedDeltaToFinal signals that an invalid allowed delta to final value has been provided
var ErrInvalidAllowedDeltaToFinal = errors.New("invalid allowed delta to final value")

// ErrNilHTTPClientWrapper signals that a nil HTTP client wrapper was provided
var ErrNilHTTPClientWrapper = errors.New("nil HTTP client wrapper")

// ErrInvalidNonceCrossCheckValueFormat signals that an invalid nonce cross-check value has been provided
var ErrInvalidNonceCrossCheckValueFormat = errors.New("invalid nonce cross check value format")
