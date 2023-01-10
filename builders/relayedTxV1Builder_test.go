package builders

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain/cryptoProvider"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/stretchr/testify/require"
)

const (
	testRelayerMnemonic     = "bid involve twenty cave offer life hello three walnut travel rare bike edit canyon ice brave theme furnace cotton swing wear bread fine latin"
	testInnerSenderMnemonic = "acid twice post genre topic observe valid viable gesture fortune funny dawn around blood enemy page update reduce decline van bundle zebra rookie real"
)

func TestRelayedTxV1Builder(t *testing.T) {
	proxyProvider, err := blockchain.NewElrondProxy(blockchain.ArgsElrondProxy{
		ProxyURL:            "https://testnet-gateway.elrond.com",
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	})
	require.NoError(t, err)

	netConfig, err := proxyProvider.GetNetworkConfig(context.Background())
	require.NoError(t, err)

	relayerAcc, relayerPrivKey := getAccount(t, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, testInnerSenderMnemonic)

	innerTx := &data.Transaction{
		Nonce:    innerSenderAcc.Nonce,
		Value:    "100000000",
		RcvAddr:  "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		SndAddr:  innerSenderAcc.Address,
		GasPrice: netConfig.MinGasPrice,
		GasLimit: netConfig.MinGasLimit,
		Data:     nil,
		ChainID:  netConfig.ChainID,
		Version:  netConfig.MinTransactionVersion,
		Options:  0,
	}

	innerTxSig := signTx(t, innerSenderPrivKey, innerTx)
	innerTx.Signature = hex.EncodeToString(innerTxSig)

	txJson, _ := json.Marshal(innerTx)
	require.Equal(t,
		`{"nonce":37,"value":"100000000","receiver":"erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th","sender":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","gasPrice":1000000000,"gasLimit":50000,"signature":"c315bc927bac80996e77b5c2f4b76506ce98b12bb0511583e8e20ef9e2bbbf58a07687e341023cf2a69bab5d116b116e32cb67bf697ce451bc5bed55b24c410a","chainID":"T","version":1}`,
		string(txJson),
	)

	relayedV1Builder := NewRelayedTxV1Builder()
	relayedV1Builder.SetInnerTransaction(innerTx)
	relayedV1Builder.SetRelayerAccount(relayerAcc)
	relayedV1Builder.SetNetworkConfig(netConfig)

	relayedTx, err := relayedV1Builder.Build()
	require.NoError(t, err)

	relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
	relayedTx.Signature = hex.EncodeToString(relayedTxSig)

	txJson, _ = json.Marshal(relayedTx)
	require.Equal(t,
		`{"nonce":37,"value":"100000000","receiver":"erd1mlh7q3fcgrjeq0et65vaaxcw6m5ky8jhu296pdxpk9g32zga6uhsemxx2a","sender":"erd1h692scsz3um6e5qwzts4yjrewxqxwcwxzavl5n9q8sprussx8fqsu70jf5","gasPrice":1000000000,"gasLimit":1060000,"data":"cmVsYXllZFR4QDdiMjI2ZTZmNmU2MzY1MjIzYTMzMzcyYzIyNzY2MTZjNzU2NTIyM2EzMTMwMzAzMDMwMzAzMDMwMzAyYzIyNzI2NTYzNjU2OTc2NjU3MjIyM2EyMjQxNTQ2YzQ4NGM3NjM5NmY2ODZlNjM2MTZkNDMzODc3NjczOTcwNjQ1MTY4Mzg2Yjc3NzA0NzQyMzU2YTY5NDk0OTZmMzM0OTQ4NGI1OTRlNjE2NTQ1M2QyMjJjMjI3MzY1NmU2NDY1NzIyMjNhMjIzMzJiMmY2NzUyNTQ2ODQxMzU1YTQxMmY0YjM5NTU1YTMzNzA3MzRmMzE3NTZjNjk0ODZjNjY2OTY5MzY0MzMwNzc2MjQ2NTI0NjUxNmI2NDMxNzkzODNkMjIyYzIyNjc2MTczNTA3MjY5NjM2NTIyM2EzMTMwMzAzMDMwMzAzMDMwMzAzMDJjMjI2NzYxNzM0YzY5NmQ2OTc0MjIzYTM1MzAzMDMwMzAyYzIyNjM2ODYxNjk2ZTQ5NDQyMjNhMjI1NjQxM2QzZDIyMmMyMjc2NjU3MjczNjk2ZjZlMjIzYTMxMmMyMjczNjk2NzZlNjE3NDc1NzI2NTIyM2EyMjc3Nzg1NzM4NmI2ZTc1NzM2NzRhNmM3NTY0Mzc1ODQzMzk0YzY0NmM0MjczMzY1OTczNTM3NTc3NTU1MjU3NDQzNjRmNDk0ZjJiNjU0YjM3NzYzMTY5Njc2NDZmNjY2YTUxNTE0OTM4Mzg3MTYxNjI3MTMxMzA1MjYxNzg0Njc1NGQ3Mzc0NmU3NjMyNmMzODM1NDY0NzM4NTcyYjMxNTY3MzZiNzg0MjQzNjczZDNkMjI3ZA==","signature":"49ff32dcfef469d14887a9940c2943a80e1d50cecf1e611a2175e4cbc1de1b900e5e41d54d9beb5760af80a87720c9099100d5b1d5ff2b73a7e37e6a364fc401","chainID":"T","version":1}`,
		string(txJson),
	)
}

func getAccount(t *testing.T, mnemonic string) (*data.Account, []byte) {
	wallet := interactors.NewWallet()

	privKey := wallet.GetPrivateKeyFromMnemonic(data.Mnemonic(mnemonic), 0, 0)
	address, err := wallet.GetAddressFromPrivateKey(privKey)
	require.NoError(t, err)

	account := &data.Account{
		Nonce:   37,
		Address: address.AddressAsBech32String(),
	}

	return account, privKey
}

func signTx(t *testing.T, privKeyBytes []byte, tx *data.Transaction) []byte {
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	privKey, err := keyGen.PrivateKeyFromByteArray(privKeyBytes)
	require.NoError(t, err)
	signer := cryptoProvider.NewSigner()

	signature, err := signer.SignTransaction(tx, privKey)
	require.NoError(t, err)

	return signature
}
