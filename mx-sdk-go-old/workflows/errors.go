package workflows

import "errors"

// ErrInvalidTransactionValue signals that an invalid transaction value was provided
var ErrInvalidTransactionValue = errors.New("invalid transaction value")

// ErrInvalidAvailableBalanceValue signals that an invalid available balance value was provided
var ErrInvalidAvailableBalanceValue = errors.New("invalid available balance value")

// ErrNilTrackableAddressesProvider signals that a nil trackable address provider was used
var ErrNilTrackableAddressesProvider = errors.New("nil trackable address provider")

// ErrNilProxy signals that a nil proxy has been provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNilLastProcessedNonceHandler signals that a nil last processed nonce handler was provided
var ErrNilLastProcessedNonceHandler = errors.New("nil last processed nonce handler")

// ErrNilMinimumBalance signals that a nil minimum balance was provided
var ErrNilMinimumBalance = errors.New("nil minimum balance")

// ErrNilTransactionInteractor signals that a nil transaction interactor was provided
var ErrNilTransactionInteractor = errors.New("nil transaction interactor")
