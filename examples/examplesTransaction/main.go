package main

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
)

var log = logger.GetOrCreate("examples/examplesTransaction")

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	ep := blockchain.NewElrondProxy(examples.TestnetGateway, nil)

	privateKey, err := erdgo.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
		return
	}
	// Generate address from private key
	addressString, err := erdgo.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Error("unable to load the address from the private key", "error", err)
		return
	}

	address, err := data.NewAddressFromBech32String(addressString)
	if err != nil {
		log.Error("error converting the address string", "error", err)
		return
	}

	//netConfigs can be used multiple times (eg. when sending multiple transactions) as to improve the
	//responsiveness of the system
	netConfigs, err := ep.GetNetworkConfig()
	if err != nil {
		log.Error("unable to get the network configs", "error", err)
		return
	}

	transactionArguments, err := ep.GetDefaultTransactionArguments(address, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	transactionArguments.RcvAddr = addressString       //send to self
	transactionArguments.Value = "1000000000000000000" //1EGLD

	ti, err := interactors.NewTransactionInteractor(ep, blockchain.NewTxSigner())
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return
	}

	tx, err := ti.ApplySignatureAndGenerateTransaction(privateKey, transactionArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return
	}
	ti.AddTransaction(tx)

	//a new transaction with the signature done on the hash of the transaction
	//it's ok to reuse the arguments here, they will be copied, anyway
	transactionArguments.Version = 2
	transactionArguments.Options = 1
	transactionArguments.Nonce++ //do not forget to increment the nonce, otherwise you will get 2 transactions
	// with the same nonce (only one of them will get executed)
	txSigOnHash, err := ti.ApplySignatureAndGenerateTransaction(privateKey, transactionArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return
	}
	ti.AddTransaction(txSigOnHash)

	hashes, err := ti.SendTransactionsAsBunch(100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
		return
	}

	log.Info("transactions sent", "hashes", hashes)
}
