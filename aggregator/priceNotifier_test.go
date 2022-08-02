package aggregator_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockArgsPriceNotifier() aggregator.ArgsPriceNotifier {
	return aggregator.ArgsPriceNotifier{
		Pairs: []*aggregator.ArgsPair{
			{
				Base:                      "BASE",
				Quote:                     "QUOTE",
				PercentDifferenceToNotify: 1,
				TrimPrecision:             0.01,
				DenominationFactor:        100,
			},
		},
		Aggregator:       &mock.PriceFetcherStub{},
		Notifee:          &mock.PriceNotifeeStub{},
		AutoSendInterval: time.Minute,
	}
}

func TestNewPriceNotifier(t *testing.T) {
	t.Parallel()

	t.Run("empty pair arguments slice should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs = nil

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, aggregator.ErrEmptyArgsPairsSlice, err)
	})
	t.Run("nil pair argument should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs = append(args.Pairs, nil)

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, aggregator.ErrNilArgsPair))
		assert.True(t, strings.Contains(err.Error(), "index 1"))
	})
	t.Run("0 trim precision", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].TrimPrecision = 0

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, aggregator.ErrInvalidTrimPrecision))
	})
	t.Run("0 denomination factor", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].DenominationFactor = 0

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, aggregator.ErrInvalidDenominationFactor))
	})
	t.Run("invalid auto send interval", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.AutoSendInterval = time.Second - time.Nanosecond

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.True(t, errors.Is(err, aggregator.ErrInvalidAutoSendInterval))
	})
	t.Run("nil notifee", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Notifee = nil

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, aggregator.ErrNilPriceNotifee, err)
	})
	t.Run("nil aggregator", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Aggregator = nil

		pn, err := aggregator.NewPriceNotifier(args)
		assert.True(t, check.IfNil(pn))
		assert.Equal(t, aggregator.ErrNilPriceAggregator, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()

		pn, err := aggregator.NewPriceNotifier(args)
		assert.False(t, check.IfNil(pn))
		assert.Nil(t, err)
	})
	t.Run("should work with 0 percentage to notify", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].PercentDifferenceToNotify = 0

		pn, err := aggregator.NewPriceNotifier(args)
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
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 0, expectedErr
			},
		}
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				assert.Fail(t, "should have not called notifee.PriceChanged")
				return nil
			},
		}

		pn, _ := aggregator.NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.True(t, errors.Is(err, expectedErr))
	})
	t.Run("first time should notify", func(t *testing.T) {
		t.Parallel()

		var startTimestamp, endTimestamp, receivedTimestamp int64
		args := createMockArgsPriceNotifier()
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		wasCalled := false
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				require.Equal(t, 1, len(args))
				for _, arg := range args {
					assert.Equal(t, arg.Base, "BASE")
					assert.Equal(t, arg.Quote, "QUOTE")
					assert.Equal(t, uint64(199), arg.DenominatedPrice)
					assert.Equal(t, uint64(100), arg.DenominationFactor)
					receivedTimestamp = arg.Timestamp
				}
				wasCalled = true

				return nil
			},
		}

		pn, _ := aggregator.NewPriceNotifier(args)
		startTimestamp = time.Now().Unix()
		err := pn.Execute(context.Background())
		endTimestamp = time.Now().Unix()
		assert.Nil(t, err)
		assert.True(t, wasCalled)
		assert.True(t, startTimestamp <= receivedTimestamp)
		assert.True(t, endTimestamp >= receivedTimestamp)
	})
	t.Run("double call should notify once", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				require.Equal(t, 1, len(args))
				for _, arg := range args {
					assert.Equal(t, arg.Base, "BASE")
					assert.Equal(t, arg.Quote, "QUOTE")
					assert.Equal(t, uint64(199), arg.DenominatedPrice)
					assert.Equal(t, uint64(100), arg.DenominationFactor)
				}
				numCalled++

				return nil
			},
		}

		pn, _ := aggregator.NewPriceNotifier(args)
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
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				require.Equal(t, 1, len(args))
				for _, arg := range args {
					assert.Equal(t, arg.Base, "BASE")
					assert.Equal(t, arg.Quote, "QUOTE")
					assert.Equal(t, uint64(199), arg.DenominatedPrice)
					assert.Equal(t, uint64(100), arg.DenominationFactor)
				}
				numCalled++

				return nil
			},
		}

		pn, err := aggregator.NewPriceNotifier(args)
		require.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, 2, numCalled)
	})
	t.Run("no price changes should not notify", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].PercentDifferenceToNotify = 1
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				require.Fail(t, "should have not called pricesChanged")

				return nil
			},
		}

		pn, err := aggregator.NewPriceNotifier(args)
		require.Nil(t, err)
		pn.SetLastNotifiedPrices([]float64{1.987654321})

		err = pn.Execute(context.Background())
		assert.Nil(t, err)
	})
	t.Run("no price changes but auto send duration exceeded", func(t *testing.T) {
		t.Parallel()

		startTime := time.Now()

		args := createMockArgsPriceNotifier()
		args.Pairs[0].PercentDifferenceToNotify = 1
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				return 1.987654321, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				numCalled++
				return nil
			},
		}

		time.Sleep(time.Second)
		pn, err := aggregator.NewPriceNotifier(args)
		require.Nil(t, err)
		pn.SetLastNotifiedPrices([]float64{1.987654321})

		lastTimeAutoSent := pn.LastTimeAutoSent()
		assert.True(t, lastTimeAutoSent.Sub(startTime) > 0)

		pn.SetTimeSinceHandler(func(providedTime time.Time) time.Duration {
			assert.Equal(t, pn.LastTimeAutoSent(), providedTime)

			return time.Second * time.Duration(10000)
		})

		time.Sleep(time.Second)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, 1, numCalled)
		assert.True(t, pn.LastTimeAutoSent().Sub(lastTimeAutoSent) > 0)
	})
	t.Run("price changed over the limit should notify twice", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceNotifier()
		price := 1.987654321
		args.Aggregator = &mock.PriceFetcherStub{
			FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
				price = price * 1.012 // due to rounding errors, we need this slightly higher increase

				fmt.Printf("new price: %v\n", price)

				return price, nil
			},
		}
		numCalled := 0
		args.Notifee = &mock.PriceNotifeeStub{
			PriceChangedCalled: func(ctx context.Context, args []*aggregator.ArgsPriceChanged) error {
				require.Equal(t, 1, len(args))
				for _, arg := range args {
					assert.Equal(t, arg.Base, "BASE")
					assert.Equal(t, arg.Quote, "QUOTE")
				}
				numCalled++

				return nil
			},
		}

		pn, _ := aggregator.NewPriceNotifier(args)
		err := pn.Execute(context.Background())
		assert.Nil(t, err)

		err = pn.Execute(context.Background())
		assert.Nil(t, err)

		assert.Equal(t, 2, numCalled)
	})
}
