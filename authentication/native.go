package authentication

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthClient is the DTO used in the native auth client constructor
type ArgsNativeAuthClient struct {
	TxSigner             builders.TxSigner
	ExtraInfo            interface{}
	Proxy                workflows.ProxyHandler
	PrivateKey           crypto.PrivateKey
	TokenExpiryInSeconds uint64
	Host                 string
}

type nativeAuthClient struct {
	txSigner             builders.TxSigner
	encodedExtraInfo     string
	proxy                workflows.ProxyHandler
	skBytes              []byte
	tokenExpiryInSeconds uint64
	encodedAddress       string
	encodedHost          string
	token                string
	tokenExpire          time.Time
}

// NewNativeAuthClient will create a new native client able to create authentication tokens
func NewNativeAuthClient(args ArgsNativeAuthClient) (AuthClient, error) {
	if check.IfNil(args.TxSigner) {
		return nil, ErrNilTxSigner
	}

	extraInfoBytes, err := json.Marshal(args.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("%w while marshaling args.ExtraInfo", err)
	}

	if check.IfNil(args.Proxy) {
		return nil, ErrNilProxy
	}

	if check.IfNil(args.PrivateKey) {
		return nil, ErrNilPrivateKey
	}

	skBytes, err := args.PrivateKey.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("%w while getting skBytes from args.PrivateKey", err)
	}

	publicKey := args.PrivateKey.GeneratePublic()
	pkBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("%w while getting pkBytes from publicKey", err)
	}

	address := data.NewAddressFromBytes(pkBytes)

	encodedAddress := base64.StdEncoding.EncodeToString(address.AddressBytes())
	encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))
	encodedExtraInfo := base64.StdEncoding.EncodeToString(extraInfoBytes)

	return &nativeAuthClient{
		txSigner:             args.TxSigner,
		encodedExtraInfo:     encodedExtraInfo,
		proxy:                args.Proxy,
		skBytes:              skBytes,
		encodedHost:          encodedHost,
		encodedAddress:       encodedAddress,
		tokenExpiryInSeconds: args.TokenExpiryInSeconds,
	}, nil
}

// GetAccessToken returns an access token used for authentication into different elrond services
func (nac *nativeAuthClient) GetAccessToken() (string, error) {
	now := time.Now()
	if nac.tokenExpire.IsZero() || nac.tokenExpire.After(now) {
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

	signature, err := nac.txSigner.SignMessage([]byte(token), nac.skBytes)
	if err != nil {
		return err
	}

	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))

	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	nac.token = fmt.Sprintf("%s.%s.%s", nac.encodedAddress, encodedToken, encodedSignature)
	nac.tokenExpire = time.Now().Add(time.Duration(nac.tokenExpiryInSeconds))
	return nil
}
