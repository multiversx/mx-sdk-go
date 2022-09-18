package authentication

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestNativeAuthClient_NewNativeAuthClient(t *testing.T) {
	t.Parallel()

	t.Run("nil txsigner should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TxSigner = nil
		authClient, err := NewNativeAuthClient(args)
		require.Nil(t, authClient)
		require.Equal(t, ErrNilTxSigner, err)
	})
	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Proxy = nil
		authClient, err := NewNativeAuthClient(args)
		require.Nil(t, authClient)
		require.Equal(t, ErrNilProxy, err)
	})
	t.Run("nil private key should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.PrivateKey = nil
		authClient, err := NewNativeAuthClient(args)
		require.Nil(t, authClient)
		require.Equal(t, ErrNilPrivateKey, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		authClient, err := NewNativeAuthClient(args)
		require.NotNil(t, authClient)
		require.Nil(t, err)
	})
}

func TestNativeAuthClient_GetAccessToken(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("proxy give errors", func(t *testing.T) {
		t.Parallel()

		t.Run("for GetLatestHyperBlockNonce", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsNativeAuthClient()
			args.Proxy = &testsCommon.ProxyStub{
				GetLatestHyperBlockNonceCalled: func(ctx context.Context) (uint64, error) {
					return 0, expectedErr
				}}
			authClient, _ := NewNativeAuthClient(args)

			token, err := authClient.GetAccessToken()
			require.Equal(t, "", token)
			require.Equal(t, expectedErr, err)
		})
		t.Run("for GetHyperBlockByNonce", func(t *testing.T) {
			t.Parallel()

			args := createMockArgsNativeAuthClient()
			args.Proxy = &testsCommon.ProxyStub{
				GetHyperBlockByNonceCalled: func(ctx context.Context, nonce uint64) (*data.HyperBlock, error) {
					return &data.HyperBlock{}, expectedErr
				},
			}
			authClient, _ := NewNativeAuthClient(args)

			token, err := authClient.GetAccessToken()
			require.Equal(t, "", token)
			require.Equal(t, expectedErr, err)
		})
	})

	t.Run("txSigner errors when sign message", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.TxSigner = &testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
				return make([]byte, 0), expectedErr
			},
		}
		authClient, _ := NewNativeAuthClient(args)

		token, err := authClient.GetAccessToken()
		require.Equal(t, "", token)
		require.Equal(t, expectedErr, err)

	})

	t.Run("should work", func(t *testing.T) {
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
		authClient, _ := NewNativeAuthClient(args)

		encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))
		encodedExtraInfo := base64.StdEncoding.EncodeToString([]byte("null"))
		internalToken := fmt.Sprintf("%s.%s.%d.%s", encodedHost, expectedHash, args.TokenExpiryInSeconds, encodedExtraInfo)
		encodedInternalToken := base64.StdEncoding.EncodeToString([]byte(internalToken))
		encodedAddress := base64.StdEncoding.EncodeToString(publicKeyBytes)
		encodedSignature := base64.StdEncoding.EncodeToString([]byte(expectedSignature))
		token, err := authClient.GetAccessToken()
		require.Nil(t, err)
		require.Equal(t, fmt.Sprintf("%s.%s.%s", encodedAddress, encodedInternalToken, encodedSignature), token)

	})
}

func createMockArgsNativeAuthClient() ArgsNativeAuthClient {
	return ArgsNativeAuthClient{
		TxSigner:             &testsCommon.TxSignerStub{},
		ExtraInfo:            nil,
		Proxy:                &testsCommon.ProxyStub{},
		PrivateKey:           &testsCommon.PrivateKeyStub{},
		TokenExpiryInSeconds: 0,
		Host:                 "",
	}
}
