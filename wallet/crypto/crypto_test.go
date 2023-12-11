package crypto

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPasswordBasedEncyptDecrypt(t *testing.T) {
	message := "Thisisateststring"
	password := "randompassword"

	bytes := []byte(message)

	p := NewPasswordBaseEncryptorDecryptor()

	encrypt, err := p.Encrypt(bytes, password)
	require.NoError(t, err, "failed to encrypt message")

	decryptedMessage, err := p.Decrypt(encrypt, password)
	require.NoError(t, err, "failed to decrypt message")

	require.Equal(t, string(decryptedMessage), message)
}
