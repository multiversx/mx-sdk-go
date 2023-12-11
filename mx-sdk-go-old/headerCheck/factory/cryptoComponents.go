package factory

import (
	"github.com/multiversx/mx-chain-core-go/hashing/blake2b"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	disabledSig "github.com/multiversx/mx-chain-crypto-go/signing/disabled/singlesig"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	mclMultiSig "github.com/multiversx/mx-chain-crypto-go/signing/mcl/multisig"
	"github.com/multiversx/mx-chain-crypto-go/signing/multisig"
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
