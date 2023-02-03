package gin

import (
	"net/http"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewWebServerHandler(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ws, err := NewWebServerHandler("127.0.0.1:8080")
		assert.Nil(t, err)
		assert.False(t, check.IfNil(ws))
	})
}

func TestWebServer_StartHttpServer(t *testing.T) {
	t.Run("upgrade on get returns error", func(t *testing.T) {
		ws, _ := NewWebServerHandler("127.0.0.1:8080")
		assert.False(t, check.IfNil(ws))

		err := ws.StartHttpServer()
		assert.Nil(t, err)

		time.Sleep(2 * time.Second)

		resp, err := http.Get("http://127.0.0.1:8080/log")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // Bad request

		time.Sleep(2 * time.Second)
		err = ws.Close()
		assert.Nil(t, err)
	})
	t.Run("should work", func(t *testing.T) {
		ws, _ := NewWebServerHandler("127.0.0.1:8080")
		assert.False(t, check.IfNil(ws))

		err := ws.StartHttpServer()
		assert.Nil(t, err)

		time.Sleep(2 * time.Second)

		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://127.0.0.1:8080/log", nil)
		assert.Nil(t, err)

		req.Header.Set("Sec-Websocket-Version", "13")
		req.Header.Set("Connection", "upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-Websocket-Key", "key")

		resp, err := client.Do(req)
		assert.Nil(t, err)

		err = resp.Body.Close()
		assert.Nil(t, err)

		time.Sleep(2 * time.Second)
		err = ws.Close()
		assert.Nil(t, err)
	})
}
