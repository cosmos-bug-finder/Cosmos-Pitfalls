package keeper

import (
	"context"
	"strconv"

	"github.com/onomyprotocol/haven/x/haven/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) TipPost(goCtx context.Context, msg *types.MsgTipPost) (*types.MsgTipPostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	uid, _ := strconv.ParseUint(msg.Uid, 10, 64)

	tipAmount, _ := sdk.NewIntFromString(msg.Amount)

	tip := sdk.NewCoin("kudos", tipAmount)

	// Test if post exists
	post, found := k.GetPost(ctx, uid)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrPostDoesNotExist, "%s", msg.Uid)
	}

	post.Tips = post.Tips.Add(tip)

	k.SetPost(ctx, post)

	haven, _ := k.GetHaven(ctx, post.Haven)

	havenTip := (tip.Amount.Mul(sdk.NewIntFromUint64(haven.Rake))).Quo(sdk.NewInt(100))
	postTip := tip.Amount.Sub(havenTip)

	haven.Earnings = haven.Earnings.AddAmount(havenTip)

	k.SetHaven(ctx, haven)

	havenCoins := sdk.NewCoin("kudos", havenTip)
	postCoins := sdk.NewCoin("kudos", postTip)

	// moduleAcc := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	havenOwner, _ := sdk.AccAddressFromBech32(haven.Owner)
	postOwner, _ := sdk.AccAddressFromBech32(post.Owner)

	k.bankKeeper.SendCoins(ctx, creator, havenOwner, sdk.NewCoins(havenCoins))
	k.bankKeeper.SendCoins(ctx, creator, postOwner, sdk.NewCoins(postCoins))

	return &types.MsgTipPostResponse{}, nil
}