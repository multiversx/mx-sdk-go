package native

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/data"
)

// ArgsNativeAuthServer is the DTO used in the native auth server constructor
type ArgsNativeAuthServer struct {
	ApiNetworkAddress string
	HttpClientWrapper authentication.HttpClientWrapper
	TokenHandler      authentication.AuthTokenHandler
	Signer            builders.Signer
	PubKeyConverter   core.PubkeyConverter
	KeyGenerator      crypto.KeyGenerator
}

type authServer struct {
	httpClientWrapper authentication.HttpClientWrapper
	apiNetworkAddress string
	tokenHandler      authentication.AuthTokenHandler
	signer            builders.Signer
	keyGenerator      crypto.KeyGenerator
	pubKeyConverter   core.PubkeyConverter
	getTimeHandler    func() time.Time
}

// NewNativeAuthServer returns a native authentication server that verifies
// authentication tokens:
// 1. Checks whether the provided signature from tokens corresponds with the provided address for the provided body
// 2. Checks the token expiration status
func NewNativeAuthServer(args ArgsNativeAuthServer) (*authServer, error) {
	if len(args.ApiNetworkAddress) == 0 {
		return nil, authentication.ErrEmptyApiNetworkAddress
	}

	if check.IfNil(args.HttpClientWrapper) {
		return nil, authentication.ErrNilHttpClientWrapper
	}

	if check.IfNil(args.Signer) {
		return nil, authentication.ErrNilSigner
	}

	if check.IfNil(args.KeyGenerator) {
		return nil, crypto.ErrNilKeyGenerator
	}

	if check.IfNil(args.PubKeyConverter) {
		return nil, core.ErrNilPubkeyConverter
	}

	if check.IfNil(args.TokenHandler) {
		return nil, authentication.ErrNilTokenHandler
	}

	return &authServer{
		apiNetworkAddress: args.ApiNetworkAddress,
		httpClientWrapper: args.HttpClientWrapper,
		tokenHandler:      args.TokenHandler,
		signer:            args.Signer,
		keyGenerator:      args.KeyGenerator,
		pubKeyConverter:   args.PubKeyConverter,
		getTimeHandler:    time.Now,
	}, nil

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
	block, err := server.getBlockByHash(context.Background(), token.GetBlockHash())
	if err != nil {
		return err
	}

	expires := int64(block.Timestamp) + token.GetTtl()

	isTokenExpired := server.getTimeHandler().After(time.Unix(expires, 0))

	if isTokenExpired {
		return authentication.ErrTokenExpired
	}
	return nil
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
	return server.signer.VerifyMessage(signableMessage, pubkey, token.GetSignature())
}

// IsInterfaceNil returns true if there is no value under the interface
func (server *authServer) IsInterfaceNil() bool {
	return server == nil
}

func (server *authServer) getBlockByHash(ctx context.Context, hash string) (*data.Block, error) {
	var block data.Block
	buff, code, err := server.httpClientWrapper.GetHTTP(ctx, "blocks/"+hash)
	if err != nil || code != http.StatusOK {
		return nil, authentication.CreateHTTPStatusError(code, err)
	}

	err = json.Unmarshal(buff, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}
