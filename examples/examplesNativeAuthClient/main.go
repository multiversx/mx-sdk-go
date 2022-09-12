package examplesNativeAuthClient

import (
	"time"

	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/authentication"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/prometheus/common/log"
)

var suite = ed25519.NewEd25519()
var keyGen = signing.NewKeyGenerator(suite)

const (
	networkAddress = "https://testnet-gateway.elrond.com"
	dataApiUrl     = "https://tools.elrond.com/data-api/graphql"
	maiarPriceUrl  = "query MaiarPriceUrl($base: String!, $quote: String!) { trading { pair(first_token: $base, second_token: $quote) { price { last time } } } }"
)

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

	w := interactors.NewWallet()
	privateKeyBytes, err := w.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
		return err
	}
	privateKey, err := keyGen.PrivateKeyFromByteArray(privateKeyBytes)
	if err != nil {
		return err
	}

	argsProxy := blockchain.ArgsElrondProxy{
		ProxyURL:            networkAddress,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		AllowedDeltaToFinal: 1,
		CacheExpirationTime: time.Second,
		EntityType:          core.RestAPIEntityType("Proxy"),
	}

	proxy, err := blockchain.NewElrondProxy(argsProxy)
	if err != nil {
		return err
	}

	args := authentication.ArgsNativeAuthClient{
		TxSigner:             blockchain.NewTxSigner(),
		ExtraInfo:            nil,
		Proxy:                proxy,
		PrivateKey:           privateKey,
		TokenExpiryInSeconds: 60 * 60 * 24,
		Host:                 "oracle",
	}

	authClient, err := authentication.NewNativeAuthClient(args)
	if err != nil {
		return err
	}

	token, err := authClient.GetAccessToken()
	if err != nil {
		return err
	}

	return nil
}
