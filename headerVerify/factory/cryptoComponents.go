package factory

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go-core/hashing/sha256"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	disabledCrypto "github.com/ElrondNetwork/elrond-go-crypto/signing/disabled"
	disabledSig "github.com/ElrondNetwork/elrond-go-crypto/signing/disabled/singlesig"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/mcl"
	mclMultiSig "github.com/ElrondNetwork/elrond-go-crypto/signing/mcl/multisig"
	mclSig "github.com/ElrondNetwork/elrond-go-crypto/signing/mcl/singlesig"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/multisig"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/consensus"
	"github.com/ElrondNetwork/elrond-go/errors"
)

var log = logger.GetOrCreate("example/leveldb/factory")

var ErrInvalidConsensusConfig = fmt.Errorf("invalid consensus type provided in config file")

const disabledSigChecking = "disabled"

type cryptoComponents struct {
	KeyGen            crypto.KeyGenerator
	MultiSigVerifier  crypto.MultiSigner
	SingleSigVerifier crypto.SingleSigner
}

func CreateCryptoComponents() (*cryptoComponents, error) {
	consensusType := "bls"
	multisigHasherType := "blake2b"

	suite, err := getSuite(consensusType)
	if err != nil {
		return nil, err
	}

	blockSignKeyGen := signing.NewKeyGenerator(suite)

	interceptSingleSigner, err := createSingleSigner(consensusType)
	if err != nil {
		return nil, err
	}

	multisigHasher, err := getMultiSigHasher(consensusType, multisigHasherType)
	if err != nil {
		return nil, err
	}

	privateKey, publicKey := blockSignKeyGen.GeneratePair()

	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	//multiSigner, err := createMultiSigner(consensusType, multisigHasher, cp, blockSignKeyGen)
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

func getSuite(consensusType string) (crypto.Suite, error) {
	switch consensusType {
	case consensus.BlsConsensusType:
		return mcl.NewSuiteBLS12(), nil
	case disabledSigChecking:
		return disabledCrypto.NewDisabledSuite(), nil
	default:
		return nil, ErrInvalidConsensusConfig
	}
}

func createSingleSigner(consensusType string) (crypto.SingleSigner, error) {
	switch consensusType {
	case consensus.BlsConsensusType:
		return &mclSig.BlsSingleSigner{}, nil
	case disabledSigChecking:
		log.Warn("using disabled single signer")
		return &disabledSig.DisabledSingleSig{}, nil
	default:
		return nil, errors.ErrInvalidConsensusConfig
	}
}

func getMultiSigHasher(consensusType string, multisigHasherType string) (hashing.Hasher, error) {
	if consensusType == consensus.BlsConsensusType && multisigHasherType != "blake2b" {
		return nil, errors.ErrMultiSigHasherMissmatch
	}

	switch multisigHasherType {
	case "sha256":
		return sha256.NewSha256(), nil
	case "blake2b":
		if consensusType == consensus.BlsConsensusType {
			return blake2b.NewBlake2bWithSize(multisig.BlsHashSize)
		}
		return blake2b.NewBlake2b(), nil
	}

	return nil, errors.ErrMissingMultiHasherConfig
}
