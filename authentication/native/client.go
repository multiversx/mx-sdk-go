package native

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthClient is the DTO used in the native auth client constructor
type ArgsNativeAuthClient struct {
	Signer               crypto.SingleSigner
	ExtraInfo            struct{}
	Proxy                workflows.ProxyHandler
	PrivateKey           crypto.PrivateKey
	TokenHandler         authentication.AuthTokenHandler
	TokenExpiryInSeconds int64
	Host                 string
}

type authClient struct {
	signer               crypto.SingleSigner
	extraInfo            string
	proxy                workflows.ProxyHandler
	privateKey           crypto.PrivateKey
	tokenExpiryInSeconds int64
	address              []byte
	host                 string
	token                string
	tokenHandler         authentication.AuthTokenHandler
	tokenExpire          time.Time
	getTimeHandler       func() time.Time
}

// NewNativeAuthClient will create a new native client able to create authentication tokens
func NewNativeAuthClient(args ArgsNativeAuthClient) (*authClient, error) {
	if check.IfNil(args.Signer) {
		return nil, builders.ErrNilTxSigner
	}

	extraInfoBytes, err := json.Marshal(args.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("%w while marshaling args.extraInfo", err)
	}
	encodedExtraInfo := base64.StdEncoding.EncodeToString(extraInfoBytes)
	encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))
	if check.IfNil(args.Proxy) {
		return nil, workflows.ErrNilProxy
	}

	if check.IfNil(args.TokenHandler) {
		return nil, authentication.ErrNilTokenHandler
	}

	if check.IfNil(args.PrivateKey) {
		return nil, crypto.ErrNilPrivateKey
	}

	publicKey := args.PrivateKey.GeneratePublic()
	pkBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("%w while getting pkBytes from publicKey", err)
	}

	address := data.NewAddressFromBytes(pkBytes)

	return &authClient{
		signer:               args.Signer,
		extraInfo:            encodedExtraInfo,
		proxy:                args.Proxy,
		privateKey:           args.PrivateKey,
		host:                 encodedHost,
		address:              []byte(address.AddressAsBech32String()),
		tokenHandler:         args.TokenHandler,
		tokenExpiryInSeconds: args.TokenExpiryInSeconds,
		getTimeHandler:       time.Now,
	}, nil
}

// GetAccessToken returns an access token used for authentication into different elrond services
func (nac *authClient) GetAccessToken() (string, error) {
	now := nac.getTimeHandler()
	noToken := nac.tokenExpire.IsZero()
	tokenExpired := now.After(nac.tokenExpire)
	if noToken || tokenExpired {
		err := nac.createNewToken()
		if err != nil {
			return "", err
		}
	}
	return nac.token, nil
}

func (nac *authClient) createNewToken() error {
	nonce, err := nac.proxy.GetLatestHyperBlockNonce(context.Background())
	if err != nil {
		return err
	}

	lastHyperblock, err := nac.proxy.GetHyperBlockByNonce(context.Background(), nonce)
	if err != nil {
		return err
	}

	token := &AuthToken{
		ttl:       nac.tokenExpiryInSeconds,
		host:      nac.host,
		extraInfo: nac.extraInfo,
		blockHash: lastHyperblock.Hash,
		address:   nac.address,
	}

	token.signature, err = nac.signer.Sign(nac.privateKey, token.GetBody())
	if err != nil {
		return err
	}

	nac.token, err = nac.tokenHandler.Encode(token)
	if err != nil {
		return err
	}
	nac.tokenExpire = nac.getTimeHandler().Add(time.Duration(nac.tokenExpiryInSeconds))
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nac *authClient) IsInterfaceNil() bool {
	return nac == nil
}
