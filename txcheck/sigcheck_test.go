package txcheck_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/hashing/keccak"
	marshallerFactory "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/txcheck"
	"github.com/stretchr/testify/require"
)

func Test_VerifyTransactionSignature(t *testing.T) {
	signer := blockchain.NewTxSigner()

	sk, pk, skGuardian, pkGuardian := createUserAndGuardianKeys(t, signer)

	builder, err := builders.NewTxBuilder(signer)
	require.Nil(t, err)

	txArg := createTransactionArgs()
	tx, signature, signatureGuardian, err := createSignedTransaction(sk, skGuardian, &txArg, builder)
	require.Nil(t, err)

	txArgHashSigning := txArg
	txArgHashSigning.Options |= transaction.MaskSignedWithHash
	txArgHashSigning.Version = 2

	txHashSign, signatureOnHash, signatureGuardianOnHash, err := createSignedTransaction(sk, skGuardian, &txArg, builder)
	require.Nil(t, err)

	marshaller, err := marshallerFactory.NewMarshalizer(marshallerFactory.JsonMarshalizer)
	require.Nil(t, err)

	hasher := keccak.NewKeccak()
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
	t.Run("verify user signature OK", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pk, signature, signer, marshaller, hasher)
		require.Nil(t, err)
	})
	t.Run("verify user signature OK with hashSigning", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(txHashSign, pk, signatureOnHash, signer, marshaller, hasher)
		require.Nil(t, err)
	})
	t.Run("verify guardian signature OK", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(tx, pkGuardian, signatureGuardian, signer, marshaller, hasher)
		require.Nil(t, err)
	})
	t.Run("verify guardian signature OK with hashSigning", func(t *testing.T) {
		err := txcheck.VerifyTransactionSignature(txHashSign, pkGuardian, signatureGuardianOnHash, signer, marshaller, hasher)
		require.Nil(t, err)
	})
}

func createUserAndGuardianKeys(t *testing.T, signer builders.TxSigner) (skUser, pkUser, skGuardian, pkGuardian []byte) {
	var err error

	skUser, err = hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c4")
	require.Nil(t, err)

	pkUser, err = signer.GeneratePkBytes(skUser)
	require.Nil(t, err)

	skGuardian, err = hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
	require.Nil(t, err)

	pkGuardian, err = signer.GeneratePkBytes(skGuardian)
	require.Nil(t, err)

	return
}

func createTransactionArgs() data.ArgCreateTransaction {
	value := big.NewInt(999)
	txArg := data.ArgCreateTransaction{
		Value:        value.String(),
		RcvAddr:      "erd1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zsf9esqq",
		GasPrice:     10,
		GasLimit:     100000,
		Data:         []byte(""),
		ChainID:      "chain id",
		Version:      uint32(1),
		GuardianAddr: "erd1lta2vgd0tkeqqadkvgef73y0efs6n3xe5ss589ufhvmt6tcur8kq34qkwr",
		Options:      transaction.MaskGuardedTransaction,
	}

	return txArg
}

func createSignedTransaction(
	skUser []byte,
	skGuardian []byte,
	arg *data.ArgCreateTransaction,
	builder interactors.GuardedTxBuilder,
) (tx *data.Transaction, sigUser []byte, sigGuardian []byte, err error) {
	tx, err = builder.ApplyUserSignatureAndGenerateTx(skUser, *arg)
	if err != nil {
		return nil, nil, nil, err
	}

	err = builder.ApplyGuardianSignature(skGuardian, tx)
	if err != nil {
		return nil, nil, nil, err
	}

	signatureGuardian, err := hex.DecodeString(tx.GuardianSignature)
	if err != nil {
		return nil, nil, nil, err
	}

	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return nil, nil, nil, err
	}

	return tx, signature, signatureGuardian, nil
}
