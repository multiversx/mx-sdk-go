package main

import (
	"bytes"
	"fmt"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
)

var log = logger.GetOrCreate("examples/examplesWallet")

func main() {
	ep := blockchain.NewElrondProxy(examples.TestnetGateway, nil)

	// Generating new mnemonic wallet
	mnemonic, err := erdgo.GenerateNewMnemonic()
	if err != nil {
		log.Error("error generating mnemonic", "error", err)
		return
	}
	log.Info("generated mnemonics", "mnemonics", mnemonic)

	// Generating private key for account 0, address index 0 (based on the previous mnemonic)
	account := uint8(0)
	addressIndex := uint8(0)
	privateKey := erdgo.GetPrivateKeyFromMnemonic(mnemonic, account, addressIndex)
	log.Info("generated private key", "private key", privateKey)
	// Generating wallet address from the private key
	addressString, err := erdgo.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Error("error getting address from private key", "error", err)
		return
	}

	address, err := data.NewAddressFromBech32String(addressString)
	if err != nil {
		log.Error("error converting address to pubkey", "error", err)
		return
	}
	log.Info("computed public key", "public key", address.AddressBytes())

	// Retrieving network configuration parameters
	networkConfig, err := ep.GetNetworkConfig()
	if err != nil {
		log.Error("error getting network config", "error", err)
		return
	}

	shardCoordinator, err := blockchain.NewShardCoordinator(networkConfig.NumShardsWithoutMeta, 0)
	if err != nil {
		log.Error("error creating shard coordinator", "error", err)
		return
	}
	// Computing the shard ID of the address
	shard, err := shardCoordinator.ComputeShardId(address)
	if err != nil {
		log.Error("error computing shard ID", "error", err)
		return
	}

	log.Info("address on shard",
		"address (bech32)", address.AddressAsBech32String(),
		"address (hex)", address.AddressBytes(),
		"shard ID", shard)

	// Save the private key to a .PEM file and reload it
	walletFilename := fmt.Sprintf("test%v.pem", time.Now().Unix())
	log.Info("pem file saved", "filename", walletFilename)

	_ = erdgo.SavePrivateKeyToPemFile(privateKey, walletFilename)
	privateKey, err = erdgo.LoadPrivateKeyFromPemFile(walletFilename)
	// Generate the address from the loaded private key
	address2String, err := erdgo.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Error("error getting address from private key", "error", err)
		return
	}

	address2, err := data.NewAddressFromBech32String(address2String)
	if err != nil {
		log.Error("error converting address to pubkey", "error", err)
		return
	}

	// Compare the old and new addresses (should match)
	if !bytes.Equal(address.AddressBytes(), address2.AddressBytes()) {
		log.Error("different address error encountered. Something went wrong",
			"address1", address.AddressAsBech32String(),
			"address2", address2.AddressAsBech32String(),
		)
		return
	}
}
