package cryptoProvider

import (
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

type cryptoComponentsHolder struct {
	publicKey      crypto.PublicKey
	privateKey     crypto.PrivateKey
	addressHandler core.AddressHandler
}

// NewCryptoComponentsHolder returns a new cryptoComponentsHolder instance
func NewCryptoComponentsHolder(skBytes []byte) (*cryptoComponentsHolder, error) {
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

	return &cryptoComponentsHolder{
		privateKey:     privateKey,
		publicKey:      publicKey,
		addressHandler: addressHandler,
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
	return holder.addressHandler.AddressAsBech32String()
}

// GetAddressHandler returns the held bech32 address
func (holder *cryptoComponentsHolder) GetAddressHandler() core.AddressHandler {
	return holder.addressHandler
}

// IsInterfaceNil returns true if there is no value under the interface
func (holder *cryptoComponentsHolder) IsInterfaceNil() bool {
	return holder == nil
}
