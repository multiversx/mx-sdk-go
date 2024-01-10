package main

import (
	"context"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
	log    = logger.GetOrCreate("mx-sdk-go/examples/examplesTransaction")
)

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	args := blockchain.ArgsProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
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

	tx, _, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	receiverAsBech32String, err := address.AddressAsBech32String()
	if err != nil {
		log.Error("unable to get receiver address as bech 32 string", "error", err)
		return
	}

	tx.Receiver = receiverAsBech32String // send to self
	tx.Value = "1000000000000000000"     // 1EGLD

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return
	}

	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		log.Error("error signing transaction", "error", err)
		return
	}
	ti.AddTransaction(&tx)

	// a new transaction with the signature done on the hash of the transaction
	// it's ok to reuse the arguments here, they will be copied, anyway
	tx.Version = 2
	tx.Options = 1
	tx.Nonce++ // do not forget to increment the nonce, otherwise you will get 2 transactions
	// with the same nonce (only one of them will get executed)
	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		log.Error("error creating transaction", "error", err)
		return
	}
	ti.AddTransaction(&tx)

	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
		return
	}

	log.Info("transactions sent", "hashes", hashes)
}
