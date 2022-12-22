package cryptoProvider

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/stretchr/testify/require"
)

var (
	suite         = ed25519.NewEd25519()
	keyGen        = signing.NewKeyGenerator(suite)
	expectedError = errors.New("expected error")
)

func TestSigner_SignMessage_VerifyMessage(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		signer := NewSigner()
		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		holder, err := NewCryptoComponentsHolder(keyGen, sk)
		require.Nil(t, err)
		sig, err := signer.SignMessage([]byte("msg"), holder.GetPrivateKey())
		require.NotNil(t, sig)
		require.Nil(t, err)
		err = signer.VerifyMessage([]byte("msg"), holder.GetPublicKey(), sig)
		require.Nil(t, err)
	})
}

func TestSigner_SignByteSlice_VerifyByteSlice(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		msg := []byte("msg")
		signer := NewSigner()
		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		privateKey, err := keyGen.PrivateKeyFromByteArray(sk)
		require.Nil(t, err)
		sig, err := signer.SignByteSlice(msg, privateKey)
		require.Nil(t, err)

		publicKey := privateKey.GeneratePublic()
		err = signer.VerifyByteSlice(msg, publicKey, sig)
		require.Nil(t, err)
	})
}

func TestSigner_SignTransaction(t *testing.T) {
	t.Parallel()

	sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
	privateKey, _ := keyGen.PrivateKeyFromByteArray(sk)

	t.Run("tx already signed", func(t *testing.T) {
		t.Parallel()

		signer := NewSigner()
		require.False(t, check.IfNil(signer))
		tx := &data.Transaction{Signature: "sig"}
		sig, err := signer.SignTransaction(tx, nil)
		require.Nil(t, sig)

		require.Equal(t, ErrTxAlreadySigned, err)

	})

	t.Run("should work if all the tx is signed", func(t *testing.T) {
		t.Parallel()

		signer := NewSigner()
		require.False(t, check.IfNil(signer))
		tx := &data.Transaction{Version: 1}
		sig, err := signer.SignTransaction(tx, privateKey)
		require.NotNil(t, sig)
		require.Nil(t, err)
	})
	t.Run("should work if only txHash is signed", func(t *testing.T) {
		t.Parallel()

		signer := NewSigner()
		require.False(t, check.IfNil(signer))
		tx := &data.Transaction{Version: 2, Options: 1}
		sig, err := signer.SignTransaction(tx, privateKey)
		require.NotNil(t, sig)
		require.Nil(t, err)
	})
}
