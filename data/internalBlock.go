package data

// RawBlockRespone holds the raw blocks endpoint response
type RawBlockRespone struct {
	Data struct {
		Block []byte `json:"block"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

// RawMiniBlockRespone holds the raw miniblock endpoint respone
type RawMiniBlockRespone struct {
	Data struct {
		MiniBlock []byte `json:"miniblock"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}
