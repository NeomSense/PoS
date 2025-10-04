package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NeomSense/PoS/x/pos/types"
)

// Record queries a single record by ID
func (qs queryServer) Record(ctx context.Context, req *types.QueryRecordRequest) (*types.QueryRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "record id cannot be empty")
	}

	record, err := qs.k.GetRecord(ctx, req.Id)
	if err != nil {
		if types.ErrRecordNotFound.Is(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRecordResponse{Record: record}, nil
}

// Records queries all records with pagination
func (qs queryServer) Records(ctx context.Context, req *types.QueryRecordsRequest) (*types.QueryRecordsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var records []types.Record
	// Simple iteration without pagination for now
	err := qs.k.Records.Walk(ctx, nil, func(key string, value types.Record) (bool, error) {
		records = append(records, value)
		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRecordsResponse{
		Records:    records,
		Pagination: nil,
	}, nil
}

// ValidatorRecords queries all records for a specific validator
func (qs queryServer) ValidatorRecords(ctx context.Context, req *types.QueryValidatorRecordsRequest) (*types.QueryValidatorRecordsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	// Get all records for this validator
	var validatorRecords []types.Record
	err := qs.k.Records.Walk(ctx, nil, func(key string, record types.Record) (bool, error) {
		if record.ValidatorAddress == req.ValidatorAddress {
			validatorRecords = append(validatorRecords, record)
		}
		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Implement proper pagination for filtered results
	// For now, returning all records without pagination

	return &types.QueryValidatorRecordsResponse{
		Records:    validatorRecords,
		Pagination: nil,
	}, nil
}

// ValidatorStats queries statistics for a specific validator
func (qs queryServer) ValidatorStats(ctx context.Context, req *types.QueryValidatorStatsRequest) (*types.QueryValidatorStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	stats, err := qs.k.GetValidatorStats(ctx, req.ValidatorAddress)
	if err != nil {
		if err != nil && err.Error() == "collections: not found" || types.ErrValidatorStatsNotFound.Is(err) {
			// Return empty stats if not found
			stats = types.ValidatorRecordStats{
				ValidatorAddress: req.ValidatorAddress,
				TotalRecords:     0,
				VerifiedRecords:  0,
				RejectedRecords:  0,
				LastRecordTime:   0,
				IsEligible:       false,
			}
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryValidatorStatsResponse{Stats: stats}, nil
}
