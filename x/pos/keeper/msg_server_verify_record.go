package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NeomSense/PoS/x/pos/types"
)

// VerifyRecord handles the MsgVerifyRecord message
func (ms msgServer) VerifyRecord(ctx context.Context, msg *types.MsgVerifyRecord) (*types.MsgVerifyRecordResponse, error) {
	// Authorization: Only bonded validators can verify records
	verifierAddr, err := sdk.ValAddressFromBech32(msg.Verifier)
	if err != nil {
		// Not a validator address, return error
		return nil, fmt.Errorf("invalid verifier address: %w", err)
	}

	validator, err := ms.k.stakingKeeper.GetValidator(ctx, verifierAddr)
	if err != nil {
		return nil, fmt.Errorf("verifier is not a validator: %w", err)
	}

	if !validator.IsBonded() {
		return nil, fmt.Errorf("only bonded validators can verify records")
	}

	// Verify the record
	err = ms.k.VerifyRecord(ctx, msg.RecordId, msg.Approved)
	if err != nil {
		return nil, err
	}

	// If rejected, optionally slash the validator for invalid record
	if !msg.Approved {
		record, err := ms.k.GetRecord(ctx, msg.RecordId)
		if err != nil {
			return nil, err
		}

		// Slash for invalid record
		if err := ms.k.SlashValidatorForInvalidRecord(ctx, record.ValidatorAddress, msg.RecordId); err != nil {
			// Log error but don't fail the transaction
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Error(
				"failed to slash validator for invalid record",
				"validator", record.ValidatorAddress,
				"record_id", msg.RecordId,
				"error", err,
			)
		}
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecordVerified,
			sdk.NewAttribute(types.AttributeKeyRecordID, msg.RecordId),
			sdk.NewAttribute(types.AttributeKeyVerifier, msg.Verifier),
			sdk.NewAttribute(types.AttributeKeyApproved, fmt.Sprintf("%t", msg.Approved)),
		),
	)

	return &types.MsgVerifyRecordResponse{}, nil
}
