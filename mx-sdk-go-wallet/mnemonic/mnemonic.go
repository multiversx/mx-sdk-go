package mnemonic

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

const (
	mnemonicBitSize = 256
	egldCoinType    = uint32(508)
	hardened        = uint32(0x80000000)
	keystoreVersion = 4
	keyHeaderKDF    = "scrypt"
	scryptN         = 4096
	scryptR         = 8
	scryptP         = 1
	scryptDKLen     = 32
	addressLen      = 32
	mnemonicKind    = "mnemonic"
)

type bip32Path []uint32

type bip32 struct {
	Key       []byte
	ChainCode []byte
}

type Mnemonic struct {
	Text string
}

func NewMnemonicFromText(text string) (*Mnemonic, error) {
	if !bip39.IsMnemonicValid(text) {
		return nil, errors.New(fmt.Sprintf("failed to create mnemonic from text: %s", text))
	}

	return &Mnemonic{Text: text}, nil
}

func (m *Mnemonic) DeriveKey(addressIndex uint32, password string) []byte {
	seed := bip39.NewSeed(m.Text, password)
	var egldPath = bip32Path{
		44 | hardened,
		egldCoinType | hardened,
		hardened, // account
		hardened,
		hardened, // addressIndex
	}

	egldPath[2] = 0 | hardened
	egldPath[4] = addressIndex | hardened
	keyData := derivePrivateKey(seed, egldPath)

	return keyData.Key
}

func (m *Mnemonic) GetWords() []string {
	return strings.Split(m.Text, " ")
}

func (m *Mnemonic) String() string { return m.Text }

func derivePrivateKey(seed []byte, path bip32Path) *bip32 {
	b := &bip32{}
	digest := hmac.New(sha512.New, []byte("ed25519 seed"))
	digest.Write(seed)
	intermediary := digest.Sum(nil)
	serializedKeyLen := 32
	serializedChildIndexLen := 4
	hardenedChildPadding := byte(0x00)
	b.Key = intermediary[:serializedKeyLen]
	b.ChainCode = intermediary[serializedKeyLen:]
	for _, childIdx := range path {
		buff := make([]byte, 1+serializedKeyLen+4)
		buff[0] = hardenedChildPadding
		copy(buff[1:1+serializedKeyLen], b.Key)
		binary.BigEndian.PutUint32(buff[1+serializedKeyLen:1+serializedKeyLen+serializedChildIndexLen], childIdx)
		digest = hmac.New(sha512.New, b.ChainCode)
		digest.Write(buff)
		intermediary = digest.Sum(nil)
		b.Key = intermediary[:serializedKeyLen]
		b.ChainCode = intermediary[serializedKeyLen:]
	}

	return b
}
