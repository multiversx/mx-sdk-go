package core

import (
	"github.com/multiversx/mx-chain-core-go/core"
)

// RestAPIEntityType defines the entity that can resolve REST API requests
type RestAPIEntityType string

const (
	// ObserverNode the entity queried is an observer
	ObserverNode RestAPIEntityType = "observer"
	// Proxy the entity queried is a proxy
	Proxy RestAPIEntityType = "proxy"
)

type FilterQuery struct {
	BlockHash []byte              // return logs only from block with this hash
	FromBlock core.OptionalUint64 // beginning of the queried range, no value set means genesis block
	ToBlock   core.OptionalUint64 // end of the range, no value set means latest block
	Addresses []string            // restricts matches to events created by specific contracts
	ShardID   uint32              // identifies the shard to query

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Events are only returned if they match all topics. The order of the topics is not important.
	Topics [][]byte // Topics is a slice of arrays of 32 bytes each
}
