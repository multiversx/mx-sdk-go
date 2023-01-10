package builders

import (
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

type relayedTxV1Builder struct {
	innerTransaction *data.Transaction
	relayerAccount   *data.Account
	networkConfig    *data.NetworkConfig
}

// NewRelayedTxV1Builder creates a new relayed transaction v1 builder
func NewRelayedTxV1Builder() *relayedTxV1Builder {
	return &relayedTxV1Builder{
		innerTransaction: nil,
		relayerAccount:   nil,
		networkConfig:    nil,
	}
}

// SetInnerTransaction sets the inner transaction to be relayed
func (rtb *relayedTxV1Builder) SetInnerTransaction(tx *data.Transaction) *relayedTxV1Builder {
	rtb.innerTransaction = tx

	return rtb
}

// SetRelayerAccount sets the relayer account
func (rtb *relayedTxV1Builder) SetRelayerAccount(account *data.Account) *relayedTxV1Builder {
	rtb.relayerAccount = account

	return rtb
}

// SetNetworkConfig sets the network config
func (rtb *relayedTxV1Builder) SetNetworkConfig(config *data.NetworkConfig) *relayedTxV1Builder {
	rtb.networkConfig = config

	return rtb
}

// Build builds the relayed transaction v1
// The returned transaction will not be signed
func (rtb *relayedTxV1Builder) Build() (*data.Transaction, error) {
	if rtb.innerTransaction == nil {
		return nil, ErrNilInnerTransaction
	}
	if len(rtb.innerTransaction.Signature) == 0 {
		return nil, ErrNilInnerTransactionSignature
	}
	if rtb.relayerAccount == nil {
		return nil, ErrNilRelayerAccount
	}
	if rtb.networkConfig == nil {
		return nil, ErrNilNetworkConfig
	}

	innerTxHex, err := prepareInnerTxForRelayV1(rtb.innerTransaction)
	if err != nil {
		return nil, err
	}

	payload := []byte("relayedTx@" + innerTxHex)
	gasLimit := rtb.networkConfig.MinGasLimit + rtb.networkConfig.GasPerDataByte*uint64(len(payload)) + rtb.innerTransaction.GasLimit

	innerTxValue, ok := big.NewInt(0).SetString(rtb.innerTransaction.Value, 10)
	if !ok {
		return nil, ErrInvalidValue
	}

	relayedTx := &data.Transaction{
		Nonce:    rtb.relayerAccount.Nonce,
		Value:    innerTxValue.String(),
		RcvAddr:  rtb.innerTransaction.SndAddr,
		SndAddr:  rtb.relayerAccount.Address,
		GasPrice: rtb.innerTransaction.GasPrice,
		GasLimit: gasLimit,
		Data:     payload,
		ChainID:  rtb.networkConfig.ChainID,
		Version:  rtb.networkConfig.MinTransactionVersion,
	}

	return relayedTx, nil
}

func prepareInnerTxForRelayV1(tx *data.Transaction) (string, error) {
	txValue, ok := big.NewInt(0).SetString(tx.Value, 10)
	if !ok {
		return "", ErrInvalidValue
	}

	addressConverter := core.AddressPublicKeyConverter
	receiverAddress, err := addressConverter.Decode(tx.RcvAddr)
	if err != nil {
		return "", err
	}

	senderAddress, err := addressConverter.Decode(tx.SndAddr)
	if err != nil {
		return "", err
	}

	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return "", err
	}

	coreTx := &transaction.Transaction{
		Nonce:     tx.Nonce,
		Value:     txValue,
		RcvAddr:   receiverAddress,
		SndAddr:   senderAddress,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      tx.Data,
		ChainID:   []byte(tx.ChainID),
		Version:   tx.Version,
		Signature: signatureBytes,
		Options:   tx.Options,
	}

	serializedTx, err := json.Marshal(coreTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(serializedTx), nil
}
