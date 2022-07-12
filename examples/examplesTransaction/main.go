package main

import (
	"context"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/examples/examplesTransaction")

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	args := blockchain.ArgsElrondProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewElrondProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return
	}

	w := interactors.NewWallet()

	privateKey, err := w.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
		return
	}
	// Generate address from private key
	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Error("unable to load the address from the private key", "error", err)
		return
	}

	// netConfigs can be used multiple times (for example when sending multiple transactions) as to improve the
	// responsiveness of the system
	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		log.Error("unable to get the network configs", "error", err)
		return
	}

	transactionArguments, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	transactionArguments.RcvAddr = address.AddressAsBech32String() // send to self
	transactionArguments.Value = "1000000000000000000"             // 1EGLD

	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return
	}

	tx, err := ti.ApplyUserSignatureAndGenerateTx(privateKey, transactionArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return
	}
	ti.AddTransaction(tx)

	// a new transaction with the signature done on the hash of the transaction
	// it's ok to reuse the arguments here, they will be copied, anyway
	transactionArguments.Version = 2
	transactionArguments.Options = 1
	transactionArguments.Nonce++ // do not forget to increment the nonce, otherwise you will get 2 transactions
	// with the same nonce (only one of them will get executed)
	txSigOnHash, err := ti.ApplyUserSignatureAndGenerateTx(privateKey, transactionArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return
	}
	ti.AddTransaction(txSigOnHash)

	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
		return
	}

	log.Info("transactions sent", "hashes", hashes)
}
