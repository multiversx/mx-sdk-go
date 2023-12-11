package data

// HyperBlock holds a hyper block's details
type HyperBlock struct {
	Nonce         uint64 `json:"nonce"`
	Round         uint64 `json:"round"`
	Hash          string `json:"hash"`
	PrevBlockHash string `json:"prevBlockHash"`
	Epoch         uint64 `json:"epoch"`
	NumTxs        uint64 `json:"numTxs"`
	ShardBlocks   []struct {
		Hash  string `json:"hash"`
		Nonce uint64 `json:"nonce"`
		Shard uint32 `json:"shard"`
	} `json:"shardBlocks"`
	Timestamp    uint64 `json:"timestamp"`
	Transactions []TransactionOnNetwork
}

// HyperBlockResponse holds a hyper block info response from the network
type HyperBlockResponse struct {
	Data struct {
		HyperBlock HyperBlock `json:"hyperblock"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}
