package interactors

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-go/trie"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/disabled"
)

type merkleProofVerifier struct {
	marshalizer marshal.Marshalizer
	hasher      hashing.Hasher
}

func NewMerkleProofVerifier(marshalizer marshal.Marshalizer, hasher hashing.Hasher) (*merkleProofVerifier, error) {
	if check.IfNil(marshalizer) {
		return nil, ErrNilMarshalizer
	}
	if check.IfNil(hasher) {
		return nil, ErrNilHasher
	}

	return &merkleProofVerifier{
		marshalizer: marshalizer,
		hasher:      hasher,
	}, nil
}

func (mpv *merkleProofVerifier) VerifyProof(rootHash string, address string, proof [][]byte) (bool, error) {
	rootHashBytes, err := hex.DecodeString(rootHash)
	if err != nil {
		return false, err
	}

	key, err := getKeyBytes(address)
	if err != nil {
		return false, err
	}

	tr, err := trie.NewTrie(&disabled.StorageManager{}, mpv.marshalizer, mpv.hasher, 5)
	if err != nil {
		return false, err
	}

	return tr.VerifyProof(rootHashBytes, key, proof)
}

func getKeyBytes(key string) ([]byte, error) {
	addressBytes, err := core.AddressPublicKeyConverter.Decode(key)
	if err == nil {
		return addressBytes, nil
	}

	return hex.DecodeString(key)
}
