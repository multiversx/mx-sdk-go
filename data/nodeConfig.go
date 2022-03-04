package data

// GenesisNodesResponse
type GenesisNodesResponse struct {
	Data struct {
		Config *GenesisNodes `json:"nodesconfig"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type GenesisNodes struct {
	Eligible map[uint32][][]byte `json:"eligible"`
	Waiting  map[uint32][][]byte `json:"waiting"`
}
