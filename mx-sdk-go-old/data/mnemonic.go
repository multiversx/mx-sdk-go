package data

import "strings"

const spaceChar = " "

// Mnemonic will hold the mnemonic info
type Mnemonic string

// ToSplitMnemonicWords splits the complete strings in words
func (m Mnemonic) ToSplitMnemonicWords() []string {
	return strings.Split(string(m), spaceChar)
}
