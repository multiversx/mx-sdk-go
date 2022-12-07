package native

import (
	"context"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsAuthServer is the DTO used in the native auth server constructor
type ArgsAuthServer struct {
	Proxy         workflows.ProxyHandler
	TokenHandler  authentication.AuthTokenHandler
	Signer        crypto.SingleSigner
	KeyGenerator  crypto.KeyGenerator
	AcceptedHosts map[string]struct{}
}

type authServer struct {
	proxy          workflows.ProxyHandler
	tokenHandler   authentication.AuthTokenHandler
	signer         crypto.SingleSigner
	keyGenerator   crypto.KeyGenerator
	acceptedHosts  map[string]struct{}
	getTimeHandler func() time.Time
}

// NewNativeAuthServer returns a native authentication server
func NewNativeAuthServer(args ArgsAuthServer) (*authServer, error) {
	if check.IfNil(args.Proxy) {
		return nil, workflows.ErrNilProxy
	}

	if check.IfNil(args.Signer) {
		return nil, authentication.ErrNilSigner
	}

	if check.IfNil(args.KeyGenerator) {
		return nil, crypto.ErrNilKeyGenerator
	}

	if check.IfNil(args.TokenHandler) {
		return nil, authentication.ErrNilTokenHandler
	}

	if len(args.AcceptedHosts) == 0 {
		return nil, authentication.ErrNilAcceptedHosts
	}

	return &authServer{
		proxy:          args.Proxy,
		tokenHandler:   args.TokenHandler,
		signer:         args.Signer,
		keyGenerator:   args.KeyGenerator,
		acceptedHosts:  args.AcceptedHosts,
		getTimeHandler: time.Now,
	}, nil

}

// Validate validates the given accessToken
func (server *authServer) Validate(accessToken string) error {
	token, err := server.tokenHandler.Decode(accessToken)
	if err != nil {
		return err
	}

	nativeToken, ok := token.(NativeAuthToken)
	if !ok {
		return authentication.ErrCannotConvertToken
	}

	_, exists := server.acceptedHosts[nativeToken.Host]
	if !exists {
		return authentication.ErrHostNotAccepted
	}

	hyperblock, err := server.proxy.GetHyperBlockByHash(context.Background(), nativeToken.BlockHash)
	if err != nil {
		return err
	}

	expires := int64(hyperblock.Timestamp) + nativeToken.Ttl

	isTokenExpired := server.getTimeHandler().After(time.Unix(expires, 0))

	if isTokenExpired {
		return authentication.ErrTokenExpired
	}

	pubkey, err := server.keyGenerator.PublicKeyFromByteArray(nativeToken.Body())
	if err != nil {
		return err
	}

	return server.signer.Verify(pubkey, nativeToken.Body(), nativeToken.Signature)
}

// IsInterfaceNil returns true if there is no value under the interface
func (server *authServer) IsInterfaceNil() bool {
	return server == nil
}
