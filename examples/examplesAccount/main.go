package main

import (
	"context"
	"fmt"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/examples/examplesAccount")

func main() {
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

	// Retrieving network configuration parameters
	networkConfig, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		log.Error("error getting network config", "error", err)
		return
	}

	addressAsBech32 := "erd1adfmxhyczrl2t97yx92v5nywqyse0c7qh4xs0p4artg2utnu90pspgvqty"
	address, err := data.NewAddressFromBech32String(addressAsBech32)
	if err != nil {
		log.Error("invalid address", "error", err)
		return
	}

	// Retrieve account info from the network (balance, nonce)
	accountInfo, err := ep.GetAccount(context.Background(), address)
	if err != nil {
		log.Error("error retrieving account info", "error", err)
		return
	}
	floatBalance, err := accountInfo.GetBalance(networkConfig.Denomination)
	if err != nil {
		log.Error("unable to compute balance", "error", err)
		return
	}

	log.Info("account details",
		"address", addressAsBech32,
		"nonce", accountInfo.Nonce,
		"balance as float", fmt.Sprintf("%.6f", floatBalance),
		"balance as int", accountInfo.Balance,
	)
}
