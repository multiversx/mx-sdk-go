package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/fetchers"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core/polling"
)

var log = logger.GetOrCreate("elrond-sdk-erdgo/examples/examplesPriceAggregator")

const base = "ETH"
const quote = "USD"
const percentDifferenceToNotify = 1 // 0 will notify on each fetch
const trimPrecision = 0.01          // will round the price to the hundredth
const denominationFactor = 100

const minResultsNum = 3
const pollInterval = time.Second * 2
const autoSendInterval = time.Second * 10

func main() {
	_ = logger.SetLogLevel("*:DEBUG")

	log.Info("examplesPriceAggregator will fetch the price of a defined pair from a bunch of exchanges, and will" +
		" notify a printer if the price changed")
	log.Info("application started, press CTRL+C to stop the app...")

	err := runApp()
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("application gracefully closed")
	}
}

func runApp() error {
	priceFetchers, err := createPriceFetchers()
	if err != nil {
		return err
	}

	argsPriceAggregator := aggregator.ArgsPriceAggregator{
		PriceFetchers: priceFetchers,
		MinResultsNum: minResultsNum,
	}
	aggregatorInstance, err := aggregator.NewPriceAggregator(argsPriceAggregator)
	if err != nil {
		return err
	}

	printNotifee := &mock.PriceNotifeeStub{
		PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
			for _, arg := range args {
				log.Info("Notified about the price changed",
					"pair", fmt.Sprintf("%s-%s", arg.Base, arg.Quote),
					"denominated price", arg.DenominatedPrice,
					"denomination factor", arg.DenominationFactor,
					"timestamp", arg.Timestamp)
			}

			return nil
		},
	}

	pairs := []*aggregator.ArgsPair{
		{
			Base:                      base,
			Quote:                     quote,
			PercentDifferenceToNotify: percentDifferenceToNotify,
			TrimPrecision:             trimPrecision,
			DenominationFactor:        denominationFactor,
			Exchanges:                 fetchers.ImplementedFetchers,
		},
	}
	argsPriceNotifier := aggregator.ArgsPriceNotifier{
		Pairs:            pairs,
		Aggregator:       aggregatorInstance,
		Notifee:          printNotifee,
		AutoSendInterval: autoSendInterval,
	}

	priceNotifier, err := aggregator.NewPriceNotifier(argsPriceNotifier)
	if err != nil {
		return err
	}

	addPairsToFetchers(pairs, priceFetchers)

	argsPollingHandler := polling.ArgsPollingHandler{
		Log:              log,
		Name:             "price notifier polling handler",
		PollingInterval:  pollInterval,
		PollingWhenError: pollInterval,
		Executor:         priceNotifier,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPollingHandler)
	if err != nil {
		return err
	}

	defer func() {
		errClose := pollingHandler.Close()
		log.LogIfError(errClose)
	}()

	err = pollingHandler.StartProcessingLoop()
	if err != nil {
		return err
	}

	chStop := make(chan os.Signal)
	signal.Notify(chStop, os.Interrupt)
	<-chStop

	return nil
}

func addPairsToFetchers(pairs []*aggregator.ArgsPair, priceFetchers []aggregator.PriceFetcher) {
	for _, pair := range pairs {
		for _, fetcher := range priceFetchers {
			_, ok := pair.Exchanges[fetcher.Name()]
			if ok {
				fetcher.AddPair(pair.Base, pair.Quote)
			}
		}
	}
}

func createMaiarMap() map[string]fetchers.MaiarTokensPair {
	return map[string]fetchers.MaiarTokensPair{
		"ETH-USD": {
			// for tests only until we have an ETH id
			// the price will be dropped as it is extreme compared to real price
			Base:  "WEGLD-bd4d79",
			Quote: "USDC-c76f1f",
		},
	}
}

func createPriceFetchers() ([]aggregator.PriceFetcher, error) {
	exchanges := fetchers.ImplementedFetchers
	priceFetchers := make([]aggregator.PriceFetcher, 0, len(exchanges))
	for exchangeName := range exchanges {
		priceFetcher, err := fetchers.NewPriceFetcher(exchangeName, &aggregator.HttpResponseGetter{}, createMaiarMap())
		if err != nil {
			return nil, err
		}

		priceFetchers = append(priceFetchers, priceFetcher)
	}

	return priceFetchers, nil
}
