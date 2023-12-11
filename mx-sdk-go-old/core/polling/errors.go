package polling

import "errors"

// ErrNilLogger signals that a nil logger was provided
var ErrNilLogger = errors.New("nil logger")

// ErrInvalidValue signals that an invalid value was provided
var ErrInvalidValue = errors.New("invalid value")

// ErrNilExecutor signals that a nil executor instance has been provided
var ErrNilExecutor = errors.New("nil executor")

// ErrLoopAlreadyStarted signals that a loop has already been started
var ErrLoopAlreadyStarted = errors.New("loop already started")
