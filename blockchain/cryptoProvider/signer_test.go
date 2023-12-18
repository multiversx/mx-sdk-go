package cryptoProvider

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
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

		signerInstance := NewSigner()
		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		holder, err := NewCryptoComponentsHolder(keyGen, sk)
		require.Nil(t, err)
		sig, err := signerInstance.SignMessage([]byte("msg"), holder.GetPrivateKey())
		require.NotNil(t, sig)
		require.Nil(t, err)
		err = signerInstance.VerifyMessage([]byte("msg"), holder.GetPublicKey(), sig)
		require.Nil(t, err)
	})
}

func TestSigner_SignByteSlice_VerifyByteSlice(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		msg := []byte("msg")
		signerInstance := NewSigner()
		sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
		privateKey, err := keyGen.PrivateKeyFromByteArray(sk)
		require.Nil(t, err)
		sig, err := signerInstance.SignByteSlice(msg, privateKey)
		require.Nil(t, err)

		publicKey := privateKey.GeneratePublic()
		err = signerInstance.VerifyByteSlice(msg, publicKey, sig)
		require.Nil(t, err)
	})
}

func TestSigner_SignTransaction(t *testing.T) {
	t.Parallel()

	sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
	privateKey, _ := keyGen.PrivateKeyFromByteArray(sk)

	t.Run("tx already signed", func(t *testing.T) {
		t.Parallel()

		signerInstance := NewSigner()
		require.False(t, check.IfNil(signerInstance))
		tx := &transaction.FrontendTransaction{Signature: "sig"}
		sig, err := signerInstance.SignTransaction(tx, nil)
		require.Nil(t, sig)

		require.Equal(t, ErrTxAlreadySigned, err)

	})

	t.Run("should work if all the tx is signed", func(t *testing.T) {
		t.Parallel()

		signerInstance := NewSigner()
		require.False(t, check.IfNil(signerInstance))
		tx := &transaction.FrontendTransaction{Version: 1}
		sig, err := signerInstance.SignTransaction(tx, privateKey)
		require.NotNil(t, sig)
		require.Nil(t, err)
	})
	t.Run("should work if only txHash is signed", func(t *testing.T) {
		t.Parallel()

		signerInstance := NewSigner()
		require.False(t, check.IfNil(signerInstance))
		tx := &transaction.FrontendTransaction{Version: 2, Options: 1}
		sig, err := signerInstance.SignTransaction(tx, privateKey)
		require.NotNil(t, sig)
		require.Nil(t, err)
	})
	// {
	// 	 "address":"erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th", //alice
	//	 "message":"0x6d657373616765",
	//	 "signature":"0x546c6b6d6487852f54571ab2da81b48ff8f09bef71ba07b116fcf7203538cd64ea5f9bffcc13a0279a75ca3b1b0a1e478d23e1771d381011f8135e4372a9dd00",
	//	 "version":1,
	//	 "signer":"sdk-js"
	// }
	t.Run("should work with signature generated using sdk-js", func(t *testing.T) {
		t.Parallel()

		signerInstance := NewSigner()

		w := interactors.NewWallet()
		privateKeyInstance, err := w.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
		require.Nil(t, err)
		holder, err := NewCryptoComponentsHolder(keyGen, privateKeyInstance)
		require.Nil(t, err)
		msg := []byte("message")
		sig, _ := hex.DecodeString("546c6b6d6487852f54571ab2da81b48ff8f09bef71ba07b116fcf7203538cd64ea5f9bffcc13a0279a75ca3b1b0a1e478d23e1771d381011f8135e4372a9dd00")
		err = signerInstance.VerifyMessage(msg, holder.GetPublicKey(), sig)
		require.Nil(t, err)
	})
}
