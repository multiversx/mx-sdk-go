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
