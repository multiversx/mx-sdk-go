package testsCommon

import (
	"context"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

// TxNonceHandlerStub -
type TxNonceHandlerStub struct {
	GetNonceCalled          func(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransactionCalled   func(ctx context.Context, tx *data.Transaction) (string, error)
	ForceNonceReFetchCalled func(address core.AddressHandler) error
	CloseCalled             func() error
}

// GetNonce -
func (stub *TxNonceHandlerStub) GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error) {
	if stub.GetNonceCalled != nil {
		return stub.GetNonceCalled(ctx, address)
	}

	return 0, nil
}

// SendTransaction -
func (stub *TxNonceHandlerStub) SendTransaction(ctx context.Context, tx *data.Transaction) (string, error) {
	if stub.SendTransactionCalled != nil {
		return stub.SendTransactionCalled(ctx, tx)
	}

	return "", nil
}

// ForceNonceReFetch -
func (stub *TxNonceHandlerStub) ForceNonceReFetch(address core.AddressHandler) error {
	if stub.ForceNonceReFetchCalled != nil {
		return stub.ForceNonceReFetchCalled(address)
	}

	return nil
}

// Close -
func (stub *TxNonceHandlerStub) Close() error {
	if stub.CloseCalled != nil {
		return stub.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *TxNonceHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
