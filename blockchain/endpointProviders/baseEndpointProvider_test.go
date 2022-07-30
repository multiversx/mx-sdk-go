package endpointProviders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseEndpointProvider(t *testing.T) {
	t.Parallel()

	base := &baseEndpointProvider{}
	assert.Equal(t, networkConfig, base.GetNetworkConfig())
	assert.Equal(t, networkEconomics, base.GetNetworkEconomics())
	assert.Equal(t, ratingsConfig, base.GetRatingsConfig())
	assert.Equal(t, enableEpochsConfig, base.GetEnableEpochsConfig())
	assert.Equal(t, "address/addressAsBech32", base.GetAccount("addressAsBech32"))
	assert.Equal(t, "address/addressAsBech32/keys", base.GetAccountKeys("addressAsBech32"))
	assert.Equal(t, costTransaction, base.GetCostTransaction())
	assert.Equal(t, sendTransaction, base.GetSendTransaction())
	assert.Equal(t, sendMultipleTransactions, base.GetSendMultipleTransactions())
	assert.Equal(t, "transaction/hex/status", base.GetTransactionStatus("hex"))
	assert.Equal(t, "transaction/hex", base.GetTransactionInfo("hex"))
	assert.Equal(t, "hyperblock/by-nonce/4", base.GetHyperBlockByNonce(4))
	assert.Equal(t, "hyperblock/by-hash/hex", base.GetHyperBlockByHash("hex"))
	assert.Equal(t, vmValues, base.GetVmValues())
	assert.Equal(t, genesisNodesConfig, base.GetGenesisNodesConfig())
	assert.Equal(t, "internal/raw/startofepoch/metablock/by-epoch/5", base.GetRawStartOfEpochMetaBlock(5))
}
