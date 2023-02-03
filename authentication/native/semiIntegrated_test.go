package native

import (
	"context"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/multiversx/mx-sdk-go/workflows"
	"github.com/stretchr/testify/require"
)

var keyGen = signing.NewKeyGenerator(ed25519.NewEd25519())

func TestNativeserver_ClientServer(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		lastHyperBlock := &data.HyperBlock{
			Timestamp: uint64(time.Now().Unix()),
			Hash:      "hash",
		}
		proxy := &testsCommon.ProxyStub{
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				return lastHyperBlock, nil
			},
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return lastHyperBlock, nil
			},
		}
		tokenHandler := NewAuthTokenHandler()
		server := createNativeServer(proxy, tokenHandler)
		alice := createNativeClient(examples.AlicePemContents, proxy, tokenHandler, "host")

		authToken, _ := alice.GetAccessToken()

		tokenDecoded, err := tokenHandler.Decode(authToken)
		require.Nil(t, err)
		err = server.Validate(tokenDecoded)
		require.Nil(t, err)
	})
}

func createNativeClient(pem string, proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler, host string) *authClient {
	w := interactors.NewWallet()
	privateKeyBytes, _ := w.LoadPrivateKeyFromPemData([]byte(pem))
	cryptoCompHolder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKeyBytes)

	clientArgs := ArgsNativeAuthClient{
		Signer:                 cryptoProvider.NewSigner(),
		ExtraInfo:              struct{}{},
		Proxy:                  proxy,
		CryptoComponentsHolder: cryptoCompHolder,
		TokenExpiryInSeconds:   60 * 60 * 24,
		TokenHandler:           tokenHandler,
		Host:                   host,
	}

	client, _ := NewNativeAuthClient(clientArgs)
	return client
}

func createNativeServer(proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler) *authServer {
	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, logger.GetOrCreate("testscommon"))

	serverArgs := ArgsNativeAuthServer{
		Proxy:           proxy,
		TokenHandler:    tokenHandler,
		Signer:          &testsCommon.SignerStub{},
		KeyGenerator:    keyGen,
		PubKeyConverter: converter,
	}
	server, _ := NewNativeAuthServer(serverArgs)

	return server
}
