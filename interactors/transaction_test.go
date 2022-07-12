package interactors

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionInteractor_SendTransactionsAsBunch_OneTransaction(t *testing.T) {
	t.Parallel()

	proxy := &testsCommon.ProxyStub{
		SendTransactionsCalled: func(tx []*data.Transaction) ([]string, error) {
			var msgs []string
			for i := 0; i < len(tx); i++ {
				msgs = append(msgs, "SUCCESS")
			}
			return msgs, nil
		},
	}

	txBuilder, _ := builders.NewTxBuilder(&testsCommon.TxSignerStub{})

	ti, err := NewTransactionInteractor(proxy, txBuilder)
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
	tx, err := ti.ApplyUserSignatureAndGenerateTx(make([]byte, 0), args)
	require.Nil(t, err)
	ti.AddTransaction(tx)

	msg, err := ti.SendTransactionsAsBunch(context.Background(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestTransactionInteractor_SendTransactionsAsBunch_MultipleTransactions(t *testing.T) {
	t.Parallel()

	proxy := &testsCommon.ProxyStub{
		SendTransactionsCalled: func(tx []*data.Transaction) ([]string, error) {
			var msgs []string
			for i := 0; i < len(tx); i++ {
				msgs = append(msgs, "SUCCESS")
			}
			return msgs, nil
		},
	}

	txBuilder, _ := builders.NewTxBuilder(&testsCommon.TxSignerStub{})

	ti, err := NewTransactionInteractor(proxy, txBuilder)
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
		tx, errGenerate := ti.ApplyUserSignatureAndGenerateTx(make([]byte, 0), args)
		require.Nil(t, errGenerate)
		ti.AddTransaction(tx)
		nonce++
	}

	msg, err := ti.SendTransactionsAsBunch(context.Background(), 1000)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
}

func TestTransactionInteractor_SendTransactionsAsBunch(t *testing.T) {
	t.Parallel()

	sendCalled := 0
	proxy := &testsCommon.ProxyStub{
		SendTransactionsCalled: func(txs []*data.Transaction) ([]string, error) {
			sendCalled++

			return make([]string, len(txs)), nil
		},
	}
	txBuilder, _ := builders.NewTxBuilder(&testsCommon.TxSignerStub{})
	ti, _ := NewTransactionInteractor(proxy, txBuilder)
	ti.SetTimeBetweenBunches(time.Millisecond)

	ti.AddTransaction(&data.Transaction{})
	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 0)
	assert.Nil(t, hashes)
	assert.Equal(t, ErrInvalidValue, err)

	hashes, err = ti.SendTransactionsAsBunch(context.Background(), 1)
	assert.Equal(t, 1, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	hashes, err = ti.SendTransactionsAsBunch(context.Background(), 2)
	assert.Equal(t, 0, len(hashes))
	assert.Equal(t, 0, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	ti.AddTransaction(&data.Transaction{})
	hashes, err = ti.SendTransactionsAsBunch(context.Background(), 2)
	assert.Equal(t, 1, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	numTxs := 2
	for i := 0; i < numTxs; i++ {
		ti.AddTransaction(&data.Transaction{})
	}
	hashes, err = ti.SendTransactionsAsBunch(context.Background(), 2)
	assert.Equal(t, numTxs, len(hashes))
	assert.Equal(t, 1, sendCalled)
	assert.Nil(t, err)

	sendCalled = 0
	numTxs = 101
	for i := 0; i < numTxs; i++ {
		ti.AddTransaction(&data.Transaction{})
	}
	hashes, err = ti.SendTransactionsAsBunch(context.Background(), 2)
	assert.Equal(t, numTxs, len(hashes))
	assert.Equal(t, 51, sendCalled)
	assert.Nil(t, err)
}
