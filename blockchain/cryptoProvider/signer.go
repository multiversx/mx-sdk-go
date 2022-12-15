package cryptoProvider

import (
	"encoding/json"
	"strconv"

	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519/singlesig"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

var (
	signerLog     = logger.GetOrCreate("elrond-sdk-erdgo/signer")
	hasher        = keccak.NewKeccak()
	singleSigner  = &singlesig.Ed25519Signer{}
	messagePrefix = []byte("\x17Elrond Signed Message:\n")
)

// signer contains the primitives used to correctly sign a transaction
type signer struct{}

// NewXSigner will create a new instance of signer
func NewXSigner() *signer {
	return &signer{}
}

// SignMessage will generate the signature providing the private key bytes and the message bytes
// prepending the standard messagePrefix
func (xs *signer) SignMessage(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	serializedMessage, err := xs.serializeForSigning(msg)
	if err != nil {
		return nil, err
	}

	return singleSigner.Sign(privateKey, serializedMessage)
}

// VerifyMessage will verify the signature providing the public key bytes and the message bytes
func (xs *signer) VerifyMessage(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
	serializedMessage, err := xs.serializeForSigning(msg)
	if err != nil {
		return err
	}

	return singleSigner.Verify(publicKey, serializedMessage, sig)
}

// SignByteSlice will generate the signature providing the private key bytes and some arbitrary message
func (xs *signer) SignByteSlice(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	return singleSigner.Sign(privateKey, msg)
}

// SignTransaction will generate the signature providing the private key bytes and the serialized form of the transaction
func (xs *signer) SignTransaction(tx *data.Transaction, privateKey crypto.PrivateKey) ([]byte, error) {
	if len(tx.Signature) > 0 {
		return nil, ErrTxAlreadySigned
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	shouldSignOnTxHash := tx.Version >= 2 && tx.Options&1 > 0
	if shouldSignOnTxHash {
		signerLog.Debug("signing the transaction using the hash of the message")
		txBytes = hasher.Compute(string(txBytes))
	}

	return singleSigner.Sign(privateKey, txBytes)
}

func (xs *signer) serializeForSigning(msg []byte) ([]byte, error) {
	msgSize := strconv.FormatInt(int64(len(msg)), 10)
	msg = append([]byte(msgSize), msg...)
	msg = append(messagePrefix, msg...)

	return hasher.Compute(string(msg)), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (xs *signer) IsInterfaceNil() bool {
	return xs == nil
}
