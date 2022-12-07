package native

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication/native/mock"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNativeAuthClient_NewNativeAuthClient(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("nil txsigner should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TxSigner = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, builders.ErrNilTxSigner, err)
	})
	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Proxy = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, workflows.ErrNilProxy, err)
	})
	t.Run("nil private key should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.PrivateKey = nil
		client, err := NewNativeAuthClient(args)
		require.Nil(t, client)
		require.Equal(t, crypto.ErrNilPrivateKey, err)
	})
	t.Run("private key returns error for ToByteArray", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.PrivateKey = &testsCommon.PrivateKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return make([]byte, 0), expectedErr
			}}
		client, err := NewNativeAuthClient(args)

		require.Nil(t, client)
		assert.True(t, errors.Is(err, expectedErr))
		assert.True(t, strings.Contains(err.Error(), "while getting skBytes from args.PrivateKey"))
	})
	t.Run("public key returns error for ToByteArray", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.PrivateKey = &testsCommon.PrivateKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("privateKey"), nil
			},
			GeneratePublicCalled: func() crypto.PublicKey {
				return &testsCommon.PublicKeyStub{
					ToByteArrayCalled: func() ([]byte, error) {
						return make([]byte, 0), expectedErr
					},
				}
			},
		}
		client, err := NewNativeAuthClient(args)

		require.Nil(t, client)
		assert.True(t, errors.Is(err, expectedErr))
		assert.True(t, strings.Contains(err.Error(), "while getting pkBytes from publicKey"))
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
	t.Run("txSigner errors when sign message", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TxSigner = &testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
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
		args.Host = "test.host"
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		publicKeyBytes := []byte("publicKey")
		args.PrivateKey = &testsCommon.PrivateKeyStub{GeneratePublicCalled: func() crypto.PublicKey {
			return &testsCommon.PublicKeyStub{ToByteArrayCalled: func() ([]byte, error) {
				return publicKeyBytes, nil
			}}
		}}
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.TxSigner = &testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
				return []byte(expectedSignature), nil
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
		args.Host = "test.host"
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		publicKeyBytes := []byte("publicKey")
		expectedToken := "token"
		args.PrivateKey = &testsCommon.PrivateKeyStub{GeneratePublicCalled: func() crypto.PublicKey {
			return &testsCommon.PublicKeyStub{ToByteArrayCalled: func() ([]byte, error) {
				return publicKeyBytes, nil
			}}
		}}
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.TxSigner = &testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
				return []byte(expectedSignature), nil
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
		args.Host = "test.host"
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		publicKeyBytes := []byte("publicKey")
		args.PrivateKey = &testsCommon.PrivateKeyStub{GeneratePublicCalled: func() crypto.PublicKey {
			return &testsCommon.PublicKeyStub{ToByteArrayCalled: func() ([]byte, error) {
				return publicKeyBytes, nil
			}}
		}}
		args.Proxy = &testsCommon.ProxyStub{
			GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
				return expectedNonce, nil
			},
			GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
				require.Equal(t, expectedNonce, nonce)
				return &data.HyperBlock{Hash: expectedHash}, nil
			},
		}
		args.TxSigner = &testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
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
		TxSigner:             &testsCommon.TxSignerStub{},
		ExtraInfo:            nil,
		Proxy:                &testsCommon.ProxyStub{},
		PrivateKey:           &testsCommon.PrivateKeyStub{},
		TokenExpiryInSeconds: 0,
		TokenHandler:         &mock.AuthTokenHandlerStub{},
		Host:                 "",
	}
}
