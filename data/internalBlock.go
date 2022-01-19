package data

type RawBlockRespone struct {
	Data struct {
		Block []byte `json:"block"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

type RawMiniBlockRespone struct {
	Data struct {
		MiniBlock []byte `json:"miniblock"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}
