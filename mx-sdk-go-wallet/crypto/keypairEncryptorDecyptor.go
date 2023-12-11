package crypto

//
//import (
//	crypto "github.com/multiversx/mx-chain-crypto-go"
//	"github.com/multiversx/mx-chain-crypto-go/encryption/x25519"
//	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
//)
//
//func KeyPairBasedEncrypt(data []byte, recipientPublicKey crypto.PublicKey, authSecretKey crypto.PrivateKey) *PublicKeyEncryptedData {
//	suite := ed25519.NewEd25519()
//	ephemeralEdScalar, ephemeralEdPoint := suite.CreateKeyPair()
//
//	recipientPubKeyBytes, err := recipientPublicKey.ToByteArray()
//	if err != nil {
//		return nil
//	}
//
//	ed := x25519.EncryptedData{}
//
//	nonce, err := ed.generateEncryptionNonce(data)
//	if err != nil {
//		return err
//	}
//
//	ciphertext, err := ed.createCiphertext(data, ephemeralEdScalar, recipientPubKey, nonce)
//	if err != nil {
//		return err
//	}
//
//	ephemeralEdPointBytes, err := ephemeralEdPoint.MarshalBinary()
//	if err != nil {
//		return err
//	}
//	mac, err := ed.generateMAC(senderPrivateKey, append(ciphertext, ephemeralEdPointBytes...))
//
//	senderPubKey, err := senderPrivateKey.GeneratePublic().ToByteArray()
//	if err != nil {
//		return err
//	}
//}
