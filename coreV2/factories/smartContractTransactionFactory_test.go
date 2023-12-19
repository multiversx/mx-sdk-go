package factories

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/multiversx/mx-sdk-go/coreV2"
)

var (
	Dconf = &core.Config{
		ChainID: "D",

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
	}
)

func TestSmartContractTransactionsFactory_CreateTransactionForDeploy(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := os.ReadFile("./testdata/adder.wasm")
	require.NoError(t, err, "failed to read wasm file")

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(0))}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForDeploy(
		sender,
		contract,
		args,
		big.NewInt(0),
		true,
		true,
		false,
		true,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6gq4hu",
		Value:    big.NewInt(0),
		Data:     []byte("0061736d010000000129086000006000017f60027f7f017f60027f7f0060017f0060037f7f7f017f60037f7f7f0060017f017f0290020b03656e7619626967496e74476574556e7369676e6564417267756d656e74000303656e760f6765744e756d417267756d656e7473000103656e760b7369676e616c4572726f72000303656e76126d42756666657253746f726167654c6f6164000203656e76176d427566666572546f426967496e74556e7369676e6564000203656e76196d42756666657246726f6d426967496e74556e7369676e6564000203656e76136d42756666657253746f7261676553746f7265000203656e760f6d4275666665725365744279746573000503656e760e636865636b4e6f5061796d656e74000003656e7614626967496e7446696e697368556e7369676e6564000403656e7609626967496e744164640006030b0a010104070301000000000503010003060f027f0041a080080b7f0041a080080b074607066d656d6f7279020004696e697400110667657453756d00120361646400130863616c6c4261636b00140a5f5f646174615f656e6403000b5f5f686561705f6261736503010aca010a0e01017f4100100c2200100020000b1901017f419c8008419c800828020041016b220036020020000b1400100120004604400f0b4180800841191002000b16002000100c220010031a2000100c220010041a20000b1401017f100c2202200110051a2000200210061a0b1301017f100c220041998008410310071a20000b1401017f10084101100d100b210010102000100f0b0e0010084100100d1010100e10090b2201037f10084101100d100b210110102202100e220020002001100a20022000100f0b0300010b0b2f0200418080080b1c77726f6e67206e756d626572206f6620617267756d656e747373756d00419c80080b049cffffff@0500@0504@"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecuteNoTransfer(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4")
	require.NoError(t, err, "failed to create contract address")
	function := "add"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		nil,
		nil,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4",
		Value:    big.NewInt(0),
		Data:     []byte("add@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecuteAndTransferNativeToken(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4")
	require.NoError(t, err, "failed to create contract address")
	function := "add"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}
	egldAmount := big.NewInt(1000000000000000000)

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		egldAmount,
		nil,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4",
		Value:    big.NewInt(1000000000000000000),
		Data:     []byte("add@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecuteAndSendSingleESDT(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4")
	require.NoError(t, err, "failed to create contract address")
	function := "dummy"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}
	token := core.Token{Identifier: "FOO-6ce17b"}
	transfer := &core.TokenTransfer{Token: token, Amount: big.NewInt(10)}
	transfers := []*core.TokenTransfer{transfer}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		nil,
		transfers,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Value:    big.NewInt(0),
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4",
		Data:     []byte("ESDTTransfer@464f4f2d366365313762@0a@64756d6d79@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestCreateTransactionForExecuteAndSendMultipleESDTs(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqak8zt22wl2ph4tswtyc39namqx6ysa2sd8ss4xmlj3")
	require.NoError(t, err, "failed to create contract address")
	function := "dummy"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}
	fooToken := core.Token{Identifier: "FOO-6ce17b"}
	fooTransfer := &core.TokenTransfer{Token: fooToken, Amount: big.NewInt(10)}
	barToken := core.Token{Identifier: "BAR-5bc08f"}
	barTransfer := &core.TokenTransfer{Token: barToken, Amount: big.NewInt(3140)}
	transfers := []*core.TokenTransfer{fooTransfer, barTransfer}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		nil,
		transfers,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Value:    big.NewInt(0),
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Data:     []byte("MultiESDTNFTTransfer@00000000000000000500ed8e25a94efa837aae0e593112cfbb01b448755069e1@02@464f4f2d366365313762@@0a@4241522d356263303866@@0c44@64756d6d79@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecuteAndSendSingleNFT(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqak8zt22wl2ph4tswtyc39namqx6ysa2sd8ss4xmlj3")
	require.NoError(t, err, "failed to create contract address")
	function := "dummy"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}
	token := core.Token{Identifier: "NFT-123456", Nonce: 1}
	transfer := &core.TokenTransfer{Token: token, Amount: big.NewInt(1)}
	transfers := []*core.TokenTransfer{transfer}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		nil,
		transfers,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Value:    big.NewInt(0),
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Data:     []byte("ESDTNFTTransfer@4e46542d313233343536@01@01@00000000000000000500b9353fe8407f87310c87e12fa1ac807f0485da39d152@64756d6d79@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecuteAndSendMultipleNFTS(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := core.NewAddressFromBech32("erd1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fqgtz9l4")
	require.NoError(t, err, "failed to create contract address")
	function := "dummy"

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(7))}
	firstToken := core.Token{Identifier: "NFT-123456", Nonce: 1}
	firstTokenTransfer := &core.TokenTransfer{Token: firstToken, Amount: big.NewInt(1)}
	secondToken := core.Token{Identifier: "NFT-123456", Nonce: 42}
	secondTokenTransfer := &core.TokenTransfer{Token: secondToken, Amount: big.NewInt(1)}
	transfers := []*core.TokenTransfer{firstTokenTransfer, secondTokenTransfer}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForExecute(
		sender,
		contract,
		function,
		args,
		nil,
		transfers,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Value:    big.NewInt(0),
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Data:     []byte("MultiESDTNFTTransfer@00000000000000000500b9353fe8407f87310c87e12fa1ac807f0485da39d152@02@4e46542d313233343536@01@01@4e46542d313233343536@2a@01@64756d6d79@07"),
	}
	require.Equal(t, expectedTransaction, transaction)
}

