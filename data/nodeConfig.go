package data

import "github.com/ElrondNetwork/elrond-go/state"

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

// ValidatorsInfoResponse holds the validators info endpoint reponse
type ValidatorsInfoResponse struct {
	Data struct {
		ValidatorsInfo []*state.ShardValidatorInfo `json:"validators"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
