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

func TestSendTransactionsOneByOne(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
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
		}(tt)
	}
	wg.Wait()
}

func TestSendTransactionsBulk(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
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
	require.Equal(t, mockedHashes(1000), txHashes)
}

func TestSendTransactionsCloseInstant(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(createMockArgsNonceTransactionsHandlerV3(&getAccountCalled))
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	// Create 1k transactions.
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

	// Apply nonce to them in a bulk.
	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")

	// We only do this once, we check if the bool has been modified.
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		hashes, err := transactionHandler.SendTransactions(context.Background(), txs...)

		var counter int
		// Since the close is almost instant there should be none or very few transactions that have been processed.
		for _, h := range hashes {
			if h != "" {
				counter++
			}
		}

		require.Less(t, counter, 5, "more than 5 transactions have been processed.")
		require.Equal(t, "context canceled while sending transaction for address erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht", err.Error())
		wg.Done()
	}()

	// Close the processes related to the transaction handler.
	transactionHandler.Close()

	wg.Wait()
	require.NoError(t, err, "failed to send transactions as bulk")
}

func TestSendTransactionsCloseDelay(t *testing.T) {
	t.Parallel()

	var getAccountCalled bool

	// Create another proxyStub that adds some delay when sending transactions.
	mockArgs := ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
				// Presume this operation is taking roughly 100 ms. Meaning 10 operations / second.
				time.Sleep(100 * time.Millisecond)
				return strconv.FormatUint(tx.Nonce, 10), nil
			},
			GetAccountCalled: func(address core.AddressHandler) (*data.Account, error) {
				getAccountCalled = true
				return &data.Account{}, nil
			},
		},
		PollingInterval: time.Second * 5,
	}

	// Since the endpoint to send workers for the nonce-management-service has the same definition as the one
	// in the gateway, we can create a proxy instance that points towards the nonce-management-service instead.
	// The nonce-management-service will then, in turn send the workers to the gateway.
	transactionHandler, err := NewNonceTransactionHandlerV3(mockArgs)
	require.NoError(t, err, "failed to create transaction handler")

	var txs []*transaction.FrontendTransaction

	// Create 1k transactions.
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

	// Apply nonce to them in a bulk.
	err = transactionHandler.ApplyNonceAndGasPrice(context.Background(), txs...)
	require.NoError(t, err, "failed to apply nonce")

	// We only do this once, we check if the bool has been modified.
	require.True(t, getAccountCalled, "get account was not called")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		hashes, err := transactionHandler.SendTransactions(context.Background(), txs...)

		// Since the close is not instant. There should be some hashes that have been processed.
		require.NotEmpty(t, hashes, "no transaction should be processed")
		require.Equal(t, "context canceled while sending transaction for address erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht", err.Error())
		wg.Done()
	}()

	// Close the processes related to the transaction handler with a delay.
	time.AfterFunc(2*time.Second, func() {
		transactionHandler.Close()
	})

	wg.Wait()
	require.NoError(t, err, "failed to send transactions as bulk")
}

func createMockArgsNonceTransactionsHandlerV3(getAccountCalled *bool) ArgsNonceTransactionsHandlerV3 {
	return ArgsNonceTransactionsHandlerV3{
		Proxy: &testsCommon.ProxyStub{
			SendTransactionCalled: func(tx *transaction.FrontendTransaction) (string, error) {
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

func mockedHashes(index int) []string {
	mock := make([]string, index)
	for i := 0; i < index; i++ {
		mock[i] = strconv.Itoa(i)
	}

	return mock
}
