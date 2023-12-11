package core

// RestAPIEntityType defines the entity that can resolve REST API requests
type RestAPIEntityType string

const (
	// ObserverNode the entity queried is an observer
	ObserverNode RestAPIEntityType = "observer"
	// Proxy the entity queried is a proxy
	Proxy RestAPIEntityType = "proxy"
)
