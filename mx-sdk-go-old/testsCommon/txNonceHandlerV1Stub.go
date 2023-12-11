package testsCommon

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/core"
)

// TxNonceHandlerV1Stub -
type TxNonceHandlerV1Stub struct {
	GetNonceCalled          func(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransactionCalled   func(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	ForceNonceReFetchCalled func(address core.AddressHandler) error
	CloseCalled             func() error
}

// GetNonce -
func (stub *TxNonceHandlerV1Stub) GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error) {
	if stub.GetNonceCalled != nil {
		return stub.GetNonceCalled(ctx, address)
	}

	return 0, nil
}

// SendTransaction -
func (stub *TxNonceHandlerV1Stub) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if stub.SendTransactionCalled != nil {
		return stub.SendTransactionCalled(ctx, tx)
	}

	return "", nil
}

// ForceNonceReFetch -
func (stub *TxNonceHandlerV1Stub) ForceNonceReFetch(address core.AddressHandler) error {
	if stub.ForceNonceReFetchCalled != nil {
		return stub.ForceNonceReFetchCalled(address)
	}

	return nil
}

// Close -
func (stub *TxNonceHandlerV1Stub) Close() error {
	if stub.CloseCalled != nil {
		return stub.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxNonceHandlerV1Stub) IsInterfaceNil() bool {
	return stub == nil
}
