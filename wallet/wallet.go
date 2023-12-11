package wallet

import (
	crypto "github.com/multiversx/mx-chain-crypto-go"
)

type Provider interface {
	GenerateKeyPair() (secretKey crypto.PrivateKey, publicKey crypto.PublicKey)
	Sign(data []byte, secretKey crypto.PrivateKey) ([]byte, error)
	Verify(data []byte, signature []byte, publicKey crypto.PublicKey) (bool, error)
	CreateSecretKeyFromBytes(data []byte) (crypto.PrivateKey, error)
	CreatePublicKeyFromBytes(data []byte) (crypto.PublicKey, error)
	ComputePublicKeyFromSecretKey(secretKey crypto.PrivateKey) (crypto.PublicKey, error)
}

type UserProvider interface {
	Provider
	GenerateMnemonic() *Mnemonic
}

type ValidatorProvider interface {
	Provider
}
