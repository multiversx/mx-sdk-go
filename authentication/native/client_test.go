package native

import (
	"context"
	"testing"
	"time"

	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/authentication/native/mock"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/multiversx/mx-sdk-go/workflows"
	"github.com/stretchr/testify/require"
)

func TestNativeAuthClient_NewNativeAuthClient(t *testing.T) {
	t.Parallel()

	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Signer = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, authentication.ErrNilSigner, err)
	})
	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Proxy = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, workflows.ErrNilProxy, err)
	})
	t.Run("nil crypto components holder should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.CryptoComponentsHolder = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, authentication.ErrNilCryptoComponentsHolder, err)
	})
	t.Run("nil token handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TokenHandler = nil
		server, err := NewNativeAuthClient(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilTokenHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		client, err := NewNativeAuthClient(args)
		require.NotNil(t, client)
		require.Nil(t, err)
		require.False(t, client.IsInterfaceNil())
	})
}

func TestNativeAuthClient_GetAccessToken(t *testing.T) {
	t.Parallel()

	t.Run("proxy returns error for GetLatestHyperBlockNonce", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return 0, expectedErr
			}}
		client, _ := NewNativeAuthClient(args)

		token, err := client.GetAccessToken()
		require.Equal(t, "", token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("proxy returns error for GetHyperBlockByNonce", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				return &data.HyperBlock{}, expectedErr
			},
		}
		client, _ := NewNativeAuthClient(args)

		token, err := client.GetAccessToken()
		require.Equal(t, "", token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("signer errors when sign message", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Signer = &testsCommon.SignerStub{
			SignMessageCalled: func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
				return make([]byte, 0), expectedErr
			},
		}
		client, _ := NewNativeAuthClient(args)

		token, err := client.GetAccessToken()
		require.Equal(t, "", token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("token handler returns error should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		expectedAddr := "addr"
		args.Signer = &testsCommon.SignerStub{
			SignMessageCalled: func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
				return []byte(expectedSignature), nil
			},
		}
		args.CryptoComponentsHolder = &testsCommon.CryptoComponentsHolderStub{
			GetBech32Called: func() string {
				return expectedAddr
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			EncodeCalled: func(authToken authentication.AuthToken) (string, error) {
				return "", expectedErr
			},
		}
		client, _ := NewNativeAuthClient(args)
		client.token = ""

		token, err := client.GetAccessToken()
		require.Equal(t, expectedErr, err)
		require.Equal(t, "", token)
	})
	t.Run("should work, nil token", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		expectedAddr := "addr"
		expectedToken := "token"
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			SignMessageCalled: func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
				return []byte(expectedSignature), nil
			},
		}
		args.CryptoComponentsHolder = &testsCommon.CryptoComponentsHolderStub{
			GetBech32Called: func() string {
				return expectedAddr
			},
		}
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			EncodeCalled: func(authToken authentication.AuthToken) (string, error) {
				return expectedToken, nil
			},
		}
		client, _ := NewNativeAuthClient(args)
		client.token = ""
		client.tokenExpire = time.Time{}

		token, err := client.GetAccessToken()
		require.Nil(t, err)
		require.Equal(t, expectedToken, token)
	})
	t.Run("should work, token expired should generate new one", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			SignMessageCalled: func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
				return []byte(expectedSignature), nil
			},
		}

		client, _ := NewNativeAuthClient(args)

		currentTime := time.Now()
		client.getTimeHandler = func() time.Time {
			return currentTime.Add(time.Second * 1120)
		}
		storedToken := "token"
		client.token = storedToken
		client.tokenExpire = currentTime.Add(time.Second * 1000)
		token, err := client.GetAccessToken()
		require.NotEqual(t, storedToken, token)
		require.Nil(t, err)
	})
}

func createMockArgsNativeAuthClient() ArgsNativeAuthClient {
	return ArgsNativeAuthClient{
		Signer:                 &testsCommon.SignerStub{},
		ExtraInfo:              struct{}{},
		Proxy:                  &testsCommon.ProxyStub{},
		CryptoComponentsHolder: &testsCommon.CryptoComponentsHolderStub{},
		TokenExpiryInSeconds:   0,
		TokenHandler:           &mock.AuthTokenHandlerStub{},
		Host:                   "",
	}
}
