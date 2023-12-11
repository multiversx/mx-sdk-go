package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"io"
)

type passwordBasedEncryptorDecryptor struct{}

func NewPasswordBaseEncryptorDecryptor() PasswordBasedEncryptorDecryptor {
	return &passwordBasedEncryptorDecryptor{}
}

func (p *passwordBasedEncryptorDecryptor) Decrypt(encryptedData *PasswordEncryptedData, password string) ([]byte, error) {
	if encryptedData.KDF != KeyDerivationFunctionScrypt {
		return nil, errors.New(fmt.Sprintf("Unknown derivation function: %v", encryptedData.KDF))
	}

	mac, err := hex.DecodeString(encryptedData.MAC)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MAC: %v", err)
	}
	salt, err := hex.DecodeString(encryptedData.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}
	cipherText, err := hex.DecodeString(encryptedData.CipherText)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cipherText: %v", err)
	}
	iv, err := hex.DecodeString(encryptedData.Iv)
	if err != nil {
		return nil, fmt.Errorf("failed to decode iv: %v", err)
	}

	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		encryptedData.KDFParams.N,
		encryptedData.KDFParams.R,
		encryptedData.KDFParams.P,
		encryptedData.KDFParams.DKLen,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}

	hash := hmac.New(sha256.New, derivedKey[16:32])
	_, err = hash.Write(cipherText)
	if err != nil {
		return nil, err
	}

	sha := hash.Sum(nil)
	if !bytes.Equal(sha, mac) {
		return nil, errors.New("wrong password")
	}

	aesBlock, err := aes.NewCipher(derivedKey[:16])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(aesBlock, iv)
	decryptedData := make([]byte, len(cipherText))
	stream.XORKeyStream(decryptedData, cipherText)

	return decryptedData, nil
}

func (p *passwordBasedEncryptorDecryptor) Encrypt(data []byte, password string) (*PasswordEncryptedData, error) {
	kdfParams := KeyDerivationParams{N: 4096, R: 8, P: 1, DKLen: 32}

	salt := make([]byte, RandomSaltLength)
	iv := make([]byte, RandomIvLength)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %v", err)
	}
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random iv: %v", err)
	}

	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		kdfParams.N,
		kdfParams.R,
		kdfParams.P,
		kdfParams.DKLen,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %v", err)
	}

	aesBlock, err := aes.NewCipher(derivedKey[:16])
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	cipherText := make([]byte, len(data))
	stream.XORKeyStream(cipherText, data)

	hash := hmac.New(sha256.New, derivedKey[16:32])
	_, err = hash.Write(cipherText)
	if err != nil {
		return nil, err
	}

	sha := hash.Sum(nil)

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generare uuid: %v", err)
	}

	return &PasswordEncryptedData{
		ID:         newUUID.String(),
		Version:    EncryptorVersion,
		Cipher:     CipherAlgorithmAes128Ctr,
		CipherText: hex.EncodeToString(cipherText),
		Iv:         hex.EncodeToString(iv),
		KDF:        KeyDerivationFunctionScrypt,
		KDFParams:  &kdfParams,
		Salt:       hex.EncodeToString(salt),
		MAC:        hex.EncodeToString(sha),
	}, nil
}
