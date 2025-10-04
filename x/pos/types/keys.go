package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "pos"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"

	// Event types
	EventTypeRecordSubmitted = "record_submitted"
	EventTypeRecordVerified  = "record_verified"
	EventTypeValidatorSlashed = "validator_slashed"

	// Event attributes
	AttributeKeyRecordID    = "record_id"
	AttributeKeyValidator   = "validator"
	AttributeKeyVerifier    = "verifier"
	AttributeKeyApproved    = "approved"
	AttributeKeySlashAmount = "slash_amount"
	AttributeKeyReason      = "reason"
)

// Store key prefixes
var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("p_pos")

	// RecordsKey is the prefix for storing records
	RecordsKey = collections.NewPrefix("r_pos")

	// ValidatorStatsKey is the prefix for validator statistics
	ValidatorStatsKey = collections.NewPrefix("vs_pos")
)
