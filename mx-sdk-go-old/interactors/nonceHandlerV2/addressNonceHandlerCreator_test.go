package nonceHandlerV2

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/data"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/testsCommon"
	"github.com/stretchr/testify/require"
)

func TestAddressNonceHandlerCreator_Create(t *testing.T) {
	t.Parallel()

	creator := AddressNonceHandlerCreator{}
	require.False(t, creator.IsInterfaceNil())
	pubkey := make([]byte, 32)
	_, _ = rand.Read(pubkey)
	addressHandler := data.NewAddressFromBytes(pubkey)

	create, err := creator.Create(&testsCommon.ProxyStub{}, addressHandler)
	require.Nil(t, err)
	require.NotNil(t, create)
	require.Equal(t, "*nonceHandlerV2.addressNonceHandler", fmt.Sprintf("%T", create))

}
