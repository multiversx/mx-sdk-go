package codec

import (
	"fmt"
	"io"
)

func readBytesExactly(reader io.Reader, numBytes int) ([]byte, error) {
	if numBytes == 0 {
		return []byte{}, nil
	}

	data := make([]byte, numBytes)
	n, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != numBytes {
		return nil, fmt.Errorf("cannot read exactly %d bytes", numBytes)
	}

	return data, err
}

func checkPubKeyLength(pubkey []byte) error {
	if len(pubkey) != pubKeyLength {
		return fmt.Errorf("public key (address) has invalid length: %d", len(pubkey))
	}

	return nil
}
