package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gauss/gauss/v6/x/auction/types"
)

func (k msgServer) BidOrder(goCtx context.Context, msg *types.MsgBidOrder) (*types.MsgBidOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	order, err := k.getOrder(ctx, msg.PoolAddress, msg.TokenId)
	if err != nil {
		return nil, err
	}
	// check expired
	err = k.checkOrder(ctx, *order)
	if err != nil {
		return nil, err
	}
	// check balance
	err = k.checkBalance(ctx, msg.Sender, *msg.Price)
	if err != nil {
		return nil, err
	}
	// check price
	if msg.Price.Denom != order.Price.Denom {
		return nil, types.ErrDenomMatch
	}
	bidPrice := msg.Price
	var currentPrice, oldBidderPrice sdk.Coin
	if order.Bidder == "" {
		currentPrice = *order.Price
		oldBidderPrice = sdk.NewCoin(order.Price.Denom, sdk.NewInt(0))
	} else {
		currentPrice = order.BidPrice.Add(order.MinStep)
		oldBidderPrice = order.BidPrice
	}

	if bidPrice.IsLT(currentPrice) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBidPrice, "bid price(%s) should greater than %s", bidPrice, currentPrice)
	}

	// 更新订单
	oldBidder := order.Bidder

	order.BidPrice = *bidPrice
	order.Bidder = msg.Sender
	order.MinEndTime = ctx.BlockTime().Add(order.AutoAgreePeriod)
	err = k.saveOrder(ctx, *order)
	if err != nil {
		return nil, err
	}
	// 更新账户余额
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{*bidPrice})
	// 退回上一个bidder的token
	if oldBidderPrice.Amount.IsPositive() {
		oldBidderAddr, _ := sdk.AccAddressFromBech32(oldBidder)
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, oldBidderAddr, sdk.Coins{oldBidderPrice})
	}

	return &types.MsgBidOrderResponse{}, nil
}

// check balance
func (k Keeper) checkBalance(ctx sdk.Context, senderAddr string, bidPrice sdk.Coin) error {
	sender, _ := sdk.AccAddressFromBech32(senderAddr)
	balance := k.bankKeeper.GetBalance(ctx, sender, bidPrice.Denom)
	if balance.IsLT(bidPrice) {
		return sdkerrors.Wrapf(types.ErrInsufficientFunds, "balance is smaller than %s", bidPrice)
	}
	return nil
}