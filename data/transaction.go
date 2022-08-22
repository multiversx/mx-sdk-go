package data

import "github.com/ElrondNetwork/elrond-go-core/data/transaction"

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

// Transaction holds the fields of a transaction to be broadcasted to the network
type Transaction struct {
	Nonce             uint64 `json:"nonce"`
	Value             string `json:"value"`
	RcvAddr           string `json:"receiver"`
	SndAddr           string `json:"sender"`
	GasPrice          uint64 `json:"gasPrice,omitempty"`
	GasLimit          uint64 `json:"gasLimit,omitempty"`
	Data              []byte `json:"data,omitempty"`
	Signature         string `json:"signature,omitempty"`
	ChainID           string `json:"chainID"`
	Version           uint32 `json:"version"`
	Options           uint32 `json:"options,omitempty"`
	GuardianAddr      string `json:"guardian,omitempty"`
	GuardianSignature string `json:"guardianSignature,omitempty"`
}

// TransactionStatus holds a transaction's status response from the network
type TransactionStatus struct {
	Data struct {
		Status string `json:"status"`
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
	Type             string                                `json:"type"`
	Hash             string                                `json:"hash"`
	Nonce            uint64                                `json:"nonce"`
	Value            string                                `json:"value"`
	Receiver         string                                `json:"receiver"`
	Sender           string                                `json:"sender"`
	GasPrice         uint64                                `json:"gasPrice"`
	GasLimit         uint64                                `json:"gasLimit"`
	Data             []byte                                `json:"data"`
	Signature        string                                `json:"signature"`
	SourceShard      uint32                                `json:"sourceShard"`
	DestinationShard uint32                                `json:"destinationShard"`
	BlockNonce       uint64                                `json:"blockNonce"`
	BlockHash        string                                `json:"blockHash"`
	MiniblockType    string                                `json:"miniblockType"`
	MiniblockHash    string                                `json:"miniblockHash"`
	Timestamp        uint64                                `json:"timestamp"`
	Status           string                                `json:"status"`
	HyperBlockNonce  uint64                                `json:"hyperblockNonce"`
	HyperBlockHash   string                                `json:"hyperblockHash"`
	ScResults        []*transaction.ApiSmartContractResult `json:"smartContractResults,omitempty"`
	Logs             *transaction.ApiLogs                  `json:"logs,omitempty"`
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

// ArgCreateTransaction will hold the transaction fields
type ArgCreateTransaction struct {
	Nonce             uint64
	Value             string
	RcvAddr           string
	SndAddr           string
	GasPrice          uint64
	GasLimit          uint64
	Data              []byte
	Signature         string
	ChainID           string
	Version           uint32
	Options           uint32
	AvailableBalance  string
	GuardianAddr      string
	GuardianSignature string
}
