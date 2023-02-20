package data

type Block struct {
	Hash             string   `json:"hash"`
	Epoch            int      `json:"epoch"`
	Nonce            int      `json:"nonce"`
	PrevHash         string   `json:"prevHash"`
	Proposer         string   `json:"proposer"`
	PubKeyBitmap     string   `json:"pubKeyBitmap"`
	Round            int      `json:"round"`
	Shard            int      `json:"shard"`
	Size             int      `json:"size"`
	SizeTxs          int      `json:"sizeTxs"`
	StateRootHash    string   `json:"stateRootHash"`
	Timestamp        int      `json:"timestamp"`
	TxCount          int      `json:"txCount"`
	GasConsumed      int      `json:"gasConsumed"`
	GasRefunded      int      `json:"gasRefunded"`
	GasPenalized     int      `json:"gasPenalized"`
	MaxGasLimit      int64    `json:"maxGasLimit"`
	MiniBlocksHashes []string `json:"miniBlocksHashes"`
	Validators       []string `json:"validators"`
}
