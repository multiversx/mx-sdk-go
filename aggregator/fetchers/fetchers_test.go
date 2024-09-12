package fetchers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/aggregator"
	"github.com/multiversx/mx-sdk-go/aggregator/mock"
	"github.com/multiversx/mx-sdk-go/authentication"
	"github.com/multiversx/mx-sdk-go/authentication/native"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errShouldSkipTest = errors.New("should skip test")

const networkAddress = "https://testnet-gateway.multiversx.com"

func createMockMap() map[string]XExchangeTokensPair {
	return map[string]XExchangeTokensPair{
		"ETH-USD": {
			Base:  "WETH-b4ca29",
			Quote: "USDC-c76f1f",
		},
		"EGLD-USD": {
			Base:  "WEGLD-bd4d79",
			Quote: "USDC-c76f1f",
		},
		"BTC-USD": {
			Base:  "BTC-test1",
			Quote: "USD-test1",
		},
	}
}

func createAuthClient() (authentication.AuthClient, error) {
	w := interactors.NewWallet()
	privateKeyBytes, err := w.LoadPrivateKeyFromPemData([]byte(examples.AlicePemContents))
	if err != nil {
		return nil, err
	}

	argsProxy := blockchain.ArgsProxy{
		ProxyURL:            networkAddress,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		AllowedDeltaToFinal: 1,
		CacheExpirationTime: time.Second,
		EntityType:          core.Proxy,
	}

	proxy, err := blockchain.NewProxy(argsProxy)
	if err != nil {
		return nil, err
	}

	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKeyBytes)
	args := native.ArgsNativeAuthClient{
		Signer:                 cryptoProvider.NewSigner(),
		ExtraInfo:              struct{}{},
		Proxy:                  proxy,
		CryptoComponentsHolder: holder,
		TokenExpiryInSeconds:   60 * 60 * 24,
		Host:                   "oracle",
		TokenHandler:           native.NewAuthTokenHandler(),
	}

	authClient, err := native.NewNativeAuthClient(args)
	if err != nil {
		return nil, err
	}

	return authClient, nil
}

