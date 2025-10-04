package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(
	minRecordSize uint64,
	maxRecordSize uint64,
	recordsPerEpoch uint64,
	epochLength uint64,
	slashFractionMissingRecord math.LegacyDec,
	slashFractionInvalidRecord math.LegacyDec,
	minVerifiedRecordsForEligibility uint64,
) Params {
	return Params{
		MinRecordSize:                    minRecordSize,
		MaxRecordSize:                    maxRecordSize,
		RecordsPerEpoch:                  recordsPerEpoch,
		EpochLength:                      epochLength,
		SlashFractionMissingRecord:       slashFractionMissingRecord,
		SlashFractionInvalidRecord:       slashFractionInvalidRecord,
		MinVerifiedRecordsForEligibility: minVerifiedRecordsForEligibility,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		100,                            // MinRecordSize: 100 bytes
		1024*1024,                      // MaxRecordSize: 1MB
		10,                             // RecordsPerEpoch: 10 records per epoch
		100,                            // EpochLength: 100 blocks (~10 minutes with 6s blocks)
		math.LegacyNewDecWithPrec(1, 2), // SlashFractionMissingRecord: 0.01 (1%)
		math.LegacyNewDecWithPrec(5, 2), // SlashFractionInvalidRecord: 0.05 (5%)
		5,                              // MinVerifiedRecordsForEligibility: 5 verified records
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MinRecordSize == 0 {
		return fmt.Errorf("min record size must be positive")
	}
	if p.MaxRecordSize == 0 {
		return fmt.Errorf("max record size must be positive")
	}
	if p.MinRecordSize > p.MaxRecordSize {
		return fmt.Errorf("min record size cannot be greater than max record size")
	}
	if p.RecordsPerEpoch == 0 {
		return fmt.Errorf("records per epoch must be positive")
	}
	if p.EpochLength == 0 {
		return fmt.Errorf("epoch length must be positive")
	}
	if p.SlashFractionMissingRecord.IsNegative() || p.SlashFractionMissingRecord.GT(math.LegacyOneDec()) {
		return fmt.Errorf("slash fraction for missing record must be between 0 and 1")
	}
	if p.SlashFractionInvalidRecord.IsNegative() || p.SlashFractionInvalidRecord.GT(math.LegacyOneDec()) {
		return fmt.Errorf("slash fraction for invalid record must be between 0 and 1")
	}

	return nil
}
