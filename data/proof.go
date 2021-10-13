package data

type ProofResponse struct {
	Data struct {
		Proof    []string `json:"proof"`
		Value    string   `json:"value"`
		RootHash string   `json:"rootHash"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
