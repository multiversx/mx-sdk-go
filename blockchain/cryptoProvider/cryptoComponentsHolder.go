package cryptoProvider

import (
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

type cryptoComponentsHolder struct {
	publicKey      crypto.PublicKey
	privateKey     crypto.PrivateKey
	addressHandler core.AddressHandler
	bech32Address  string
}

// NewCryptoComponentsHolder returns a new cryptoComponentsHolder instance
func NewCryptoComponentsHolder(keyGen crypto.KeyGenerator, skBytes []byte) (*cryptoComponentsHolder, error) {
	privateKey, err := keyGen.PrivateKeyFromByteArray(skBytes)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}
	addressHandler := data.NewAddressFromBytes(publicKeyBytes)
	bech32Address, err := addressHandler.AddressAsBech32String()
	if err != nil {
		return nil, err
	}

	return &cryptoComponentsHolder{
		privateKey:     privateKey,
		publicKey:      publicKey,
		addressHandler: addressHandler,
		bech32Address:  bech32Address,
	}, nil
}

// GetPublicKey returns the held publicKey
func (holder *cryptoComponentsHolder) GetPublicKey() crypto.PublicKey {
	return holder.publicKey
}

// GetPrivateKey returns the held privateKey
func (holder *cryptoComponentsHolder) GetPrivateKey() crypto.PrivateKey {
	return holder.privateKey
}

// GetBech32 returns the held bech32 address
func (holder *cryptoComponentsHolder) GetBech32() string {
	return holder.bech32Address
}

// GetAddressHandler returns the held bech32 address
func (holder *cryptoComponentsHolder) GetAddressHandler() core.AddressHandler {
	return holder.addressHandler
}

// IsInterfaceNil returns true if there is no value under the interface
func (holder *cryptoComponentsHolder) IsInterfaceNil() bool {
	return holder == nil
}
