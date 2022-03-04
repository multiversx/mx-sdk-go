package data

import "github.com/ElrondNetwork/elrond-go/sharding"

// GenesisNodesConfigResponse
type GenesisNodesConfigResponse struct {
	Data struct {
		Config *GenesisNodesConfig `json:"nodesconfig"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

type GenesisNodesConfig struct {
	Eligible map[uint32][]*sharding.InitialNode `json:"eligible"`
	Waiting  map[uint32][]*sharding.InitialNode `json:"waiting"`
}
