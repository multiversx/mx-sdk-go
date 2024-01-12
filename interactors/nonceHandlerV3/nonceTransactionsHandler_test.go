package nonceHandlerV3

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
)

var testAddressAsBech32String = "erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht"
var testAddress, _ = data.NewAddressFromBech32String(testAddressAsBech32String)

func TestDefault(t *testing.T) {
	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3())
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	for i := 0; i < 1000; i++ {
		tx := &transaction.FrontendTransaction{
			Sender:   testAddressAsBech32String,
			Receiver: testAddressAsBech32String,
			GasLimit: 50000,
			ChainID:  "T",
			Value:    "5000000000000000000",
			Nonce:    uint64(i),
			GasPrice: 1000000000,
			Version:  2,
		}
		txs = append(txs, tx)
	}

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), testAddress, txs...)
	require.NoError(t, err, "failed to apply nonce")

	var wg sync.WaitGroup
	for _, tt := range txs {
		wg.Add(1)
		go func(tt *transaction.FrontendTransaction) {
			defer wg.Done()
			h, err := transactionHandler.SendTransaction(context.Background(), tt)
			require.NoError(t, err, "failed to send transaction")
			require.Equal(t, strconv.FormatUint(tt.Nonce, 10), h)
		}(tt)
	}
	wg.Wait()
}

func createMockArgsNonceTransactionsHandlerV3() ArgsNonceTransactionsHandlerV3 {
	return ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
				return strconv.FormatUint(tx.Nonce, 10), nil
			},
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				return &data.Account{}, nil
			},
		},
		PollingInterval: time.Second * 5,
	}
}
