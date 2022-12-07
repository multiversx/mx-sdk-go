package native

import (
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/stretchr/testify/require"
)

func TestNativeserver_Decode(t *testing.T) {
	t.Parallel()

	//
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
			return []byte("host.blockHash.Ttl.ExtraInfo"), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.Nil(t, token)
		require.NotNil(t, err)
	})
	t.Run("parseIntHandler errors should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.decodeHandler = func(s string) ([]byte, error) {
			return []byte("host.blockHash.10.ExtraInfo"), nil
		}
		token, err := handler.Decode("address.body.signature")
		require.NotNil(t, token)
		require.Nil(t, err)
	})
}

func TestNativeserver_Encode(t *testing.T) {
	t.Parallel()

	//
	t.Run("decodeHandler errors for address should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		token, err := handler.Encode(nil)
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrCannotConvertToken, err)
	})
	t.Run("nil signature should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		token, err := handler.Encode(&NativeAuthToken{
			Signature: nil,
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
		token, err := handler.Encode(&NativeAuthToken{
			Signature: []byte("signature"),
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
		token, err := handler.Encode(&NativeAuthToken{
			Signature: []byte("signature"),
			Address:   []byte("address"),
		})
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrNilBody, err)
	})
	t.Run("encodeHandler return empty string for signature should error", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.encodeHandler = func(src []byte) string {
			if string(src) == "address" ||
				string(src) == "host.blockHash.10.extraInfo" {
				return "addr"
			}
			return ""
		}
		token, err := handler.Encode(&NativeAuthToken{
			Signature: []byte("signature"),
			Address:   []byte("address"),
			Host:      "host",
			BlockHash: "blockHash",
			Ttl:       10,
			ExtraInfo: "extraInfo",
		})
		require.Equal(t, "", token)
		require.Equal(t, authentication.ErrNilSignature, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		handler := NewAuthTokenHandler()
		handler.encodeHandler = func(src []byte) string {
			return "a"
		}
		token, err := handler.Encode(&NativeAuthToken{
			Signature: []byte("signature"),
		})
		require.Equal(t, "a.a.a", token)
		require.Nil(t, err)
	})
}
