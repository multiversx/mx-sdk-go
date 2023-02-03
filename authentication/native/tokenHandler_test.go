package native

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNativeserver_Decode(t *testing.T) {
	t.Parallel()

	t.Run("decodeHandler errors for address should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			return make([]byte, 0), expectedErr
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("decodeHandler errors for body should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			if s == "address" {
				return nil, nil
			}
			return make([]byte, 0), expectedErr
		}
		token, err := handler.Decode("address.body.signature")
		require.True(t, check.IfNil(token))
		require.Equal(t, expectedErr, err)
	})
	t.Run("hexDecodeHandler errors should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			return []byte("host.blockHash.ttl.extraInfo"), nil
		}
		handler.hexDecodeHandler = func(s string) ([]byte, error) {
			return nil, expectedErr
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("decodeHandler errors for host should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			if s == "body" {
				return []byte("host.blockHash.ttl.extraInfo"), nil
			}
			if s == "host" {
				return make([]byte, 0), expectedErr
			}
			return nil, nil
		}
		handler.hexDecodeHandler = func(s string) ([]byte, error) {
			return []byte(s), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("ParseInt errors for ttl should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			if s == "body" {
				return []byte("host.blockHash.ttl.extraInfo"), nil
			}
			return nil, nil
		}
		handler.hexDecodeHandler = func(s string) ([]byte, error) {
			return []byte(s), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		assert.Equal(t, "strconv.ParseInt: parsing \"ttl\": invalid syntax", err.Error())
	})
	t.Run("decodeHandler errors for extraInfo should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			if s == "body" {
				return []byte("host.blockHash.1234.extraInfo"), nil
			}
			if s == "extraInfo" {
				return make([]byte, 0), expectedErr
			}
			return nil, nil
		}
		handler.hexDecodeHandler = func(s string) ([]byte, error) {
			return []byte(s), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.Equal(t, expectedErr, err)
	})
}

func TestNativeserver_Encode(t *testing.T) {
	t.Parallel()

	t.Run("nil signature should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		token, err := handler.Encode(&AuthToken{
			signature: nil,
		})
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrNilSignature, err)
	})
	t.Run("encodeHandler return empty string for address should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.encodeHandler = func(src []byte) string {
			return ""
		}
		token, err := handler.Encode(&AuthToken{
			signature: []byte("signature"),
		})
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrNilAddress, err)
	})
	t.Run("encodeHandler return empty string for body should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.encodeHandler = func(src []byte) string {
			if string(src) == "address" {
				return "addr"
			}
			return ""
		}
		token, err := handler.Encode(&AuthToken{
			signature: []byte("signature"),
			address:   []byte("address"),
		})
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrNilBody, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.encodeHandler = func(src []byte) string {
			return "a"
		}
		token, err := handler.Encode(&AuthToken{
			signature: []byte("signature"),
		})
		require.Equal(t, "a.a.7369676e6174757265", token)
		require.Nil(t, err)
	})
	t.Run("should work with real components", func(t *testing.T) {
		t.Parallel()

		authToken := &AuthToken{
			ttl:       110,
			address:   []byte("address"),
			host:      []byte("host"),
			extraInfo: []byte("extra info"),
			signature: []byte("sig"),
			blockHash: "block hash",
		}
		assert.False(t, check.IfNil(authToken))
		handler := NewAuthTokenHandler()
		token, err := handler.Encode(authToken)
		assert.Nil(t, err)

		recoveredToken, err := handler.Decode(token)
		assert.Nil(t, err)

		assert.Equal(t, authToken, recoveredToken)
	})
}
