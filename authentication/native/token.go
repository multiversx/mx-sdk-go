package native

import (
	"fmt"
)

type NativeAuthToken struct {
	Ttl       int64
	Address   []byte
	Host      string
	ExtraInfo string
	Signature []byte
	BlockHash string
}

// Body returns the authentication token body as string
func (token NativeAuthToken) Body() []byte {
	return []byte(fmt.Sprintf("%s.%s.%d.%s", token.Host, token.BlockHash, token.Ttl, token.ExtraInfo))
}
