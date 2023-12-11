package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMnemonic_ToSplitMnemonicWords(t *testing.T) {
	t.Parallel()

	m := Mnemonic("tag volcano eight thank tide danger coast health above argue embrace heavy")
	expected := []string{"tag", "volcano", "eight", "thank", "tide", "danger", "coast", "health", "above", "argue", "embrace", "heavy"}

	assert.Equal(t, expected, m.ToSplitMnemonicWords())
}
