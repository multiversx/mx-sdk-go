package disabled

// EnableEpochsHandler is a disabled implementation of EnableEpochsHandler interface
type EnableEpochsHandler struct {
}

// IsGuardAccountEnabled returns false
func (eeh *EnableEpochsHandler) IsGuardAccountEnabled() bool {
	return false
}

// IsSetGuardianEnabled returns false
func (eeh *EnableEpochsHandler) IsSetGuardianEnabled() bool {
	return false
}

// IsWipeSingleNFTLiquidityDecreaseEnabled returns false
func (eeh *EnableEpochsHandler) IsWipeSingleNFTLiquidityDecreaseEnabled() bool {
	return false
}

// BlockGasAndFeesReCheckEnableEpoch returns 0
func (eeh *EnableEpochsHandler) BlockGasAndFeesReCheckEnableEpoch() uint32 {
	return 0
}

// StakingV2EnableEpoch returns 0
func (eeh *EnableEpochsHandler) StakingV2EnableEpoch() uint32 {
	return 0
}

// ScheduledMiniBlocksEnableEpoch returns 0
func (eeh *EnableEpochsHandler) ScheduledMiniBlocksEnableEpoch() uint32 {
	return 0
}

// SwitchJailWaitingEnableEpoch returns 0
func (eeh *EnableEpochsHandler) SwitchJailWaitingEnableEpoch() uint32 {
	return 0
}

// BalanceWaitingListsEnableEpoch returns WaitingListFixEnableEpochField
func (eeh *EnableEpochsHandler) BalanceWaitingListsEnableEpoch() uint32 {
	return 0
}

// WaitingListFixEnableEpoch returns WaitingListFixEnableEpochField
func (eeh *EnableEpochsHandler) WaitingListFixEnableEpoch() uint32 {
	return 0
}

// MultiESDTTransferAsyncCallBackEnableEpoch returns 0
func (eeh *EnableEpochsHandler) MultiESDTTransferAsyncCallBackEnableEpoch() uint32 {
	return 0
}

// FixOOGReturnCodeEnableEpoch returns 0
func (eeh *EnableEpochsHandler) FixOOGReturnCodeEnableEpoch() uint32 {
	return 0
}

// RemoveNonUpdatedStorageEnableEpoch returns 0
func (eeh *EnableEpochsHandler) RemoveNonUpdatedStorageEnableEpoch() uint32 {
	return 0
}

// CreateNFTThroughExecByCallerEnableEpoch returns 0
func (eeh *EnableEpochsHandler) CreateNFTThroughExecByCallerEnableEpoch() uint32 {
	return 0
}

// FixFailExecutionOnErrorEnableEpoch returns 0
func (eeh *EnableEpochsHandler) FixFailExecutionOnErrorEnableEpoch() uint32 {
	return 0
}

// ManagedCryptoAPIEnableEpoch returns 0
func (eeh *EnableEpochsHandler) ManagedCryptoAPIEnableEpoch() uint32 {
	return 0
}

// DisableExecByCallerEnableEpoch returns 0
func (eeh *EnableEpochsHandler) DisableExecByCallerEnableEpoch() uint32 {
	return 0
}

// RefactorContextEnableEpoch returns 0
func (eeh *EnableEpochsHandler) RefactorContextEnableEpoch() uint32 {
	return 0
}

// CheckExecuteReadOnlyEnableEpoch returns 0
func (eeh *EnableEpochsHandler) CheckExecuteReadOnlyEnableEpoch() uint32 {
	return 0
}

// StorageAPICostOptimizationEnableEpoch returns 0
func (eeh *EnableEpochsHandler) StorageAPICostOptimizationEnableEpoch() uint32 {
	return 0
}

// MiniBlockPartialExecutionEnableEpoch returns 0
func (eeh *EnableEpochsHandler) MiniBlockPartialExecutionEnableEpoch() uint32 {
	return 0
}

// RefactorPeersMiniBlocksEnableEpoch returns 0
func (eeh *EnableEpochsHandler) RefactorPeersMiniBlocksEnableEpoch() uint32 {
	return 0
}

// IsSCDeployFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSCDeployFlagEnabled() bool {
	return false
}

// IsBuiltInFunctionsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsBuiltInFunctionsFlagEnabled() bool {
	return false
}

// IsRelayedTransactionsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRelayedTransactionsFlagEnabled() bool {
	return false
}

// IsPenalizedTooMuchGasFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsPenalizedTooMuchGasFlagEnabled() bool {
	return false
}

// ResetPenalizedTooMuchGasFlag does nothing
func (eeh *EnableEpochsHandler) ResetPenalizedTooMuchGasFlag() {
}

// IsSwitchJailWaitingFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSwitchJailWaitingFlagEnabled() bool {
	return false
}

// IsBelowSignedThresholdFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsBelowSignedThresholdFlagEnabled() bool {
	return false
}

// IsSwitchHysteresisForMinNodesFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSwitchHysteresisForMinNodesFlagEnabled() bool {
	return false
}

// IsSwitchHysteresisForMinNodesFlagEnabledForCurrentEpoch returns false
func (eeh *EnableEpochsHandler) IsSwitchHysteresisForMinNodesFlagEnabledForCurrentEpoch() bool {
	return false
}

// IsTransactionSignedWithTxHashFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsTransactionSignedWithTxHashFlagEnabled() bool {
	return false
}

// IsMetaProtectionFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsMetaProtectionFlagEnabled() bool {
	return false
}

// IsAheadOfTimeGasUsageFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsAheadOfTimeGasUsageFlagEnabled() bool {
	return false
}

// IsGasPriceModifierFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsGasPriceModifierFlagEnabled() bool {
	return false
}

// IsRepairCallbackFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRepairCallbackFlagEnabled() bool {
	return false
}

// IsBalanceWaitingListsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsBalanceWaitingListsFlagEnabled() bool {
	return false
}

// IsReturnDataToLastTransferFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsReturnDataToLastTransferFlagEnabled() bool {
	return false
}

// IsSenderInOutTransferFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSenderInOutTransferFlagEnabled() bool {
	return false
}

// IsStakeFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsStakeFlagEnabled() bool {
	return false
}

// IsStakingV2FlagEnabled returns false
func (eeh *EnableEpochsHandler) IsStakingV2FlagEnabled() bool {
	return false
}

// IsStakingV2OwnerFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsStakingV2OwnerFlagEnabled() bool {
	return false
}

// IsStakingV2FlagEnabledForActivationEpochCompleted returns false
func (eeh *EnableEpochsHandler) IsStakingV2FlagEnabledForActivationEpochCompleted() bool {
	return false
}

// IsDoubleKeyProtectionFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDoubleKeyProtectionFlagEnabled() bool {
	return false
}

// IsESDTFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTFlagEnabled() bool {
	return false
}

// IsESDTFlagEnabledForCurrentEpoch returns false
func (eeh *EnableEpochsHandler) IsESDTFlagEnabledForCurrentEpoch() bool {
	return false
}

// IsGovernanceFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsGovernanceFlagEnabled() bool {
	return false
}

// IsGovernanceFlagEnabledForCurrentEpoch returns false
func (eeh *EnableEpochsHandler) IsGovernanceFlagEnabledForCurrentEpoch() bool {
	return false
}

// IsDelegationManagerFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDelegationManagerFlagEnabled() bool {
	return false
}

// IsDelegationSmartContractFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDelegationSmartContractFlagEnabled() bool {
	return false
}

// IsDelegationSmartContractFlagEnabledForCurrentEpoch returns false
func (eeh *EnableEpochsHandler) IsDelegationSmartContractFlagEnabledForCurrentEpoch() bool {
	return false
}

// IsCorrectLastUnJailedFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCorrectLastUnJailedFlagEnabled() bool {
	return false
}

// IsCorrectLastUnJailedFlagEnabledForCurrentEpoch returns false
func (eeh *EnableEpochsHandler) IsCorrectLastUnJailedFlagEnabledForCurrentEpoch() bool {
	return false
}

// IsRelayedTransactionsV2FlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRelayedTransactionsV2FlagEnabled() bool {
	return false
}

// IsUnBondTokensV2FlagEnabled returns false
func (eeh *EnableEpochsHandler) IsUnBondTokensV2FlagEnabled() bool {
	return false
}

// IsSaveJailedAlwaysFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSaveJailedAlwaysFlagEnabled() bool {
	return false
}

// IsReDelegateBelowMinCheckFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsReDelegateBelowMinCheckFlagEnabled() bool {
	return false
}

// IsValidatorToDelegationFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsValidatorToDelegationFlagEnabled() bool {
	return false
}

// IsWaitingListFixFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsWaitingListFixFlagEnabled() bool {
	return false
}

// IsIncrementSCRNonceInMultiTransferFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsIncrementSCRNonceInMultiTransferFlagEnabled() bool {
	return false
}

// IsESDTMultiTransferFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTMultiTransferFlagEnabled() bool {
	return false
}

// IsGlobalMintBurnFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsGlobalMintBurnFlagEnabled() bool {
	return false
}

// IsESDTTransferRoleFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTTransferRoleFlagEnabled() bool {
	return false
}

// IsBuiltInFunctionOnMetaFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsBuiltInFunctionOnMetaFlagEnabled() bool {
	return false
}

// IsComputeRewardCheckpointFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsComputeRewardCheckpointFlagEnabled() bool {
	return false
}

// IsSCRSizeInvariantCheckFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSCRSizeInvariantCheckFlagEnabled() bool {
	return false
}

// IsBackwardCompSaveKeyValueFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsBackwardCompSaveKeyValueFlagEnabled() bool {
	return false
}

// IsESDTNFTCreateOnMultiShardFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTNFTCreateOnMultiShardFlagEnabled() bool {
	return false
}

// IsMetaESDTSetFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsMetaESDTSetFlagEnabled() bool {
	return false
}

// IsAddTokensToDelegationFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsAddTokensToDelegationFlagEnabled() bool {
	return false
}

// IsMultiESDTTransferFixOnCallBackFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsMultiESDTTransferFixOnCallBackFlagEnabled() bool {
	return false
}

// IsOptimizeGasUsedInCrossMiniBlocksFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsOptimizeGasUsedInCrossMiniBlocksFlagEnabled() bool {
	return false
}

// IsCorrectFirstQueuedFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCorrectFirstQueuedFlagEnabled() bool {
	return false
}

// IsDeleteDelegatorAfterClaimRewardsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDeleteDelegatorAfterClaimRewardsFlagEnabled() bool {
	return false
}

// IsFixOOGReturnCodeFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsFixOOGReturnCodeFlagEnabled() bool {
	return false
}

// IsRemoveNonUpdatedStorageFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRemoveNonUpdatedStorageFlagEnabled() bool {
	return false
}

// IsOptimizeNFTStoreFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsOptimizeNFTStoreFlagEnabled() bool {
	return false
}

// IsCreateNFTThroughExecByCallerFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCreateNFTThroughExecByCallerFlagEnabled() bool {
	return false
}

// IsStopDecreasingValidatorRatingWhenStuckFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsStopDecreasingValidatorRatingWhenStuckFlagEnabled() bool {
	return false
}

// IsFrontRunningProtectionFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsFrontRunningProtectionFlagEnabled() bool {
	return false
}

// IsPayableBySCFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsPayableBySCFlagEnabled() bool {
	return false
}

// IsCleanUpInformativeSCRsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCleanUpInformativeSCRsFlagEnabled() bool {
	return false
}

// IsStorageAPICostOptimizationFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsStorageAPICostOptimizationFlagEnabled() bool {
	return false
}

// IsESDTRegisterAndSetAllRolesFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTRegisterAndSetAllRolesFlagEnabled() bool {
	return false
}

// IsScheduledMiniBlocksFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsScheduledMiniBlocksFlagEnabled() bool {
	return false
}

// IsCorrectJailedNotUnStakedEmptyQueueFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCorrectJailedNotUnStakedEmptyQueueFlagEnabled() bool {
	return false
}

// IsDoNotReturnOldBlockInBlockchainHookFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDoNotReturnOldBlockInBlockchainHookFlagEnabled() bool {
	return false
}

// IsAddFailedRelayedTxToInvalidMBsFlag returns false
func (eeh *EnableEpochsHandler) IsAddFailedRelayedTxToInvalidMBsFlag() bool {
	return false
}

