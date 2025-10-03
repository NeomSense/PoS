package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/you/pos/x/blog/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	// Validate creator bech32
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get new ID from counter and store the post
	id := k.GetPostCount(ctx)

	post := types.Post{
		Id:      id,
		Creator: msg.Creator,
		Title:   msg.Title,
		Body:    msg.Content, // Post has Body; Msg has Content
	}

	k.SetPost(ctx, post)
	k.SetPostCount(ctx, id+1)

	return &types.MsgCreatePostResponse{}, nil
}
