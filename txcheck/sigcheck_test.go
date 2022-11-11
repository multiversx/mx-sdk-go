package txcheck_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	hasherFactory "github.com/ElrondNetwork/elrond-go-core/hashing/factory"
	marshallerFactory "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/txcheck"
	"github.com/stretchr/testify/require"
)

func Test_VerifyTransactionSignature(t *testing.T) {
	signer := blockchain.NewTxSigner()

	sk, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)

	pk, err := signer.GeneratePkBytes(sk)
	require.Nil(t, err)

	builder, err := builders.NewTxBuilder(signer)
	require.Nil(t, err)

	value := big.NewInt(999)
	txArg := data.ArgCreateTransaction{
		Value:    value.String(),
		RcvAddr:  "erd1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zsf9esqq",
		GasPrice: 10,
		GasLimit: 100000,
		Data:     []byte(""),
		ChainID:  "chain id",
		Version:  uint32(1),
	}

	tx, err :=builder.ApplyUserSignatureAndGenerateTx(sk, txArg)
	require.Nil(t, err)

	signature, err := hex.DecodeString(tx.Signature)
	require.Nil(t, err)

	marshaller, err := marshallerFactory.NewMarshalizer(marshallerFactory.JsonMarshalizer)
	require.Nil(t, err)

	hasher, err := hasherFactory.NewHasher("sha256")
	require.Nil(t, err)

	t.Run("nil transaction should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(nil, pk, signature, signer, marshaller, hasher)
		require.Equal(t, txcheck.ErrNilTransaction, err)
	})

	t.Run("nil public key should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, nil, signature, signer, marshaller, hasher)
		require.Equal(t, txcheck.ErrNilPubKey, err)
	})
	t.Run("nil signature should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, nil, signer, marshaller, hasher)
		require.Equal(t, txcheck.ErrNilSignature, err)
	})
	t.Run("nil verifier should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, signature, nil, marshaller, hasher)
		require.Equal(t, txcheck.ErrNilSignatureVerifier, err)
	})
	t.Run("nil marshaller should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, signature, signer, nil, hasher)
		require.Equal(t, txcheck.ErrNilMarshaller, err)
	})
	t.Run("nil hasher should err", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, signature, signer, marshaller, nil)
		require.Equal(t, txcheck.ErrNilHasher, err)
	})
	t.Run("verify OK", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, signature, signer, marshaller, hasher)
		require.Nil(t, err)
	})
}
