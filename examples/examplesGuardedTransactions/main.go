package main

import (
	"context"
	"errors"
	"os"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
	"github.com/urfave/cli"
)

var (
	helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
	setGuardian = cli.BoolFlag{
		Name:        "setGuardian",
		Usage:       "Should be set in order to construct data field for setGuardian. Could fail in combination with some other flags, e.g. with dataField flag",
		Destination: &argsConfig.setGuardian,
	}

	guardedTxBy = cli.StringFlag{
		Name:        "guardedTxBy",
		Usage:       "If used, will set the guardian to the given one. Options: alice, bob, eve, charlie",
		Destination: &argsConfig.guardedTxBy,
	}

	guardian = cli.StringFlag{
		Name:        "guardian",
		Usage:       "If used, it replaces the default guardian with the given one. Options: alice, bob, eve, charlie, erd1..... Could fail with custom address and guardedTx set",
		Destination: &argsConfig.guardian,
	}

	sender = cli.StringFlag{
		Name:        "sender",
		Usage:       "If used, it replaces the default sender with the given one. Options: alice, bob, eve, charlie",
		Destination: &argsConfig.sender,
	}

	receiver = cli.StringFlag{
		Name:        "receiver",
		Usage:       "If used, it replaces the default receiver with the given one. Options: alice, bob, eve, charlie, erd1...",
		Destination: &argsConfig.receiver,
	}

	dataField = cli.StringFlag{
		Name:        "dataField",
		Usage:       "If used, it replaces the data field with this one. Could fail in combination with some other flags, e.g. setGuardian flag set",
		Destination: &argsConfig.dataField,
	}

	value = cli.StringFlag{
		Name:        "value",
		Usage:       "If set it replaces the default transaction value with this value. If might fail in combination with other flags, e.g setGuardian (which requires zero value)",
		Destination: &argsConfig.value,
	}

	withFunding = cli.BoolFlag{
		Name:        "withFunding",
		Usage:       "If set the default accounts will be funded with 10 egld",
		Destination: &argsConfig.withFunding,
	}

	argsConfig = &cfg{}
	log        = logger.GetOrCreate("elrond-sdk-erdgo/examples/examplesGuardedTransaction")
)

var HOME = os.Getenv("HOME")
var pathGeneratedWallets = HOME + "/Elrond/testnet/filegen/output/walletKey.pem"

const (
	alice              = "alice"
	bob                = "bob"
	charlie            = "charlie"
	eve                = "eve"
	setGuardianGasCost = 250000
	maskGuardedTx      = 1 << 1
)

type cfg struct {
	setGuardian bool
	withFunding bool
	guardedTxBy string
	guardian    string
	sender      string
	receiver    string
	dataField   string
	value       string
}

type selectedOptions struct {
	skSender        []byte
	skGuardian      []byte
	guardianAddress core.AddressHandler
	txArguments     data.ArgCreateTransaction
}

type testData struct {
	skFunding      []byte
	skAlice        []byte
	skBob          []byte
	skCharlie      []byte
	skEve          []byte
	addressFunding core.AddressHandler
	addressAlice   core.AddressHandler
	addressBob     core.AddressHandler
	addressCharlie core.AddressHandler
	addressEve     core.AddressHandler
	selected       selectedOptions
}

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	app := cli.NewApp()
	cli.AppHelpTemplate = helpTemplate

	app.Name = "guarded transactions cli"
	app.Version = "v1.0.0"
	app.Usage = "This binary enables sending and receiving transactions to the configured net"

	app.Flags = []cli.Flag{
		setGuardian,
		guardedTxBy,
		guardian,
		sender,
		receiver,
		value,
		dataField,
		withFunding,
	}

	app.Action = func(_ *cli.Context) error {
		return process()
	}

	_ = app.Run(os.Args)
}

func process() error {
	td, err := loadPemFiles()
	if err != nil {
		return err
	}

	ep, err := blockchain.NewElrondProxy(createElrondProxyArgs())
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return err
	}

	// netConfigs can be used multiple times (for example when sending multiple transactions) as to improve the
	// responsiveness of the system
	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		log.Error("unable to get the network configs", "error", err)
		return err
	}

	options, err := processCommand(td, ep, netConfigs)
	if err != nil {
		return err
	}

	if argsConfig.withFunding {
		err = fundWallets(td, ep, netConfigs)
		if err != nil {
			return err
		}
	}

	return generateAndSendTransaction(options, ep)
}

