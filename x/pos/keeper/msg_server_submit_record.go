package keeper

import (
	"context"

	"github.com/NeomSense/PoS/x/pos/types"
)

// SubmitRecord handles the MsgSubmitRecord message
func (ms msgServer) SubmitRecord(ctx context.Context, msg *types.MsgSubmitRecord) (*types.MsgSubmitRecordResponse, error) {
	// Create the record
	recordID, err := ms.k.CreateRecord(ctx, msg.ValidatorAddress, msg.Data, msg.MerkleRoot)
	if err != nil {
		return nil, err
	}

	// Get the timestamp
	record, err := ms.k.GetRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitRecordResponse{
		RecordId:  recordID,
		Timestamp: record.Timestamp,
	}, nil
}
