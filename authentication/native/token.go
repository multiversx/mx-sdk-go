package native

// AuthToken is the native authentication token implementation
type AuthToken struct {
	ttl       int64
	address   []byte
	extraInfo []byte
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

// GetSignature is the getter to signature member
func (token AuthToken) GetSignature() []byte {
	return token.signature
}

// GetBlockHash is the getter to blockHash member
func (token AuthToken) GetBlockHash() string {
	return token.blockHash
}

// GetExtraInfo is the getter to extraInfo member
func (token AuthToken) GetExtraInfo() []byte {
	return token.extraInfo
}
