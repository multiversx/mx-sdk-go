package main

import (
	"context"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
)

var log = logger.GetOrCreate("examples/examplesVMQuery")

func main() {
	ep := blockchain.NewElrondProxy(examples.TestnetGateway, nil)

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
