package native

import (
	"context"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/builders"
)

// ArgsNativeAuthServer is the DTO used in the native auth server constructor
type ArgsNativeAuthServer struct {
	BlockhashHandler authentication.BlockhashHandler
	TokenHandler     authentication.AuthTokenHandler
	Signer           builders.Signer
	PubKeyConverter  core.PubkeyConverter
	KeyGenerator     crypto.KeyGenerator
}

type authServer struct {
	blockhashHandler authentication.BlockhashHandler
	tokenHandler     authentication.AuthTokenHandler
	signer           builders.Signer
	keyGenerator     crypto.KeyGenerator
	pubKeyConverter  core.PubkeyConverter
	getTimeHandler   func() time.Time
}

// NewNativeAuthServer returns a native authentication server that verifies
// authentication tokens:
// 1. Checks whether the provided signature from tokens corresponds with the provided address for the provided body
// 2. Checks the token expiration status
func NewNativeAuthServer(args ArgsNativeAuthServer) (*authServer, error) {
	if check.IfNil(args.BlockhashHandler) {
		return nil, authentication.ErrNilBlockhashHandler
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
		blockhashHandler: args.BlockhashHandler,
		tokenHandler:     args.TokenHandler,
		signer:           args.Signer,
		keyGenerator:     args.KeyGenerator,
		pubKeyConverter:  args.PubKeyConverter,
		getTimeHandler:   time.Now,
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
	block, err := server.blockhashHandler.GetBlockByHash(context.Background(), token.GetBlockHash())
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
