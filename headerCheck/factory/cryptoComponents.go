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
	KeyGen    crypto.KeyGenerator
	MultiSig  crypto.MultiSigner
	SingleSig crypto.SingleSigner
	PublicKey crypto.PublicKey
}

// CreateCryptoComponents creates crypto components needed for header verification
func CreateCryptoComponents() (*cryptoComponents, error) {
	blockSignKeyGen := signing.NewKeyGenerator(mcl.NewSuiteBLS12())

	interceptSingleSigner := &disabledSig.DisabledSingleSig{}

	multisigHasher, err := blake2b.NewBlake2bWithSize(mclMultiSig.HasherOutputSize)
	if err != nil {
		return nil, err
	}

	// dummy key
	_, publicKey := blockSignKeyGen.GeneratePair()

	multiSigner, err := multisig.NewBLSMultisig(
		&mclMultiSig.BlsMultiSigner{Hasher: multisigHasher},
		blockSignKeyGen,
	)
	if err != nil {
		return nil, err
	}

	return &cryptoComponents{
		KeyGen:    blockSignKeyGen,
		SingleSig: interceptSingleSigner,
		MultiSig:  multiSigner,
		PublicKey: publicKey,
	}, nil
}
