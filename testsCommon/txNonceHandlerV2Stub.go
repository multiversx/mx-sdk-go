package testsCommon

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxNonceHandlerV2Stub -
type TxNonceHandlerV2Stub struct {
	ApplyNonceCalled        func(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error
	SendTransactionCalled   func(ctx context.Context, tx *data.Transaction) (string, error)
	ForceNonceReFetchCalled func(address core.AddressHandler) error
	CloseCalled             func() error
}

// ApplyNonce -
func (stub *TxNonceHandlerV2Stub) ApplyNonce(ctx context.Context, address core.AddressHandler, txArgs *data.ArgCreateTransaction) error {
	if stub.ApplyNonceCalled != nil {
		return stub.ApplyNonceCalled(ctx, address, txArgs)
	}

	return nil
}

// SendTransaction -
func (stub *TxNonceHandlerV2Stub) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	if stub.SendTransactionCalled != nil {
		return stub.SendTransactionCalled(ctx, tx)
	}

	return "", nil
}

// Close -
func (stub *TxNonceHandlerV2Stub) Close() error {
	if stub.CloseCalled != nil {
		return stub.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxNonceHandlerV2Stub) IsInterfaceNil() bool {
	return stub == nil
}
