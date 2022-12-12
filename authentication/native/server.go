package native

import (
	"context"
	"time"

	goCore "github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthServer is the DTO used in the native auth server constructor
type ArgsNativeAuthServer struct {
	Proxy           workflows.ProxyHandler
	TokenHandler    authentication.AuthTokenHandler
	Signer          crypto.SingleSigner
	PubKeyConverter goCore.PubkeyConverter
	KeyGenerator    crypto.KeyGenerator
}

type authServer struct {
	proxy           workflows.ProxyHandler
	tokenHandler    authentication.AuthTokenHandler
	signer          crypto.SingleSigner
	keyGenerator    crypto.KeyGenerator
	pubKeyConverter goCore.PubkeyConverter
	getTimeHandler  func() time.Time
}

// NewNativeAuthServer returns a native authentication server that verifies
// authentication tokens:
// 1. Checks whether the provided signature from tokens corresponds with the provided address for the provided body
// 2. Checks the token expiration status
func NewNativeAuthServer(args ArgsNativeAuthServer) (*authServer, error) {
	if check.IfNil(args.Proxy) {
		return nil, workflows.ErrNilProxy
	}

	if check.IfNil(args.Signer) {
		return nil, authentication.ErrNilSigner
	}

	if check.IfNil(args.KeyGenerator) {
		return nil, crypto.ErrNilKeyGenerator
	}

	if check.IfNil(args.PubKeyConverter) {
		return nil, goCore.ErrNilPubkeyConverter
	}

	if check.IfNil(args.TokenHandler) {
		return nil, authentication.ErrNilTokenHandler
	}

	return &authServer{
		proxy:           args.Proxy,
		tokenHandler:    args.TokenHandler,
		signer:          args.Signer,
		keyGenerator:    args.KeyGenerator,
		pubKeyConverter: args.PubKeyConverter,
		getTimeHandler:  time.Now,
	}, nil

}

// Validate validates the given accessToken
func (server *authServer) Validate(accessToken string) (string, error) {
	token, err := server.tokenHandler.Decode(accessToken)
	if err != nil {
		return "", err
	}

	err = server.validateExpiration(token)
	if err != nil {
		return "", err
	}

	err = server.validateSignature(token)
	if err != nil {
		return "", err
	}

	return string(token.GetAddress()), nil
}

func (server *authServer) validateExpiration(token authentication.AuthToken) error {
	hyperblock, err := server.proxy.GetHyperBlockByHash(context.Background(), token.GetBlockHash())
	if err != nil {
		return err
	}

	expires := int64(hyperblock.Timestamp) + token.GetTtl()

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

	err = server.signer.Verify(pubkey, token.GetBody(), token.GetSignature())
	if err != nil {
		return err
	}
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (server *authServer) IsInterfaceNil() bool {
	return server == nil
}
