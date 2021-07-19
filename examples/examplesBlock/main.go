package main

import (
	"encoding/json"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
)

var log = logger.GetOrCreate("examples/examplesBlock")

func main() {
	ep := blockchain.NewElrondProxy(examples.TestnetGateway, nil)

	// Get latest hyper block (metachain) nonce
	nonce, err := ep.GetLatestHyperBlockNonce()
	if err != nil {
		log.Error("error retrieving latest block nonce", "error", err)
		return
	}
	log.Info("latest hyper block", "nonce", nonce)

	// Get block info
	block, errGet := ep.GetHyperBlockByNonce(nonce)
	if errGet != nil {
		log.Error("error retrieving hyper block", "error", err)
		return
	}
	data, errMarshal := json.MarshalIndent(block, "", "    ")
	if errMarshal != nil {
		log.Error("error serializing block", "error", errMarshal)
		return
	}
	log.Info("\n" + string(data))
}
