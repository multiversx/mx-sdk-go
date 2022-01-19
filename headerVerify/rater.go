package headerVerify

import (
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/process/rating"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

// NewRater
func NewRater(args rating.RatingsDataArg) (sharding.ChanceComputer, error) {
	ratingsConfig := createRatingsConfig()

	args.Config = ratingsConfig

	ratingsData, err := rating.NewRatingsData(args)
	if err != nil {
		return nil, err
	}

	rater, err := rating.NewBlockSigningRater(ratingsData)
	if err != nil {
		return nil, err
	}

	return rater, nil
}

// TODO: read economics data from config file or from proxy (if available)
func createRatingsConfig() config.RatingsConfig {
	general := config.General{
		StartRating:           5000001,
		MaxRating:             10000000,
		MinRating:             1,
		SignedBlocksThreshold: 0.01,
		SelectionChances: []*config.SelectionChance{
			{MaxThreshold: 0, ChancePercent: 5},
			{MaxThreshold: 1000000, ChancePercent: 0},
			{MaxThreshold: 2000000, ChancePercent: 16},
			{MaxThreshold: 3000000, ChancePercent: 17},
			{MaxThreshold: 4000000, ChancePercent: 18},
			{MaxThreshold: 5000000, ChancePercent: 19},
			{MaxThreshold: 6000000, ChancePercent: 20},
			{MaxThreshold: 7000000, ChancePercent: 21},
			{MaxThreshold: 8000000, ChancePercent: 22},
			{MaxThreshold: 9000000, ChancePercent: 23},
			{MaxThreshold: 10000000, ChancePercent: 24},
		},
	}

	shardChain := config.ShardChain{
		RatingSteps: config.RatingSteps{
			HoursToMaxRatingFromStartRating: 72,
			ProposerValidatorImportance:     1,
			ProposerDecreaseFactor:          -4,
			ValidatorDecreaseFactor:         -4,
			ConsecutiveMissedBlocksPenalty:  1.5,
		},
	}

	metaChain := config.MetaChain{
		RatingSteps: config.RatingSteps{
			HoursToMaxRatingFromStartRating: 55,
			ProposerValidatorImportance:     1,
			ProposerDecreaseFactor:          -4,
			ValidatorDecreaseFactor:         -4,
			ConsecutiveMissedBlocksPenalty:  1.5,
		},
	}

	peerHonesty := config.PeerHonestyConfig{
		DecayCoefficient:             0.9779,
		DecayUpdateIntervalInSeconds: 10,
		MaxScore:                     100,
		MinScore:                     -100,
		BadPeerThreshold:             -80,
		UnitValue:                    1,
	}

	ratingsConfig := config.RatingsConfig{
		General:     general,
		ShardChain:  shardChain,
		MetaChain:   metaChain,
		PeerHonesty: peerHonesty,
	}

	return ratingsConfig
}
