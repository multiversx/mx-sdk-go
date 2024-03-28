package native

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-go/storage"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	blockByHashEndpoint = "blocks/%s"
	int64Size           = 64
)

// ArgsNativeAuthServer is the DTO used in the native auth server constructor
type ArgsNativeAuthServer struct {
	HttpClientWrapper authentication.HttpClientWrapper
	TokenHandler      authentication.AuthTokenHandler
	Signer            builders.Signer
	PubKeyConverter   core.PubkeyConverter
	KeyGenerator      crypto.KeyGenerator
	TimestampsCacher  storage.Cacher
}

type authServer struct {
	httpClientWrapper authentication.HttpClientWrapper
	tokenHandler      authentication.AuthTokenHandler
	signer            builders.Signer
	keyGenerator      crypto.KeyGenerator
	pubKeyConverter   core.PubkeyConverter
	getTimeHandler    func() time.Time
	cacher            storage.Cacher
}

// NewNativeAuthServer returns a native authentication server that verifies
// authentication tokens:
// 1. Checks whether the provided signature from tokens corresponds with the provided address for the provided body
// 2. Checks the token expiration status
func NewNativeAuthServer(args ArgsNativeAuthServer) (*authServer, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &authServer{
		httpClientWrapper: args.HttpClientWrapper,
		tokenHandler:      args.TokenHandler,
		signer:            args.Signer,
		keyGenerator:      args.KeyGenerator,
		pubKeyConverter:   args.PubKeyConverter,
		getTimeHandler:    time.Now,
		cacher:            args.TimestampsCacher,
	}, nil
}

func checkArgs(args ArgsNativeAuthServer) error {
	if check.IfNil(args.HttpClientWrapper) {
		return authentication.ErrNilHttpClientWrapper
	}
	if check.IfNil(args.Signer) {
		return authentication.ErrNilSigner
	}
	if check.IfNil(args.KeyGenerator) {
		return crypto.ErrNilKeyGenerator
	}
	if check.IfNil(args.PubKeyConverter) {
		return core.ErrNilPubkeyConverter
	}
	if check.IfNil(args.TokenHandler) {
		return authentication.ErrNilTokenHandler
	}
	if check.IfNil(args.TimestampsCacher) {
		return authentication.ErrNilCacher
	}

	return nil
}

// Validate validates the given accessToken
func (server *authServer) Validate(authToken authentication.AuthToken) error {
	err := server.validateExpiration(authToken)
	if err != nil {
		return err
	}

	err = server.validateSignature(authToken)
	if err != nil {
		return err
	}

	return nil
}

func (server *authServer) validateExpiration(token authentication.AuthToken) error {
	blockHash := token.GetBlockHash()
	timestamp, err := server.getBlockTimestamp(blockHash)
	if err != nil {
		return err
	}

	expireTime := timestamp + token.GetTtl()

	isTokenExpired := server.getTimeHandler().After(time.Unix(expireTime, 0))

	if isTokenExpired {
		return authentication.ErrTokenExpired
	}
	return nil
}

func (server *authServer) getBlockTimestamp(blockHash string) (int64, error) {
	cachedTimestampValue, found := server.cacher.Get([]byte(blockHash))
	if found {
		timestamp, ok := cachedTimestampValue.(int64)
		if !ok {
			return 0, fmt.Errorf("%w while casting timestamp value: %v", authentication.ErrInvalidValue, cachedTimestampValue)
		}

		return timestamp, nil
	}

	block, err := server.getBlockByHash(context.Background(), blockHash)
	if err != nil {
		return 0, err
	}

	intTimestamp := int64(block.Timestamp)

	server.cacher.Put([]byte(blockHash), intTimestamp, int64Size)

	return intTimestamp, nil
}

func (server *authServer) validateSignature(token authentication.AuthToken) error {
	address, err := server.pubKeyConverter.Decode(string(token.GetAddress()))
	if err != nil {
		return err
	}

	pubkey, err := server.keyGenerator.PublicKeyFromByteArray(address)
	if err != nil {
		return err
	}

	unsignedToken := server.tokenHandler.GetUnsignedToken(token)
	signableMessage := server.tokenHandler.GetSignableMessage(token.GetAddress(), unsignedToken)

	err = server.signer.VerifyMessage(signableMessage, pubkey, token.GetSignature())
	if err != nil {
		signableMessageLegacy := server.tokenHandler.GetSignableMessageLegacy(token.GetAddress(), unsignedToken)
		return server.signer.VerifyMessage(signableMessageLegacy, pubkey, token.GetSignature())
	}
	return err
}

// IsInterfaceNil returns true if there is no value under the interface
func (server *authServer) IsInterfaceNil() bool {
	return server == nil
}

func (server *authServer) getBlockByHash(ctx context.Context, hash string) (*data.Block, error) {
	var block data.Block
	buff, code, err := server.httpClientWrapper.GetHTTP(ctx, fmt.Sprintf(blockByHashEndpoint, hash))
	if err != nil || code != http.StatusOK {
		return nil, authentication.CreateHTTPStatusError(code, err)
	}

	err = json.Unmarshal(buff, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}
