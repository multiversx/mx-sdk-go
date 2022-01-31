package headerCheck_test

// func TestNewHeaderVerifier(t *testing.T) {
// 	t.Parallel()

// 	t.Run("nil ratings config", func(t *testing.T) {
// 		t.Parallel()

// 		args := createMockArgHeaderVerifier()
// 		args.RatingsConfig = nil
// 		hv, err := headerCheck.NewHeaderVerifier(args)

// 		assert.True(t, check.IfNil(hv))
// 		assert.True(t, errors.Is(err, headerCheck.ErrNilRatingsConfig))
// 	})
// 	t.Run("nil network config", func(t *testing.T) {
// 		t.Parallel()

// 		args := createMockArgHeaderVerifier()
// 		args.NetworkConfig = nil
// 		hv, err := headerCheck.NewHeaderVerifier(args)

// 		assert.True(t, check.IfNil(hv))
// 		assert.True(t, errors.Is(err, headerCheck.ErrNilNetworkConfig))
// 	})
// 	t.Run("nil enableEpochs config", func(t *testing.T) {
// 		t.Parallel()

// 		args := createMockArgHeaderVerifier()
// 		args.EnableEpochsConfig = nil
// 		hv, err := headerCheck.NewHeaderVerifier(args)

// 		assert.True(t, check.IfNil(hv))
// 		assert.True(t, errors.Is(err, headerCheck.ErrNilEnableEpochsConfig))
// 	})
// 	t.Run("should work", func(t *testing.T) {
// 		t.Parallel()

// 		args := createMockArgHeaderVerifier()
// 		hv, err := headerCheck.NewHeaderVerifier(args)

// 		assert.False(t, check.IfNil(hv))
// 		assert.Nil(t, err)
// 	})
// }

// func TestHeaderVerifier_IsInCache(t *testing.T) {
// 	t.Parallel()

// 	args := createMockArgHeaderVerifier()
// 	hv, err := headerCheck.NewHeaderVerifier(args)
// 	require.Nil(t, err)

// 	isInCache := hv.IsInCache(uint32(0))
// 	assert.True(t, isInCache)

// 	isInCache = hv.IsInCache(uint32(1))
// 	assert.False(t, isInCache)
// }

// func TestHeaderVerifier_SetNodesConfigPerEpoch(t *testing.T) {
// 	t.Parallel()

// 	args := createMockArgHeaderVerifier()
// 	hv, err := headerCheck.NewHeaderVerifier(args)
// 	require.Nil(t, err)
// }

// func createMockArgHeaderVerifier() headerCheck.ArgHeaderVerifier {
// 	return headerCheck.ArgHeaderVerifier{
// 		RatingsConfig: createDummyRatingsConfig(),
// 		NetworkConfig: createDummyNetworkConfig(),
// 		EnableEpochsConfig: &data.EnableEpochsConfig{
// 			BalanceWaitingListsEnableEpoch: 0,
// 			WaitingListFixEnableEpoch:      0,
// 			MaxNodesChangeEnableEpoch:      []data.MaxNodesChangeConfig{},
// 		},
// 	}
// }

// func createDummyRatingsConfig() *data.RatingsConfig {
// 	selectionChances := []*data.SelectionChances{
// 		{
// 			ChancePercent: 5,
// 			MaxThreshold:  0,
// 		},
// 		{
// 			ChancePercent: 20,
// 			MaxThreshold:  10000000,
// 		},
// 	}

// 	return &data.RatingsConfig{
// 		GeneralMaxRating:                          10000000,
// 		GeneralMinRating:                          1,
// 		GeneralSignedBlocksThreshold:              "0.025",
// 		GeneralStartRating:                        5000001,
// 		GeneralSelectionChances:                   selectionChances,
// 		MetachainConsecutiveMissedBlocksPenalty:   "1.1",
// 		MetachainHoursToMaxRatingFromStartRating:  2,
// 		MetachainProposerDecreaseFactor:           "-4",
// 		MetachainProposerValidatorImportance:      "1.0",
// 		MetachainValidatorDecreaseFactor:          "-4",
// 		PeerhonestyBadPeerThreshold:               "1.0",
// 		PeerhonestyDecayCoefficient:               "1.0",
// 		PeerhonestyDecayUpdateIntervalInseconds:   0,
// 		PeerhonestyMaxScore:                       "1.0",
// 		PeerhonestyMinScore:                       "1.0",
// 		PeerhonestyUnitValue:                      "1.0",
// 		ShardchainConsecutiveMissedBlocksPenalty:  "1.1",
// 		ShardchainHoursToMaxRatingFromStartRating: 2,
// 		ShardchainProposerDecreaseFactor:          "-4",
// 		ShardchainProposerValidatorImportance:     "1.0",
// 		ShardchainValidatorDecreaseFactor:         "-4"}
// }

// func createDummyNetworkConfig() *data.NetworkConfig {
// 	return &data.NetworkConfig{
// 		ChainID:                  "test",
// 		Denomination:             1,
// 		GasPerDataByte:           2,
// 		LatestTagSoftwareVersion: "test",
// 		MetaConsensusGroup:       3,
// 		MinGasLimit:              4,
// 		MinGasPrice:              5,
// 		MinTransactionVersion:    6,
// 		NumMetachainNodes:        3,
// 		NumNodesInShard:          3,
// 		NumShardsWithoutMeta:     2,
// 		RoundDuration:            10,
// 		ShardConsensusGroupSize:  3,
// 		StartTime:                12,
// 		Adaptivity:               "true",
// 		Hysteresys:               "0.1",
// 	}
// }
