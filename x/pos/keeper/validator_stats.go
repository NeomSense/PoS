package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NeomSense/PoS/x/pos/types"
)

// GetValidatorStats retrieves statistics for a validator
func (k Keeper) GetValidatorStats(ctx context.Context, validatorAddr string) (types.ValidatorRecordStats, error) {
	stats, err := k.ValidatorStats.Get(ctx, validatorAddr)
	if err != nil {
		if err != nil && err.Error() == "collections: not found" {
			return types.ValidatorRecordStats{}, types.ErrValidatorStatsNotFound.Wrapf(
				"stats for validator %s not found",
				validatorAddr,
			)
		}
		return types.ValidatorRecordStats{}, err
	}
	return stats, nil
}

// SetValidatorStats stores validator statistics
func (k Keeper) SetValidatorStats(ctx context.Context, validatorAddr string, stats types.ValidatorRecordStats) error {
	return k.ValidatorStats.Set(ctx, validatorAddr, stats)
}

// GetAllValidatorStats returns statistics for all validators
func (k Keeper) GetAllValidatorStats(ctx context.Context) ([]types.ValidatorRecordStats, error) {
	var allStats []types.ValidatorRecordStats
	err := k.ValidatorStats.Walk(ctx, nil, func(key string, stats types.ValidatorRecordStats) (bool, error) {
		allStats = append(allStats, stats)
		return false, nil
	})
	return allStats, err
}

// InitializeValidatorStats creates initial stats for a new validator
func (k Keeper) InitializeValidatorStats(ctx context.Context, validatorAddr string) error {
	// Check if stats already exist
	if has, err := k.ValidatorStats.Has(ctx, validatorAddr); err != nil {
		return err
	} else if has {
		return nil // Already initialized
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	stats := types.ValidatorRecordStats{
		ValidatorAddress:        validatorAddr,
		TotalRecords:           0,
		VerifiedRecords:        0,
		RejectedRecords:        0,
		LastRecordTime:         0,
		IsEligible:             true, // Start as eligible
		NextRequiredRecordTime: currentTime + int64(params.EpochLength),
	}

	return k.SetValidatorStats(ctx, validatorAddr, stats)
}

// CheckValidatorEligibility checks if a validator meets record requirements
func (k Keeper) CheckValidatorEligibility(ctx context.Context, validatorAddr string) (bool, error) {
	stats, err := k.GetValidatorStats(ctx, validatorAddr)
	if err != nil {
		if types.ErrValidatorStatsNotFound.Is(err) {
			// No stats = not eligible
			return false, nil
		}
		return false, err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	// Check if validator has minimum verified records
	if stats.VerifiedRecords < params.MinVerifiedRecordsForEligibility {
		return false, nil
	}

	// Check if validator submitted record in current epoch
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()
	blockHeight := uint64(sdkCtx.BlockHeight())
	currentEpoch := blockHeight / params.EpochLength

	// If we're past the required time and no recent record, not eligible
	if currentTime > stats.NextRequiredRecordTime {
		lastRecordEpoch := uint64(stats.LastRecordTime) / params.EpochLength
		if lastRecordEpoch < currentEpoch {
			return false, nil
		}
	}

	return stats.IsEligible, nil
}

// UpdateValidatorEligibility updates a validator's eligibility status
func (k Keeper) UpdateValidatorEligibility(ctx context.Context, validatorAddr string, eligible bool) error {
	stats, err := k.GetValidatorStats(ctx, validatorAddr)
	if err != nil {
		return err
	}

	stats.IsEligible = eligible
	return k.SetValidatorStats(ctx, validatorAddr, stats)
}

// SlashValidatorForMissingRecords slashes a validator for not submitting required records
func (k Keeper) SlashValidatorForMissingRecords(ctx context.Context, validatorAddr string) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Get validator
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return err
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	// Get consensus address for slashing
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Slash the validator
	power := validator.GetConsensusPower(sdk.DefaultPowerReduction)
	_, err = k.stakingKeeper.Slash(
		ctx,
		consAddr,
		sdkCtx.BlockHeight(),
		power,
		params.SlashFractionMissingRecord,
	)
	if err != nil {
		return err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_slashed",
			sdk.NewAttribute("validator", validatorAddr),
			sdk.NewAttribute("reason", "missing_records"),
			sdk.NewAttribute("slash_fraction", params.SlashFractionMissingRecord.String()),
		),
	)

	return nil
}

// SlashValidatorForInvalidRecord slashes a validator for submitting invalid record
func (k Keeper) SlashValidatorForInvalidRecord(ctx context.Context, validatorAddr string, recordID string) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Get validator
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return err
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	// Get consensus address for slashing
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Slash the validator
	power := validator.GetConsensusPower(sdk.DefaultPowerReduction)
	_, err = k.stakingKeeper.Slash(
		ctx,
		consAddr,
		sdkCtx.BlockHeight(),
		power,
		params.SlashFractionInvalidRecord,
	)
	if err != nil {
		return err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_slashed",
			sdk.NewAttribute("validator", validatorAddr),
			sdk.NewAttribute("reason", "invalid_record"),
			sdk.NewAttribute("record_id", recordID),
			sdk.NewAttribute("slash_fraction", params.SlashFractionInvalidRecord.String()),
		),
	)

	return nil
}

// CheckAllValidatorsEligibility checks eligibility for all validators and slashes if needed
func (k Keeper) CheckAllValidatorsEligibility(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := uint64(sdkCtx.BlockHeight())

	// Only check at epoch boundaries
	if blockHeight%params.EpochLength != 0 {
		return nil
	}

	// Get all validators from staking module
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, validator := range validators {
		// Skip if not bonded
		if !validator.IsBonded() {
			continue
		}

		validatorAddr := validator.GetOperator()

		// Check eligibility
		eligible, err := k.CheckValidatorEligibility(ctx, validatorAddr)
		if err != nil {
			// Log error but continue with other validators
			sdkCtx.Logger().Error(
				"failed to check validator eligibility",
				"validator", validatorAddr,
				"error", err,
			)
			continue
		}

		// If not eligible, slash and update status
		if !eligible {
			if err := k.SlashValidatorForMissingRecords(ctx, validatorAddr); err != nil {
				sdkCtx.Logger().Error(
					"failed to slash validator",
					"validator", validatorAddr,
					"error", err,
				)
			}

			// Update eligibility status
			if err := k.UpdateValidatorEligibility(ctx, validatorAddr, false); err != nil {
				sdkCtx.Logger().Error(
					"failed to update validator eligibility",
					"validator", validatorAddr,
					"error", err,
				)
			}

			sdkCtx.Logger().Info(
				"validator marked ineligible due to insufficient records",
				"validator", validatorAddr,
				"epoch", blockHeight/params.EpochLength,
			)
		} else {
			// Update next required record time
			stats, err := k.GetValidatorStats(ctx, validatorAddr)
			if err == nil {
				stats.NextRequiredRecordTime = sdkCtx.BlockTime().Unix() + int64(params.EpochLength)
				_ = k.SetValidatorStats(ctx, validatorAddr, stats)
			}
		}
	}

	return nil
}
