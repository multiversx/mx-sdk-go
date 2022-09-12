package authentication

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

// ArgsNativeAuthClient -
type ArgsNativeAuthClient struct {
	TxSigner             blockchain.TxSigner
	ExtraInfo            interface{}
	Proxy                workflows.ProxyHandler
	PrivateKey           crypto.PrivateKey
	TokenExpiryInSeconds uint64
	Host                 string
}

type nativeAuthClient struct {
	txSigner             blockchain.TxSigner
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

	publicKey := args.PrivateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	address := data.NewAddressFromBytes(publicKeyBytes)

	extraInfoBytes, err := json.Marshal(args.ExtraInfo)
	if err != nil {
		return nil, err
	}
	encodedExtraInfo := base64.StdEncoding.EncodeToString(extraInfoBytes)

	encodedAddress := base64.StdEncoding.EncodeToString(address.AddressBytes())
	skBytes, err := args.PrivateKey.ToByteArray()

	encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))

	return &nativeAuthClient{
		txSigner:         args.TxSigner,
		encodedExtraInfo: encodedExtraInfo,
		proxy:            args.Proxy,
		skBytes:          skBytes,
		encodedHost:      encodedHost,
		encodedAddress:   encodedAddress,
	}, nil
}

// GetAccessToken -
func (nac *nativeAuthClient) GetAccessToken() (string, error) {
	now := time.Now()
	if nac.tokenExpire.After(now) {
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

	nac.token = fmt.Sprintf("%s.%s.%s", nac.encodedAddress, encodedToken, signature)
	nac.tokenExpire = time.Now().Add(time.Duration(nac.tokenExpiryInSeconds))
	return nil
}
