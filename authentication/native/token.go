package native

import (
	"fmt"
)

// AuthToken is the native authentication token implementation
type AuthToken struct {
	ttl       int64
	address   []byte
	host      string
	extraInfo string
	signature []byte
	blockHash string
}

// GetTtl is the getter to ttl member
func (token AuthToken) GetTtl() int64 {
	return token.ttl
}

// GetAddress is the getter to address member
func (token AuthToken) GetAddress() []byte {
	return token.address
}

// GetHost is the getter to host member
func (token AuthToken) GetHost() string {
	return token.host
}

// GetSignature is the getter to signature member
func (token AuthToken) GetSignature() []byte {
	return token.signature
}

// GetBlockHash is the getter to blockHash member
func (token AuthToken) GetBlockHash() string {
	return token.blockHash
}

// GetBody returns the authentication token body as string
func (token AuthToken) GetBody() []byte {
	return []byte(fmt.Sprintf("%s.%s.%d.%s", token.host, token.blockHash, token.ttl, token.extraInfo))
}
