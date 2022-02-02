package data

// RatingsConfigResponse holds the ratings config endpoint response
type RatingsConfigResponse struct {
	Data struct {
		Config *RatingsConfig `json:"config"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// SelectionChances holds the selection chances parameters
type SelectionChances struct {
	ChancePercent uint32 `json:"erd_chance_percent"`
	MaxThreshold  uint32 `json:"erd_max_threshold"`
}

// RatingsConfig holds the ratings configuration parameters
type RatingsConfig struct {
	GeneralMaxRating                          uint32              `json:"erd_ratings_general_max_rating"`
	GeneralMinRating                          uint32              `json:"erd_ratings_general_min_rating"`
	GeneralSignedBlocksThreshold              string              `json:"erd_ratings_general_signed_blocks_threshold"`
	GeneralStartRating                        uint32              `json:"erd_ratings_general_start_rating"`
	GeneralSelectionChances                   []*SelectionChances `json:"erd_ratings_general_selection_chances"`
	MetachainConsecutiveMissedBlocksPenalty   string              `json:"erd_ratings_metachain_consecutive_missed_blocks_penalty"`
	MetachainHoursToMaxRatingFromStartRating  uint32              `json:"erd_ratings_metachain_hours_to_max_rating_from_start_rating"`
	MetachainProposerDecreaseFactor           string              `json:"erd_ratings_metachain_proposer_decrease_factor"`
	MetachainProposerValidatorImportance      string              `json:"erd_ratings_metachain_proposer_validator_importance"`
	MetachainValidatorDecreaseFactor          string              `json:"erd_ratings_metachain_validator_decrease_factor"`
	PeerhonestyBadPeerThreshold               string              `json:"erd_ratings_peerhonesty_bad_peer_threshold"`
	PeerhonestyDecayCoefficient               string              `json:"erd_ratings_peerhonesty_decay_coefficient"`
	PeerhonestyDecayUpdateIntervalInseconds   uint32              `json:"erd_ratings_peerhonesty_decay_update_interval_inseconds"`
	PeerhonestyMaxScore                       string              `json:"erd_ratings_peerhonesty_max_score"`
	PeerhonestyMinScore                       string              `json:"erd_ratings_peerhonesty_min_score"`
	PeerhonestyUnitValue                      string              `json:"erd_ratings_peerhonesty_unit_value"`
	ShardchainConsecutiveMissedBlocksPenalty  string              `json:"erd_ratings_shardchain_consecutive_missed_blocks_penalty"`
	ShardchainHoursToMaxRatingFromStartRating uint32              `json:"erd_ratings_shardchain_hours_to_max_rating_from_start_rating"`
	ShardchainProposerDecreaseFactor          string              `json:"erd_ratings_shardchain_proposer_decrease_factor"`
	ShardchainProposerValidatorImportance     string              `json:"erd_ratings_shardchain_proposer_validator_importance"`
	ShardchainValidatorDecreaseFactor         string              `json:"erd_ratings_shardchain_validator_decrease_factor"`
}

func (rc *RatingsConfig) GetSignedBlockThreshold() float32 {
	return strToFloat32(rc.GeneralSignedBlocksThreshold)
}

func (rc *RatingsConfig) GetMetachainProposerDecreaseFactor() float32 {
	return strToFloat32(rc.MetachainProposerDecreaseFactor)
}

func (rc *RatingsConfig) GetMetachainProposerValidatorImportance() float32 {
	return strToFloat32(rc.MetachainProposerValidatorImportance)
}

func (rc *RatingsConfig) GetMetachainValidatorDecreaseFactor() float32 {
	return strToFloat32(rc.MetachainValidatorDecreaseFactor)
}

func (rc *RatingsConfig) GetPeerhonestyBadPeerThreshold() float64 {
	return strToFloat64(rc.PeerhonestyBadPeerThreshold)
}

func (rc *RatingsConfig) GetPeerhonestyDecayCoefficient() float64 {
	return strToFloat64(rc.PeerhonestyDecayCoefficient)
}

func (rc *RatingsConfig) GetPeerhonestyMaxScore() float64 {
	return strToFloat64(rc.PeerhonestyMaxScore)
}

func (rc *RatingsConfig) GetPeerhonestyMinScore() float64 {
	return strToFloat64(rc.PeerhonestyMinScore)
}

func (rc *RatingsConfig) GetPeerhonestyUnitValue() float64 {
	return strToFloat64(rc.PeerhonestyUnitValue)
}

func (rc *RatingsConfig) GetShardchainConsecutiveMissedBlocksPenalty() float32 {
	return strToFloat32(rc.ShardchainConsecutiveMissedBlocksPenalty)
}

func (rc *RatingsConfig) GetMetachainConsecutiveMissedBlocksPenalty() float32 {
	return strToFloat32(rc.MetachainConsecutiveMissedBlocksPenalty)
}

func (rc *RatingsConfig) GetShardchainProposerDecreaseFactor() float32 {
	return strToFloat32(rc.ShardchainProposerDecreaseFactor)
}

func (rc *RatingsConfig) GetShardchainProposerValidatorImportance() float32 {
	return strToFloat32(rc.ShardchainProposerValidatorImportance)
}

func (rc *RatingsConfig) GetShardchainValidatorDecreaseFactor() float32 {
	return strToFloat32(rc.ShardchainValidatorDecreaseFactor)
}
