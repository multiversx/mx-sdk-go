package factories

import (
	"encoding/hex"
	"fmt"
	"math/big"

	core "github.com/multiversx/mx-sdk-go/coreV2"
)

var (
	contractDeployAddress, _ = core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6gq4hu")
)

type SmartContractTransactionsFactory interface {
	CreateTransactionForDeploy(
		sender core.Address,
		bytecode []byte,
		arguments []any,
		nativeTransferAmount *core.Amount,
		isUpgradeable bool,
		isReadable bool,
		isPayable bool,
		isPayableBySc bool,
		gasLimit uint32,
	) (*core.Transaction, error)

	CreateTransactionForExecute(
		sender core.Address,
		contract core.Address,
		function string,
		arguments []any,
		nativeTransferAmount *core.Amount,
		tokenTransfers []*core.TokenTransfer,
		gasLimit uint32,
	) (*core.Transaction, error)

	CreateTransactionForUpgrade(
		sender core.Address,
		contract core.Address,
		bytecode []byte,
		arguments []any,
		nativeTransferAmount *core.Amount,
		isUpgradeable bool,
		isReadable bool,
		isPayable bool,
		isPayableBySc bool,
		gasLimit uint32,
	) (*core.Transaction, error)
}

type smartContractTransactionsFactory struct {
	Config          *core.Config
	TokenComputer   core.TokenComputer
	dataArgsBuilder core.TokenTransfersDataBuilder
}

func NewSmartContractTransactionsFactory(
	config *core.Config,
	tokenComputer core.TokenComputer,
) SmartContractTransactionsFactory {
	return &smartContractTransactionsFactory{
		Config:          config,
		TokenComputer:   tokenComputer,
		dataArgsBuilder: core.NewTokenTransferDataBuilder(tokenComputer),
	}
}

func (s *smartContractTransactionsFactory) CreateTransactionForDeploy(
	sender core.Address,
	bytecode []byte,
	arguments []any,
	nativeTransferAmount *core.Amount,
	isUpgradeable bool,
	isReadable bool,
	isPayable bool,
	isPayableBySc bool,
	gasLimit uint32,
) (*core.Transaction, error) {
	if nativeTransferAmount == nil {
		nativeTransferAmount = big.NewInt(0)
	}

	metadata := core.NewCodeMetadata(
		core.WithUpgradeable(isUpgradeable),
		core.WithReadable(isReadable),
		core.WithPayable(isPayable),
		core.WithPayableByContract(isPayableBySc),
	)

	parts := []string{
		hex.EncodeToString(bytecode),
		hex.EncodeToString(core.GetVMTypeWASMVM()),
		metadata.String(),
	}

	for _, arg := range arguments {
		if _, ok := arg.([]byte); !ok {
			return nil, fmt.Errorf("failed to parse argument list element to byte array - []byte")
		}
		parts = append(parts, hex.EncodeToString(arg.([]byte)))
	}

	return core.NewTransactionBuilder().
		WithConfig(s.Config).
		WithSender(sender).
		WithReceiver(contractDeployAddress).
		WithAmount(nativeTransferAmount).
		WithProvidedGasLimit(uint64(gasLimit)).
		WithAddDataMovementGas(false).
		WithDataParts(parts).
		Build()
}

func (s *smartContractTransactionsFactory) CreateTransactionForExecute(
	sender core.Address,
	contract core.Address,
	function string,
	arguments []any,
	nativeTransferAmount *core.Amount,
	tokenTransfers []*core.TokenTransfer,
	gasLimit uint32,
) (*core.Transaction, error) {
	var (
		err       error
		dataParts []string
	)
	noOfTokens := len(tokenTransfers)
	receiver := contract

	if nativeTransferAmount == nil {
		nativeTransferAmount = big.NewInt(0)
	}

	if nativeTransferAmount.Cmp(big.NewInt(0)) > 0 && noOfTokens > 0 {
		return nil, fmt.Errorf("can't send both native token and custom tokens(ESDT/NFT)")
	}

	if noOfTokens == 1 {
		transfer := tokenTransfers[0]

		if s.TokenComputer.IsFungible(transfer.Token) {
			dataParts, err = s.dataArgsBuilder.BuildArgsForESDTTransfer(transfer)
			if err != nil {
				return nil, fmt.Errorf("failed to build args for esdt transfer")
			}
		} else {
			dataParts, err = s.dataArgsBuilder.BuildArgsForSingleESDTNFTTransfer(transfer, receiver)
			if err != nil {
				return nil, fmt.Errorf("failed to build args for esdt-nft transfer")
			}
			receiver = sender
		}
	}
	if noOfTokens > 1 {
		dataParts, err = s.dataArgsBuilder.BuildArgsForMultiESDTNFTTransfer(receiver, tokenTransfers)
		if err != nil {
			return nil, fmt.Errorf("failed to build args for multi-esdt-nft transfer")
		}
		dataParts = append(dataParts)
		receiver = sender
	}

	if len(dataParts) == 0 {
		dataParts = append(dataParts, function)
	} else {
		dataParts = append(dataParts, hex.EncodeToString([]byte(function)))
	}

	args, err := encodeArguments(arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments")
	}
	dataParts = append(dataParts, args...)

	return core.NewTransactionBuilder().
		WithConfig(s.Config).
		WithSender(sender).
		WithReceiver(receiver).
		WithAmount(nativeTransferAmount).
		WithProvidedGasLimit(uint64(gasLimit)).
		WithAddDataMovementGas(false).
		WithDataParts(dataParts).
		Build()
}

func (s *smartContractTransactionsFactory) CreateTransactionForUpgrade(
	sender core.Address,
	contract core.Address,
	bytecode []byte,
	arguments []any,
	nativeTransferAmount *core.Amount,
	isUpgradeable bool,
	isReadable bool,
	isPayable bool,
	isPayableBySc bool,
	gasLimit uint32,
) (*core.Transaction, error) {
	if nativeTransferAmount == nil {
		nativeTransferAmount = big.NewInt(0)
	}

	metadata := core.NewCodeMetadata(
		core.WithUpgradeable(isUpgradeable),
		core.WithReadable(isReadable),
		core.WithPayable(isPayable),
		core.WithPayableByContract(isPayableBySc),
	)

	parts := []string{
		"upgradeContract",
		hex.EncodeToString(bytecode),
		metadata.String(),
	}

	args, err := encodeArguments(arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments")
	}
	parts = append(parts, args...)

	return core.NewTransactionBuilder().
		WithConfig(s.Config).
		WithSender(sender).
		WithReceiver(contract).
		WithAmount(nativeTransferAmount).
		WithProvidedGasLimit(uint64(gasLimit)).
		WithAddDataMovementGas(false).
		WithDataParts(parts).
		Build()
}

func encodeArguments(args []any) ([]string, error) {
	var parts []string
	for _, arg := range args {
		if _, ok := arg.([]byte); !ok {
			return nil, fmt.Errorf("failed to parse argument list element to byte array - []byte")
		}
		parts = append(parts, hex.EncodeToString(arg.([]byte)))
	}

	return parts, nil
}
