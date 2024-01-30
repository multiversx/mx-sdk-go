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

func TestSendTransactionsOneByOne(t *testing.T) {
	t.Parallel()

	var sendTransactionCalled bool
	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(
		&sendTransactionCalled, &getAccountCalled))
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

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	for _, tt := range txs {
		wg.Add(1)
		go func(tt *transaction.FrontendTransaction) {
			defer wg.Done()
			h, err := transactionHandler.SendTransactions(context.Background(), tt)
			require.NoError(t, err, "failed to send transaction")
			require.Equal(t, []string{strconv.FormatUint(tt.Nonce, 10)}, h)
			require.True(t, sendTransactionCalled, "send transaction was not called")
		}(tt)
	}
	wg.Wait()
}

func TestSendTransactionsBulk(t *testing.T) {
	t.Parallel()

	var sendTransactionCalled bool
	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(
		&sendTransactionCalled, &getAccountCalled))
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

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")
	require.True(t, getAccountCalled, "get account was not called")

	txHashes, err := transactionHandler.SendTransactions(context.Background(), txs...)
	require.NoError(t, err, "failed to send transactions as bulk")
	require.Equal(t, mockedNonces(1000), txHashes)
	require.True(t, sendTransactionCalled, "send transaction was not called")
}

func TestSendTransactionsClose(t *testing.T) {
	t.Parallel()

	var sendTransactionCalled bool
	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(
		&sendTransactionCalled, &getAccountCalled))
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

	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		hashes, err := transactionHandler.SendTransactions(context.Background(), txs...)
		require.Empty(t, hashes, "no transaction should be processed")
		require.Equal(t, "context canceled while sending transaction for address erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht", err.Error())
		wg.Done()
	}()
	transactionHandler.Close()
	wg.Wait()
	require.NoError(t, err, "failed to send transactions as bulk")
}

func createMockArgsNonceTransactionsHandlerV3(sendTransactionWasCalled, getAccountCalled *bool) ArgsNonceTransactionsHandlerV3 {
	return ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
				*sendTransactionWasCalled = true
				return strconv.FormatUint(tx.Nonce, 10), nil
			},
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				*getAccountCalled = true
				return &data.Account{}, nil
			},
		},
		PollingInterval: time.Second * 5,
	}
}

func mockedNonces(index int) []string {
	mock := make([]string, index)
	for i := 0; i < index; i++ {
		mock[i] = strconv.Itoa(i)
	}

	return mock
}
