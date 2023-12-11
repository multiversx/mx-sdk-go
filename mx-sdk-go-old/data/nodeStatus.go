package data

// NodeStatusResponse holds the node status response
type NodeStatusResponse struct {
	Data struct {
		Status *NetworkStatus `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
