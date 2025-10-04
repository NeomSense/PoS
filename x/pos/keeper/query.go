package keeper

import (
	"github.com/NeomSense/PoS/x/pos/types"
)

var _ types.QueryServer = &queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

type queryServer struct {
	types.UnimplementedQueryServer
	k Keeper
}
