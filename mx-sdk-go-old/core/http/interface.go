package http

import "net/http"

// Client is the interface we expect to call in order to do the HTTP requests
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
