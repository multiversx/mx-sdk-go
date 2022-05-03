package endpointProviders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseEndpointProvider(t *testing.T) {
	t.Parallel()

	base := &baseEndpointProvider{}
	assert.Equal(t, networkConfigEndpoint, base.GetNetworkConfigEndpoint())
	assert.Equal(t, networkEconomicsEndpoint, base.GetNetworkEconomicsEndpoint())
	assert.Equal(t, ratingsConfigEndpoint, base.GetRatingsConfigEndpoint())
	assert.Equal(t, enableEpochsConfigEndpoint, base.GetEnableEpochsConfigEndpoint())
	assert.Equal(t, "address/addressAsBech32", base.GetAccountEndpoint("addressAsBech32"))
	assert.Equal(t, costTransactionEndpoint, base.GetCostTransactionEndpoint())
	assert.Equal(t, sendTransactionEndpoint, base.GetSendTransactionEndpoint())
	assert.Equal(t, sendMultipleTransactionsEndpoint, base.GetSendMultipleTransactionsEndpoint())
	assert.Equal(t, "transaction/hex/status", base.GetTransactionStatusEndpoint("hex"))
	assert.Equal(t, "transaction/hex", base.GetTransactionInfoEndpoint("hex"))
	assert.Equal(t, "hyperblock/by-nonce/4", base.GetHyperBlockByNonceEndpoint(4))
	assert.Equal(t, "hyperblock/by-hash/hex", base.GetHyperBlockByHashEndpoint("hex"))
	assert.Equal(t, vmValuesEndpoint, base.GetVmValuesEndpoint())
	assert.Equal(t, genesisNodesConfigEndpoint, base.GetGenesisNodesConfigEndpoint())
	assert.Equal(t, "internal/raw/startofepoch/metablock/by-epoch/5", base.GetRawStartOfEpochMetaBlock(5))
}
