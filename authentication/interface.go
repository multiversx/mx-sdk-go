package authentication

import (
	"context"

	"github.com/multiversx/mx-sdk-go/data"
)

// AuthClient defines the behavior of an authentication client
type AuthClient interface {
	GetAccessToken() (string, error)
	IsInterfaceNil() bool
}

// AuthServer defines the behavior of an authentication server
type AuthServer interface {
	Validate(accessToken AuthToken) error
	IsInterfaceNil() bool
}

// AuthTokenHandler defines the behavior of an authentication token handler
type AuthTokenHandler interface {
	Decode(accessToken string) (AuthToken, error)
	Encode(authToken AuthToken) (string, error)
	GetUnsignedToken(authToken AuthToken) []byte
	GetSignableMessage(address, unsignedToken []byte) []byte
	GetSignableMessageLegacy(address, unsignedToken []byte) []byte
	IsInterfaceNil() bool
}

// AuthToken defines the behavior of an authentication token
type AuthToken interface {
	GetTtl() int64
	GetAddress() []byte
	GetHost() []byte
	GetSignature() []byte
	GetBlockHash() string
	GetExtraInfo() []byte
	IsInterfaceNil() bool
}

// HttpClientWrapper defines the behavior of http client able to make http requests
type HttpClientWrapper interface {
	GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
	IsInterfaceNil() bool
}

// BlockhashHandler defines the behavior of a blockhash handler
type BlockhashHandler interface {
	GetBlockByHash(ctx context.Context, hash string) (*data.Block, error)
	GetBlockByNonce(ctx context.Context, nonce uint64) (*data.Block, error)
	IsInterfaceNil() bool
}
