package blockchain

import (
	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519/singlesig"
)

var (
	suite         = ed25519.NewEd25519()
	hasher        = keccak.NewKeccak()
	singleSigner  = &singlesig.Ed25519Signer{}
	messagePrefix = []byte("\x17Elrond Signed Message:\n")
)

// xSigner contains the primitives used to correctly sign a transaction
type xSigner struct {
	keyGen crypto.KeyGenerator
}

// NewXSigner will create a new instance of xSigner
func NewXSigner() *xSigner {
	return &xSigner{
		keyGen: signing.NewKeyGenerator(suite),
	}
}

// SignMessage will generate the signature providing the private key bytes and the message bytes
// prepending the standard messagePrefix
func (xs *xSigner) SignMessage(msg []byte, skBytes []byte) ([]byte, error) {
	sk, err := xs.keyGen.PrivateKeyFromByteArray(skBytes)
	if err != nil {
		return nil, err
	}

	serializedMessage, err := xs.serializeForSigning(msg)
	if err != nil {
		return nil, err
	}

	return singleSigner.Sign(sk, serializedMessage)
}

// SignTransaction will generate the signature providing the private key bytes and the serialized form of the transaction
func (xs *xSigner) SignTransaction(tx []byte, skBytes []byte) ([]byte, error) {
	sk, err := xs.keyGen.PrivateKeyFromByteArray(skBytes)
	if err != nil {
		return nil, err
	}

	return singleSigner.Sign(sk, tx)
}

// GeneratePkBytes will generate the public key bytes out of the provided private key bytes
func (xs *xSigner) GeneratePkBytes(skBytes []byte) ([]byte, error) {
	sk, err := xs.keyGen.PrivateKeyFromByteArray(skBytes)
	if err != nil {
		return nil, err
	}

	pk := sk.GeneratePublic()

	return pk.ToByteArray()
}

func (xs *xSigner) serializeForSigning(msg []byte) ([]byte, error) {
	msg = append(messagePrefix, msg...)

	return hasher.Compute(string(msg)), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (xs *xSigner) IsInterfaceNil() bool {
	return xs == nil
}
