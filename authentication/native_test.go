package authentication

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestNativeAuthClient_NewNativeAuthClient(t *testing.T) {
	t.Parallel()

	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Signer = nil
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
	t.Run("nil crypto components holder should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.CryptoComponentsHolder = nil
		authClient, err := NewNativeAuthClient(args)
		require.Nil(t, authClient)
		require.Equal(t, ErrNilCryptoComponentsHolder, err)
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
	t.Run("proxy returns error for GetLatestHyperBlockNonce", func(t *testing.T) {
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
	t.Run("proxy returns error for GetHyperBlockByNonce", func(t *testing.T) {
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
	t.Run("signer errors when sign message", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Signer = &testsCommon.SignerStub{
			SignMessageCalled: func(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
				return make([]byte, 0), expectedErr
			},
		}
		authClient, _ := NewNativeAuthClient(args)

		token, err := authClient.GetAccessToken()
		require.Equal(t, "", token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("should work, nil token", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Host = "test.host"
		args.TokenExpiryInSeconds = 120
		expectedNonce := uint64(100)
		expectedHash := "hash"
		expectedSignature := "signature"
		expectedAddr := "addr"
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
		authClient, _ := NewNativeAuthClient(args)
		authClient.token = ""
		authClient.tokenExpire = time.Time{}
		encodedHost := base64.StdEncoding.EncodeToString([]byte(args.Host))
		encodedExtraInfo := base64.StdEncoding.EncodeToString([]byte("null"))
		internalToken := fmt.Sprintf("%s.%s.%d.%s", encodedHost, expectedHash, args.TokenExpiryInSeconds, encodedExtraInfo)
		encodedInternalToken := base64.StdEncoding.EncodeToString([]byte(internalToken))
		encodedAddress := base64.StdEncoding.EncodeToString([]byte(expectedAddr))
		encodedSignature := base64.StdEncoding.EncodeToString([]byte(expectedSignature))
		token, err := authClient.GetAccessToken()
		require.Nil(t, err)
		require.Equal(t, fmt.Sprintf("%s.%s.%s", encodedAddress, encodedInternalToken, encodedSignature), token)
	})
	t.Run("should work, token expired should generate new one", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthClient()
		args.Host = "test.host"
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

		authClient, _ := NewNativeAuthClient(args)

		currentTime := time.Now()
		authClient.getTimeHandler = func() time.Time {
			return currentTime.Add(time.Second * 1120)
		}
		storedToken := "token"
		authClient.token = storedToken
		authClient.tokenExpire = currentTime.Add(time.Second * 1000)
		token, err := authClient.GetAccessToken()
		require.NotEqual(t, storedToken, token)
		require.Nil(t, err)
	})
}

func createMockArgsNativeAuthClient() ArgsNativeAuthClient {
	return ArgsNativeAuthClient{
		Signer:                 &testsCommon.SignerStub{},
		ExtraInfo:              nil,
		Proxy:                  &testsCommon.ProxyStub{},
		CryptoComponentsHolder: &testsCommon.CryptoComponentsHolderStub{},
		TokenExpiryInSeconds:   0,
		Host:                   "",
	}
}
