package aggregator

import (
	"context"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func createMockArgsPriceAggregator() ArgsPriceAggregator {
	return ArgsPriceAggregator{
		PriceFetchers: []PriceFetcher{&mock.PriceFetcherStub{}},
		MinResultsNum: 1,
	}
}

func TestNewPriceAggregator(t *testing.T) {
	t.Parallel()

	t.Run("invalid MinResultsNum should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.MinResultsNum = 0
		pa, err := NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, errInvalidMinNumberOfResults))
	})
	t.Run("invalid number of price fetchers should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.PriceFetchers = make([]PriceFetcher, 0)
		pa, err := NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, errInvalidNumberOfPriceFetchers))
	})
	t.Run("nil price fetcher should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.PriceFetchers = append(args.PriceFetchers, nil)
		pa, err := NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, errNilPriceFetcher))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		pa, err := NewPriceAggregator(args)

		assert.Equal(t, "price aggregator", pa.Name())
		assert.False(t, check.IfNil(pa))
		assert.Nil(t, err)
	})
}

func TestPriceAggregator_FetchPrice(t *testing.T) {
	t.Parallel()

	t.Run("all stubs return valid values and no errors", func(t *testing.T) {
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []PriceFetcher{
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 1.0045, nil
				},
			},
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 1.0047, nil
				},
			},
		}
		pa, _ := NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Nil(t, err)
		assert.Equal(t, 1.0046, value)
	})
	t.Run("one stub returns an error", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []PriceFetcher{
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 1.0045, nil
				},
			},
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 0, expectedErr
				},
			},
		}
		pa, _ := NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Nil(t, err)
		assert.Equal(t, 1.0045, value)
	})
	t.Run("all stubs return errors", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []PriceFetcher{
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 0, expectedErr
				},
			},
			&mock.PriceFetcherStub{
				FetchPriceCalled: func(ctx context.Context, base string, quote string) (float64, error) {
					return 0, expectedErr
				},
			},
		}
		pa, _ := NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Equal(t, errNotEnoughResponses, err)
		assert.Equal(t, 0.00, value)
	})
}
