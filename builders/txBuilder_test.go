package builders

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTxBuilder(t *testing.T) {
	t.Parallel()

	t.Run("nil txSigner should error", func(t *testing.T) {
		t.Parallel()

		tb, err := NewTxBuilder(nil)
		assert.True(t, check.IfNil(tb))
		assert.Equal(t, ErrNilTxSigner, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		tb, err := NewTxBuilder(&testsCommon.TxSignerStub{})
		assert.False(t, check.IfNil(tb))
		assert.Nil(t, err)
	})
}

func TestTxBuilder_ApplySignatureAndGenerateTx(t *testing.T) {
	t.Parallel()

	sk, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)
	value := big.NewInt(999)
	args := data.ArgCreateTransaction{
		Value:    value.String(),
		RcvAddr:  "erd1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zsf9esqq",
		GasPrice: 10,
		GasLimit: 100000,
		Data:     []byte(""),
		ChainID:  "integration test chain id",
		Version:  uint32(1),
	}

	t.Run("tx signer errors when generating public key should error", func(t *testing.T) {
		t.Parallel()

		argsCopy := args
		expectedErr := errors.New("expected error")
		tb, _ := NewTxBuilder(&testsCommon.TxSignerStub{
			GeneratePkBytesCalled: func(skBytes []byte) ([]byte, error) {
				return nil, expectedErr
			},
		})

		tx, errGenerate := tb.ApplyUserSignatureAndGenerateTx(sk, argsCopy)
		assert.Nil(t, tx)
		assert.Equal(t, expectedErr, errGenerate)
	})
	t.Run("tx signer errors when signing should error", func(t *testing.T) {
		t.Parallel()

		argsCopy := args
		expectedErr := errors.New("expected error")
		tb, _ := NewTxBuilder(&testsCommon.TxSignerStub{
			SignMessageCalled: func(msg []byte, skBytes []byte) ([]byte, error) {
				return nil, expectedErr
			},
		})

		tx, errGenerate := tb.ApplyUserSignatureAndGenerateTx(sk, argsCopy)
		assert.Nil(t, tx)
		assert.Equal(t, expectedErr, errGenerate)
	})

	txSigner := blockchain.NewTxSigner()
	tb, err := NewTxBuilder(txSigner)
	require.Nil(t, err)

	t.Run("sign on all tx bytes should work", func(t *testing.T) {
		t.Parallel()

		argsCopy := args
		tx, errGenerate := tb.ApplyUserSignatureAndGenerateTx(sk, argsCopy)
		require.Nil(t, errGenerate)

		assert.Equal(t, "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x", tx.SndAddr)
		assert.Equal(t, "80e1b5476c5ea9567614d9c364e1a7380b7990b53e7b6fd8431bf8536d174c8b3e73cc354b783a03e5ae0a53b128504a6bcf32c3b9bbc06f284afe1fac179e0d",
			tx.Signature)
	})
	t.Run("sign on tx hash should work", func(t *testing.T) {
		t.Parallel()

		argsCopy := args
		argsCopy.Version = 2
		argsCopy.Options = 1

		tx, errGenerate := tb.ApplyUserSignatureAndGenerateTx(sk, argsCopy)
		require.Nil(t, errGenerate)

		assert.Equal(t, "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x", tx.SndAddr)
		assert.Equal(t, "1761bcac651a65839b53e89f6b0738e0956cb12e8624826b98bfc577c9f8d5e36a2544a9c5445ce7d5059972b2c5f42e25f3ad9f59255465a2ba128f0764b90e",
			tx.Signature)
	})
}

func TestTxBuilder_ApplySignatureAndGenerateTxHash(t *testing.T) {
	t.Parallel()

	t.Run("fails if the signature is missing", func(t *testing.T) {
		t.Parallel()

		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		txHash, errGenerate := tb.ComputeTxHash(&data.Transaction{})
		assert.Nil(t, txHash)
		assert.Equal(t, ErrMissingSignature, errGenerate)
	})

	t.Run("should generate tx hash", func(t *testing.T) {
		t.Parallel()

		sk, err := hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
		require.Nil(t, err)

		args := data.ArgCreateTransaction{
			Nonce:    1,
			Value:    "11500313000000000000",
			RcvAddr:  "erd1p72ru5zcdsvgkkcm9swtvw2zy5epylwgv8vwquptkw7ga7pfvk7qz7snzw",
			GasPrice: 1000000000,
			GasLimit: 60000,
			Data:     []byte(""),
			ChainID:  "T",
			Version:  uint32(1),
		}
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())

		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)
		assert.Equal(t, "725c6aa7def724c60f02ee481734807038fef125e453242bf4dc570fc4a4f2ff1b78e996a2ec67ef8be03f9b98b0251d419cfc72c6e6c5c9e33f879af938f008", tx.Signature)

		txHash, errGenerate := tb.ComputeTxHash(tx)
		assert.Nil(t, errGenerate)
		assert.Equal(t, "8bbb2b7474deb2e67fa8f9db1eccef57ec14aa93710452a5de5ff52e5a369144", hex.EncodeToString(txHash))
	})
}

func TestTxBuilder_ApplyUserSignatureAndGenerateWithTxGuardian(t *testing.T) {
	t.Parallel()

	guardianAddress := "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x"
	skGuardian, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")
	require.Nil(t, err)

	senderAddress := "erd1lta2vgd0tkeqqadkvgef73y0efs6n3xe5ss589ufhvmt6tcur8kq34qkwr"
	sk, err := hex.DecodeString("28654d9264f55f18d810bb88617e22c117df94fa684dfe341a511a72dfbf2b68")
	require.Nil(t, err)

	args := data.ArgCreateTransaction{
		Nonce:        1,
		Value:        "11500313000000000000",
		RcvAddr:      "erd1p72ru5zcdsvgkkcm9swtvw2zy5epylwgv8vwquptkw7ga7pfvk7qz7snzw",
		GasPrice:     1000000000,
		GasLimit:     60000,
		Data:         []byte(""),
		ChainID:      "T",
		Version:      uint32(2),
		Options:      transaction.MaskGuardedTransaction,
		GuardianAddr: guardianAddress,
	}

	t.Run("no guardian option should fail", func(t *testing.T) {
		args := args
		args.Options = 0
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)

		err = tb.ApplyGuardianSignature(skGuardian, tx)
		require.Equal(t, ErrMissingGuardianOption, err)
	})

	t.Run("no guardian address should fail", func(t *testing.T) {
		args := args
		args.GuardianAddr = ""
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)

		err = tb.ApplyGuardianSignature(skGuardian, tx)
		require.NotNil(t, err)
	})

	t.Run("different guardian address should fail", func(t *testing.T) {
		args := args
		args.GuardianAddr = senderAddress
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)

		err = tb.ApplyGuardianSignature(skGuardian, tx)
		require.Equal(t, ErrGuardianDoesNotMatch, err)
	})
	t.Run("correct guardian ok", func(t *testing.T) {
		args := args
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)

		err = tb.ApplyGuardianSignature(skGuardian, tx)
		require.Nil(t, err)
	})
	t.Run("correct guardian and sign with hash ok", func(t *testing.T) {
		args := args
		args.Options |= transaction.MaskSignedWithHash
		tb, _ := NewTxBuilder(blockchain.NewTxSigner())
		tx, _ := tb.ApplyUserSignatureAndGenerateTx(sk, args)

		err = tb.ApplyGuardianSignature(skGuardian, tx)
		require.Nil(t, err)
	})
}
