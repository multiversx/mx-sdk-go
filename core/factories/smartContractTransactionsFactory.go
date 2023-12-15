package factories

import "github.com/multiversx/mx-sdk-go/core"

type SmartContractTransactionsFactory interface {
	CreateTransactionForDeploy(
		sender core.Address,
		bytecode []byte,
		arguments []int,

	)
}

type smartContractTransactionsFactory struct {
	Config          *core.Config
	TokenComputer   core.TokenComputer
	dataArgsBuilder core.TokenTransfersDataBuilder
}
