package native

import (
	"context"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519/singlesig"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
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
	privateKey, _ := keyGen.PrivateKeyFromByteArray(privateKeyBytes)

	clientArgs := ArgsNativeAuthClient{
		Signer:               blockchain.NewTxSigner(),
		ExtraInfo:            struct{}{},
		Proxy:                proxy,
		PrivateKey:           privateKey,
		TokenExpiryInSeconds: 60 * 60 * 24,
		TokenHandler:         tokenHandler,
		Host:                 host,
	}

	client, _ := NewNativeAuthClient(clientArgs)
	return client
}

func createNativeServer(proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler) *authServer {
	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, logger.GetOrCreate("testscommon"))

	serverArgs := ArgsNativeAuthServer{
		Proxy:           proxy,
		TokenHandler:    tokenHandler,
		Signer:          &singlesig.Ed25519Signer{},
		KeyGenerator:    keyGen,
		PubKeyConverter: converter,
	}
	server, _ := NewNativeAuthServer(serverArgs)

	return server
}
