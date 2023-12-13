package core

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageComputer_ComputeBytesForSigning(t *testing.T) {
	message := Message{Data: []byte("test message")}
	msgComputer := NewMessageComputer()
	serialized := msgComputer.ComputeBytesForSigning(message)
	require.Equal(t, hex.EncodeToString(serialized), "2162d6271208429e6d3e664139e98ba7c5f1870906fb113e8903b1d3f531004d")
}
