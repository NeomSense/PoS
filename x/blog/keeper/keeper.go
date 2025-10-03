package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"

	"github.com/you/pos/x/blog/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          sdkcodec.Codec
	addressCodec address.Codec

	// Address capable of executing a MsgUpdateParams message (usually the gov module account).
	authority string

	Schema  collections.Schema
	Params  collections.Item[types.Params]
	PostSeq collections.Sequence
	Post    collections.Map[uint64, types.Post]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc sdkcodec.Codec,
	addressCodec address.Codec,
	authority string,
) Keeper {
	// Validate authority is a proper bech32 address for the configured address codec.
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		// NOTE: CollValue is provided by github.com/cosmos/cosmos-sdk/codec in your setup.
		Params:  collections.NewItem(sb, types.ParamsKey, "params", sdkcodec.CollValue[types.Params](cdc)),
		Post:    collections.NewMap(sb, types.PostKey, "post", collections.Uint64Key, sdkcodec.CollValue[types.Post](cdc)),
		PostSeq: collections.NewSequence(sb, types.PostCountKey, "postSequence"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority (bech32 string).
func (k Keeper) GetAuthority() string {
	return k.authority
}
