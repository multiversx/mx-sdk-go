package cryptoProvider

import "errors"

// ErrTxAlreadySigned signals that the provided transaction is already signed
var ErrTxAlreadySigned = errors.New("tx already signed")
