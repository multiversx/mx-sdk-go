package native

import (
	"context"
	"time"

	goCore "github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthServer is the DTO used in the native auth server constructor
type ArgsNativeAuthServer struct {
	Proxy           workflows.ProxyHandler
	TokenHandler    authentication.AuthTokenHandler
	Signer          crypto.SingleSigner
	PubKeyConverter goCore.PubkeyConverter
	KeyGenerator    crypto.KeyGenerator
	AcceptedHosts   map[string]struct{}
}

type authServer struct {
	proxy           workflows.ProxyHandler
	tokenHandler    authentication.AuthTokenHandler
	signer          crypto.SingleSigner
	keyGenerator    crypto.KeyGenerator
	pubKeyConverter goCore.PubkeyConverter
	acceptedHosts   map[string]struct{}
	getTimeHandler  func() time.Time
}

// NewNativeAuthServer returns a native authentication server
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

	if len(args.AcceptedHosts) == 0 {
		return nil, authentication.ErrNilAcceptedHosts
	}

	return &authServer{
		proxy:           args.Proxy,
		tokenHandler:    args.TokenHandler,
		signer:          args.Signer,
		keyGenerator:    args.KeyGenerator,
		acceptedHosts:   args.AcceptedHosts,
		pubKeyConverter: args.PubKeyConverter,
		getTimeHandler:  time.Now,
	}, nil

}

// Validate validates the given accessToken
func (server *authServer) Validate(accessToken string) (core.AddressHandler, error) {
	token, err := server.tokenHandler.Decode(accessToken)
	if err != nil {
		return nil, err
	}

	_, exists := server.acceptedHosts[token.GetHost()]
	if !exists {
		return nil, authentication.ErrHostNotAccepted
	}

	hyperblock, err := server.proxy.GetHyperBlockByHash(context.Background(), token.GetBlockHash())
	if err != nil {
		return nil, err
	}

	expires := int64(hyperblock.Timestamp) + token.GetTtl()

	isTokenExpired := server.getTimeHandler().After(time.Unix(expires, 0))

	if isTokenExpired {
		return nil, authentication.ErrTokenExpired
	}
	address, err := server.pubKeyConverter.Decode(string(token.GetAddress()))
	if err != nil {
		return nil, err
	}

	pubkey, err := server.keyGenerator.PublicKeyFromByteArray(address)
	if err != nil {
		return nil, err
	}

	err = server.signer.Verify(pubkey, token.GetBody(), token.GetSignature())
	if err != nil {
		return nil, err
	}

	return data.NewAddressFromBytes(address), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (server *authServer) IsInterfaceNil() bool {
	return server == nil
}
