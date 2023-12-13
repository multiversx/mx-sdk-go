package core

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	TokenRandomSequenceLength = 6
)

var (
	alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

type Token struct {
	Identifier string
	Nonce      uint64
}

type TokenIdentifierParts struct {
	Ticker         string
	RandomSequence string
	Nonce          uint64
}

type TokenTransfer struct {
	Token  Token
	Amount Amount
}

type TokenComputer interface {
	IsFungible(token Token) bool
	ExtractNonceFromExtendedIdentifier(extendedIdentifier string) (uint64, error)
	ExtractIdentifierFromExtendedIdentifier(extendedIdentifier string) (string, error)
	ExtractTickerFromIdentifier(identifier string) (string, error)
	ParseExtendedIdentifierParts(identifier string) (*TokenIdentifierParts, error)
	ComputeExtendedIdentifierFromIdentifierAndNonce(identifier string, nonce uint64) (string, error)
	ComputeExtendedIdentifierFromParts(ticker string, randomSequence string, nonce uint64) (string, error)
}

type tokenComputer struct{}

func NewTokenComputer() TokenComputer {
	return &tokenComputer{}
}

func (t *tokenComputer) IsFungible(token Token) bool {
	return token.Nonce == 0
}

func (t *tokenComputer) ExtractNonceFromExtendedIdentifier(extendedIdentifier string) (uint64, error) {
	parts := strings.Split(extendedIdentifier, "-")

	err := checkIfExtendedIdentifierWasProvided(parts)
	if err != nil {
		return 0, fmt.Errorf("extended identifier was not provided: %v", err)
	}
	err = checkLengthOfRandomSequence(parts[1])
	if err != nil {
		return 0, fmt.Errorf("extended identifier was not provided: %v", err)
	}

	//# in case the identifier of a fungible token is provided
	if len(parts) == 2 {
		return 0, nil
	}

	hexNonce, err := hex.DecodeString(parts[2])
	if err != nil {
		return 0, fmt.Errorf("failed to decode string to hex: %v", err)
	}

	return uint64(hexNonce[0]), nil
}

func (t *tokenComputer) ExtractIdentifierFromExtendedIdentifier(extendedIdentifier string) (string, error) {
	parts := strings.Split(extendedIdentifier, "-")

	err := checkIfExtendedIdentifierWasProvided(parts)
	if err != nil {
		return "", fmt.Errorf("extended identifier was not provided: %v", err)
	}
	err = ensureTokenTickerValidity(parts[0])
	if err != nil {
		return "", fmt.Errorf("token ticker is not valid: %v", err)
	}
	err = checkLengthOfRandomSequence(parts[1])
	if err != nil {
		return "", fmt.Errorf("length of random sequence is incorrect: %v", err)
	}

	return parts[0] + "-" + parts[1], nil
}

func (t *tokenComputer) ExtractTickerFromIdentifier(identifier string) (string, error) {
	parts := strings.Split(identifier, "-")

	err := checkLengthOfRandomSequence(parts[1])
	if err != nil {
		return "", fmt.Errorf("length of random sequence is incorrect: %v", err)
	}
	err = ensureTokenTickerValidity(parts[0])
	if err != nil {
		return "", fmt.Errorf("token ticker is not valid: %v", err)
	}

	return parts[0], nil
}

func (t *tokenComputer) ParseExtendedIdentifierParts(identifier string) (*TokenIdentifierParts, error) {
	var (
		nonce    uint64
		hexNonce []byte
	)

	parts := strings.Split(identifier, "-")

	err := checkIfExtendedIdentifierWasProvided(parts)
	if err != nil {
		return nil, fmt.Errorf("extended identifier was not provided: %v", err)
	}
	err = checkLengthOfRandomSequence(parts[1])
	if err != nil {
		return nil, fmt.Errorf("length of random sequence is incorrect: %v", err)
	}
	err = ensureTokenTickerValidity(parts[0])
	if err != nil {
		return nil, fmt.Errorf("token ticker is not valid: %v", err)
	}

	if len(parts) == 3 {
		hexNonce, err = hex.DecodeString(parts[2])
		if err != nil {
			return nil, fmt.Errorf("failed to decode string to hex: %v", err)
		}
		nonce = uint64(hexNonce[0])
	}

	return &TokenIdentifierParts{parts[0], parts[1], nonce}, nil
}

func (t *tokenComputer) ComputeExtendedIdentifierFromIdentifierAndNonce(identifier string, nonce uint64) (string, error) {
	parts := strings.Split(identifier, "-")

	err := checkLengthOfRandomSequence(parts[1])
	if err != nil {
		return "", fmt.Errorf("length of random sequence is incorrect: %v", err)
	}
	err = ensureTokenTickerValidity(parts[0])
	if err != nil {
		return "", fmt.Errorf("token ticker is not valid: %v", err)
	}

	if nonce == 0 {
		return identifier, nil
	}

	encodedNonce := encodeUnsignedNumber(nonce)

	return fmt.Sprintf("%s-%s", identifier, hex.EncodeToString(encodedNonce)), nil
}

func (t *tokenComputer) ComputeExtendedIdentifierFromParts(ticker string, randomSequence string, nonce uint64) (string, error) {
	identifier := ticker + randomSequence
	return t.ComputeExtendedIdentifierFromIdentifierAndNonce(identifier, nonce)
}

func checkIfExtendedIdentifierWasProvided(tokenParts []string) error {
	//this is for the identifiers of fungible tokens
	minExtendedIdentifierLengthIfSplit := 2
	//# this is for the identifiers of nft, sft and meta-esdt
	maxExtendedIdentifierLengthIfSplit := 3

	if len(tokenParts) < minExtendedIdentifierLengthIfSplit || len(tokenParts) > maxExtendedIdentifierLengthIfSplit {
		return errors.New("invalid extended token identifier provided")
	}

	return nil
}

func checkLengthOfRandomSequence(randomSequence string) error {
	if len(randomSequence) != TokenRandomSequenceLength {
		return errors.New("the identifier is not valid. The random sequence does not have the right length")
	}

	return nil
}

func ensureTokenTickerValidity(ticker string) error {
	MinTickerLength := 3
	MaxTickerLength := 10

	if len(ticker) < MinTickerLength || len(ticker) > MaxTickerLength {
		return fmt.Errorf("the token ticker should be between %d and %d characters", MinTickerLength, MaxTickerLength)
	}

	if !isAlphaNumeric(ticker) {
		return fmt.Errorf("the token ticker should only contain alphanumeric characters")
	}

	if !isUpper(ticker) {
		fmt.Println("the token ticker should be upper case")
	}

	return nil
}

func isAlphaNumeric(str string) bool {
	return alphanumeric.MatchString(str)
}

func isUpper(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func encodeUnsignedNumber(arg uint64) []byte {
	// Determine the maximum number of bytes needed based on the size of int
	const IntegerMaxNumBytes = 8 // Assuming a 64-bit integer

	// Convert int to bytes
	argBytes := make([]byte, IntegerMaxNumBytes)
	binary.BigEndian.PutUint64(argBytes, arg)

	// Remove leading zero bytes
	argBytes = removeLeadingZeros(argBytes)

	return argBytes
}

// removeLeadingZeros removes leading zero bytes from a byte slice
func removeLeadingZeros(data []byte) []byte {
	for i, b := range data {
		if b != 0 {
			return data[i:]
		}
	}
	return []byte{0} // If all bytes are zero, return a single zero byte
}
