package blockchain

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
)

const txCompleted = "completedTxEvent"
const txFailed = "signalError"

// ProcessTransactionStatus will parse the provided transaction info and return its status accordingly
func ProcessTransactionStatus(txInfo data.TransactionInfo) transaction.TxStatus {
	if txInfo.Data.Transaction.Status != string(transaction.TxStatusSuccess) {
		return transaction.TxStatus(txInfo.Data.Transaction.Status)
	}
	if findIdentifierInLogs(txInfo, txFailed) {
		return transaction.TxStatusFail
	}
	if findIdentifierInLogs(txInfo, txCompleted) {
		return transaction.TxStatusSuccess
	}

	return transaction.TxStatusPending
}

func findIdentifierInLogs(txInfo data.TransactionInfo, identifier string) bool {
	for _, event := range txInfo.Data.Transaction.Logs.Events {
		if event.Identifier == identifier {
			return true
		}
	}

	return false
}
