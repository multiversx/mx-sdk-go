package wallet

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	userWallet := NewUserWalletProvider()
	secretKey, publicKey := userWallet.GenerateKeyPair()
	message := "someRandomMessage"
	signature, err := userWallet.Sign([]byte(message), secretKey)
	require.NoError(t, err)
	result, err := userWallet.Verify([]byte(message), signature, publicKey)
	require.NoError(t, err)
	require.Equal(t, result, true)

	validatorWallet := NewValidatorWalletProvider()
	secretKey, publicKey = validatorWallet.GenerateKeyPair()
	signature, err = validatorWallet.Sign([]byte(message), secretKey)
	require.NoError(t, err)
	result, err = validatorWallet.Verify([]byte(message), signature, publicKey)
	require.NoError(t, err)
	require.Equal(t, result, true)
}

func TestComputePublicKeyFromSecretKey(t *testing.T) {
	userWallet := NewUserWalletProvider()
	secretKey, publicKey := userWallet.GenerateKeyPair()
	key, err := userWallet.ComputePublicKeyFromSecretKey(secretKey)
	require.NoError(t, err)
	require.Equal(t, publicKey, key)

	validatorWallet := NewValidatorWalletProvider()
	secretKey, publicKey = validatorWallet.GenerateKeyPair()
	key, err = validatorWallet.ComputePublicKeyFromSecretKey(secretKey)
	require.NoError(t, err)
	require.Equal(t, publicKey, key)
}