func Test_FunctionalTesting(t *testing.T) {
	t.Parallel()

	t.Skip("this test should be run only when doing debugging work on the component")

	responseGetter, err := aggregator.NewHttpResponseGetter()
	require.Nil(t, err)

	authClient, err := createAuthClient()
	require.Nil(t, err)

	graphqlGetter, err := aggregator.NewGraphqlResponseGetter(authClient)
	require.Nil(t, err)

	wg := sync.WaitGroup{}
	wg.Add(len(ImplementedFetchers))
	for name := range ImplementedFetchers {
		go func(fetcherName string) {
			args := ArgsPriceFetcher{
				FetcherName:        fetcherName,
				ResponseGetter:     responseGetter,
				GraphqlGetter:      graphqlGetter,
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			ethTicker := "ETH"
			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, fetchErr := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			require.Nil(t, fetchErr)
			fmt.Printf("price between %s and %s is: %v from %s\n", ethTicker, quoteUSDFiat, price, fetcherName)
			require.True(t, price > 0)
			wg.Done()
		}(name)
	}

	wg.Wait()
}

func Test_FunctionalTestingForEVMGasPrice(t *testing.T) {
	t.Parallel()

	t.Skip("this test should be run only when doing debugging work on the component")

	responseGetter, err := aggregator.NewHttpResponseGetter()
	require.Nil(t, err)

	authClient, err := createAuthClient()
	require.Nil(t, err)

	graphqlGetter, err := aggregator.NewGraphqlResponseGetter(authClient)
	require.Nil(t, err)

	// IMPORTANT: on the API URL value we should append &apikey=<APIKEY>
	//            with api keys created on the bscscan.io & etherscan.io.
	//            free plan should suffice for intended purpose

	args := ArgsPriceFetcher{
		FetcherName:        EVMGasPriceStation,
		ResponseGetter:     responseGetter,
		GraphqlGetter:      graphqlGetter,
		XExchangeTokensMap: createMockMap(),
		EVMGasConfig: EVMGasPriceFetcherConfig{
			ApiURL:   "https://api.bscscan.com/api?module=gastracker&action=gasoracle",
			Selector: "SafeGasPrice",
		},
	}
	fetcher, _ := NewPriceFetcher(args)
	fetcher.AddPair("BSC", "gas")
	price, fetchErr := fetcher.FetchPrice(context.Background(), "BSC", "gas")
	require.Nil(t, fetchErr)
	fmt.Printf("gas price for %s and is: %v from %s\n", "BSC-gas", price, fetcher.Name())
	require.True(t, price > 0)

	args.EVMGasConfig.ApiURL = "https://api.etherscan.io/api?module=gastracker&action=gasoracle"
	fetcher, _ = NewPriceFetcher(args)
	fetcher.AddPair("ETH", "gas")
	price, fetchErr = fetcher.FetchPrice(context.Background(), "ETH", "gas")
	require.Nil(t, fetchErr)
	fmt.Printf("gas price for %s and is: %v from %s\n", "ETH-gas", price, fetcher.Name())
	require.True(t, price > 0)
}

func Test_FetchPriceErrors(t *testing.T) {
	t.Parallel()

	ethTicker := "ETH"
	pair := ethTicker + quoteUSDFiat

	expectedError := errors.New("expected error")
	for f := range ImplementedFetchers {
		fetcherName := f

		t.Run("response getter errors should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := ""
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, expectedError),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, expectedError),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Equal(t, expectedError, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("empty string for price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := ""
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("negative price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := "-1"
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Equal(t, errInvalidResponseData, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("invalid string for price should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := "not a number"
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.NotNil(t, err)
			require.Equal(t, float64(0), price)
			require.IsType(t, err, &strconv.NumError{})
		})
		t.Run("xExchange: missing key from map should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			if fetcherName != XExchangeName {
				return
			}

			returnPrice := "4714.05000000"
			args := ArgsPriceFetcher{
				FetcherName:    fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			missingTicker := "missing ticker"
			fetcher.AddPair(missingTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), missingTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			assert.Equal(t, errInvalidPair, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("xExchange: invalid graphql response should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			if fetcherName != XExchangeName {
				return
			}

			args := ArgsPriceFetcher{
				FetcherName:    fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: func(ctx context.Context, url string, query string, variables string) ([]byte, error) {
						return make([]byte, 0), nil
					},
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			assert.Equal(t, errInvalidGraphqlResponse, err)
			require.Equal(t, float64(0), price)
		})
		t.Run("pair not added should error "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := ""
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Equal(t, aggregator.ErrPairNotSupported, err)
			require.Equal(t, float64(0), price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
		t.Run("should work eth-usd "+fetcherName, func(t *testing.T) {
			t.Parallel()

			returnPrice := "4714.05000000"
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, pair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(ethTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), ethTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Nil(t, err)
			require.Equal(t, 4714.05, price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
		t.Run("should work btc-usd "+fetcherName, func(t *testing.T) {
			t.Parallel()

			btcTicker := "BTC"
			btcUsdPair := btcTicker + quoteUSDFiat
			returnPrice := "4714.05000000"
			args := ArgsPriceFetcher{
				FetcherName: fetcherName,
				ResponseGetter: &mock.HttpResponseGetterStub{
					GetCalled: getFuncGetCalled(fetcherName, returnPrice, btcUsdPair, nil),
				},
				GraphqlGetter: &mock.GraphqlResponseGetterStub{
					GetCalled: getFuncQueryCalled(fetcherName, returnPrice, nil),
				},
				XExchangeTokensMap: createMockMap(),
				EVMGasConfig:       EVMGasPriceFetcherConfig{},
			}
			fetcher, _ := NewPriceFetcher(args)
			assert.False(t, check.IfNil(fetcher))

			fetcher.AddPair(btcTicker, quoteUSDFiat)
			price, err := fetcher.FetchPrice(context.Background(), btcTicker, quoteUSDFiat)
			if errors.Is(err, errShouldSkipTest) {
				return
			}
			require.Nil(t, err)
			require.Equal(t, 4714.05, price)
			assert.Equal(t, fetcherName, fetcher.Name())
		})
	}
}

func getFuncQueryCalled(name, returnPrice string, returnErr error) func(ctx context.Context, url string, query string, variables string) ([]byte, error) {
	switch name {
	case XExchangeName:
		return func(ctx context.Context, url string, query string, variables string) ([]byte, error) {
			priceArray := make([]priceResponse, 0)
			var p priceResponse

			var err error
			p.Last, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return nil, errShouldSkipTest
			}
			p.Time = time.Now()

			priceArray = append(priceArray, p)

			var response graphqlResponse
			response.Data.Trading.Pair.Price = priceArray
			responseBytes, _ := json.Marshal(response)

			return responseBytes, returnErr
		}
	}
	return nil
}

func getFuncGetCalled(name, returnPrice, pair string, returnErr error) func(ctx context.Context, url string, response interface{}) error {
	switch name {
	case BinanceName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*binancePriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case BitfinexName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*bitfinexPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case CryptocomName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*cryptocomPriceRequest)
			cast.Result.Data = []cryptocomPair{
				{
					Price: returnPrice,
				},
			}
			return returnErr
		}
	case GeminiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*geminiPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case HitbtcName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*hitbtcPriceRequest)
			cast.Price = returnPrice
			return returnErr
		}
	case HuobiName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*huobiPriceRequest)
			var err error
			cast.Ticker.Price, err = strconv.ParseFloat(returnPrice, 64)
			if err != nil {
				return errShouldSkipTest
			}
			return returnErr
		}
	case KrakenName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*krakenPriceRequest)
			cast.Result = map[string]krakenPricePair{
				pair: {[]string{returnPrice, ""}},
			}
			return returnErr
		}
	case OkxName:
		return func(ctx context.Context, url string, response interface{}) error {
			cast, _ := response.(*okxPriceRequest)
			cast.Data = []okxTicker{{returnPrice}}
			return returnErr
		}
	}

	return nil
}
