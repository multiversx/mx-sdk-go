package data

// EnableEpochsConfigResponse holds the enable epochs config endpoint response
type EnableEpochsConfigResponse struct {
	Data struct {
		Config *EnableEpochsConfig `json:"enableEpochs"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// MaxNodesChangeConfig holds the max nodes change config
type MaxNodesChangeConfig struct {
	EpochEnable            uint32 `json:"erd_epoch_enable"`
	MaxNumNodes            uint32 `json:"erd_max_num_nodes"`
	NodesToShufflePerShard uint32 `json:"erd_nodes_to_shuffle_per_shard"`
}

// EnableEpochsConfig holds the enable epochs configuration parameters
type EnableEpochsConfig struct {
	BalanceWaitingListsEnableEpoch uint32                 `json:"erd_balance_waiting_lists_enable_epoch"`
	WaitingListFixEnableEpoch      uint32                 `json:"erd_waiting_list_fix_enable_epoch"`
	MaxNodesChangeEnableEpoch      []MaxNodesChangeConfig `json:"erd_max_nodes_change_enable_epoch"`
}
