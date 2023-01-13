package blockchain

import (
	"errors"
	"fmt"
	"net/http"
)

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

// ErrHTTPStatusCodeIsNotOK signals that the returned HTTP status code is not OK
var ErrHTTPStatusCodeIsNotOK = errors.New("HTTP status code is not OK")

// ErrNilEndpointProvider signals that a nil endpoint provider was provided
var ErrNilEndpointProvider = errors.New("nil endpoint provider")

// ErrInvalidEndpointProvider signals that an invalid endpoint provider was provided
var ErrInvalidEndpointProvider = errors.New("invalid endpoint provider")

// ErrShardIDMismatch signals that a shard ID mismatch has occurred
var ErrShardIDMismatch = errors.New("shard ID mismatch")

// ErrNilNetworkStatus signals that nil network status was received
var ErrNilNetworkStatus = errors.New("nil network status")

// ErrAddressLengthNotCorrect signals that an account does not have the correct address
var ErrAddressLengthNotCorrect = errors.New("address length is not correct")

// ErrVMTypeLengthIsNotCorrect signals that the vm type length is not correct
var ErrVMTypeLengthIsNotCorrect = errors.New("vm type length is not correct")

func createHTTPStatusError(httpStatusCode int, err error) error {
	if err == nil {
		err = ErrHTTPStatusCodeIsNotOK
	}

	return fmt.Errorf("%w, returned http status: %d, %s",
		err, httpStatusCode, http.StatusText(httpStatusCode))
}
