package main

import (
	"context"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
)

var log = logger.GetOrCreate("mx-sdk-go/examples/examplesVMQuery")

func main() {
	args := blockchain.ArgsMultiversXProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewMultiversXProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return
	}

	vmRequest := &data.VmValueRequest{
		Address:    "erd1qqqqqqqqqqqqqpgqp699jngundfqw07d8jzkepucvpzush6k3wvqyc44rx",
		FuncName:   "version",
		CallerAddr: "erd1rh5ws22jxm9pe7dtvhfy6j3uttuupkepferdwtmslms5fydtrh5sx3xr8r",
		CallValue:  "",
		Args:       nil,
	}
	response, err := ep.ExecuteVMQuery(context.Background(), vmRequest)
	if err != nil {
		log.Error("error executing vm query", "error", err)
		return
	}

	contractVersion := string(response.Data.ReturnData[0])
	log.Info("response", "contract version", contractVersion)
}
