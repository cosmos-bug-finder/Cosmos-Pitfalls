package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stateset/core/x/purchaseorder/types"
)

func (k msgServer) FinancePurchaseorder(goCtx context.Context, msg *types.MsgFinancePurchaseorder) (*types.MsgFinancePurchaseorderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	purchaseorder, found := k.GetPurchaseorder(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	if purchaseorder.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongPurchaseOrderState, "%v", purchaseorder.State)
	}

	financer, _ := sdk.AccAddressFromBech32(msg.Creator)
	seller, _ := sdk.AccAddressFromBech32(purchaseorder.Seller)
	amount, _ := sdk.ParseCoinsNormalized(purchaseorder.Amount)

	k.bankKeeper.SendCoins(ctx, financer, seller, amount)

	purchaseorder.Financer = msg.Creator
	purchaseorder.State = "financed"

	k.SetPurchaseorder(ctx, purchaseorder)

	return &types.MsgFinancePurchaseorderResponse{}, nil
}