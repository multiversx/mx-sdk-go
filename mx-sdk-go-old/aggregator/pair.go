package aggregator

import (
	"fmt"
	"math"
)

const (
	minDecimals = 1
	maxDecimals = 18
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
	base                      string
	quote                     string
	percentDifferenceToNotify uint32
	decimals                  uint64
	trimPrecision             float64
	denominationFactor        uint64
	exchanges                 map[string]struct{}
}

func newPair(args *ArgsPair) (*pair, error) {
	err := checkPairArgs(args)
	if err != nil {
		return nil, err
	}

	denominationFactorAsFloat64 := math.Pow(10, float64(args.Decimals))
	return &pair{
		base:                      args.Base,
		quote:                     args.Quote,
		percentDifferenceToNotify: args.PercentDifferenceToNotify,
		decimals:                  args.Decimals,
		trimPrecision:             float64(1) / denominationFactorAsFloat64,
		denominationFactor:        uint64(denominationFactorAsFloat64),
		exchanges:                 args.Exchanges,
	}, nil
}

func checkPairArgs(args *ArgsPair) error {
	if len(args.Base) == 0 {
		return ErrNilBaseName
	}
	if len(args.Quote) == 0 {
		return ErrNilQuoteName
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

// IsInterfaceNil returns true if there is no value under the interface
func (p *pair) IsInterfaceNil() bool {
	return p == nil
}
