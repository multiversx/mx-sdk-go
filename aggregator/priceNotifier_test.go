package aggregator

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsPriceNotifier() ArgsPriceNotifier {
	return ArgsPriceNotifier{
		Pairs: []*ArgsPair{
			{
				Base:                      "BASE",
				Quote:                     "QUOTE",
				PercentDifferenceToNotify: 1,
				TrimPrecision:             0.01,
			},
		},
		Fetcher: &mock.PriceFetcherStub{},
		Notifee: &mock.PriceNotifeeStub{},
	}
}

func TestNewPriceNotifier(t *testing.T) {
	t.Parallel()

	t.Run("empty pair arguments slice should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs = nil

		pn, err := NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, errEmptyArgsPairsSlice, err)
	})
	t.Run("nil pair argument should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs = append(args.Pairs, nil)

		pn, err := NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, errNilArgsPair))
		assert.True(t, strings.Contains(err.Error(), "index 1"))
	})
	t.Run("0 trim precision", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].TrimPrecision = 0

		pn, err := NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, errInvalidTrimPrecision))
	})
	t.Run("nil notifee", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Notifee = nil

		pn, err := NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, errNilPriceNotifee, err)
	})
	t.Run("nil fetcher", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Fetcher = nil

		pn, err := NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, errNilPriceFetcher, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()

		pn, err := NewPriceNotifier(args)
		assert.False(t, check.IfNil(pn))
		assert.Nil(t, err)
	})
	t.Run("should work with 0 percentage to notify", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].PercentDifferenceToNotify = 0

		pn, err := NewPriceNotifier(args)
		assert.False(t, check.IfNil(pn))
		assert.Nil(t, err)
	})
}

func TestPriceNotifier_Execute(t *testing.T) {
	t.Parallel()

	t.Run("price fetch errors should error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockArgsPriceNotifier()
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 0, expectedErr
			},
		}
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Fail(t, "should have not called notifee.PriceChanged")
				return nil
			},
		}

		pn, _ := NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.True(t, errors.Is(err, expectedErr))
	})
	t.Run("first time should notify", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		wasCalled := false
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Equal(t, base, "BASE")
				assert.Equal(t, quote, "QUOTE")
				assert.Equal(t, 1.99, price)
				wasCalled = true

				return nil
			},
		}

		pn, _ := NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.Nil(t, err)
		assert.True(t, wasCalled)
	})
	t.Run("double call should notify once", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Equal(t, base, "BASE")
				assert.Equal(t, quote, "QUOTE")
				assert.Equal(t, 1.99, price)
				numCalled++

				return nil
			},
		}

		pn, _ := NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, 1, numCalled)
	})
	t.Run("double call should notify twice if the percentage value is 0", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].PercentDifferenceToNotify = 0
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Equal(t, base, "BASE")
				assert.Equal(t, quote, "QUOTE")
				assert.Equal(t, 1.99, price)
				numCalled++

				return nil
			},
		}

		pn, err := NewPriceNotifier(args)
		require.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, 2, numCalled)
	})
	t.Run("one notify fails should try to notify all", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs = []*ArgsPair{
			{
				Base:                      "BASE1",
				Quote:                     "QUOTE1",
				PercentDifferenceToNotify: 1,
				TrimPrecision:             0.01,
			},
			{
				Base:                      "BASE2",
				Quote:                     "QUOTE2",
				PercentDifferenceToNotify: 1,
				TrimPrecision:             0.01,
			},
		}
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		expectedErr := errors.New("expected error")
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Equal(t, 1.99, price)
				numCalled++

				if base == "BASE1" {
					return expectedErr
				}

				return nil
			},
		}

		pn, _ := NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.True(t, errors.Is(err, expectedErr))

		assert.Equal(t, 2, numCalled)
	})
	t.Run("price changed over the limit should notify twice", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		price := 1.987654321
		args.Fetcher = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				price = price * 1.012 // due to rounding errors, we need this slightly higher increase

				fmt.Printf("new price: %v\n", price)

				return price, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, base string, quote string, price float64) error {
				assert.Equal(t, base, "BASE")
				assert.Equal(t, quote, "QUOTE")
				numCalled++

				return nil
			},
		}

		pn, _ := NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, 2, numCalled)
	})
}
