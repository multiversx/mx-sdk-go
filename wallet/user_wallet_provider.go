package wallet

import (
	"errors"
	"fmt"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519/singlesig"
	"github.com/tyler-smith/go-bip39"
)

type userWalletProvider struct {
	suite  crypto.Suite
	keyGen crypto.KeyGenerator
	signer singlesig.Ed25519Signer
}

func NewUserWalletProvider() UserProvider {
	suite := ed25519.NewEd25519()
	return &userWalletProvider{
		suite,
		signing.NewKeyGenerator(suite),
		singlesig.Ed25519Signer{},
	}
}

func (u *userWalletProvider) GenerateKeyPair() (secretKey crypto.PrivateKey, publicKey crypto.PublicKey) {
	return u.keyGen.GeneratePair()
}

func (u *userWalletProvider) Sign(data []byte, secretKey crypto.PrivateKey) ([]byte, error) {
	sign, err := u.signer.Sign(secretKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}
	return sign, nil
}

func (u *userWalletProvider) Verify(data []byte, signature []byte, publicKey crypto.PublicKey) (bool, error) {
	err := u.signer.Verify(publicKey, data, signature)
	if err != nil {
		return false, fmt.Errorf("failed to verify signature: %v", err)
	}
	return true, nil
}

func (u *userWalletProvider) CreateSecretKeyFromBytes(data []byte) (crypto.PrivateKey, error) {
	pk, err := u.keyGen.PrivateKeyFromByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret key: %v", err)
	}
	return pk, nil
}

func (u *userWalletProvider) CreatePublicKeyFromBytes(data []byte) (crypto.PublicKey, error) {
	pk, err := u.keyGen.PublicKeyFromByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %v", err)
	}
	return pk, nil
}

func (u *userWalletProvider) ComputePublicKeyFromSecretKey(secretKey crypto.PrivateKey) (crypto.PublicKey, error) {
	public := secretKey.GeneratePublic()
	if public == nil {
		return nil, errors.New("failed to compute public key")
	}
	return public, nil
}

func (u *userWalletProvider) GenerateMnemonic() *Mnemonic {
	entropy, _ := bip39.NewEntropy(256)
	m, _ := bip39.NewMnemonic(entropy)
	return &Mnemonic{Text: m}
}