// IsSCRSizeInvariantOnBuiltInResultFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSCRSizeInvariantOnBuiltInResultFlagEnabled() bool {
	return false
}

// IsCheckCorrectTokenIDForTransferRoleFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCheckCorrectTokenIDForTransferRoleFlagEnabled() bool {
	return false
}

// IsFailExecutionOnEveryAPIErrorFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsFailExecutionOnEveryAPIErrorFlagEnabled() bool {
	return false
}

// IsMiniBlockPartialExecutionFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsMiniBlockPartialExecutionFlagEnabled() bool {
	return false
}

// IsManagedCryptoAPIsFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsManagedCryptoAPIsFlagEnabled() bool {
	return false
}

// IsESDTMetadataContinuousCleanupFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTMetadataContinuousCleanupFlagEnabled() bool {
	return false
}

// IsDisableExecByCallerFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsDisableExecByCallerFlagEnabled() bool {
	return false
}

// IsRefactorContextFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRefactorContextFlagEnabled() bool {
	return false
}

// IsCheckFunctionArgumentFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCheckFunctionArgumentFlagEnabled() bool {
	return false
}

// IsCheckExecuteOnReadOnlyFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCheckExecuteOnReadOnlyFlagEnabled() bool {
	return false
}

// IsFixAsyncCallbackCheckFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsFixAsyncCallbackCheckFlagEnabled() bool {
	return false
}

// IsSaveToSystemAccountFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSaveToSystemAccountFlagEnabled() bool {
	return false
}

// IsCheckFrozenCollectionFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCheckFrozenCollectionFlagEnabled() bool {
	return false
}

// IsSendAlwaysFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsSendAlwaysFlagEnabled() bool {
	return false
}

// IsValueLengthCheckFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsValueLengthCheckFlagEnabled() bool {
	return false
}

// IsCheckTransferFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsCheckTransferFlagEnabled() bool {
	return false
}

// IsTransferToMetaFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsTransferToMetaFlagEnabled() bool {
	return false
}

// IsESDTNFTImprovementV1FlagEnabled returns false
func (eeh *EnableEpochsHandler) IsESDTNFTImprovementV1FlagEnabled() bool {
	return false
}

// IsSetSenderInEeiOutputTransferFlagEnabled -
func (eeh *EnableEpochsHandler) IsSetSenderInEeiOutputTransferFlagEnabled() bool {
	return false
}

// IsChangeDelegationOwnerFlagEnabled -
func (eeh *EnableEpochsHandler) IsChangeDelegationOwnerFlagEnabled() bool {
	return false
}

// IsRefactorPeersMiniBlocksFlagEnabled returns false
func (eeh *EnableEpochsHandler) IsRefactorPeersMiniBlocksFlagEnabled() bool {
	return false
}

// IsFixAsyncCallBackArgsListFlagEnabled -
func (eeh *EnableEpochsHandler) IsFixAsyncCallBackArgsListFlagEnabled() bool {
	return false
}

// IsFixOldTokenLiquidityEnabled -
func (eeh *EnableEpochsHandler) IsFixOldTokenLiquidityEnabled() bool {
	return false
}

// IsMaxBlockchainHookCountersFlagEnabled -
func (eeh *EnableEpochsHandler) IsMaxBlockchainHookCountersFlagEnabled() bool {
	return false
}

// IsRuntimeMemStoreLimitEnabled -
func (eeh *EnableEpochsHandler) IsRuntimeMemStoreLimitEnabled() bool {
	return false
}

// IsAlwaysSaveTokenMetaDataEnabled -
func (eeh *EnableEpochsHandler) IsAlwaysSaveTokenMetaDataEnabled() bool {
	return false
}

// IsRuntimeCodeSizeFixEnabled -
func (eeh *EnableEpochsHandler) IsRuntimeCodeSizeFixEnabled() bool {
	return false
}

// IsRuntimeCodeSizeFixEnabled -
func (eeh *EnableEpochsHandler) IsRuntimeCodeSizeFixEnabled() bool {
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (eeh *EnableEpochsHandler) IsInterfaceNil() bool {
	return eeh == nil
}
