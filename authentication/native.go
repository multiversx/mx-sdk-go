package authentication

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthClient is the DTO used in the native auth client constructor
type ArgsNativeAuthClient struct {
	Signer                 builders.Signer
	ExtraInfo              interface{}
	Proxy                  workflows.ProxyHandler
	CryptoComponentsHolder core.CryptoComponentsHolder
	TokenExpiryInSeconds   uint64
	Host                   string
}

type nativeAuthClient struct {
	signer                 builders.Signer
	encodedExtraInfo       string
	proxy                  workflows.ProxyHandler
	tokenExpiryInSeconds   uint64
	cryptoComponentsHolder core.CryptoComponentsHolder
	encodedHost            string
	token                  string
	tokenExpire            time.Time
	getTimeHandler         func() time.Time
}

// NewNativeAuthClient will create a new native client able to create authentication tokens
func NewNativeAuthClient(args ArgsNativeAuthClient) (*nativeAuthClient, error) {
	if check.IfNil(args.Signer) {
		return nil, ErrNilTxSigner
	}

	extraInfoBytes, err := json.Marshal(args.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("%w while marshaling args.ExtraInfo", err)
	}

	if check.IfNil(args.Proxy) {
		return nil, ErrNilProxy
	}

	if check.IfNil(args.CryptoComponentsHolder) {
		return nil, ErrNilCryptoComponentsHolder
	}

	encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))
	encodedExtraInfo := base64.StdEncoding.EncodeToString(extraInfoBytes)

	return &nativeAuthClient{
		signer:                 args.Signer,
		encodedExtraInfo:       encodedExtraInfo,
		proxy:                  args.Proxy,
		cryptoComponentsHolder: args.CryptoComponentsHolder,
		encodedHost:            encodedHost,
		tokenExpiryInSeconds:   args.TokenExpiryInSeconds,
		getTimeHandler:         time.Now,
	}, nil
}

// GetAccessToken returns an access token used for authentication into different elrond services
func (nac *nativeAuthClient) GetAccessToken() (string, error) {
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

func (nac *nativeAuthClient) createNewToken() error {
	nonce, err := nac.proxy.GetLatestHyperBlockNonce(context.Background())
	if err != nil {
		return err
	}

	lastHyperblock, err := nac.proxy.GetHyperBlockByNonce(context.Background(), nonce)
	if err != nil {
		return err
	}

	token := fmt.Sprintf("%s.%s.%d.%s", nac.encodedHost, lastHyperblock.Hash, nac.tokenExpiryInSeconds, nac.encodedExtraInfo)

	signature, err := nac.signer.SignMessage([]byte(token), nac.cryptoComponentsHolder.GetPrivateKey())
	if err != nil {
		return err
	}

	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))

	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	encodedAddress := base64.StdEncoding.EncodeToString([]byte(nac.cryptoComponentsHolder.GetBech32()))
	nac.token = fmt.Sprintf("%s.%s.%s", encodedAddress, encodedToken, encodedSignature)
	nac.tokenExpire = nac.getTimeHandler().Add(time.Duration(nac.tokenExpiryInSeconds))
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nac *nativeAuthClient) IsInterfaceNil() bool {
	return nac == nil
}
