package factories

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/multiversx/mx-sdk-go/coreV2"
)

var (
	contractAddress, _ = core.NewAddressFromBech32("erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u")
	Tconf              = &core.Config{
		ChainID: "T",

		GasLimitIssue:                   60_000_000,
		GasLimitToggleBurnRoleGlobally:  60_000_000,
		GasLimitESDTLocalMint:           300_000,
		GasLimitESDTLocalBurn:           300_000,
		GasLimitSetSpecialRole:          60_000_000,
		GasLimitPausing:                 60_000_000,
		GasLimitFreezing:                60_000_000,
		GasLimitWiping:                  60_000_000,
		GasLimitESDTNFTCreate:           3_000_000,
		GasLimitESDTNFTUpdateAttributes: 1_000_000,
		GasLimitESDTNFTAddQuantity:      1_000_000,
		GasLimitESDTNFTBurn:             1_000_000,
		GasLimitStorePerByte:            50_000,
		IssueCost:                       50_000_000_000_000_000,
		ESDTContractAddress:             contractAddress,

		MinGasLimit:                  50_000,
		GasLimitPerByte:              1_500,
		GasLimitESDTTransfer:         200_000,
		GasLimitESDTNFTTransfer:      200_000,
		GasLimitMultiESDTNFTTransfer: 200_000,

		GasLimitStake:                        5_000_000,
		GasLimitUnstake:                      5_000_000,
		GasLimitUnbound:                      5_000_000,
		GasLimitCreateDelegationContract:     50_000_000,
		GasLimitDelegationOperations:         1_000_000,
		AdditionalGasLimitPerValidatorNode:   6_000_000,
		AdditionalGasForDelegationOperations: 10_000_000,
	}
)

