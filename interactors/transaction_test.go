package interactors

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionInteractor_ApplySignatureAndSenderWithRealTxSigner(t *testing.T) {
	t.Parallel()

	txSigner := blockchain.NewTxSigner()
	proxy := &mock.ProxyStub{}

	sk, err := hex.DecodeString("6ae10fed53a84029e53e35afdbe083688eea0917a09a9431951dd42fd4da14c40d248169f4dd7c90537f05be1c49772ddbf8f7948b507ed17fb23284cf218b7d")

	assert.Nil(t, err)

	ti, err := NewTransactionInteractor(proxy, txSigner)
	assert.Nil(t, err)
	assert.NotNil(t, ti)

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

	tx, err := ti.ApplySignatureAndGenerateTransaction(sk, args)
	require.Nil(t, err)

	assert.Equal(t, "erd1p5jgz605m47fq5mlqklpcjth9hdl3au53dg8a5tlkgegfnep3d7stdk09x", tx.SndAddr)
	assert.Equal(t, "80e1b5476c5ea9567614d9c364e1a7380b7990b53e7b6fd8431bf8536d174c8b3e73cc354b783a03e5ae0a53b128504a6bcf32c3b9bbc06f284afe1fac179e0d",
		tx.Signature)
}

func TestTransactionInteractor_SendTransactionsAsBunch_OneTransaction(t *testing.T) {
	t.Parallel()

	proxy := &mock.ProxyStub{
		SendTransactionsCalled: func(tx []*data.Transaction) ([]string, error) {
			var msgs []string
			for i := 0; i < len(tx); i++ {
				msgs = append(msgs, "SUCCESS")
			}
			return msgs, nil
		},
	}

	var signer TxSigner = &mock.TxSignerStub{}

	ti, err := NewTransactionInteractor(proxy, signer)
	assert.Nil(t, err, "Error on transaction interactor constructor")

	value := big.NewInt(999)
	args := data.ArgCreateTransaction{
		Value:     value.String(),
		RcvAddr:   "erd12dnfhej64s6c56ka369gkyj3hwv5ms0y5rxgsk2k7hkd2vuk7rvqxkalsa",
		SndAddr:   "erd1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zsf9esqq",
		GasPrice:  10,
		GasLimit:  100000,
		Data:      []byte(""),
		Signature: "394c6f1375f6511dd281465fb9dd7caf013b6512a8f8ac278bbe2151cbded89da28bd539bc1c1c7884835742712c826900c092edb24ac02de9015f0f494f6c0a",
		ChainID:   "integration test chain id",
		Version:   uint32(1),
	}
	tx := ti.createTransaction(args)
	ti.AddTransaction(tx)

	msg, err := ti.SendTransactionsAsBunch(1)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestTransactionInteractor_SendTransactionsAsBunch_MultipleTransactions(t *testing.T) {
	t.Parallel()

	proxy := &mock.ProxyStub{
		SendTransactionsCalled: func(tx []*data.Transaction) ([]string, error) {
			var msgs []string
			for i := 0; i < len(tx); i++ {
				msgs = append(msgs, "SUCCESS")
			}
			return msgs, nil
		},
	}

	var signer TxSigner = &mock.TxSignerStub{}

	ti, err := NewTransactionInteractor(proxy, signer)
	assert.Nil(t, err, "Error on transaction interactor constructor")

	value := big.NewInt(999)
	nonce := uint64(0)
	ti.SetTimeBetweenBunches(time.Millisecond)
	for nonce < 10000 {
		args := data.ArgCreateTransaction{
			Nonce:     nonce,
			Value:     value.String(),
			RcvAddr:   "erd12dnfhej64s6c56ka369gkyj3hwv5ms0y5rxgsk2k7hkd2vuk7rvqxkalsa",
			SndAddr:   "erd1l20m7kzfht5rhdnd4zvqr82egk7m4nvv3zk06yw82zqmrt9kf0zsf9esqq",
			GasPrice:  10,
			GasLimit:  100000,
			Data:      []byte(""),
			Signature: "394c6f1375f6511dd281465fb9dd7caf013b6512a8f8ac278bbe2151cbded89da28bd539bc1c1c7884835742712c826900c092edb24ac02de9015f0f494f6c0a",
			ChainID:   "integration test chain id",
			Version:   uint32(1),
		}
		tx := ti.createTransaction(args)
		ti.AddTransaction(tx)
		nonce++
	}

	msg, err := ti.SendTransactionsAsBunch(1000)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestTransactionInteractor_SendTransactionsAsBunch(t *testing.T) {
	t.Parallel()

	sendCalled := 0
	proxy := &mock.ProxyStub{
		SendTransactionsCalled: func(txs []*data.Transaction) ([]string, error) {
			sendCalled++

			return make([]string, len(txs)), nil
		},
	}
	txSigner := &mock.TxSignerStub{}
	ti, _ := NewTransactionInteractor(proxy, txSigner)
	ti.SetTimeBetweenBunches(time.Millisecond)

	ti.AddTransaction(&data.Transaction{})
	hashes, err := ti.SendTransactionsAsBunch(0)
	assert.Nil(t, hashes)
	assert.Equal(t, ErrInvalidValue, err)

	hashes, err = ti.SendTransactionsAsBunch(1)
	assert.Equal(t, 1, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	hashes, err = ti.SendTransactionsAsBunch(2)
	assert.Equal(t, 0, len(hashes))
	assert.Equal(t, 0, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	ti.AddTransaction(&data.Transaction{})
	hashes, err = ti.SendTransactionsAsBunch(2)
	assert.Equal(t, 1, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	numTxs := 2
	for i := 0; i < numTxs; i++ {
		ti.AddTransaction(&data.Transaction{})
	}
	hashes, err = ti.SendTransactionsAsBunch(2)
	assert.Equal(t, numTxs, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	numTxs = 101
	for i := 0; i < numTxs; i++ {
		ti.AddTransaction(&data.Transaction{})
	}
	hashes, err = ti.SendTransactionsAsBunch(2)
	assert.Equal(t, numTxs, len(hashes))
	assert.Equal(t, 51, sendCalled)
	assert.Nil(t, err)
}
