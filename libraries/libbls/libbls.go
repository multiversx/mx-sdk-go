package main

import (
	"encoding/hex"
	"log"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl/singlesig"
)

var (
	keyGenerator = signing.NewKeyGenerator(mcl.NewSuiteBLS12())
	blsSigner    = singlesig.BlsSingleSigner{}
)

func doGeneratePrivateKeyAsHex() string {
	privateKey, _ := keyGenerator.GeneratePair()
	privateKeyBytes, err := privateKey.ToByteArray()
	if err != nil {
		log.Println("doGeneratePrivateKey(): error when decoding the private key", err)
		return ""
	}

	return hex.EncodeToString(privateKeyBytes)
}

func doGeneratePublicKeyAsHex(privateKeyHex string) string {
	privateKeyBytes, ok := decodeInputParameter("private key", privateKeyHex)
	if !ok {
		return ""
	}

	privateKey, err := keyGenerator.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		log.Println("doGeneratePublicKeyAsHex(): error when creating the private key", err)
		return ""
	}

	publicKey := privateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		log.Println("doGeneratePublicKeyAsHex(): error when decoding the public key", err)
		return ""
	}

	return hex.EncodeToString(publicKeyBytes)
}

func doComputeMessageSignatureAsHex(messageHex string, privateKeyHex string) string {
	message, ok := decodeInputParameter("message", messageHex)
	if !ok {
		return ""
	}

	privateKeyBytes, ok := decodeInputParameter("private key", privateKeyHex)
	if !ok {
		return ""
	}

	privateKey, err := keyGenerator.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		log.Println("doComputeMessageSignatureAsHex(): error when creating the private key", err)
		return ""
	}

	signature, err := blsSigner.Sign(privateKey, message)
	if err != nil {
		log.Println("doComputeMessageSignatureAsHex(): error when signing the message", err)
		return ""
	}

	return hex.EncodeToString(signature)
}

func doVerifyMessageSignature(publicKeyHex string, messageHex string, signatureHex string) bool {
	publicKeyBytes, ok := decodeInputParameter("public key", publicKeyHex)
	if !ok {
		return false
	}

	message, ok := decodeInputParameter("message", messageHex)
	if !ok {
		return false
	}

	signature, ok := decodeInputParameter("signature", signatureHex)
	if !ok {
		return false
	}

	publicKey, err := keyGenerator.PublicKeyFromByteArray(publicKeyBytes)
	if err != nil {
		log.Println("doVerifyMessageSignature(): error when creating the public key", err)
		return false
	}

	err = blsSigner.Verify(publicKey, message, signature)
	return err == nil
}

func decodeInputParameter(parameterName string, parameterValueHex string) ([]byte, bool) {
	data, err := hex.DecodeString(parameterValueHex)
	if err != nil {
		log.Println("cannot decode input parameter", parameterName, err)
		return nil, false
	}

	return data, true
}
