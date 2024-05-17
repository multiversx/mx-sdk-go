package builders

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
)

type relayedTxV3Builder struct {
	innerTransactions []*transaction.FrontendTransaction
	relayerAccount    *data.Account
	networkConfig     *data.NetworkConfig
}

// NewRelayedTxV3Builder creates a new relayed transaction v2 builder
func NewRelayedTxV3Builder() *relayedTxV3Builder {
	return &relayedTxV3Builder{
		innerTransactions: nil,
		relayerAccount:    nil,
		networkConfig:     nil,
	}
}

// SetInnerTransactions sets the inner transactions to be relayed
func (rtb *relayedTxV3Builder) SetInnerTransactions(innerTxs []*transaction.FrontendTransaction) *relayedTxV3Builder {
	rtb.innerTransactions = innerTxs

	return rtb
}

// SetRelayerAccount sets the relayer account (that will send the wrapped transaction)
func (rtb *relayedTxV3Builder) SetRelayerAccount(account *data.Account) *relayedTxV3Builder {
	rtb.relayerAccount = account

	return rtb
}

// SetNetworkConfig sets the network config
func (rtb *relayedTxV3Builder) SetNetworkConfig(config *data.NetworkConfig) *relayedTxV3Builder {
	rtb.networkConfig = config

	return rtb
}

// Build builds the relayed transaction v3
// The returned transaction will not be signed
func (rtb *relayedTxV3Builder) Build() (*transaction.FrontendTransaction, error) {
	if len(rtb.innerTransactions) == 0 {
		return nil, ErrEmptyInnerTransactions
	}
	innerTxsGasLimit := uint64(0)
	for _, innerTx := range rtb.innerTransactions {
		if len(innerTx.Signature) == 0 {
			return nil, ErrNilInnerTransactionSignature
		}
		if len(innerTx.Relayer) == 0 {
			return nil, ErrEmptyRelayerOnInnerTransaction
		}

		innerTxsGasLimit += innerTx.GasLimit
	}

	if rtb.relayerAccount == nil {
		return nil, ErrNilRelayerAccount
	}
	if rtb.networkConfig == nil {
		return nil, ErrNilNetworkConfig
	}

	minGasLimit := rtb.networkConfig.MinGasLimit
	moveBalancesGas := minGasLimit * uint64(len(rtb.innerTransactions))
	gasLimit := minGasLimit + moveBalancesGas + innerTxsGasLimit

	relayedTx := &transaction.FrontendTransaction{
		Nonce:             rtb.relayerAccount.Nonce,
		Value:             "0",
		Receiver:          rtb.relayerAccount.Address,
		Sender:            rtb.relayerAccount.Address,
		GasPrice:          rtb.innerTransactions[0].GasPrice,
		GasLimit:          gasLimit,
		Data:              []byte(""),
		ChainID:           rtb.networkConfig.ChainID,
		Version:           rtb.networkConfig.MinTransactionVersion,
		InnerTransactions: rtb.innerTransactions,
	}

	return relayedTx, nil
}
