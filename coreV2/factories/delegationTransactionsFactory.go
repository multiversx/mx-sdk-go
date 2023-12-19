package factories

import (
	"encoding/hex"
	"fmt"

	crypto "github.com/multiversx/mx-chain-crypto-go"

	core "github.com/multiversx/mx-sdk-go/coreV2"
)

const (
	DelegationManagerScAddress = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllslmq6y6"
)

var delegationAddress, _ = core.NewAddressFromBech32(DelegationManagerScAddress)

type DelegationTransactionsFactory interface {
	CreateTransactionForNewDelegationContract(
		sender core.Address,
		totalDelegationCap *core.Amount,
		serviceFee int,
		amount *core.Amount,
	) (*core.Transaction, error)

	CreateTransactionForAddingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
		signedMessages [][]byte,
	) (*core.Transaction, error)

	CreateTransactionForRemovingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
	) (*core.Transaction, error)

	CreateTransactionForStakingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
	) (*core.Transaction, error)

	CreateTransactionForUnboundingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
	) (*core.Transaction, error)

	CreateTransactionForUnstakingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
	) (*core.Transaction, error)

	CreateTransactionForUnjailingNodes(
		sender core.Address,
		delegationContract core.Address,
		publicKeys []crypto.PublicKey,
	) (*core.Transaction, error)

	CreateTransactionForChangingServiceFee(
		sender core.Address,
		delegationContract core.Address,
		serviceFee int,
	) (*core.Transaction, error)

	CreateTransactionForModifyingDelegationCap(
		sender core.Address,
		delegationContract core.Address,
		delegationCap *core.Amount,
	) (*core.Transaction, error)

	CreateTransactionForSettingAutomaticActivation(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)

	CreateTransactionForSettingCapCheckOnRedelegateRewards(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)

	CreateTransactionForUnsettingCapCheckOnRedelegateRewards(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)

	CreateTransactionForSettingMetadata(
		sender core.Address,
		delegationContract core.Address,
		name string,
		website string,
		identifier string,
	) (*core.Transaction, error)

	CreateTransactionForDelegating(
		sender core.Address,
		delegationContract core.Address,
		amount *core.Amount,
	) (*core.Transaction, error)

	CreateTransactionForClaimingRewards(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)

	CreateTransactionForRedelegatingRewards(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)

	CreateTransactionForUndelegating(
		sender core.Address,
		delegationContract core.Address,
		amount *core.Amount,
	) (*core.Transaction, error)

	CreateTransactionForWithdrawing(
		sender core.Address,
		delegationContract core.Address,
	) (*core.Transaction, error)
}

type delegationTransactionsFactory struct {
	config *core.Config
}

func NewDelegationTransactionsFactory(config *core.Config) DelegationTransactionsFactory {
	return &delegationTransactionsFactory{config}
}

func (d *delegationTransactionsFactory) CreateTransactionForNewDelegationContract(
	sender core.Address,
	totalDelegationCap *core.Amount,
	serviceFee int,
	amount *core.Amount,
) (*core.Transaction, error) {
	parts := []string{
		"createNewDelegationContract",
		hex.EncodeToString(totalDelegationCap.Bytes()),
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(serviceFee))),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationAddress).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitCreateDelegationContract +
			d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		WithAmount(amount).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForAddingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
	signedMessages [][]byte,
) (*core.Transaction, error) {
	if len(publicKeys) != len(signedMessages) {
		return nil, fmt.Errorf("the number of public keys %q should match the number of signed messages %q",
			len(publicKeys), len(signedMessages))
	}

	parts := []string{"addNodes"}

	for i, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
		parts = append(parts, hex.EncodeToString(signedMessages[i]))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.computeExecutionGasLimitForNodesManagement(numNodes))).
		WithAddDataMovementGas(true).
		Build()

}

