package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/runtime"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/you/pos/x/blog/types"
)

// ===== Post count (total) =====

func (k Keeper) GetPostCount(ctx sdk.Context) uint64 {
	kv := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := kv.Get(types.KeyPrefix(types.PostCountKey))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetPostCount(ctx sdk.Context, count uint64) {
	kv := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	kv.Set(types.KeyPrefix(types.PostCountKey), bz)
}

// ===== CRUD for Post =====

func (k Keeper) SetPost(ctx sdk.Context, post types.Post) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PostKeyPrefix))
	b := k.cdc.MustMarshal(&post)
	store.Set(postIDToBz(post.Id), b)
}

func (k Keeper) GetPost(ctx sdk.Context, id uint64) (val types.Post, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PostKeyPrefix))
	b := store.Get(postIDToBz(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePost(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PostKeyPrefix))
	store.Delete(postIDToBz(id))
}

func postIDToBz(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}
