package interactors

import (
	"encoding/hex"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/hashing/sha256"

	"github.com/ElrondNetwork/elrond-go-core/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/assert"
)

func TestNewMerkleProofVerifier_NilMarshalizer(t *testing.T) {
	t.Parallel()

	mpv, err := NewMerkleProofVerifier(nil, blake2b.NewBlake2b())
	assert.Nil(t, mpv)
	assert.Equal(t, ErrNilMarshalizer, err)
}

func TestNewMerkleProofVerifier_NilHasher(t *testing.T) {
	t.Parallel()

	mpv, err := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, nil)
	assert.Nil(t, mpv)
	assert.Equal(t, ErrNilHasher, err)
}

func TestNewMerkleProofVerifier(t *testing.T) {
	t.Parallel()

	mpv, err := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, blake2b.NewBlake2b())
	assert.Nil(t, err)
	assert.NotNil(t, mpv)
}

func TestMerkleProofVerifier_VerifyProofInvalidRootHash(t *testing.T) {
	t.Parallel()

	mpv, _ := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, blake2b.NewBlake2b())

	rootHash := "not hex characters"
	address := "erd1adfmxhyczrl2t97yx92v5nywqyse0c7qh4xs0p4artg2utnu90pspgvqty"
	p, _ := hex.DecodeString("0a410a0f0202020f07090f030401020607060901050c0903020f060605020b080d080e090900050e02030c0a0e0a0a090a06040500000f0e0a0f0409080f050f090c10124c120200002220d70d8c3ed2b3ff31a6935aa2b6d0ea4cbeeba2fd44922810c0bc5ea9b626adae2a20c9f5f894faef00546a9aaeac32e5099e8d8b2566f239c5196762143f97f222fa3202000001")
	proof := [][]byte{p}

	ok, err := mpv.VerifyProof(rootHash, address, proof)
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestMerkleProofVerifier_VerifyProofInvalidAddress(t *testing.T) {
	t.Parallel()

	mpv, _ := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, blake2b.NewBlake2b())

	rootHash := "79cc4abb4eafd2977a30c4b3c5f4b76a6893d528edc5c7490314b2bb631a4030"
	address := "address"
	p, _ := hex.DecodeString("0a410a0f0202020f07090f030401020607060901050c0903020f060605020b080d080e090900050e02030c0a0e0a0a090a06040500000f0e0a0f0409080f050f090c10124c120200002220d70d8c3ed2b3ff31a6935aa2b6d0ea4cbeeba2fd44922810c0bc5ea9b626adae2a20c9f5f894faef00546a9aaeac32e5099e8d8b2566f239c5196762143f97f222fa3202000001")
	proof := [][]byte{p}

	ok, err := mpv.VerifyProof(rootHash, address, proof)
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestMerkleProofVerifier_VerifyProofWrongProof(t *testing.T) {
	t.Parallel()

	mpv, _ := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, blake2b.NewBlake2b())

	rootHash := "79cc4abb4eafd2977a30c4b3c5f4b76a6893d528edc5c7490314b2bb631a4030"
	address := "erd1adfmxhyczrl2t97yx92v5nywqyse0c7qh4xs0p4artg2utnu90pspgvqty"
	p, _ := hex.DecodeString("0a410a0f0202020f07090f0304010")
	proof := [][]byte{p}

	ok, err := mpv.VerifyProof(rootHash, address, proof)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestMerkleProofVerifier_VerifyProofOk(t *testing.T) {
	t.Parallel()

	mpv, _ := NewMerkleProofVerifier(&marshal.GogoProtoMarshalizer{}, sha256.NewSha256())

	rootHash := "bc2e549d98c31ffe6e9419b933d03b37e84f74c42601412302799d277651a6d8"
	address := "erd1hapzzd68d9lfmmzzz8h4pwnqvx65w2d48wsvfx2ffr9tg7903p2qm2rku5"
	p, _ := hex.DecodeString("0a41040508080f0a0807040b0a0c080409040909040c000a0b03050b09020704050b010600060a0b00050f0e010102040c0e0d090e07090607040703010202040f0b10124c1202000022206182d14320be95434f5508acad9478d3b6cf837bfce7ebfe47c2e860d1b98ca72a20bf42213747697e9dec4211ef50ba6061b54729b53ba0c4994948cab478af88543202000001")
	proof := [][]byte{p}

	ok, err := mpv.VerifyProof(rootHash, address, proof)
	assert.Nil(t, err)
	assert.True(t, ok)
}