func processCommand(td *testData, proxy workflows.ProxyHandler, config *data.NetworkConfig) (*selectedOptions, error) {
	options, err := getDefaultOptions(td, proxy, config)
	if err != nil {
		return nil, err
	}

	err = setSenderOption(td, options)
	if err != nil {
		return nil, err
	}

	err = setReceiverOption(td, options)
	if err != nil {
		return nil, err
	}

	err = setGuardianOption(td, options)
	if err != nil {
		return nil, err
	}

	err = setGuardedTxByOption(td, options)
	if err != nil {
		return nil, err
	}

	err = setOtherOptions(options, config)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func setSenderOption(td *testData, options *selectedOptions) error {
	selectedAddress, sk, err := selectAddressAndSkFromString(td, argsConfig.sender)
	if err != nil {
		return err
	}

	if selectedAddress != nil {
		options.txArguments.SndAddr = selectedAddress.AddressAsBech32String()
		options.skSender = sk
	}

	return nil
}

func setReceiverOption(td *testData, options *selectedOptions) error {
	selectedAddress, _, err := selectAddressAndSkFromString(td, argsConfig.receiver)
	if err != nil {
		return err
	}
	if selectedAddress != nil {
		options.txArguments.RcvAddr = selectedAddress.AddressAsBech32String()
	}

	return nil
}

func setGuardianOption(td *testData, options *selectedOptions) error {
	var err error
	selectedAddress, _, err := selectAddressAndSkFromString(td, argsConfig.guardian)
	if err != nil {
		return err
	}

	if selectedAddress != nil {
		options.guardianAddress = selectedAddress
	}
	return err
}

func setGuardedTxByOption(td *testData, options *selectedOptions) error {
	selectedAddress, sk, err := selectAddressAndSkFromString(td, argsConfig.guardedTxBy)
	if err != nil {
		return err
	}

	if selectedAddress != nil {
		options.txArguments.GuardianAddr = selectedAddress.AddressAsBech32String()
		options.skGuardian = sk
		options.txArguments.Options = maskGuardedTx
	}

	return nil
}

func selectAddressAndSkFromString(td *testData, option string) (core.AddressHandler, []byte, error) {
	var selectedAddress core.AddressHandler
	var sk []byte

	switch option {
	case alice:
		selectedAddress = td.addressAlice
		sk = td.skAlice
	case bob:
		selectedAddress = td.addressBob
		sk = td.skBob
	case charlie:
		selectedAddress = td.addressCharlie
		sk = td.skCharlie
	case eve:
		selectedAddress = td.addressEve
		sk = td.skEve
	default:
		var err error
		if len(option) > 0 {
			selectedAddress, err = data.NewAddressFromBech32String(option)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return selectedAddress, sk, nil
}

func setOtherOptions(options *selectedOptions, config *data.NetworkConfig) error {
	if argsConfig.setGuardian && len(argsConfig.dataField) > 0 {
		return errors.New("dataField and setGuardian cannot be set together")
	}

	err := treatDataIfNeeded(options, config)
	if err != nil {
		return err
	}

	if len(argsConfig.value) > 0 {
		options.txArguments.Value = argsConfig.value
	}

	return nil
}

func treatDataIfNeeded(options *selectedOptions, config *data.NetworkConfig) error {
	var err error
	if argsConfig.setGuardian {
		options.txArguments.Data, err = createSetGuardianData(options.guardianAddress)
		if err != nil {
			return err
		}
		options.txArguments.GasLimit += setGuardianGasCost
		options.txArguments.RcvAddr = options.txArguments.SndAddr
	}
	if len(argsConfig.dataField) > 0 {
		options.txArguments.Data = []byte(argsConfig.dataField)
	}
	options.txArguments.Version = 2
	options.txArguments.GasLimit += uint64(len(options.txArguments.Data)) * config.GasPerDataByte

	return nil
}

func createSetGuardianData(guardianAddress core.AddressHandler) ([]byte, error) {
	builder := builders.NewTxDataBuilder()
	builder.Function("SetGuardian").ArgAddress(guardianAddress)
	return builder.ToDataBytes()
}

func getDefaultOptions(td *testData, ep workflows.ProxyHandler, netConfigs *data.NetworkConfig) (*selectedOptions, error) {
	transactionArguments, err := ep.GetDefaultTransactionArguments(context.Background(), td.addressAlice, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return nil, err
	}
	transactionArguments.Value = "0"
	transactionArguments.RcvAddr = td.addressBob.AddressAsBech32String()

	return &selectedOptions{
		skSender:        td.skAlice, // default if nothing provided
		skGuardian:      td.skBob,   // default if nothing provided
		txArguments:     transactionArguments,
		guardianAddress: td.addressBob,
	}, nil
}

func generateAndSendTransaction(options *selectedOptions, proxy interactors.Proxy) error {
	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return err
	}
	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return err
	}

	tx, err := ti.ApplyUserSignatureAndGenerateTx(options.skSender, options.txArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return err
	}

	if len(argsConfig.guardedTxBy) > 0 {
		err = ti.ApplyGuardianSignature(options.skGuardian, tx)
		if err != nil {
			log.Error("error applying guardian signature", "error", err)
			return err
		}
	}

	ti.AddTransaction(tx)

	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
		return err
	}

	log.Info("transactions sent", "hashes", hashes)
	return nil
}

