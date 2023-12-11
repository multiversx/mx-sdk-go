package keystores

import (
	mxcrypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-wallet/mnemonic"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-wallet/wallet"
)

type Keystore struct{}

func NewFromSecretKey(secretKey mxcrypto.PrivateKey) (*Keystore, error) {
	return nil, nil
}

// static new_from_mnemonic(wallet_provider: IWalletProvider, mnemonic: Mnemonic): Keystore
func NewFromMnemonic(walletProvider wallet.Provider, mnemonic mnemonic.Mnemonic) {

}
