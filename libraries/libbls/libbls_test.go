package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoGeneratePrivateKey(t *testing.T) {
	privateKeyHex := doGeneratePrivateKey()
	privateKey, err := hex.DecodeString(privateKeyHex)

	require.Nil(t, err)
	require.Len(t, privateKey, 32)
}

func TestDoGeneratePublicKey(t *testing.T) {
	t.Run("with good input", func(t *testing.T) {
		publicKeyHex := doGeneratePublicKey("7cff99bd671502db7d15bc8abc0c9a804fb925406fbdd50f1e4c17a4cd774247")
		require.Equal(t, "e7beaa95b3877f47348df4dd1cb578a4f7cabf7a20bfeefe5cdd263878ff132b765e04fef6f40c93512b666c47ed7719b8902f6c922c04247989b7137e837cc81a62e54712471c97a2ddab75aa9c2f58f813ed4c0fa722bde0ab718bff382208", publicKeyHex)
	})

	t.Run("with bad input", func(t *testing.T) {
		publicKeyHex := doGeneratePublicKey("7cff99bd671502db7d15bc8abc0c9a804fb925406fbdd50f1e4c17a4cd7742")
		require.Equal(t, "", publicKeyHex)
	})
}

func TestDoSignMessage(t *testing.T) {
	t.Run("with good input", func(t *testing.T) {
		messageHex := hex.EncodeToString([]byte("hello"))
		privateKeyHex := "7cff99bd671502db7d15bc8abc0c9a804fb925406fbdd50f1e4c17a4cd774247"
		signatureHex := doSignMessage(messageHex, privateKeyHex)
		require.Equal(t, "84fd0a3a9d4f1ea2d4b40c6da67f9b786284a1c3895b7253fec7311597cda3f757862bb0690a92a13ce612c33889fd86", signatureHex)
	})

	t.Run("with bad input (not hex)", func(t *testing.T) {
		signatureHex := doSignMessage("not hex", "7cff99bd671502db7d15bc8abc0c9a804fb925406fbdd50f1e4c17a4cd774247")
		require.Equal(t, "", signatureHex)
	})

	t.Run("with bad input (bad key)", func(t *testing.T) {
		messageHex := hex.EncodeToString([]byte("hello"))
		signatureHex := doSignMessage(messageHex, "7cff99bd671502db7d15bc8abc0c9a804fb925406fbdd50f1e4c17a4cd7742")
		require.Equal(t, "", signatureHex)
	})
}

func TestDoVerifyMessage(t *testing.T) {
	t.Run("with good input", func(t *testing.T) {
		publicKeyHex := "e7beaa95b3877f47348df4dd1cb578a4f7cabf7a20bfeefe5cdd263878ff132b765e04fef6f40c93512b666c47ed7719b8902f6c922c04247989b7137e837cc81a62e54712471c97a2ddab75aa9c2f58f813ed4c0fa722bde0ab718bff382208"
		messageHex := hex.EncodeToString([]byte("hello"))
		signatureHex := "84fd0a3a9d4f1ea2d4b40c6da67f9b786284a1c3895b7253fec7311597cda3f757862bb0690a92a13ce612c33889fd86"
		require.True(t, doVerifyMessage(publicKeyHex, messageHex, signatureHex))
	})

	t.Run("with altered signature", func(t *testing.T) {
		publicKeyHex := "e7beaa95b3877f47348df4dd1cb578a4f7cabf7a20bfeefe5cdd263878ff132b765e04fef6f40c93512b666c47ed7719b8902f6c922c04247989b7137e837cc81a62e54712471c97a2ddab75aa9c2f58f813ed4c0fa722bde0ab718bff382208"
		messageHex := hex.EncodeToString([]byte("hello"))
		signatureHex := "94fd0a3a9d4f1ea2d4b40c6da67f9b786284a1c3895b7253fec7311597cda3f757862bb0690a92a13ce612c33889fd86"
		require.False(t, doVerifyMessage(publicKeyHex, messageHex, signatureHex))
	})

	t.Run("with altered message", func(t *testing.T) {
		publicKeyHex := "e7beaa95b3877f47348df4dd1cb578a4f7cabf7a20bfeefe5cdd263878ff132b765e04fef6f40c93512b666c47ed7719b8902f6c922c04247989b7137e837cc81a62e54712471c97a2ddab75aa9c2f58f813ed4c0fa722bde0ab718bff382208"
		messageHex := hex.EncodeToString([]byte("helloWorld"))
		signatureHex := "84fd0a3a9d4f1ea2d4b40c6da67f9b786284a1c3895b7253fec7311597cda3f757862bb0690a92a13ce612c33889fd86"
		require.False(t, doVerifyMessage(publicKeyHex, messageHex, signatureHex))
	})

	t.Run("with bad public key", func(t *testing.T) {
		publicKeyHex := "badbad95b3877f47348df4dd1cb578a4f7cabf7a20bfeefe5cdd263878ff132b765e04fef6f40c93512b666c47ed7719b8902f6c922c04247989b7137e837cc81a62e54712471c97a2ddab75aa9c2f58f813ed4c0fa722bde0ab718bff382208"
		messageHex := hex.EncodeToString([]byte("hello"))
		signatureHex := "84fd0a3a9d4f1ea2d4b40c6da67f9b786284a1c3895b7253fec7311597cda3f757862bb0690a92a13ce612c33889fd86"
		require.False(t, doVerifyMessage(publicKeyHex, messageHex, signatureHex))
	})
}

func TestGenerateSignAndVerify(t *testing.T) {
	messageHex := hex.EncodeToString([]byte("hello"))

	privateKeyHex := doGeneratePrivateKey()
	publicKeyHex := doGeneratePublicKey(privateKeyHex)
	signatureHex := doSignMessage(messageHex, privateKeyHex)
	isOk := doVerifyMessage(publicKeyHex, messageHex, signatureHex)

	require.True(t, isOk)
}
