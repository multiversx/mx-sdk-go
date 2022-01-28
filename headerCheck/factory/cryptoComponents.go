package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	disabledSig "github.com/ElrondNetwork/elrond-go-crypto/signing/disabled/singlesig"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/mcl"
	mclMultiSig "github.com/ElrondNetwork/elrond-go-crypto/signing/mcl/multisig"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/multisig"
)

type cryptoComponents struct {
	KeyGen            crypto.KeyGenerator
	MultiSigVerifier  crypto.MultiSigner
	SingleSigVerifier crypto.SingleSigner
}

func CreateCryptoComponents() (*cryptoComponents, error) {
	blockSignKeyGen := signing.NewKeyGenerator(mcl.NewSuiteBLS12())

	interceptSingleSigner := &disabledSig.DisabledSingleSig{}

	multisigHasher, err := blake2b.NewBlake2bWithSize(multisig.BlsHashSize)
	if err != nil {
		return nil, err
	}

	// dummy key
	privateKey, publicKey := blockSignKeyGen.GeneratePair()

	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	multiSigner, err := multisig.NewBLSMultisig(
		&mclMultiSig.BlsMultiSigner{Hasher: multisigHasher},
		[]string{string(publicKeyBytes)},
		privateKey,
		blockSignKeyGen,
		uint16(0),
	)
	if err != nil {
		return nil, err
	}

	return &cryptoComponents{
		KeyGen:            blockSignKeyGen,
		SingleSigVerifier: interceptSingleSigner,
		MultiSigVerifier:  multiSigner,
	}, nil
}
