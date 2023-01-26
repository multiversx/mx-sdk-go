package cryptoProvider

import (
	"encoding/json"
	"strconv"

	"github.com/multiversx/mx-chain-core-go/hashing/keccak"
	"github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519/singlesig"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/data"
)

var (
	signerLog     = logger.GetOrCreate("mx-sdk-go/signer")
	hasher        = keccak.NewKeccak()
	singleSigner  = &singlesig.Ed25519Signer{}
	messagePrefix = []byte("\x17Elrond Signed Message:\n")
)

// signer contains the primitives used to correctly sign a transaction
type signer struct{}

// NewSigner will create a new instance of signer
func NewSigner() *signer {
	return &signer{}
}

// SignMessage will generate the signature providing the private key bytes and the message bytes
// prepending the standard messagePrefix
func (s *signer) SignMessage(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	serializedMessage := s.serializeForSigning(msg)

	return s.SignByteSlice(serializedMessage, privateKey)
}

// VerifyMessage will verify the signature providing the public key bytes and the message bytes
func (s *signer) VerifyMessage(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
	serializedMessage := s.serializeForSigning(msg)

	return singleSigner.Verify(publicKey, serializedMessage, sig)
}

// SignByteSlice will generate the signature providing the private key bytes and some arbitrary message
func (s *signer) SignByteSlice(msg []byte, privateKey crypto.PrivateKey) ([]byte, error) {
	return singleSigner.Sign(privateKey, msg)
}

// VerifyByteSlice will verify the signature providing the public key bytes and the message bytes
func (s *signer) VerifyByteSlice(msg []byte, publicKey crypto.PublicKey, sig []byte) error {
	return singleSigner.Verify(publicKey, msg, sig)
}

// SignTransaction will generate the signature providing the private key bytes and the serialized form of the transaction
func (s *signer) SignTransaction(tx *data.Transaction, privateKey crypto.PrivateKey) ([]byte, error) {
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

	return s.SignByteSlice(txBytes, privateKey)
}

func (s *signer) serializeForSigning(msg []byte) []byte {
	msgSize := strconv.FormatInt(int64(len(msg)), 10)
	msg = append([]byte(msgSize), msg...)
	msg = append(messagePrefix, msg...)

	return hasher.Compute(string(msg))
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *signer) IsInterfaceNil() bool {
	return s == nil
}
