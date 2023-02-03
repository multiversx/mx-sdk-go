package aggregator

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewPair(t *testing.T) {
	t.Parallel()

	t.Run("invalid base name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()
		args.Base = ""

		pn, err := newPair(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, ErrNilBaseName, err)
	})
	t.Run("invalid quote name", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()
		args.Quote = ""

		pn, err := newPair(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, ErrNilQuoteName, err)
	})
	t.Run("0 decimals", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()
		args.Decimals = 0

		pn, err := newPair(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, ErrInvalidDecimals))
	})
	t.Run(">18 decimals", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()
		args.Decimals = 19

		pn, err := newPair(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, ErrInvalidDecimals))
	})
	t.Run("nil exchanges map", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()
		args.Exchanges = nil

		pn, err := newPair(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, ErrNilExchanges, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPair()

		pn, err := newPair(args)
		assert.False(t, check.IfNil(pn))
		assert.Nil(t, err)
	})
}

func createMockArgsPair() *ArgsPair {
	return &ArgsPair{
		Base:                      "BASE",
		Quote:                     "QUOTE",
		PercentDifferenceToNotify: 1,
		Decimals:                  2,
		Exchanges:                 map[string]struct{}{"Binance": {}},
	}
}
