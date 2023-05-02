package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func loadJsonIntoTransactionInfo(tb testing.TB, path string) data.TransactionInfo {
	txInfo := &data.TransactionInfo{}
	buff, err := ioutil.ReadFile(path)
	require.Nil(tb, err)

	err = json.Unmarshal(buff, txInfo)
	require.Nil(tb, err)

	return *txInfo
}

func TestProcessTransactionStatus(t *testing.T) {
	t.Parallel()

	t.Run("pending new", func(t *testing.T) {
		t.Parallel()

		txInfo := loadJsonIntoTransactionInfo(t, "./testdata/pendingNew.json")
		status := ProcessTransactionStatus(txInfo)
		require.Equal(t, transaction.TxStatusPending, status)
	})
	t.Run("pending executing", func(t *testing.T) {
		t.Parallel()

		txInfo := loadJsonIntoTransactionInfo(t, "./testdata/pendingExecuting.json")
		status := ProcessTransactionStatus(txInfo)
		require.Equal(t, transaction.TxStatusPending, status)
	})
	t.Run("tx info ok", func(t *testing.T) {
		t.Parallel()

		txInfo := loadJsonIntoTransactionInfo(t, "./testdata/finishedOK.json")
		status := ProcessTransactionStatus(txInfo)
		require.Equal(t, transaction.TxStatusSuccess, status)
	})
	t.Run("tx info failed", func(t *testing.T) {
		t.Parallel()

		txInfo := loadJsonIntoTransactionInfo(t, "./testdata/finishedFailed.json")
		status := ProcessTransactionStatus(txInfo)
		require.Equal(t, transaction.TxStatusFail, status)
	})
}
