package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/NeomSense/PoS/x/blog/types"
)

func (k msgServer) CreatePost(ctx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	// Validate creator bech32
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Get new ID from sequence
	id, err := k.PostSeq.Next(ctx)
	if err != nil {
		return nil, err
	}

	post := types.Post{
		Id:      id,
		Creator: msg.Creator,
		Title:   msg.Title,
		Body:    msg.Content, // Post has Body; Msg has Content
	}

	if err := k.Post.Set(ctx, id, post); err != nil {
		return nil, err
	}

	return &types.MsgCreatePostResponse{}, nil
}
