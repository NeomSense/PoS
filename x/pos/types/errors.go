package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/pos module sentinel errors
var (
	ErrInvalidSigner          = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrNotValidator           = errors.Register(ModuleName, 1101, "signer is not a validator")
	ErrInvalidRecordSize      = errors.Register(ModuleName, 1102, "record size is invalid")
	ErrRecordNotFound         = errors.Register(ModuleName, 1103, "record not found")
	ErrRecordAlreadyVerified  = errors.Register(ModuleName, 1104, "record already verified")
	ErrInvalidRecordStatus    = errors.Register(ModuleName, 1105, "invalid record status")
	ErrValidatorNotEligible   = errors.Register(ModuleName, 1106, "validator not eligible")
	ErrInsufficientRecords    = errors.Register(ModuleName, 1107, "insufficient verified records")
	ErrInvalidMerkleRoot      = errors.Register(ModuleName, 1108, "invalid merkle root")
	ErrDuplicateRecord        = errors.Register(ModuleName, 1109, "duplicate record submission")
	ErrEpochRecordsExceeded   = errors.Register(ModuleName, 1110, "epoch record limit exceeded")
	ErrValidatorStatsNotFound = errors.Register(ModuleName, 1111, "validator stats not found")
)