func TestTokenManagementTransactionFactory_CreateTransactionForIssuingFungible(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForIssuingFungible(
		frank,
		"FRANK",
		"FRANK",
		100,
		0,
		true,
		true,
		true,
		true,
		true,
		true,
	)
	require.NoError(t, err, "failed to issue fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60384500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(50000000000000000),
		Data:     []byte("issue@4652414e4b@4652414e4b@64@@63616e467265657a65@74727565@63616e57697065@74727565@63616e5061757365@74727565@63616e4368616e67654f776e6572@74727565@63616e55706772616465@74727565@63616e4164645370656369616c526f6c6573@74727565"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForIssuingSemiFungible(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForIssuingSemiFungible(
		frank,
		"FRANK",
		"FRANK",
		true,
		true,
		true,
		true,
		true,
		true,
		true,
	)
	require.NoError(t, err, "failed to issue semi-fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60483500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(50000000000000000),
		Data:     []byte("issueSemiFungible@4652414e4b@4652414e4b@63616e467265657a65@74727565@63616e57697065@74727565@63616e5061757365@74727565@63616e5472616e736665724e4654437265617465526f6c65@74727565@63616e4368616e67654f776e6572@74727565@63616e55706772616465@74727565@63616e4164645370656369616c526f6c6573@74727565"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForIssuingNonFungible(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForIssuingNonFungible(
		frank,
		"FRANK",
		"FRANK",
		true,
		true,
		true,
		true,
		true,
		true,
		true,
	)
	require.NoError(t, err, "failed to issue non-fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60482000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(50000000000000000),
		Data:     []byte("issueNonFungible@4652414e4b@4652414e4b@63616e467265657a65@74727565@63616e57697065@74727565@63616e5061757365@74727565@63616e5472616e736665724e4654437265617465526f6c65@74727565@63616e4368616e67654f776e6572@74727565@63616e55706772616465@74727565@63616e4164645370656369616c526f6c6573@74727565"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForRegisteringMetaESDT(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForRegisteringMetaESDT(
		frank,
		"FRANK",
		"FRANK",
		10,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
	)
	require.NoError(t, err, "failed to register meta ESDT")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60486500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(50000000000000000),
		Data:     []byte("registerMetaESDT@4652414e4b@4652414e4b@0a@63616e467265657a65@74727565@63616e57697065@74727565@63616e5061757365@74727565@63616e5472616e736665724e4654437265617465526f6c65@74727565@63616e4368616e67654f776e6572@74727565@63616e55706772616465@74727565@63616e4164645370656369616c526f6c6573@74727565"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForRegisteringAndSettingRoles(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForRegisteringAndSettingRoles(
		frank,
		"TEST",
		"TEST",
		FNG,
		2,
	)
	require.NoError(t, err, "failed to register and set roles")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60125000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(50000000000000000),
		Data:     []byte("registerAndSetAllRoles@54455354@54455354@464e47@02"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForSettingBurnRoleGlobally(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	burnRole, err := tmtf.CreateTransactionForSettingBurnRoleGlobally(
		frank,
		"FRANK-11ce3e",
	)
	require.NoError(t, err, "failed to set create transaction for setting burn role globally")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60116000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("setBurnRoleGlobally@4652414e4b2d313163653365"),
	}
	require.Equal(t, expectedTransaction, burnRole)
}

func TestTokenManagementTransactionFactory_CreateTransactionForSettingSpecialRoleOnFungibleToken(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	fungible, err := tmtf.CreateTransactionForSettingSpecialRoleOnFungibleToken(
		frank,
		grace,
		"FRANK-11ce3e",
		true,
		false,
	)
	require.NoError(t, err, "failed to set special role on fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60258500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("setSpecialRole@4652414e4b2d313163653365@1e8a8b6b49de5b7be10aaa158a5a6a4abb4b56cc08f524bb5e6cd5f211ad3e13@45534454526f6c654c6f63616c4d696e74"),
	}
	require.Equal(t, expectedTransaction, fungible)
}

func TestTokenManagementTransactionFactory_CreateTransactionForSettingSpecialRoleOnSemiFungibleToken(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	semiFungilbe, err := tmtf.CreateTransactionForSettingSpecialRoleOnSemiFungibleToken(
		frank,
		grace,
		"FRANK-11ce3e",
		true,
		true,
		true,
		true,
	)
	require.NoError(t, err, "failed to set special role on semi-fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60422000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("setSpecialRole@4652414e4b2d313163653365@1e8a8b6b49de5b7be10aaa158a5a6a4abb4b56cc08f524bb5e6cd5f211ad3e13@45534454526f6c654e4654437265617465@45534454526f6c654e46544275726e@45534454526f6c654e46544164645175616e74697479@455344545472616e73666572526f6c65"),
	}
	require.Equal(t, expectedTransaction, semiFungilbe)
}

func TestTokenManagementTransactionFactory_CreateTransactionForSettingSpecialRoleOnNonFungibleToken(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	nonFungilbe, err := tmtf.CreateTransactionForSettingSpecialRoleOnNonFungibleToken(
		frank,
		grace,
		"FRANK-11ce3e",
		true,
		false,
		true,
		true,
		false,
	)
	require.NoError(t, err, "failed to set special role on non-fungible")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60393500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("setSpecialRole@4652414e4b2d313163653365@1e8a8b6b49de5b7be10aaa158a5a6a4abb4b56cc08f524bb5e6cd5f211ad3e13@45534454526f6c654e4654437265617465@45534454526f6c654e465455706461746541747472696275746573@45534454526f6c654e4654416464555249"),
	}
	require.Equal(t, expectedTransaction, nonFungilbe)
}

func TestTokenManagementTransactionFactory_CreateTransactionForCreatingNFT(t *testing.T) {
	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	nonFungilbe, err := tmtf.CreateTransactionForCreatingNFT(
		grace,
		"FRANK-aa9e8d",
		1,
		"test",
		1000,
		"abba",
		[]byte("test"),
		[]string{"a", "b"},
	)
	require.NoError(t, err, "failed to set create transaction for creating NFT")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 4068500,
		Sender:   "erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede",
		Receiver: "erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTNFTCreate@4652414e4b2d616139653864@01@74657374@03e8@61626261@74657374@61@62"),
	}
	require.Equal(t, expectedTransaction, nonFungilbe)
}

func TestTokenManagementTransactionFactory_CreateTransactionForPausing(t *testing.T) {
	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	pausing, err := tmtf.CreateTransactionForPausing(
		grace,
		"FRANK-11ce3e",
	)
	require.NoError(t, err, "failed to set create transaction for pausing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60095000,
		Sender:   "erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("pause@4652414e4b2d313163653365"),
	}
	require.Equal(t, expectedTransaction, pausing)
}

func TestTokenManagementTransactionFactory_CreateTransactionForUnpausing(t *testing.T) {
	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	unPausing, err := tmtf.CreateTransactionForUnpausing(
		grace,
		"FRANK-11ce3e",
	)
	require.NoError(t, err, "failed to set create transaction for unpausing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60098000,
		Sender:   "erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("unPause@4652414e4b2d313163653365"),
	}
	require.Equal(t, expectedTransaction, unPausing)
}

func TestTokenManagementTransactionFactory_CreateTransactionForFreezing(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	freezing, err := tmtf.CreateTransactionForFreezing(
		frank,
		grace,
		"FRANK-11ce3e",
	)
	require.NoError(t, err, "failed to set create transaction for freezing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60194000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("freeze@4652414e4b2d313163653365@1e8a8b6b49de5b7be10aaa158a5a6a4abb4b56cc08f524bb5e6cd5f211ad3e13"),
	}
	require.Equal(t, expectedTransaction, freezing)
}

func TestTokenManagementTransactionFactory_CreateTransactionForUnfreezing(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	grace, err := core.NewAddressFromBech32("erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	unFreeze, err := tmtf.CreateTransactionForUnfreezing(
		frank,
		grace,
		"FRANK-11ce3e",
	)
	require.NoError(t, err, "failed to set create transaction for unfreezing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 60197000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("unFreeze@4652414e4b2d313163653365@1e8a8b6b49de5b7be10aaa158a5a6a4abb4b56cc08f524bb5e6cd5f211ad3e13"),
	}
	require.Equal(t, expectedTransaction, unFreeze)
}

func TestTokenManagementTransactionFactory_CreateTransactionForLocalMinting(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	localMinting, err := tmtf.CreateTransactionForLocalMinting(
		frank,
		"FRANK-11ce3e",
		10,
	)
	require.NoError(t, err, "failed to set create transaction for unfreezing")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 411500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTLocalMint@4652414e4b2d313163653365@0a"),
	}
	require.Equal(t, expectedTransaction, localMinting)
}

func TestTokenManagementTransactionFactory_CreateTransactionForLocalBurning(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	localBurning, err := tmtf.CreateTransactionForLocalBurning(
		frank,
		"FRANK-11ce3e",
		10,
	)
	require.NoError(t, err, "failed to set create transaction for local burning")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 411500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTLocalBurn@4652414e4b2d313163653365@0a"),
	}
	require.Equal(t, expectedTransaction, localBurning)
}

func TestTokenManagementTransactionFactory_CreateTransactionForUpdatingAttributes(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	updateAttributes, err := tmtf.CreateTransactionForUpdatingAttributes(
		frank,
		"FRANK-11ce3e",
		10,
		[]byte("test"),
	)
	require.NoError(t, err, "failed to set create transaction for updating attributes")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 1140000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTNFTUpdateAttributes@4652414e4b2d313163653365@0a@74657374"),
	}
	require.Equal(t, expectedTransaction, updateAttributes)
}

func TestTokenManagementTransactionFactory_CreateTransactionForAddingQuantity(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	addQuantity, err := tmtf.CreateTransactionForAddingQuantity(
		frank,
		"FRANK-11ce3e",
		10,
		10,
	)
	require.NoError(t, err, "failed to set create transaction for adding quantity")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 1123500,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTNFTAddQuantity@4652414e4b2d313163653365@0a@0a"),
	}
	require.Equal(t, expectedTransaction, addQuantity)
}

func TestTokenManagementTransactionFactory_CreateTransactionForBurningQuantity(t *testing.T) {
	frank, err := core.NewAddressFromBech32("erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv")
	require.NoError(t, err, "failed to create address from bech32")

	tmtf := NewTokenManagementTransactionFactory(Tconf)
	burnQuantity, err := tmtf.CreateTransactionForBurningQuantity(
		frank,
		"FRANK-11ce3e",
		10,
		10,
	)
	require.NoError(t, err, "failed to set create transaction for burning quantity")

	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 1113000,
		Sender:   "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Receiver: "erd1kdl46yctawygtwg2k462307dmz2v55c605737dp3zkxh04sct7asqylhyv",
		Value:    big.NewInt(0),
		Data:     []byte("ESDTNFTBurn@4652414e4b2d313163653365@0a@0a"),
	}
	require.Equal(t, expectedTransaction, burnQuantity)
}
