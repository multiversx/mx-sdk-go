package cryptoProvider

import (
	"encoding/hex"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/testscommon/cryptoMocks"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestNewCryptoComponentsHolder(t *testing.T) {
	t.Parallel()

	t.Run("invalid privateKey bytes", func(t *testing.T) {
		t.Parallel()

		keyGen := &cryptoMocks.KeyGenStub{
			PrivateKeyFromByteArrayStub: func(b []byte) (crypto.PrivateKey, error) {
				return nil, expectedError
			},
		}
		holder, err := NewCryptoComponentsHolder(keyGen, []byte(""))
		require.Nil(t, holder)
		require.Equal(t, expectedError, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		privateKey := &testsCommon.PrivateKeyStub{
			GeneratePublicCalled: func() crypto.PublicKey {
				return &testsCommon.PublicKeyStub{}
			},
		}
		keyGen := &cryptoMocks.KeyGenStub{
			PrivateKeyFromByteArrayStub: func(b []byte) (crypto.PrivateKey, error) {
				return privateKey, nil
			},
		}
		holder, err := NewCryptoComponentsHolder(keyGen, []byte(""))
		require.False(t, check.IfNil(holder))
		require.Nil(t, err)
		_ = holder.GetPublicKey()
		_ = holder.GetPrivateKey()
		_ = holder.GetBech32()
		_ = holder.GetAddressHandler()
	})
	t.Run("should work with real components", func(t *testing.T) {
		t.Parallel()

		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		holder, err := NewCryptoComponentsHolder(keyGen, sk)
		require.False(t, check.IfNil(holder))
		require.Nil(t, err)
		publicKey := holder.GetPublicKey()
		privateKey := holder.GetPrivateKey()
		require.False(t, check.IfNil(publicKey))
		require.False(t, check.IfNil(privateKey))

		bech32Address := holder.GetBech32()
		addressHandler := holder.GetAddressHandler()
		require.Equal(t, addressHandler.AddressAsBech32String(), bech32Address)
		require.Equal(t, "erd1j84k44nsqsme8r6e5aawutx0z2cd6cyx3wprkzdh73x2cf0kqvksa3snnq", bech32Address)
	})
}
