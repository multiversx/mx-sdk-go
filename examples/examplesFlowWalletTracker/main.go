package main

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples/examplesFlowWalletTracker/mock"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/workflows"
)

const timeToExecuteRequest = time.Second

var log = logger.GetOrCreate("elrond-sdk-erdgo/examples/examplesFlowWalletTracker")

type moveBalanceHandler interface {
	GenerateMoveBalanceTransactions(ctx context.Context, addresses []string)
}

type transactionInteractor interface {
	workflows.TransactionInteractor
	SendTransactionsAsBunch(ctx context.Context, bunchSize int) ([]string, error)
}

type walletTracker interface {
	GetLatestTrackedAddresses() []string
}

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	log.Info("examplesFlowWalletTracker application started, press CTRL+C to stop the app...")

	err := runApp()
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("application gracefully closed")
	}
}

func runApp() error {
	args := blockchain.ArgsElrondProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
	}
	ep, err := blockchain.NewElrondProxy(args)
	if err != nil {
		return err
	}

	tap := mock.NewTrackableAddressProviderMock()
	mnt := &mock.MemoryNonceTracker{}
	err = setTestParams(ep, tap, mnt)
	if err != nil {
		return err
	}

	minimumBalance := big.NewInt(1000000000000000) // 0.001 EGLD

	argsWalletsTracker := workflows.WalletTrackerArgs{
		TrackableAddressesProvider: tap,
		Proxy:                      ep,
		NonceHandler:               mnt,
		CheckInterval:              time.Second * 2,
		MinimumBalance:             minimumBalance,
	}
	wt, err := workflows.NewWalletTracker(argsWalletsTracker)
	if err != nil {
		return err
	}

	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		return err
	}

	receiverAddress := "erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht" // /elrond-sdk-erdgo/interactors/testdata/test.pem
	txInteractor, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		return err
	}
	argsMoveBalanceHandler := workflows.MoveBalanceHandlerArgs{
		Proxy:                      ep,
		TxInteractor:               txInteractor,
		ReceiverAddress:            receiverAddress,
		TrackableAddressesProvider: tap,
		MinimumBalance:             minimumBalance,
	}

	mbh, err := workflows.NewMoveBalanceHandler(argsMoveBalanceHandler)
	if err != nil {
		return err
	}

	ctxNetworkConfigs, cancel := context.WithTimeout(context.Background(), timeToExecuteRequest)
	// TODO: call this periodically as to avoid network config changes that can impact the transaction creation and sending
	err = mbh.CacheNetworkConfigs(ctxNetworkConfigs)
	cancel()
	if err != nil {
		return err
	}

	ctxDone, cancel := context.WithCancel(context.Background())
	// generateMoveBalanceTransactionsAndSendThem function can be either periodically triggered or manually triggered (we choose automatically)
	interval := time.Second * 20
	go func() {
		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			timer.Reset(interval)

			select {
			case <-timer.C:
				// send transaction batches each 20 seconds
				generateMoveBalanceTransactionsAndSendThem(ctxDone, wt, txInteractor, mbh)
			case <-ctxDone.Done():
				log.Debug("closing automatically send move-balance transactions go routine...")
				return
			}
		}
	}()

	log.Info("setup complete, please send tokens to the following addresses:\n\t" + strings.Join(tap.AllTrackableAddresses(), "\n\t"))

	chStop := make(chan os.Signal)
	signal.Notify(chStop, os.Interrupt)
	<-chStop

	_ = wt.Close()
	cancel()

	time.Sleep(time.Second)

	return nil
}

func generateMoveBalanceTransactionsAndSendThem(
	ctx context.Context,
	wt walletTracker,
	txInteractor transactionInteractor,
	mbh moveBalanceHandler,
) {
	ctxSendTransactions, cancelSendTransactions := context.WithTimeout(ctx, timeToExecuteRequest)

	addresses := wt.GetLatestTrackedAddresses()
	log.Debug("trying to send move balance transactions...", "num", len(addresses))
	timeToGenerateTransactions := timeToExecuteRequest * time.Duration(len(addresses))
	ctxGenerateMoveBalanceTransactions, cancelGenerateTransactions := context.WithTimeout(ctx, timeToGenerateTransactions)
	defer func() {
		cancelSendTransactions()
		cancelGenerateTransactions()
	}()

	mbh.GenerateMoveBalanceTransactions(ctxGenerateMoveBalanceTransactions, addresses)
	hashes, errSend := txInteractor.SendTransactionsAsBunch(ctxSendTransactions, 100)
	if errSend != nil {
		log.Error(errSend.Error())
	}
	if len(hashes) > 0 {
		log.Debug("sent transactions", "hashes", strings.Join(hashes, " "))
	}
}

func setTestParams(
	ep workflows.ProxyHandler,
	trackableAddresses *mock.TrackableAddressProviderMock,
	tracker *mock.MemoryNonceTracker,
) error {

	nonce, err := ep.GetLatestHyperBlockNonce(context.Background())
	if err != nil {
		return err
	}

	// since this is an example, and we are using a memory tracker, we need this to be executed each time as to not request ancient blocks
	tracker.ProcessedNonce(nonce)

	// add 2 trackable addresses for demo purposes
	sk, _ := hex.DecodeString("45f72e8b6e8d10086bacd2fc8fa1340f82a3f5d4ef31953b463ea03c606533a6")
	trackableAddresses.AddTrackableAddress("erd1j84k44nsqsme8r6e5aawutx0z2cd6cyx3wprkzdh73x2cf0kqvksa3snnq", sk)

	sk, _ = hex.DecodeString("6babe6936d8b089a1f3b464a2050376462769782239b31dca4311e379b0391f3")
	trackableAddresses.AddTrackableAddress("erd1kjjl7lssufpmml2yy4x6cklvnxdd40c4ym3dpw93vrflwchydt3q749v2z", sk)

	return nil
}
