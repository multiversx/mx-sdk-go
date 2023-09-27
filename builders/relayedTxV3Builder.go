package builders

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
)

type relayedTxV3Builder struct {
	innerTransaction *transaction.FrontendTransaction
	relayerAccount   *data.Account
	networkConfig    *data.NetworkConfig
}

// NewRelayedTxV3Builder creates a new relayed transaction v2 builder
func NewRelayedTxV3Builder() *relayedTxV3Builder {
	return &relayedTxV3Builder{
		innerTransaction: nil,
		relayerAccount:   nil,
		networkConfig:    nil,
	}
}

// SetInnerTransaction sets the inner transaction to be relayed
func (rtb *relayedTxV3Builder) SetInnerTransaction(tx *transaction.FrontendTransaction) *relayedTxV3Builder {
	rtb.innerTransaction = tx

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
	if rtb.innerTransaction == nil {
		return nil, ErrNilInnerTransaction
	}
	if len(rtb.innerTransaction.Signature) == 0 {
		return nil, ErrNilInnerTransactionSignature
	}
	if len(rtb.innerTransaction.Relayer) == 0 {
		return nil, ErrEmptyRelayerOnInnerTransaction
	}
	if rtb.relayerAccount == nil {
		return nil, ErrNilRelayerAccount
	}
	if rtb.networkConfig == nil {
		return nil, ErrNilNetworkConfig
	}

	gasLimit := rtb.networkConfig.MinGasLimit + rtb.innerTransaction.GasLimit

	relayedTx := &transaction.FrontendTransaction{
		Nonce:            rtb.relayerAccount.Nonce,
		Value:            "0",
		Receiver:         rtb.innerTransaction.Sender,
		Sender:           rtb.relayerAccount.Address,
		GasPrice:         rtb.innerTransaction.GasPrice,
		GasLimit:         gasLimit,
		Data:             []byte(""),
		ChainID:          rtb.networkConfig.ChainID,
		Version:          rtb.networkConfig.MinTransactionVersion,
		InnerTransaction: rtb.innerTransaction,
	}

	return relayedTx, nil
}
