package data

import "github.com/multiversx/mx-chain-core-go/data/transaction"

// SendTransactionResponse holds the response received from the network when broadcasting a transaction
type SendTransactionResponse struct {
	Data struct {
		TxHash string `json:"txHash"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// SendTransactionsResponse holds the response received from the network when broadcasting multiple transactions
type SendTransactionsResponse struct {
	Data struct {
		NumOfSentTxs int            `json:"numOfSentTxs"`
		TxsHashes    map[int]string `json:"txsHashes"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// TransactionStatus holds a transaction's status response from the network
type TransactionStatus struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// ProcessedTransactionStatus holds a transaction's processed status response from the network
type ProcessedTransactionStatus struct {
	Data struct {
		ProcessedStatus string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// TransactionInfo holds a transaction info response from the network
type TransactionInfo struct {
	Data struct {
		Transaction TransactionOnNetwork `json:"transaction"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// TransactionOnNetwork holds a transaction's info entry in a hyper block
type TransactionOnNetwork struct {
	Type                              string                                `json:"type"`
	ProcessingTypeOnSource            string                                `json:"processingTypeOnSource,omitempty"`
	ProcessingTypeOnDestination       string                                `json:"processingTypeOnDestination,omitempty"`
	Hash                              string                                `json:"hash"`
	Nonce                             uint64                                `json:"nonce"`
	Value                             string                                `json:"value"`
	Receiver                          string                                `json:"receiver"`
	Sender                            string                                `json:"sender"`
	GasPrice                          uint64                                `json:"gasPrice"`
	GasLimit                          uint64                                `json:"gasLimit"`
	Data                              []byte                                `json:"data"`
	Signature                         string                                `json:"signature"`
	SourceShard                       uint32                                `json:"sourceShard"`
	DestinationShard                  uint32                                `json:"destinationShard"`
	BlockNonce                        uint64                                `json:"blockNonce"`
	BlockHash                         string                                `json:"blockHash"`
	MiniblockType                     string                                `json:"miniblockType"`
	MiniblockHash                     string                                `json:"miniblockHash"`
	Timestamp                         uint64                                `json:"timestamp"`
	Status                            string                                `json:"status"`
	HyperBlockNonce                   uint64                                `json:"hyperblockNonce"`
	HyperBlockHash                    string                                `json:"hyperblockHash"`
	NotarizedAtSourceInMetaNonce      uint64                                `json:"notarizedAtSourceInMetaNonce,omitempty"`
	NotarizedAtSourceInMetaHash       string                                `json:"NotarizedAtSourceInMetaHash,omitempty"`
	NotarizedAtDestinationInMetaNonce uint64                                `json:"notarizedAtDestinationInMetaNonce,omitempty"`
	NotarizedAtDestinationInMetaHash  string                                `json:"notarizedAtDestinationInMetaHash,omitempty"`
	ScResults                         []*transaction.ApiSmartContractResult `json:"smartContractResults,omitempty"`
	Logs                              *transaction.ApiLogs                  `json:"logs,omitempty"`
}

// TxCostResponseData follows the format of the data field of a transaction cost request
type TxCostResponseData struct {
	TxCost     uint64 `json:"txGasUnits"`
	RetMessage string `json:"returnMessage"`
}

// ResponseTxCost defines a response from the node holding the transaction cost
type ResponseTxCost struct {
	Data  TxCostResponseData `json:"data"`
	Error string             `json:"error"`
	Code  string             `json:"code"`
}
