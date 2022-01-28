package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	hasherFactory "github.com/ElrondNetwork/elrond-go-core/hashing/factory"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	marshalizerFactory "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/process/rating"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
)

const marshalizerType = "gogo protobuf"
const hasherType = "blake2b"

type coreComponents struct {
	Marshalizer marshal.Marshalizer
	Hasher      hashing.Hasher
	Rater       sharding.ChanceComputer
}

func CreateCoreComponents(
	ratingsConfig *data.RatingsConfig,
	networkConfig *data.NetworkConfig,
) (*coreComponents, error) {
	marshalizer, err := marshalizerFactory.NewMarshalizer(marshalizerType)
	if err != nil {
		return nil, err
	}

	hasher, err := hasherFactory.NewHasher(hasherType)
	if err != nil {
		return nil, err
	}

	rater, err := createRater(ratingsConfig, networkConfig)
	if err != nil {
		return nil, err
	}

	return &coreComponents{
		Marshalizer: marshalizer,
		Hasher:      hasher,
		Rater:       rater,
	}, nil
}

func createRater(rc *data.RatingsConfig, nc *data.NetworkConfig) (sharding.ChanceComputer, error) {
	ratingsConfig := createRatingsConfig(rc)

	ratingDataArgs := rating.RatingsDataArg{
		Config:                   ratingsConfig,
		ShardConsensusSize:       uint32(nc.ShardConsensusGroupSize),
		MetaConsensusSize:        uint32(nc.MetaConsensusGroup),
		ShardMinNodes:            uint32(nc.NumNodesInShard),
		MetaMinNodes:             uint32(nc.NumMetachainNodes),
		RoundDurationMiliseconds: uint64(nc.RoundDuration),
	}

	ratingsData, err := rating.NewRatingsData(ratingDataArgs)
	if err != nil {
		return nil, err
	}

	rater, err := rating.NewBlockSigningRater(ratingsData)
	if err != nil {
		return nil, err
	}

	return rater, nil
}

func createRatingsConfig(rd *data.RatingsConfig) config.RatingsConfig {
	selectionChances := make([]*config.SelectionChance, len(rd.GeneralSelectionChances))
	for i, v := range rd.GeneralSelectionChances {
		selectionChance := &config.SelectionChance{
			MaxThreshold:  v.MaxThreshold,
			ChancePercent: v.ChancePercent,
		}
		selectionChances[i] = selectionChance
	}

	general := config.General{
		StartRating:           rd.GeneralStartRating,
		MaxRating:             rd.GeneralMaxRating,
		MinRating:             rd.GeneralMinRating,
		SignedBlocksThreshold: rd.GetSignedBlockThreshold(),
		SelectionChances:      selectionChances,
	}

	shardChain := config.ShardChain{
		RatingSteps: config.RatingSteps{
			HoursToMaxRatingFromStartRating: rd.ShardchainHoursToMaxRatingFromStartRating,
			ProposerValidatorImportance:     rd.GetShardchainProposerValidatorImportance(),
			ProposerDecreaseFactor:          rd.GetShardchainProposerDecreaseFactor(),
			ValidatorDecreaseFactor:         rd.GetShardchainValidatorDecreaseFactor(),
			ConsecutiveMissedBlocksPenalty:  rd.GetShardchainConsecutiveMissedBlocksPenalty(),
		},
	}

	metaChain := config.MetaChain{
		RatingSteps: config.RatingSteps{
			HoursToMaxRatingFromStartRating: rd.MetachainHoursToMaxRatingFromStartRating,
			ProposerValidatorImportance:     rd.GetMetachainProposerValidatorImportance(),
			ProposerDecreaseFactor:          rd.GetMetachainProposerDecreaseFactor(),
			ValidatorDecreaseFactor:         rd.GetMetachainValidatorDecreaseFactor(),
			ConsecutiveMissedBlocksPenalty:  rd.GetMetachainConsecutiveMissedBlocksPenalty(),
		},
	}

	peerHonesty := config.PeerHonestyConfig{
		DecayCoefficient:             rd.GetPeerhonestyDecayCoefficient(),
		DecayUpdateIntervalInSeconds: rd.PeerhonestyDecayUpdateIntervalInseconds,
		MaxScore:                     rd.GetPeerhonestyMaxScore(),
		MinScore:                     rd.GetPeerhonestyMinScore(),
		BadPeerThreshold:             rd.GetPeerhonestyBadPeerThreshold(),
		UnitValue:                    rd.GetPeerhonestyUnitValue(),
	}

	ratingsConfig := config.RatingsConfig{
		General:     general,
		ShardChain:  shardChain,
		MetaChain:   metaChain,
		PeerHonesty: peerHonesty,
	}

	return ratingsConfig
}
