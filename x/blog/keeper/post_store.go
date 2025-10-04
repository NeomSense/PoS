package keeper

import (
	"context"

	"github.com/NeomSense/PoS/x/blog/types"
)

// ===== Post count (total) =====

func (k Keeper) GetPostCount(ctx context.Context) uint64 {
	count, err := k.PostSeq.Peek(ctx)
	if err != nil {
		return 0
	}
	return count
}

func (k Keeper) SetPostCount(ctx context.Context, count uint64) {
	// With collections.Sequence, we don't manually set count
	// The sequence is auto-incremented via Next()
	// This method is kept for compatibility but does nothing
}

// ===== CRUD for Post =====

func (k Keeper) SetPost(ctx context.Context, post types.Post) {
	_ = k.Post.Set(ctx, post.Id, post)
}

func (k Keeper) GetPost(ctx context.Context, id uint64) (val types.Post, found bool) {
	post, err := k.Post.Get(ctx, id)
	if err != nil {
		return types.Post{}, false
	}
	return post, true
}

func (k Keeper) RemovePost(ctx context.Context, id uint64) {
	_ = k.Post.Remove(ctx, id)
}
