package native

import (
	"context"
	"encoding/json"
	"net/http"
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
		lastBlock := &data.HyperBlock{
			Timestamp: uint64(time.Now().Unix()),
			Hash:      "hash",
		}
		proxy := &testsCommon.ProxyStub{
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				return lastBlock, nil
			},
		}

		httpClientWrapper := &testsCommon.HTTPClientWrapperStub{
			GetHTTPCalled: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				block := &data.Block{
					Timestamp: int(lastBlock.Timestamp),
					Hash:      lastBlock.Hash,
				}
				buff, _ := json.Marshal(block)
				return buff, http.StatusOK, nil
			},
		}
		tokenHandler := NewAuthTokenHandler()
		server := createNativeServer(httpClientWrapper, tokenHandler)
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

func createNativeServer(httpClientWrapper authentication.HttpClientWrapper, tokenHandler authentication.AuthTokenHandler) *authServer {
	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, logger.GetOrCreate("testscommon"))

	serverArgs := ArgsNativeAuthServer{
		ApiNetworkAddress: "api.multiversx.com",
		HttpClientWrapper: httpClientWrapper,
		TokenHandler:      tokenHandler,
		Signer:            &testsCommon.SignerStub{},
		PubKeyConverter:   converter,
		KeyGenerator:      keyGen,
	}
	server, _ := NewNativeAuthServer(serverArgs)

	return server
}
