package data

import (
	"errors"
	"math/big"
)

var errInvalidBalance = errors.New("invalid balance")

// AccountResponse holds the account endpoint response
type AccountResponse struct {
	Data struct {
		Account *Account `json:"account"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// IsDataTrieMigratedResponse holds the IsDataTrieMigrated endpoint response
type IsDataTrieMigratedResponse struct {
	Data  map[string]bool `json:"data"`
	Error string          `json:"error"`
	Code  string          `json:"code"`
}

// Account holds an Account's information
type Account struct {
	Address         string `json:"address"`
	Nonce           uint64 `json:"nonce"`
	Balance         string `json:"balance"`
	Code            string `json:"code"`
	CodeHash        []byte `json:"codeHash"`
	RootHash        []byte `json:"rootHash"`
	CodeMetadata    []byte `json:"codeMetadata"`
	Username        string `json:"username"`
	DeveloperReward string `json:"developerReward"`
	OwnerAddress    string `json:"ownerAddress"`
}

// GetBalance computes the float representation of the balance,
// based on the provided number of decimals
func (a *Account) GetBalance(decimals int) (float64, error) {
	balance, ok := big.NewFloat(0).SetString(a.Balance)
	if !ok {
		return 0, errInvalidBalance
	}
	// Compute denominated balance to 18 decimals
	denomination := big.NewInt(int64(decimals))
	denominationMultiplier := big.NewInt(10)
	denominationMultiplier.Exp(denominationMultiplier, denomination, nil)
	floatDenomination, _ := big.NewFloat(0).SetString(denominationMultiplier.String())
	balance.Quo(balance, floatDenomination)
	floatBalance, _ := balance.Float64()

	return floatBalance, nil
}

// ESDTFungibleResponse holds the ESDT (fungible) token data endpoint response
type ESDTFungibleResponse struct {
	Data struct {
		TokenData *ESDTFungibleTokenData `json:"tokenData"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// ESDTFungibleTokenData holds the ESDT (fungible) token data definition
type ESDTFungibleTokenData struct {
	TokenIdentifier string `json:"tokenIdentifier"`
	Balance         string `json:"balance"`
	Properties      string `json:"properties"`
}

// ESDTNFTResponse holds the NFT token data endpoint response
type ESDTNFTResponse struct {
	Data struct {
		TokenData *ESDTNFTTokenData `json:"tokenData"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// ESDTNFTTokenData holds the ESDT (NDT, SFT or MetaESDT) token data definition
type ESDTNFTTokenData struct {
	TokenIdentifier string   `json:"tokenIdentifier"`
	Balance         string   `json:"balance"`
	Properties      string   `json:"properties,omitempty"`
	Name            string   `json:"name,omitempty"`
	Nonce           uint64   `json:"nonce,omitempty"`
	Creator         string   `json:"creator,omitempty"`
	Royalties       string   `json:"royalties,omitempty"`
	Hash            []byte   `json:"hash,omitempty"`
	URIs            [][]byte `json:"uris,omitempty"`
	Attributes      []byte   `json:"attributes,omitempty"`
}
