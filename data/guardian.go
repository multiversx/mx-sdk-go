package data

import "github.com/multiversx/mx-chain-core-go/data/api"

// GuardianDataResponse holds the guardian data endpoint response
type GuardianDataResponse struct {
	Data struct {
		GuardianData *api.GuardianData `json:"guardianData"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
