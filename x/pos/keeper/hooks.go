package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Hooks wrapper struct for keeper
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterValidatorCreated - Initialize validator stats when a validator is created
func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	return h.k.InitializeValidatorStats(ctx, valAddr.String())
}

// AfterValidatorRemoved - Clean up validator stats when removed (optional)
func (h Hooks) AfterValidatorRemoved(ctx context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	// Optionally, you could archive or remove validator stats here
	// For now, we'll keep the stats for historical purposes
	return nil
}

// BeforeDelegationCreated - called before a delegation is created
func (h Hooks) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// BeforeDelegationSharesModified - called before a delegation's shares are modified
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// BeforeDelegationRemoved - called before a delegation is removed
func (h Hooks) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// BeforeValidatorModified - called before a validator is modified
func (h Hooks) BeforeValidatorModified(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

// AfterDelegationModified - called after a delegation is modified
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// BeforeValidatorSlashed - called before a validator is slashed
func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	return nil
}

// AfterValidatorBeginUnbonding - called after a validator begins unbonding
func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// AfterValidatorBonded - called after a validator is bonded
func (h Hooks) AfterValidatorBonded(ctx context.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	// When a validator becomes bonded, ensure they have stats initialized
	return h.k.InitializeValidatorStats(ctx, valAddr.String())
}

// AfterUnbondingInitiated - called after unbonding has been initiated
func (h Hooks) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}

// BeforeConsensusPubKeyRotated - called before consensus key rotation
func (h Hooks) BeforeConsensusPubKeyRotated(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

// AfterConsensusPubKeyRotated - called after consensus key rotation
func (h Hooks) AfterConsensusPubKeyRotated(ctx context.Context, oldPubKey, newPubKey []byte) error {
	return nil
}