func (d *delegationTransactionsFactory) CreateTransactionForRemovingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
) (*core.Transaction, error) {
	parts := []string{"removeNodes"}

	for _, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.computeExecutionGasLimitForNodesManagement(numNodes))).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForStakingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
) (*core.Transaction, error) {
	parts := []string{"stakeNodes"}

	for _, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.GasLimitStake + numNodes*d.config.AdditionalGasLimitPerValidatorNode)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForUnboundingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
) (*core.Transaction, error) {
	parts := []string{"unBoundNodes"}

	for _, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.GasLimitUnbound + numNodes*d.config.AdditionalGasLimitPerValidatorNode)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForUnstakingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
) (*core.Transaction, error) {
	parts := []string{"unStakeNodes"}

	for _, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.GasLimitUnstake + numNodes*d.config.AdditionalGasLimitPerValidatorNode)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForUnjailingNodes(
	sender core.Address,
	delegationContract core.Address,
	publicKeys []crypto.PublicKey,
) (*core.Transaction, error) {
	parts := []string{"unJailNodes"}

	for _, pk := range publicKeys {
		bytes, err := pk.ToByteArray()
		if err != nil {
			return nil, fmt.Errorf("failed to get bytes from public key: %v", err)
		}
		parts = append(parts, hex.EncodeToString(bytes))
	}

	numNodes := len(publicKeys)

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.computeExecutionGasLimitForNodesManagement(numNodes))).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForChangingServiceFee(
	sender core.Address,
	delegationContract core.Address,
	serviceFee int,
) (*core.Transaction, error) {
	parts := []string{
		"changeServiceFee",
		hex.EncodeToString(core.EncodeUnsignedNumber(uint64(serviceFee))),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForModifyingDelegationCap(
	sender core.Address,
	delegationContract core.Address,
	delegationCap *core.Amount,
) (*core.Transaction, error) {
	parts := []string{
		"modifyTotalDelegationCap",
		hex.EncodeToString(delegationCap.Bytes()),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForSettingAutomaticActivation(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	parts := []string{
		"setAutomaticActivation",
		hex.EncodeToString([]byte("true")),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForSettingCapCheckOnRedelegateRewards(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	parts := []string{
		"setCheckCapOnReDelegateRewards",
		hex.EncodeToString([]byte("true")),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForUnsettingCapCheckOnRedelegateRewards(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	parts := []string{
		"setCheckCapOnReDelegateRewards",
		hex.EncodeToString([]byte("false")),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForSettingMetadata(
	sender core.Address,
	delegationContract core.Address,
	name string,
	website string,
	identifier string,
) (*core.Transaction, error) {
	parts := []string{
		"setMetaData",
		hex.EncodeToString([]byte(name)),
		hex.EncodeToString([]byte(website)),
		hex.EncodeToString([]byte(identifier)),
	}

	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts(parts).
		WithProvidedGasLimit(uint64(d.config.GasLimitDelegationOperations + d.config.AdditionalGasForDelegationOperations)).
		WithAddDataMovementGas(true).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForDelegating(
	sender core.Address,
	delegationContract core.Address,
	amount *core.Amount,
) (*core.Transaction, error) {
	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts([]string{"delegate"}).
		WithAmount(amount).
		WithProvidedGasLimit(12000000).
		WithAddDataMovementGas(false).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForClaimingRewards(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts([]string{"claimRewards"}).
		WithProvidedGasLimit(6000000).
		WithAddDataMovementGas(false).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForRedelegatingRewards(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts([]string{"reDelegateRewards"}).
		WithProvidedGasLimit(12000000).
		WithAddDataMovementGas(false).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForUndelegating(
	sender core.Address,
	delegationContract core.Address,
	amount *core.Amount,
) (*core.Transaction, error) {
	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts([]string{"unDelegate", hex.EncodeToString(amount.Bytes())}).
		WithProvidedGasLimit(12000000).
		WithAddDataMovementGas(false).
		Build()
}

func (d *delegationTransactionsFactory) CreateTransactionForWithdrawing(
	sender core.Address,
	delegationContract core.Address,
) (*core.Transaction, error) {
	return core.NewTransactionBuilder().
		WithConfig(d.config).
		WithSender(sender).
		WithReceiver(delegationContract).
		WithDataParts([]string{"withdraw"}).
		WithProvidedGasLimit(12000000).
		WithAddDataMovementGas(false).
		Build()
}

func (d *delegationTransactionsFactory) computeExecutionGasLimitForNodesManagement(numNodes int) int {
	return d.config.GasLimitDelegationOperations + numNodes*d.config.AdditionalGasLimitPerValidatorNode
}
