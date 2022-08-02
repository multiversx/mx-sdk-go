package aggregator

import (
	"fmt"
	"math"
)

const (
	minPercentDifferenceToNotify = 1
	minDecimals                  = 1
	maxDecimals                  = 18
)

// ArgsPair is the argument DTO for a pair
type ArgsPair struct {
	Base                      string
	Quote                     string
	PercentDifferenceToNotify uint32
	Decimals                  uint64
	Exchanges                 map[string]struct{}
}

type pair struct {
	Base                      string
	Quote                     string
	PercentDifferenceToNotify uint32
	Decimals                  uint64
	TrimPrecision             float64
	DenominationFactor        uint64
	Exchanges                 map[string]struct{}
}

// NewPair creates a new pair instance
func NewPair(args *ArgsPair) (*pair, error) {
	err := checkPairArgs(args)
	if err != nil {
		return nil, err
	}

	denominationFactorAsFloat64 := math.Pow(10, float64(args.Decimals))
	return &pair{
		Base:                      args.Base,
		Quote:                     args.Quote,
		PercentDifferenceToNotify: args.PercentDifferenceToNotify,
		Decimals:                  args.Decimals,
		TrimPrecision:             float64(1) / denominationFactorAsFloat64,
		DenominationFactor:        uint64(denominationFactorAsFloat64),
		Exchanges:                 args.Exchanges,
	}, nil
}

func checkPairArgs(args *ArgsPair) error {
	if len(args.Base) == 0 {
		return ErrNilBaseName
	}
	if len(args.Quote) == 0 {
		return ErrNilQuoteName
	}
	if args.PercentDifferenceToNotify < minPercentDifferenceToNotify {
		return ErrInvalidPercentDifferenceToNotify
	}
	if args.Decimals < minDecimals || args.Decimals > maxDecimals {
		return fmt.Errorf("%w, got %d for pair %s-%s", ErrInvalidDecimals,
			args.Decimals, args.Base, args.Quote)
	}
	if len(args.Exchanges) == 0 {
		return ErrNilExchanges
	}

	return nil
}
