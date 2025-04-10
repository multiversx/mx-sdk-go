package nonceHandlerV1

import (
	"github.com/multiversx/mx-sdk-go/testsCommon"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testsCommon.ReplaceConverter()
	os.Exit(m.Run())
}
