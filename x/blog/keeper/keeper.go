package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/funmachine/pos/x/blog/types"
)

// Keeper defines the blog module's keeper
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	addressCodec address.Codec

	// Collections for storing params, posts, and post count
	Params  collections.Item[types.Params]
	Post    collections.Map[uint64, types.Post]
	PostSeq collections.Sequence

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
	Schema    collections.Schema
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
