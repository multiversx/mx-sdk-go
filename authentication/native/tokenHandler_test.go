package native

import (
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
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
		require.Nil(t, token)
		require.Equal(t, expectedErr, err)
	})
	t.Run("parseIntHandler errors should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			return []byte("host.blockHash.ttl.extraInfo"), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.NotNil(t, err)
	})
	t.Run("parseIntHandler errors should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			return []byte("blockHash.10.extraInfo"), nil
		}
		handler.hexDecodeHandler = func(s string) ([]byte, error) {
			return []byte(s), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.NotNil(t, token)
		require.Nil(t, err)
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
}

func TestNativeserver_decodeHandler(t *testing.T) {
	t.Parallel()

	handler := NewAuthTokenHandler()
	t.Run("len % 4 == 0 should do nothing", func(t *testing.T) {
		bytes, err := handler.decodeHandler("MTIz")
		require.Nil(t, err)
		require.Equal(t, "123", string(bytes))
	})
	t.Run("len % 4 == 1 should error", func(t *testing.T) {
		bytes, err := handler.decodeHandler("e")
		require.Nil(t, bytes)
		require.Equal(t, "illegal base64 data at input byte 1", err.Error())
	})
	t.Run("len % 2 == 0 should work", func(t *testing.T) {
		bytes, err := handler.decodeHandler("MQ")
		require.Nil(t, err)
		require.Equal(t, "1", string(bytes))
	})
	t.Run("len % 2 == 0 should work", func(t *testing.T) {
		bytes, err := handler.decodeHandler("MTI")
		require.Nil(t, err)
		require.Equal(t, "12", string(bytes))
	})
}
