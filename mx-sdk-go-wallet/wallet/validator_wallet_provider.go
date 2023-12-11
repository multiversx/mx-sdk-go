package wallet

import (
	"errors"
	"fmt"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	mcl "github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	signer "github.com/multiversx/mx-chain-crypto-go/signing/mcl/singlesig"
)

type validatorWalletProvider struct {
	suite  crypto.Suite
	keyGen crypto.KeyGenerator
	signer signer.BlsSingleSigner
}

func NewValidatorWalletProvider() ValidatorProvider {
	suite := mcl.NewSuiteBLS12()
	return &validatorWalletProvider{
		suite,
		signing.NewKeyGenerator(suite),
		signer.BlsSingleSigner{},
	}
}

func (v *validatorWalletProvider) GenerateKeyPair() (secretKey crypto.PrivateKey, publicKey crypto.PublicKey) {
	return v.keyGen.GeneratePair()
}

func (v *validatorWalletProvider) Sign(data []byte, secretKey crypto.PrivateKey) ([]byte, error) {
	sign, err := v.signer.Sign(secretKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}
	return sign, nil
}

func (v *validatorWalletProvider) Verify(data []byte, signature []byte, publicKey crypto.PublicKey) (bool, error) {
	err := v.signer.Verify(publicKey, data, signature)
	if err != nil {
		return false, fmt.Errorf("failed to verify signature: %v", err)
	}
	return true, nil
}

func (v *validatorWalletProvider) CreateSecretKeyFromBytes(data []byte) (crypto.PrivateKey, error) {
	pk, err := v.keyGen.PrivateKeyFromByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret key: %v", err)
	}
	return pk, nil
}

func (v *validatorWalletProvider) CreatePublicKeyFromBytes(data []byte) (crypto.PublicKey, error) {
	pk, err := v.keyGen.PublicKeyFromByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %v", err)
	}
	return pk, nil
}

func (v *validatorWalletProvider) ComputePublicKeyFromSecretKey(secretKey crypto.PrivateKey) (crypto.PublicKey, error) {
	public := secretKey.GeneratePublic()
	if public == nil {
		return nil, errors.New("failed to compute public key")
	}
	return public, nil
}
