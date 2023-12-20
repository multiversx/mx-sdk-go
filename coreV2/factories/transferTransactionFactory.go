package factories

import (
	"fmt"
	"math/big"

	core "github.com/multiversx/mx-sdk-go/coreV2"
)

const (
	AdditionalGasForEsdtTransfer    = 100000
	AdditionalGasForEsdtNftTransfer = 800000
)

type transferTransactionFactory struct {
	Config          *core.Config
	TokenComputer   core.TokenComputer
	dataArgsBuilder core.TokenTransfersDataBuilder
}

type TransferTransactionFactory interface {
	CreateTransactionForNativeTokenTransfer(
		sender core.Address,
		receiver core.Address,
		nativeAmount int,
		data string) (*core.Transaction, error)

	CreateTransactionForESDTTokenTransfer(
		sender core.Address,
		receiver core.Address,
		tokenTransfers []*core.TokenTransfer,
	) (*core.Transaction, error)
}

func NewTransferTransactionFactory(config *core.Config, computer core.TokenComputer) TransferTransactionFactory {
	dataArgsBuilder := core.NewTokenTransferDataBuilder(computer)
	return &transferTransactionFactory{config, computer, dataArgsBuilder}
}

func (tf *transferTransactionFactory) CreateTransactionForNativeTokenTransfer(
	sender core.Address,
	receiver core.Address,
	nativeAmount int,
	data string,
) (*core.Transaction, error) {
	txBuilder := core.NewTransactionBuilder()

	return txBuilder.WithConfig(tf.Config).
		WithSender(sender).
		WithReceiver(receiver).
		WithDataParts([]string{data}).
		WithProvidedGasLimit(0).
		WithAddDataMovementGas(true).
		WithAmount(big.NewInt(int64(nativeAmount))).
		Build()
}

func (tf *transferTransactionFactory) CreateTransactionForESDTTokenTransfer(
	sender core.Address,
	receiver core.Address,
	tokenTransfers []*core.TokenTransfer,
) (*core.Transaction, error) {
	var err error
	dataParts := []string{""}
	extraGasForTransfer := 0

	switch len(tokenTransfers) {
	case 0:
		return nil, fmt.Errorf("no token transfer has been provided")

	case 1:
		transfer := tokenTransfers[0]

		if tf.TokenComputer.IsFungible(transfer.Token) {
			dataParts, _ = tf.dataArgsBuilder.BuildArgsForESDTTransfer(transfer)
			extraGasForTransfer = tf.Config.GasLimitESDTNFTTransfer + AdditionalGasForEsdtTransfer
		} else {
			dataParts, err = tf.dataArgsBuilder.BuildArgsForSingleESDTNFTTransfer(transfer, receiver)
			if err != nil {
				return nil, fmt.Errorf("failed to build args for single-ESDT NFT transfer: %v", err)
			}
			extraGasForTransfer = tf.Config.GasLimitMultiESDTNFTTransfer*len(tokenTransfers) +
				AdditionalGasForEsdtNftTransfer
		}

	default:
		dataParts, err = tf.dataArgsBuilder.BuildArgsForMultiESDTNFTTransfer(receiver, tokenTransfers)
		if err != nil {
			return nil, fmt.Errorf("failed to build args for multi-ESDT NFT transfer: %v", err)
		}
		extraGasForTransfer = tf.Config.GasLimitMultiESDTNFTTransfer*len(tokenTransfers) +
			AdditionalGasForEsdtNftTransfer
		receiver = sender
	}

	return core.NewTransactionBuilder().
		WithConfig(tf.Config).
		WithSender(sender).
		WithReceiver(receiver).
		WithAmount(big.NewInt(0)).
		WithDataParts(dataParts).
		WithProvidedGasLimit(uint64(extraGasForTransfer)).
		WithAddDataMovementGas(true).
		Build()
}
