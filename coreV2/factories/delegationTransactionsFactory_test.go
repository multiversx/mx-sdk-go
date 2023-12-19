package factories

import (
	"math/big"
	"testing"

	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/stretchr/testify/require"

	core "github.com/multiversx/mx-sdk-go/coreV2"
	"github.com/multiversx/mx-sdk-go/testsCommon"
)

func TestDelegationTransactionsFactory_CreateTransactionForNewDelegationContract(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	totalDelegationCap := new(core.Amount)
	tdc, _ := totalDelegationCap.SetString("5000000000000000000000", 10)
	amount := new(core.Amount)
	a, _ := amount.SetString("1250000000000000000000", 10)

	transaction, err := factory.CreateTransactionForNewDelegationContract(
		sender,
		tdc,
		10,
		a,
	)
	require.NoError(t, err, "failed to create transaction for delegating new contract")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60126500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllslmq6y6",
		Value:    a,
		Data:     []byte("createNewDelegationContract@010f0cf064dd59200000@0a"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

// TODO: implement this, once the wallet is merged.-
func TestDelegationTransactionsFactory_CreateTransactionForAddingNodes(t *testing.T) {

}

func TestDelegationTransactionsFactory_CreateTransactionForRemovingNodes(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	pks := []crypto.PublicKey{
		&testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("notavalidblskeyhexencoded"), nil
			},
		},
	}

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForRemovingNodes(
		sender,
		delegationContract,
		pks,
	)
	require.NoError(t, err, "failed to create transaction for removing nodes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 7143000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("removeNodes@6e6f746176616c6964626c736b6579686578656e636f646564"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForStakingNodes(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	pks := []crypto.PublicKey{
		&testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("notavalidblskeyhexencoded"), nil
			},
		},
	}

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForStakingNodes(
		sender,
		delegationContract,
		pks,
	)
	require.NoError(t, err, "failed to create transaction for staking nodes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12141500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("stakeNodes@6e6f746176616c6964626c736b6579686578656e636f646564"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForUnboundingNodes(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	pks := []crypto.PublicKey{
		&testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("notavalidblskeyhexencoded"), nil
			},
		},
	}

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForUnboundingNodes(
		sender,
		delegationContract,
		pks,
	)
	require.NoError(t, err, "failed to create transaction for unbounding nodes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12144500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("unBoundNodes@6e6f746176616c6964626c736b6579686578656e636f646564"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForUnstakingNodes(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	pks := []crypto.PublicKey{
		&testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("notavalidblskeyhexencoded"), nil
			},
		},
	}

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForUnstakingNodes(
		sender,
		delegationContract,
		pks,
	)
	require.NoError(t, err, "failed to create transaction for unstaking nodes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12144500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("unStakeNodes@6e6f746176616c6964626c736b6579686578656e636f646564"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForUnjailingNodes(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	pks := []crypto.PublicKey{
		&testsCommon.PublicKeyStub{
			ToByteArrayCalled: func() ([]byte, error) {
				return []byte("notavalidblskeyhexencoded"), nil
			},
		},
	}

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForUnjailingNodes(
		sender,
		delegationContract,
		pks,
	)
	require.NoError(t, err, "failed to create transaction for unjailing nodes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 7143000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("unJailNodes@6e6f746176616c6964626c736b6579686578656e636f646564"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForChangingServiceFee(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForChangingServiceFee(
		sender,
		delegationContract,
		10,
	)
	require.NoError(t, err, "failed to create transaction for changing service fee")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11078500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("changeServiceFee@0a"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForModifyingDelegationCap(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	amount, _ := new(big.Int).SetString("5000000000000000000000", 10)
	transaction, err := factory.CreateTransactionForModifyingDelegationCap(
		sender,
		delegationContract,
		amount,
	)
	require.NoError(t, err, "failed to create transaction for modifying delegation cap")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11117500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("modifyTotalDelegationCap@010f0cf064dd59200000"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForSettingAutomaticActivation(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForSettingAutomaticActivation(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for setting automatic activation")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11096500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("setAutomaticActivation@74727565"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForSettingCapCheckOnRedelegateRewards(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForSettingCapCheckOnRedelegateRewards(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for setting cap check")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11108500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("setCheckCapOnReDelegateRewards@74727565"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForUnsettingCapCheckOnRedelegateRewards(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForUnsettingCapCheckOnRedelegateRewards(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for unsetting cap check")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11111500,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("setCheckCapOnReDelegateRewards@66616c7365"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForSettingMetadata(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForSettingMetadata(
		sender,
		delegationContract,
		"name",
		"website",
		"identifier",
	)
	require.NoError(t, err, "failed to create transaction for setting metadata")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 11134000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("setMetaData@6e616d65@77656273697465@6964656e746966696572"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForDelegating(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	amount := new(big.Int)
	a, _ := amount.SetString("1000000000000000000", 10)
	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForDelegating(
		sender,
		delegationContract,
		a,
	)
	require.NoError(t, err, "failed to create transaction for delegating")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12000000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    a,
		Data:     []byte("delegate"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForClaimingRewards(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForClaimingRewards(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for claiming rewards")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("claimRewards"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForRedelegatingRewards(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForRedelegatingRewards(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for redelegating rewards")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12000000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("reDelegateRewards"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForUndelegating(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	amount := new(big.Int)
	a, _ := amount.SetString("1000000000000000000", 10)
	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForUndelegating(
		sender,
		delegationContract,
		a,
	)
	require.NoError(t, err, "failed to create transaction for undelegating")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12000000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("unDelegate@0de0b6b3a7640000"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestDelegationTransactionsFactory_CreateTransactionForWithdrawing(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2")
	require.NoError(t, err, "failed to create address from bech32")

	delegationContract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc")
	require.NoError(t, err, "failed to create address from bech32")

	factory := NewDelegationTransactionsFactory(Tconf)
	transaction, err := factory.CreateTransactionForWithdrawing(
		sender,
		delegationContract,
	)
	require.NoError(t, err, "failed to create transaction for withdrawing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 12000000,
		Sender:   "erd18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawq553rt2",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqtllllls002zgc",
		Value:    big.NewInt(0),
		Data:     []byte("withdraw"),
	}
	require.Equal(t, expectedTransaction, transaction)
}
