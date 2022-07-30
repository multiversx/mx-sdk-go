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

type AccountKeys map[string]string

type AccountKeysResponse struct {
	Data struct {
		Pairs AccountKeys `json:"pairs"`
	} `json:"data"`

	Error string `json:"error"`
	Code  string `json:"code"`
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
