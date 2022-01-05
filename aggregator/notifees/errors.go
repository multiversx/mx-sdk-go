package notifees

import "errors"

var (
	errNilProxy                  = errors.New("nil proxy")
	errNilTxBuilder              = errors.New("nil tx builder")
	errNilTxNonceHandler         = errors.New("nil tx nonce handler")
	errNilContractAddressHandler = errors.New("nil contract address handler")
	errInvalidContractAddress    = errors.New("invalid contract address")
	errInvalidBaseGasLimit       = errors.New("invalid base gas limit")
	errInvalidGasLimitForEach    = errors.New("invalid gas limit for each price change")
	errNilPrivateKey             = errors.New("nil private key")
)
