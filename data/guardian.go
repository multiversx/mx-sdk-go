package data

// GuardianDataResponse holds the guardian data endpoint response
type GuardianDataResponse struct {
	Data struct {
		GuardianData *GuardianData `json:"guardianData"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// Guardian holds the relevant information for an account guardian
type Guardian struct {
	Address         string `json:"address"`
	ActivationEpoch uint32 `json:"activationEpoch"`
}

// GuardianData holds data relating to the configured guardian(s) and frozen state of an account
type GuardianData struct {
	ActiveGuardian  *Guardian `json:"activeGuardian,omitempty"`
	PendingGuardian *Guardian `json:"pendingGuardian,omitempty"`
	Frozen          bool      `json:"frozen,omitempty"`
}
