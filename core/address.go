package core

import (
	"encoding/hex"
	"fmt"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"strings"
)

const (
	CHARSET = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

type address struct {
	publicKey []byte
	hrp       string
}

type AddressGetter interface {
	GetPublicKey() []byte
	GetHRP() string
}

type Address interface {
	AddressGetter
	ToHex() string
	ToBech32() (string, error)
	IsSmartContract() bool
}

func NewAddress(publicKey []byte, hrp string) (Address, error) {
	if len(publicKey) != AddressBytesLen {
		return nil, pubkeyConverter.ErrInvalidAddressLength
	}

	return &address{publicKey: publicKey, hrp: hrp}, nil
}

func NewAddressFromBech32(value string) (Address, error) {
	hrp, err := determineHRP(value)
	if err != nil {
		return nil, fmt.Errorf("failed to determine HRP: %v", err)
	}
	converter, err := pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, hrp)
	if err != nil {
		return nil, fmt.Errorf("failed to create bech32 converter: %v", err)
	}

	decode, err := converter.Decode(value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q to bech32: %v", decode, err)
	}

	return &address{publicKey: decode, hrp: hrp}, nil
}

func NewAddressFromHex(value, hrp string) (Address, error) {
	converter, err := pubkeyConverter.NewHexPubkeyConverter(AddressBytesLen)
	if err != nil {
		return nil, fmt.Errorf("failed to create hex converter: %v", err)
	}

	publicKey, err := converter.Decode(value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q to hex: %v", value, err)
	}

	return &address{publicKey: publicKey, hrp: hrp}, nil
}

func (a *address) GetPublicKey() []byte {
	return a.publicKey
}

func (a *address) GetHRP() string {
	return a.hrp
}

func (a *address) ToHex() string {
	return hex.EncodeToString(a.publicKey)
}

func (a *address) ToBech32() (string, error) {
	converter, err := pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, a.hrp)
	if err != nil {
		return "", fmt.Errorf("failed to create bech32 converter: %v", err)
	}
	bech32, err := converter.Encode(a.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert address to bech32: %v", err)
	}

	return bech32, nil
}

func (a *address) IsSmartContract() bool {
	return strings.HasPrefix(a.ToHex(), ScHexPubKeyPrefix)
}

func determineHRP(bech string) (string, error) {
	if len(bech) < 32 || len(bech) > 125 {
		return "", fmt.Errorf("address provided %q doesn't have the correct length", bech)
	}

	if strings.ToLower(bech) != bech && strings.ToUpper(bech) != bech {
		return "", fmt.Errorf("address provided %q has both uppercase and lowercase characters", bech)
	}

	bech = strings.ToLower(bech)
	pos := strings.Index(bech, "1")

	if pos < 1 || pos+7 > len(bech) || len(bech) > 90 {
		return "", fmt.Errorf("address provided %q is invalid", bech)
	}

	for _, x := range bech[pos+1:] {
		if !strings.Contains(CHARSET, string(x)) {
			return "", fmt.Errorf("address provided %q is invalid", bech)
		}
	}

	return bech[:pos], nil
}
