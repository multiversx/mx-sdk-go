package native

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	genesisMock "github.com/multiversx/mx-chain-go/genesis/mock"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/authentication/native/mock"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedErr = errors.New("expected error")

func TestNativeserver_NewNativeAuthServer(t *testing.T) {
	t.Parallel()

	t.Run("nil blockhash handler should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.BlockhashHandler = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilBlockhashHandler, err)
	})
	t.Run("nil signer should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.Signer = nil
		server, err := NewNativeAuthServer(args)
		require.Nil(t, server)
		require.Equal(t, authentication.ErrNilSigner, err)
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

	tokenTtl := int64(20)
	blockTimestamp := int64(10)
	t.Run("blockhash handler returns error should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("token expired should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				return block, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp+tokenTtl+1, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, authentication.ErrTokenExpired, err)
	})
	t.Run("pubKeyConverter errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				return block, nil
			},
		}
		args.PubKeyConverter = &genesisMock.PubkeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				return nil, expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(blockTimestamp, 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("keyGenerator errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				return block, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, expectedErr
			},
		}
		args.PubKeyConverter = &genesisMock.PubkeyConverterStub{
			DecodeCalled: func(humanReadable string) ([]byte, error) {
				return nil, nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(blockTimestamp), 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("verification errors should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			GetUnsignedTokenCalled: func(authToken authentication.AuthToken) []byte {
				return []byte("token")
			},
		}
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				return block, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			VerifyMessageCalled: func(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
				return expectedErr
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(blockTimestamp), 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		require.Equal(t, expectedErr, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsNativeAuthServer()
		args.TokenHandler = &mock.AuthTokenHandlerStub{
			GetUnsignedTokenCalled: func(authToken authentication.AuthToken) []byte {
				return []byte("token")
			},
		}
		args.BlockhashHandler = &testsCommon.BlockhashHandlerStub{
			GetBlockByHashCalled: func(ctx context.Context, hash string) (*data.Block, error) {
				block := &data.Block{
					Timestamp: int(blockTimestamp),
				}
				return block, nil
			},
		}
		args.KeyGenerator = &genesisMock.KeyGeneratorStub{
			PublicKeyFromByteArrayCalled: func(b []byte) (crypto.PublicKey, error) {
				return nil, nil
			},
		}
		args.Signer = &testsCommon.SignerStub{
			VerifyMessageCalled: func(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
				return nil
			},
		}
		server, _ := NewNativeAuthServer(args)
		server.getTimeHandler = func() time.Time {
			return time.Unix(int64(blockTimestamp), 1)
		}

		err := server.Validate(&AuthToken{
			ttl: tokenTtl,
		})
		assert.Nil(t, err)
	})
}

func createMockArgsNativeAuthServer() ArgsNativeAuthServer {
	return ArgsNativeAuthServer{
		BlockhashHandler: &testsCommon.BlockhashHandlerStub{},
		TokenHandler:     &mock.AuthTokenHandlerStub{},
		Signer:           &testsCommon.SignerStub{},
		PubKeyConverter:  &genesisMock.PubkeyConverterStub{},
		KeyGenerator:     &genesisMock.KeyGeneratorStub{},
	}
}
