package coreV2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	extendedNftIdentifier      = "TEST-123456-0a"
	extendedFungibleIdentifier = "FNG-123456"
	tkComputer                 = NewTokenComputer()
)

func TestTokenComputer_IsFungible(t *testing.T) {
	fungibleToken := Token{Identifier: "TEST-123456"}
	nonFungibleToken := Token{Identifier: "NFT-987654", Nonce: 7}

	require.True(t, tkComputer.IsFungible(fungibleToken))
	require.False(t, tkComputer.IsFungible(nonFungibleToken))
}

func TestTokenComputer_ExtractNonceFromExtendedIdentifier(t *testing.T) {
	nftNonce, err := tkComputer.ExtractNonceFromExtendedIdentifier(extendedNftIdentifier)
	require.NoError(t, err, "failed to extract nonce from extended identifier")
	require.Equal(t, uint64(10), nftNonce)

	fungibleNonce, err := tkComputer.ExtractNonceFromExtendedIdentifier(extendedFungibleIdentifier)
	require.NoError(t, err, "failed to extract nonce from extended identifier")
	require.Equal(t, uint64(0), fungibleNonce)
}

func TestTokenComputer_ExtractIdentifierFromExtendedIdentifier(t *testing.T) {
	nftIdentifier, err := tkComputer.ExtractIdentifierFromExtendedIdentifier(extendedNftIdentifier)
	require.NoError(t, err, "failed to extract identifier from extender identifier")
	require.Equal(t, "TEST-123456", nftIdentifier)

	fungibleIdentifier, err := tkComputer.ExtractIdentifierFromExtendedIdentifier(extendedFungibleIdentifier)
	require.NoError(t, err, "failed to extract identifier from extender identifier")
	require.Equal(t, "FNG-123456", fungibleIdentifier)
}

func TestTokenComputer_ExtractTickerFromIdentifier(t *testing.T) {
	fungibleIdentifier := "FNG-123456"
	nonFungibleIdentifier := "NFT-987654-0a"

	fungibleTicker, err := tkComputer.ExtractTickerFromIdentifier(fungibleIdentifier)
	require.NoError(t, err, "failed to extract ticker from identifier")
	require.Equal(t, "FNG", fungibleTicker)

	nonFungibleTicker, err := tkComputer.ExtractTickerFromIdentifier(nonFungibleIdentifier)
	require.NoError(t, err, "failed to extract ticker from identifier")
	require.Equal(t, "NFT", nonFungibleTicker)
}

func TestTokenComputer_ParseExtendedIdentifierParts(t *testing.T) {
	fungibleIdentifier := "FNG-123456"
	nonFungibleIdentifier := "NFT-987654-0a"

	fungibleParts, err := tkComputer.ParseExtendedIdentifierParts(fungibleIdentifier)
	require.NoError(t, err, "failed to parse extended identifier parts")
	require.Equal(t, &TokenIdentifierParts{Ticker: "FNG", RandomSequence: "123456", Nonce: 0}, fungibleParts)

	nonFungibleParts, err := tkComputer.ParseExtendedIdentifierParts(nonFungibleIdentifier)
	require.NoError(t, err, "failed to parse extended identifier parts")
	require.Equal(t, &TokenIdentifierParts{Ticker: "NFT", RandomSequence: "987654", Nonce: 10}, nonFungibleParts)
}

func TestTokenComputer_ComputeExtendedIdentifierFromIdentifierAndNonce(t *testing.T) {
	fungibleIdentifier := "FNG-123456"
	fungibleNonce := 0

	nonFungibleIdentifier := "NFT-987654"
	nonFungibleNonce := 10

	fungibleTokenIdentifier, err := tkComputer.ComputeExtendedIdentifierFromIdentifierAndNonce(
		fungibleIdentifier,
		uint64(fungibleNonce),
	)
	require.NoError(t, err, "failed to compute extended identifier from identifier and nonce")
	require.Equal(t, "FNG-123456", fungibleTokenIdentifier)

	nonFungibleTokenIdentifier, err := tkComputer.ComputeExtendedIdentifierFromIdentifierAndNonce(
		nonFungibleIdentifier,
		uint64(nonFungibleNonce),
	)
	require.NoError(t, err, "failed to compute extended identifier from identifier and nonce")
	require.Equal(t, "NFT-987654-0a", nonFungibleTokenIdentifier)
}
