package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NeomSense/PoS/x/pos/types"
)

// CreateRecord creates a new record submitted by a validator
func (k Keeper) CreateRecord(
	ctx context.Context,
	validatorAddr string,
	data []byte,
	merkleRoot string,
) (string, error) {
	// Get params for validation
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}

	// Validate record size
	dataSize := uint64(len(data))
	if dataSize < params.MinRecordSize || dataSize > params.MaxRecordSize {
		return "", types.ErrInvalidRecordSize.Wrapf(
			"record size %d is not within bounds [%d, %d]",
			dataSize,
			params.MinRecordSize,
			params.MaxRecordSize,
		)
	}

	// Validate merkle root
	if len(merkleRoot) == 0 {
		return "", types.ErrInvalidMerkleRoot.Wrap("merkle root cannot be empty")
	}

	// Get validator to ensure they exist
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return "", err
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return "", types.ErrNotValidator.Wrap(err.Error())
	}

	// Check if validator is bonded
	if !validator.IsBonded() {
		return "", types.ErrNotValidator.Wrap("validator must be bonded to submit records")
	}

	// Generate record ID from hash of validator + data + timestamp
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	timestamp := sdkCtx.BlockTime().Unix()
	blockHeight := uint64(sdkCtx.BlockHeight())

	recordID := generateRecordID(validatorAddr, data, timestamp)

	// Check for duplicate
	if has, err := k.Records.Has(ctx, recordID); err != nil {
		return "", err
	} else if has {
		return "", types.ErrDuplicateRecord.Wrapf("record %s already exists", recordID)
	}

	// Check epoch record limits
	stats, err := k.GetValidatorStats(ctx, validatorAddr)
	if err != nil {
		// If stats don't exist, create them
		stats = types.ValidatorRecordStats{
			ValidatorAddress:        validatorAddr,
			TotalRecords:           0,
			VerifiedRecords:        0,
			RejectedRecords:        0,
			LastRecordTime:         0,
			IsEligible:             true,
			NextRequiredRecordTime: timestamp + int64(params.EpochLength),
		}
	}

	// Calculate current epoch
	currentEpoch := blockHeight / params.EpochLength
	lastRecordEpoch := uint64(stats.LastRecordTime) / params.EpochLength

	// If same epoch, check record limit
	if currentEpoch == lastRecordEpoch {
		// Count records in current epoch
		epochRecordCount := uint64(0)
		err := k.Records.Walk(ctx, nil, func(key string, record types.Record) (bool, error) {
			if record.ValidatorAddress == validatorAddr {
				recordEpoch := record.BlockHeight / params.EpochLength
				if recordEpoch == currentEpoch {
					epochRecordCount++
				}
			}
			return false, nil
		})
		if err != nil {
			return "", err
		}

		if epochRecordCount >= params.RecordsPerEpoch {
			return "", types.ErrEpochRecordsExceeded.Wrapf(
				"validator has already submitted %d records in epoch %d",
				epochRecordCount,
				currentEpoch,
			)
		}
	}

	// Create record
	record := types.Record{
		Id:               recordID,
		ValidatorAddress: validatorAddr,
		Data:             data,
		Timestamp:        timestamp,
		Status:           types.RecordStatusPending,
		MerkleRoot:       merkleRoot,
		BlockHeight:      blockHeight,
	}

	// Store record
	if err := k.Records.Set(ctx, recordID, record); err != nil {
		return "", err
	}

	// Update validator stats
	stats.TotalRecords++
	stats.LastRecordTime = timestamp
	if err := k.SetValidatorStats(ctx, validatorAddr, stats); err != nil {
		return "", err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"record_created",
			sdk.NewAttribute("record_id", recordID),
			sdk.NewAttribute("validator", validatorAddr),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", blockHeight)),
			sdk.NewAttribute("merkle_root", merkleRoot),
		),
	)

	return recordID, nil
}

// GetRecord retrieves a record by ID
func (k Keeper) GetRecord(ctx context.Context, recordID string) (types.Record, error) {
	record, err := k.Records.Get(ctx, recordID)
	if err != nil {
		if err != nil && err.Error() == "collections: not found" {
			return types.Record{}, types.ErrRecordNotFound.Wrapf("record %s not found", recordID)
		}
		return types.Record{}, err
	}
	return record, nil
}

// VerifyRecord verifies or rejects a submitted record
func (k Keeper) VerifyRecord(ctx context.Context, recordID string, approved bool) error {
	// Get the record
	record, err := k.GetRecord(ctx, recordID)
	if err != nil {
		return err
	}

	// Check if already verified
	if record.Status != types.RecordStatusPending {
		return types.ErrRecordAlreadyVerified.Wrapf(
			"record %s is already %s",
			recordID,
			record.Status.String(),
		)
	}

	// Update status
	if approved {
		record.Status = types.RecordStatusVerified
	} else {
		record.Status = types.RecordStatusRejected
	}

	// Save record
	if err := k.Records.Set(ctx, recordID, record); err != nil {
		return err
	}

	// Update validator stats
	stats, err := k.GetValidatorStats(ctx, record.ValidatorAddress)
	if err != nil {
		return err
	}

	if approved {
		stats.VerifiedRecords++
	} else {
		stats.RejectedRecords++
	}

	// Check eligibility based on verified records
	params, _ := k.Params.Get(ctx)
	stats.IsEligible = stats.VerifiedRecords >= params.MinVerifiedRecordsForEligibility

	if err := k.SetValidatorStats(ctx, record.ValidatorAddress, stats); err != nil {
		return err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"record_verified",
			sdk.NewAttribute("record_id", recordID),
			sdk.NewAttribute("validator", record.ValidatorAddress),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", approved)),
			sdk.NewAttribute("status", record.Status.String()),
		),
	)

	return nil
}

// GetAllRecords returns all records with pagination
func (k Keeper) GetAllRecords(ctx context.Context) ([]types.Record, error) {
	var records []types.Record
	err := k.Records.Walk(ctx, nil, func(key string, record types.Record) (bool, error) {
		records = append(records, record)
		return false, nil
	})
	return records, err
}

// GetValidatorRecords returns all records for a specific validator
func (k Keeper) GetValidatorRecords(ctx context.Context, validatorAddr string) ([]types.Record, error) {
	var records []types.Record
	err := k.Records.Walk(ctx, nil, func(key string, record types.Record) (bool, error) {
		if record.ValidatorAddress == validatorAddr {
			records = append(records, record)
		}
		return false, nil
	})
	return records, err
}

// generateRecordID generates a unique ID for a record
func generateRecordID(validatorAddr string, data []byte, timestamp int64) string {
	hasher := sha256.New()
	hasher.Write([]byte(validatorAddr))
	hasher.Write(data)
	hasher.Write([]byte(fmt.Sprintf("%d", timestamp)))
	return hex.EncodeToString(hasher.Sum(nil))
}