func fundWallets(td *testData, proxy workflows.ProxyHandler, netConfigs *data.NetworkConfig) error {
	transactionArguments, err := proxy.GetDefaultTransactionArguments(context.Background(), td.addressFunding, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return err
	}

	receivers := []core.AddressHandler{td.addressAlice, td.addressBob, td.addressEve, td.addressCharlie}

	err = sendFundWalletsTxs(td, proxy, transactionArguments, receivers)
	if err != nil {
		return err
	}

	return nil
}

func sendFundWalletsTxs(td *testData, proxy workflows.ProxyHandler, txArgs data.ArgCreateTransaction, receivers []core.AddressHandler) error {
	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return err
	}

	tiProxy, ok := proxy.(interactors.Proxy)
	if !ok {
		err = errors.New("proxy assertion failure")
		log.Error("type assertion failed", "error", err)
		return err
	}

	ti, err := interactors.NewTransactionInteractor(tiProxy, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return err
	}

	var tx *data.Transaction
	txArgs.Value = "10000000000000000000" // 10EGLD
	for _, addressHandler := range receivers {
		txArgs.RcvAddr = addressHandler.AddressAsBech32String()
		tx, err = ti.ApplyUserSignatureAndGenerateTx(td.skFunding, txArgs)
		if err != nil {
			log.Error("error creating transaction", "error", err)
			return err
		}
		ti.AddTransaction(tx)
		txArgs.Nonce++
	}

	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transactions", "error", err)
		return err
	}

	log.Info("funding transactions sent", "hashes", hashes)
	return nil
}

func loadPemFiles() (*testData, error) {
	var err error
	td := &testData{}

	w := interactors.NewWallet()

	td.skFunding, err = w.LoadPrivateKeyFromPemFile(pathGeneratedWallets)
	if err != nil {
		log.Error("unable to load funding wallet", "error", err)
		return nil, err
	}

	td.skAlice, err = w.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
		return nil, err
	}

	td.skBob, err = w.LoadPrivateKeyFromPemData([]byte(examples.BobPemContents))
	if err != nil {
		log.Error("unable to load bob.pem", "error", err)
		return nil, err
	}

	td.skCharlie, err = w.LoadPrivateKeyFromPemData([]byte(examples.CharliePemContents))
	if err != nil {
		log.Error("unable to load charlie.pem", "error", err)
		return nil, err
	}

	td.skEve, err = w.LoadPrivateKeyFromPemData([]byte(examples.EvePemContents))
	if err != nil {
		log.Error("unable to load eve.pem", "error", err)
		return nil, err
	}

	td.addressFunding, err = w.GetAddressFromPrivateKey(td.skFunding)
	if err != nil {
		log.Error("unable to load funding address from private key", "error", err)
		return nil, err
	}

	// Generate address from private key
	td.addressAlice, err = w.GetAddressFromPrivateKey(td.skAlice)
	if err != nil {
		log.Error("unable to load alice address from the private key", "error", err)
		return nil, err
	}

	// Generate address from private key
	td.addressBob, err = w.GetAddressFromPrivateKey(td.skBob)
	if err != nil {
		log.Error("unable to load bob address from the private key", "error", err)
		return nil, err
	}

	// Generate address from private key
	td.addressCharlie, err = w.GetAddressFromPrivateKey(td.skCharlie)
	if err != nil {
		log.Error("unable to load charlie address from the private key", "error", err)
		return nil, err
	}

	// Generate address from private key
	td.addressEve, err = w.GetAddressFromPrivateKey(td.skEve)
	if err != nil {
		log.Error("unable to load eve address from the private key", "error", err)
		return nil, err
	}

	return td, nil
}

func createElrondProxyArgs() blockchain.ArgsElrondProxy {
	return blockchain.ArgsElrondProxy{
		ProxyURL:            examples.LocalTestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
}
