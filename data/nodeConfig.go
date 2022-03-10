package data

// GenesisNodesResponse holds the network genesis nodes endpoint reponse
type GenesisNodesResponse struct {
	Data struct {
		Nodes *GenesisNodes `json:"nodes"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// GenesisNodes holds the genesis nodes public keys per shard
type GenesisNodes struct {
	Eligible map[uint32][]string `json:"eligible"`
	Waiting  map[uint32][]string `json:"waiting"`
}
