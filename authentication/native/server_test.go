package native

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	genesisMock "github.com/ElrondNetwork/elrond-go/genesis/mock"
	"github.com/ElrondNetwork/elrond-go/testscommon/cryptoMocks"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication/native/mock"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
	"github.com/stretchr/testify/require"
)

var expectedErr = errors.New("expected error")

func TestNativeserver_NewNativeAuthServer(t *testing.T) {
	t.Parallel()

	t.Run("nil proxy should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.Proxy = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, workflows.ErrNilProxy, err)
	})
	t.Run("nil AcceptedHosts should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.AcceptedHosts = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilAcceptedHosts, err)
	})
	t.Run("empty AcceptedHosts should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.AcceptedHosts = make(map[string]struct{})
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilAcceptedHosts, err)
	})
	t.Run("nil KeyGenerator should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.KeyGenerator = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, crypto.ErrNilKeyGenerator, err)
	})
	t.Run("nil pubKeyConverter should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.PubKeyConverter = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, core.ErrNilPubkeyConverter, err)
	})
	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.Signer = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilSigner, err)
	})
	t.Run("nil token handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilTokenHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		server, err := NewNativeAuthServer(args)
		require.NotNil(t, server)
		require.False(t, server.IsInterfaceNil())
		require.Nil(t, err)
	})
}
func TestNativeserver_Validate(t *testing.T) {
	t.Parallel()
	t.Run("decode errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)

		err := server.Validate("accessToken")

		require.Equal(t, expectedErr, err)
	})
	t.Run("host not accepted should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{host: "invalidHost"}, nil
			},
		}
		server, _ := NewNativeAuthServer(args)

		err := server.Validate("accessToken")

		require.Equal(t, authentication.ErrHostNotAccepted, err)
	})
	t.Run("proxy returns error should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{host: "test.host"}, nil
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)

		err := server.Validate("accessToken")

		require.Equal(t, expectedErr, err)
	})
	t.Run("token expired should error", func(t *testing.T) {
		t.Parallel()

		hyperblockTimestamp := 10
		tokenTtl := 20
		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{
					host: "test.host",
					ttl:  int64(tokenTtl),
				}, nil
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return &data.HyperBlock{Timestamp: uint64(hyperblockTimestamp)}, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(hyperblockTimestamp+tokenTtl+1), 1)
		}

		err := server.Validate("accessToken")

		require.Equal(t, authentication.ErrTokenExpired, err)
	})
	t.Run("keyGenerator errors should error", func(t *testing.T) {
		t.Parallel()

		hyperblockTimestamp := 10
		tokenTtl := 20
		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{
					host: "test.host",
					ttl:  int64(tokenTtl),
				}, nil
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return &data.HyperBlock{Timestamp: uint64(hyperblockTimestamp)}, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, expectedErr
			},
		}
		args.PubKeyConverter = &genesisMock.PubkeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(hyperblockTimestamp), 1)
		}

		err := server.Validate("accessToken")

		require.Equal(t, expectedErr, err)
	})
	t.Run("keyGenerator errors should error", func(t *testing.T) {
		t.Parallel()

		hyperblockTimestamp := 10
		tokenTtl := 20
		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{
					host: "test.host",
					ttl:  int64(tokenTtl),
				}, nil
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return &data.HyperBlock{Timestamp: uint64(hyperblockTimestamp)}, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(hyperblockTimestamp), 1)
		}

		err := server.Validate("accessToken")

		require.Equal(t, expectedErr, err)
	})
	t.Run("verification errors should error", func(t *testing.T) {
		t.Parallel()

		hyperblockTimestamp := 10
		tokenTtl := 20
		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			DecodeCalled: func(accessToken string) (authentication.AuthToken, error) {
				return AuthToken{
					host: "test.host",
					ttl:  int64(tokenTtl),
				}, nil
			},
		}
		args.Proxy = &testsCommon.ProxyStub{
			GetHyperBlockByHashCalled: func(ctx context.Context, hash string) (*data.HyperBlock, error) {
				return &data.HyperBlock{Timestamp: uint64(hyperblockTimestamp)}, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &cryptoMocks.SignerStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(hyperblockTimestamp), 1)
		}

		err := server.Validate("accessToken")

		require.Equal(t, expectedErr, err)
	})
}

func createMockArgsNativeAuthServer() ArgsNativeAuthServer {
	acceptedHosts := make(map[string]struct{})
	acceptedHosts["test.host"] = struct{}{}
	return ArgsNativeAuthServer{
		Proxy:           &testsCommon.ProxyStub{},
		TokenHandler:    &mock.AuthTokenHandlerStub{},
		Signer:          &cryptoMocks.SignerStub{},
		KeyGenerator:    &genesisMock.KeyGeneratorStub{},
		PubKeyConverter: &genesisMock.PubkeyConverterStub{},
		AcceptedHosts:   acceptedHosts,
	}
}
