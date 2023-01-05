package builders

import (
	"context"
	"encoding/hex"
	"fmt"
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

func TestRelayedTxBuilder(t *testing.T) {
	proxyProvider, err := blockchain.NewElrondProxy(blockchain.ArgsElrondProxy{
		ProxyURL:            "https://testnet-gateway.elrond.com",
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	})
	require.NoError(t, err)

	netConfig, err := proxyProvider.GetNetworkConfig(context.Background())
	require.NoError(t, err)

	relayerAcc, relayerPrivKey := getAccount(t, proxyProvider, testRelayerMnemonic)
	innerSenderAcc, innerSenderPrivKey := getAccount(t, proxyProvider, testInnerSenderMnemonic)

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

	relayedTxV1Builder := NewRelayedTxV1Builder()
	relayedTxV1Builder.SetInnerTransaction(innerTx)
	relayedTxV1Builder.SetRelayerAccount(relayerAcc)
	relayedTxV1Builder.SetNetworkConfig(netConfig)

	relayedTx, err := relayedTxV1Builder.Build()
	require.NoError(t, err)

	relayedTxSig := signTx(t, relayerPrivKey, relayedTx)
	relayedTx.Signature = hex.EncodeToString(relayedTxSig)

	printTx(relayedTx)

	hash, err := proxyProvider.SendTransaction(context.Background(), relayedTx)
	require.NoError(t, err)
	fmt.Printf("The hash of the transaction is: %s\n", hash)
}

func getAccount(t *testing.T, proxy interactors.Proxy, mnemonic string) (*data.Account, []byte) {
	wallet := interactors.NewWallet()

	privKey := wallet.GetPrivateKeyFromMnemonic(data.Mnemonic(mnemonic), 0, 0)
	address, err := wallet.GetAddressFromPrivateKey(privKey)
	require.NoError(t, err)

	account, err := proxy.GetAccount(context.Background(), address)
	require.NoError(t, err)

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

func printTx(tx *data.Transaction) {
	fmt.Printf("Transaction: [\n\tNonce: %d\n\tSender: %s\n\tReceiver: %s\n\tValue: %s\n\tGasPrice: %d\n\tGasLimit: %d\n\tData: %s\n\tSignature: %s\n]\n",
		tx.Nonce, tx.SndAddr, tx.RcvAddr, tx.Value, tx.GasPrice, tx.GasLimit, tx.Data, tx.Signature)
}
