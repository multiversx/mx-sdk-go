package mnemonic

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	mnemonic := Mnemonic{"moral volcano peasant pass circle pen over picture flat shop clap goat never lyrics gather prepare woman film husband gravity behind test tiger improve"}

	key1 := hex.EncodeToString(mnemonic.DeriveKey(0, ""))
	require.Equal(t, key1, "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9")

	key2 := hex.EncodeToString(mnemonic.DeriveKey(1, ""))
	require.Equal(t, key2, "b8ca6f8203fb4b545a8e83c5384da033c415db155b53fb5b8eba7ff5a039d639")

	key3 := hex.EncodeToString(mnemonic.DeriveKey(2, ""))
	require.Equal(t, key3, "e253a571ca153dc2aee845819f74bcc9773b0586edead15a94cb7235a5027436")
}

func TestNewMnemonicFromText(t *testing.T) {
	badText := "bad mnemonic"
	mnemonic, err := NewMnemonicFromText(badText)
	require.Nil(t, mnemonic)
	require.Error(t, err, fmt.Sprintf("failed to create mnemonic from text: %s", badText))

	goodText := "moral volcano peasant pass circle pen over picture flat shop clap goat never lyrics gather prepare woman film husband gravity behind test tiger improve"
	mnemonic, err = NewMnemonicFromText(goodText)
	require.Equal(t, mnemonic, &Mnemonic{Text: goodText})
	require.Nil(t, err)
}

func TestGenerateAndWords(t *testing.T) {
	mnemonic, err := NewMnemonicFromText("moral volcano peasant pass circle pen over picture flat shop clap goat never lyrics gather prepare woman film husband gravity behind test tiger improve")
	require.NoError(t, err, "failed to create mnemonic")
	words := mnemonic.GetWords()
	require.Len(t, words, 24)
}
