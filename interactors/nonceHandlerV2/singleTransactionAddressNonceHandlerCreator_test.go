package nonceHandlerV2

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestSingleTransactionAddressNonceHandlerCreator_Create(t *testing.T) {
	t.Parallel()

	creator := singleTransactionAddressNonceHandlerCreator{}
	require.False(t, creator.IsInterfaceNil())
	pubkey := make([]byte, 32)
	_, _ = rand.Read(pubkey)
	addressHandler := data.NewAddressFromBytes(pubkey)

	create, err := creator.Create(&testsCommon.ProxyStub{}, addressHandler)
	require.Nil(t, err)
	require.NotNil(t, create)
	require.Equal(t, "*nonceHandlerV2.singleTransactionAddressNonceHandler", fmt.Sprintf("%T", create))

}
