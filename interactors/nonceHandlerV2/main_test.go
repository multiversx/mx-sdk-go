package nonceHandlerV2

import (
	"errors"
	sdkCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"os"
	"testing"
)

var testAddressAsBech32String string
var testAddress sdkCore.AddressHandler
var expectedErr error

func TestMain(m *testing.M) {
	testsCommon.ReplaceConverter()

	testAddressAsBech32String = "erd1zptg3eu7uw0qvzhnu009lwxupcn6ntjxptj5gaxt8curhxjqr9tsqpsnht"
	testAddress, _ = data.NewAddressFromBech32String(testAddressAsBech32String)
	expectedErr = errors.New("expected error")

	os.Exit(m.Run())
}
