package main

import (
	"encoding/hex"
	"log"

	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/mcl"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/mcl/singlesig"
)

var (
	keyGenerator = signing.NewKeyGenerator(mcl.NewSuiteBLS12())
	blsSigner    = singlesig.BlsSingleSigner{}
)

func doGeneratePrivateKey() string {
	privateKey, _ := keyGenerator.GeneratePair()
	privateKeyBytes, err := privateKey.ToByteArray()
	if err != nil {
		log.Println("doGeneratePrivateKey(): error when decoding the private key", err)
		return ""
	}

	return hex.EncodeToString(privateKeyBytes)
}

func doGeneratePublicKey(privateKeyHex string) string {
	privateKeyBytes, ok := decodeInputParameter("private key", privateKeyHex)
	if !ok {
		return ""
	}

	privateKey, err := keyGenerator.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		log.Println("doGeneratePublicKey(): error when creating the private key", err)
		return ""
	}

	publicKey := privateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		log.Println("doGeneratePublicKey(): error when decoding the public key", err)
		return ""
	}

	return hex.EncodeToString(publicKeyBytes)
}

func doSignMessage(messageHex string, privateKeyHex string) string {
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
		log.Println("doSignMessage(): error when creating the private key", err)
		return ""
	}

	signature, err := blsSigner.Sign(privateKey, message)
	if err != nil {
		log.Println("doSignMessage(): error when signing the message", err)
		return ""
	}

	return hex.EncodeToString(signature)
}

func doVerifyMessage(publicKeyHex string, messageHex string, signatureHex string) bool {
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
		log.Println("doVerifyMessage(): error when creating the public key", err)
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
