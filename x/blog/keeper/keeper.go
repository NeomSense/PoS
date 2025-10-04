package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/NeomSense/PoS/x/blog/types"
)

// Keeper defines the blog module's keeper
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	addressCodec address.Codec

	// Collections for storing params, posts, and post count
	Params collections.Item[types.Params]
	Post   collections.Map[uint64, types.Post]
	PostSeq collections.Sequence

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
	Schema    collections.Schema
}

// NewKeeper creates a new blog keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	addressCodec address.Codec,
	authority string,
) (Keeper, error) {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		addressCodec: addressCodec,
		authority:    authority,
		Params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),
		Post: collections.NewMap(
			sb,
			types.PostKey,
			"post",
			collections.Uint64Key,
			codec.CollValue[types.Post](cdc),
		),
		PostSeq: collections.NewSequence(
			sb,
			types.PostCountKey,
			"post_seq",
		),
	}

	schema, err := sb.Build()
	if err != nil {
		return Keeper{}, err
	}
	k.Schema = schema

	return k, nil
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx any) any {
	return fmt.Sprintf("blog keeper logger: %v", ctx)
}