func TestSmartContractTransactionsFactory_CreateTransactionForExecute(t *testing.T) {
	sender, err := core.NewAddressFromBech32("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.NoError(t, err, "failed to create sender address")

	contract, err := os.ReadFile("./testdata/adder.wasm")
	require.NoError(t, err, "failed to read wasm file")

	gasLimit := uint32(6000000)
	args := []any{core.EncodeUnsignedNumber(uint64(0))}

	factory := NewSmartContractTransactionsFactory(Tconf, core.NewTokenComputer())
	transaction, err := factory.CreateTransactionForUpgrade(
		sender,
		contractAddress,
		contract,
		args,
		big.NewInt(0),
		true,
		true,
		false,
		true,
		gasLimit,
	)
	require.NoError(t, err, "failed to deploy contract")
	expectedTransaction := &core.Transaction{
		ChainID:  "T",
		GasLimit: 6000000,
		Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		Receiver: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u",
		Value:    big.NewInt(0),
		Data:     []byte("upgradeContract@0061736d010000000129086000006000017f60027f7f017f60027f7f0060017f0060037f7f7f017f60037f7f7f0060017f017f0290020b03656e7619626967496e74476574556e7369676e6564417267756d656e74000303656e760f6765744e756d417267756d656e7473000103656e760b7369676e616c4572726f72000303656e76126d42756666657253746f726167654c6f6164000203656e76176d427566666572546f426967496e74556e7369676e6564000203656e76196d42756666657246726f6d426967496e74556e7369676e6564000203656e76136d42756666657253746f7261676553746f7265000203656e760f6d4275666665725365744279746573000503656e760e636865636b4e6f5061796d656e74000003656e7614626967496e7446696e697368556e7369676e6564000403656e7609626967496e744164640006030b0a010104070301000000000503010003060f027f0041a080080b7f0041a080080b074607066d656d6f7279020004696e697400110667657453756d00120361646400130863616c6c4261636b00140a5f5f646174615f656e6403000b5f5f686561705f6261736503010aca010a0e01017f4100100c2200100020000b1901017f419c8008419c800828020041016b220036020020000b1400100120004604400f0b4180800841191002000b16002000100c220010031a2000100c220010041a20000b1401017f100c2202200110051a2000200210061a0b1301017f100c220041998008410310071a20000b1401017f10084101100d100b210010102000100f0b0e0010084100100d1010100e10090b2201037f10084101100d100b210110102202100e220020002001100a20022000100f0b0300010b0b2f0200418080080b1c77726f6e67206e756d626572206f6620617267756d656e747373756d00419c80080b049cffffff@0504@"),
	}
	require.Equal(t, expectedTransaction, transaction)
}
