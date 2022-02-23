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
	GeneralSignedBlocksThreshold              float32             `json:"erd_ratings_general_signed_blocks_threshold,string"`
	GeneralStartRating                        uint32              `json:"erd_ratings_general_start_rating"`
	GeneralSelectionChances                   []*SelectionChances `json:"erd_ratings_general_selection_chances"`
	MetachainConsecutiveMissedBlocksPenalty   float32             `json:"erd_ratings_metachain_consecutive_missed_blocks_penalty,string"`
	MetachainHoursToMaxRatingFromStartRating  uint32              `json:"erd_ratings_metachain_hours_to_max_rating_from_start_rating"`
	MetachainProposerDecreaseFactor           float32             `json:"erd_ratings_metachain_proposer_decrease_factor,string"`
	MetachainProposerValidatorImportance      float32             `json:"erd_ratings_metachain_proposer_validator_importance,string"`
	MetachainValidatorDecreaseFactor          float32             `json:"erd_ratings_metachain_validator_decrease_factor,string"`
	PeerhonestyBadPeerThreshold               float64             `json:"erd_ratings_peerhonesty_bad_peer_threshold,string"`
	PeerhonestyDecayCoefficient               float64             `json:"erd_ratings_peerhonesty_decay_coefficient,string"`
	PeerhonestyDecayUpdateIntervalInseconds   uint32              `json:"erd_ratings_peerhonesty_decay_update_interval_inseconds"`
	PeerhonestyMaxScore                       float64             `json:"erd_ratings_peerhonesty_max_score,string"`
	PeerhonestyMinScore                       float64             `json:"erd_ratings_peerhonesty_min_score,string"`
	PeerhonestyUnitValue                      float64             `json:"erd_ratings_peerhonesty_unit_value,string"`
	ShardchainConsecutiveMissedBlocksPenalty  float32             `json:"erd_ratings_shardchain_consecutive_missed_blocks_penalty,string"`
	ShardchainHoursToMaxRatingFromStartRating uint32              `json:"erd_ratings_shardchain_hours_to_max_rating_from_start_rating"`
	ShardchainProposerDecreaseFactor          float32             `json:"erd_ratings_shardchain_proposer_decrease_factor,string"`
	ShardchainProposerValidatorImportance     float32             `json:"erd_ratings_shardchain_proposer_validator_importance,string"`
	ShardchainValidatorDecreaseFactor         float32             `json:"erd_ratings_shardchain_validator_decrease_factor,string"`
}
