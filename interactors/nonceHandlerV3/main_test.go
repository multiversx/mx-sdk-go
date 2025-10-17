package nonceHandlerV3

import (
	"os"
	"testing"

	"github.com/multiversx/mx-sdk-go/testsCommon"
)

func TestMain(m *testing.M) {
	testsCommon.ReplaceConverter()
	os.Exit(m.Run())
}
