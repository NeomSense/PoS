\package keeper

import (
    "github.com/NeomSense/PoS/blob/main/x/blog/keeper/keeper.go"
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/runtime"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "cosmossdk.io/collections"
)

// Keeper defines the blog module's keeper
type Keeper struct {
    cdc      codec.BinaryCodec
    storeKey sdk.StoreKey
    memStore sdk.KVStoreKey
    // Collections for storing params, posts, and post count
    Params    collections.Item[types.Params]
    Posts     collections.Map[string, types.Post]
    PostCount collections.Sequence
}

// NewKeeper creates a new blog keeper
func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey sdk.StoreKey,
    memStore sdk.KVStoreKey,
) Keeper {
    return Keeper{
        cdc:      cdc,
        storeKey: storeKey,
        memStore: memStore,
        Params: collections.NewItem(
            runtime.NewKVStoreService(storeKey),
            types.ParamsKey,
            "params",
            codec.CollValue[types.Params](cdc),
        ),
        Posts: collections.NewMap(
            runtime.NewKVStoreService(storeKey),
            types.PostKey,
            "posts",
            collections.StringKey,
            codec.CollValue[types.Post](cdc),
        ),
        PostCount: collections.NewSequence(
            runtime.NewKVStoreService(storeKey),
            types.PostCountKey,
        ),
    }
}
