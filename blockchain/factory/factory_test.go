package factory

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/testsCommon"
	"github.com/stretchr/testify/assert"
)

func TestCreateEndpointProvider(t *testing.T) {
	t.Parallel()

	t.Run("unknown type", func(t *testing.T) {
		t.Parallel()

		provider, err := CreateEndpointProvider("unknown")
		assert.True(t, check.IfNil(provider))
		assert.True(t, errors.Is(err, ErrUnknownRestAPIEntityType))
	})
	t.Run("node type", func(t *testing.T) {
		t.Parallel()

		provider, err := CreateEndpointProvider(core.ObserverNode)
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
		assert.Equal(t, "*endpointProviders.nodeEndpointProvider", fmt.Sprintf("%T", provider))
	})
	t.Run("proxy type", func(t *testing.T) {
		t.Parallel()

		provider, err := CreateEndpointProvider(core.Proxy)
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
		assert.Equal(t, "*endpointProviders.proxyEndpointProvider", fmt.Sprintf("%T", provider))
	})
}

func TestCreateFinalityProvider(t *testing.T) {
	t.Parallel()

	t.Run("unknown type", func(t *testing.T) {
		t.Parallel()

		proxyInstance := &testsCommon.ProxyStub{
			GetRestAPIEntityTypeCalled: func() core.RestAPIEntityType {
				return "unknown"
			},
		}

		provider, err := CreateFinalityProvider(proxyInstance, true)
		assert.True(t, check.IfNil(provider))
		assert.True(t, errors.Is(err, ErrUnknownRestAPIEntityType))
	})
	t.Run("disabled finality checker", func(t *testing.T) {
		t.Parallel()

		provider, err := CreateFinalityProvider(nil, false)
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
		assert.Equal(t, "*finalityProvider.disabledFinalityProvider", fmt.Sprintf("%T", provider))
	})
	t.Run("node type", func(t *testing.T) {
		t.Parallel()

		proxyInstance := &testsCommon.ProxyStub{
			GetRestAPIEntityTypeCalled: func() core.RestAPIEntityType {
				return core.ObserverNode
			},
		}

		provider, err := CreateFinalityProvider(proxyInstance, true)
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
		assert.Equal(t, "*finalityProvider.nodeFinalityProvider", fmt.Sprintf("%T", provider))
	})
	t.Run("proxy type", func(t *testing.T) {
		t.Parallel()

		proxyInstance := &testsCommon.ProxyStub{
			GetRestAPIEntityTypeCalled: func() core.RestAPIEntityType {
				return core.Proxy
			},
		}

		provider, err := CreateFinalityProvider(proxyInstance, true)
		assert.False(t, check.IfNil(provider))
		assert.Nil(t, err)
		assert.Equal(t, "*finalityProvider.proxyFinalityProvider", fmt.Sprintf("%T", provider))
	})
}
