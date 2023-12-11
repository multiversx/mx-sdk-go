package aggregator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/aggregator"
	"github.com/multiversx/mx-sdk-go/mx-sdk-go-old/aggregator/mock"
	"github.com/stretchr/testify/assert"
)

func createMockArgsPriceAggregator() aggregator.ArgsPriceAggregator {
	return aggregator.ArgsPriceAggregator{
		PriceFetchers: []aggregator.PriceFetcher{&mock.PriceFetcherStub{}},
		MinResultsNum: 1,
	}
}

func TestNewPriceAggregator(t *testing.T) {
	t.Parallel()

	t.Run("invalid MinResultsNum should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.MinResultsNum = 0
		pa, err := aggregator.NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, aggregator.ErrInvalidMinNumberOfResults))
	})
	t.Run("invalid number of price fetchers should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.PriceFetchers = make([]aggregator.PriceFetcher, 0)
		pa, err := aggregator.NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, aggregator.ErrInvalidNumberOfPriceFetchers))
	})
	t.Run("nil price fetcher should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		args.PriceFetchers = append(args.PriceFetchers, nil)
		pa, err := aggregator.NewPriceAggregator(args)

		assert.True(t, check.IfNil(pa))
		assert.True(t, errors.Is(err, aggregator.ErrNilPriceFetcher))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsPriceAggregator()
		pa, err := aggregator.NewPriceAggregator(args)

		assert.Equal(t, "price aggregator", pa.Name())
		assert.False(t, check.IfNil(pa))
		assert.Nil(t, err)
	})
}

func TestPriceAggregator_FetchPrice(t *testing.T) {
	t.Parallel()

	t.Run("all stubs return valid values and no errors", func(t *testing.T) {
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []aggregator.PriceFetcher{
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
		pa, _ := aggregator.NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Nil(t, err)
		assert.Equal(t, 1.0046, value)
	})
	t.Run("one stub returns an error", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []aggregator.PriceFetcher{
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
		pa, _ := aggregator.NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Nil(t, err)
		assert.Equal(t, 1.0045, value)
	})
	t.Run("all stubs return errors", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		args := createMockArgsPriceAggregator()
		args.PriceFetchers = []aggregator.PriceFetcher{
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
		pa, _ := aggregator.NewPriceAggregator(args)

		value, err := pa.FetchPrice(context.Background(), "", "")
		assert.Equal(t, aggregator.ErrNotEnoughResponses, err)
		assert.Equal(t, 0.00, value)
	})
}
