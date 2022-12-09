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
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
	"github.com/stretchr/testify/require"
)

var keyGen = signing.NewKeyGenerator(ed25519.NewEd25519())

func TestNativeserver_ClientServer(t *testing.T) {

	t.Run("host not accepted should error", func(t *testing.T) {
		t.Parallel()
		proxy := &testsCommon.ProxyStub{}
		tokenHandler := NewAuthTokenHandler()
		acceptedHosts := make([]string, 1)
		acceptedHosts = append(acceptedHosts, "acceptedHost")
		server := createNativeServer(proxy, tokenHandler, acceptedHosts)
		client := createNativeClient(examples.AlicePemContents, proxy, tokenHandler, "invalid")

		authToken, _ := client.GetAccessToken()

		address, err := server.Validate(authToken)
		require.Equal(t, "", address)
		require.Equal(t, authentication.ErrHostNotAccepted, err)
	})

	t.Run("host accepted should error", func(t *testing.T) {
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
		acceptedHosts := make([]string, 1)
		acceptedHosts = append(acceptedHosts, "acceptedHost")
		server := createNativeServer(proxy, tokenHandler, acceptedHosts)
		alice := createNativeClient(examples.AlicePemContents, proxy, tokenHandler, "acceptedHost")

		authToken, _ := alice.GetAccessToken()

		address, err := server.Validate(authToken)
		require.Nil(t, err)
		require.Equal(t, address, string(alice.address))
	})
}

func createNativeClient(pem string, proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler, host string) *authClient {
	w := interactors.NewWallet()
	privateKeyBytes, _ := w.LoadPrivateKeyFromPemData([]byte(pem))
	privateKey, _ := keyGen.PrivateKeyFromByteArray(privateKeyBytes)

	clientArgs := ArgsNativeAuthClient{
		Signer:               &singlesig.Ed25519Signer{},
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

func createNativeServer(proxy workflows.ProxyHandler, tokenHandler authentication.AuthTokenHandler, acceptedHostsArray []string) *authServer {
	acceptedHosts := make(map[string]struct{})

	for _, acceptedHost := range acceptedHostsArray {
		acceptedHosts[acceptedHost] = struct{}{}
	}

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, logger.GetOrCreate("testscommon"))

	serverArgs := ArgsNativeAuthServer{
		Proxy:           proxy,
		TokenHandler:    tokenHandler,
		Signer:          &singlesig.Ed25519Signer{},
		KeyGenerator:    keyGen,
		AcceptedHosts:   acceptedHosts,
		PubKeyConverter: converter,
	}
	server, _ := NewNativeAuthServer(serverArgs)

	return server
}
