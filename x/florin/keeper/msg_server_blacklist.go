package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

var _ blacklist.MsgServer = &blacklistMsgServer{}

type blacklistMsgServer struct {
	*Keeper
}

func NewBlacklistMsgServer(keeper *Keeper) blacklist.MsgServer {
	return &blacklistMsgServer{Keeper: keeper}
}

func (k blacklistMsgServer) AcceptOwnership(goCtx context.Context, msg *blacklist.MsgAcceptOwnership) (*blacklist.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingOwner := k.GetBlacklistPendingOwner(ctx)
	if pendingOwner == "" {
		return nil, blacklist.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(blacklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetBlacklistOwner(ctx)
	k.SetBlacklistOwner(ctx, msg.Signer)
	k.DeleteBlacklistPendingOwner(ctx)

	return &blacklist.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k blacklistMsgServer) AddAdminAccount(goCtx context.Context, msg *blacklist.MsgAddAdminAccount) (*blacklist.MsgAddAdminAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.SetBlacklistAdmin(ctx, msg.Account)

	return &blacklist.MsgAddAdminAccountResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.AdminAccountAdded{
		Account: msg.Account,
	})
}

func (k blacklistMsgServer) Ban(goCtx context.Context, msg *blacklist.MsgBan) (*blacklist.MsgBanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsBlacklistAdmin(ctx, msg.Signer) {
		return nil, blacklist.ErrInvalidAdmin
	}

	k.SetAdversary(ctx, msg.Adversary)

	return &blacklist.MsgBanResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.Ban{
		Adversary: msg.Adversary,
	})
}

func (k blacklistMsgServer) RemoveAdminAccount(goCtx context.Context, msg *blacklist.MsgRemoveAdminAccount) (*blacklist.MsgRemoveAdminAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.DeleteBlacklistAdmin(ctx, msg.Account)

	return &blacklist.MsgRemoveAdminAccountResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.AdminAccountRemoved{
		Account: msg.Account,
	})
}

func (k blacklistMsgServer) TransferOwnership(goCtx context.Context, msg *blacklist.MsgTransferOwnership) (*blacklist.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	if msg.NewOwner == owner {
		return nil, blacklist.ErrSameOwner
	}

	k.SetBlacklistPendingOwner(ctx, msg.NewOwner)

	return &blacklist.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k blacklistMsgServer) Unban(goCtx context.Context, msg *blacklist.MsgUnban) (*blacklist.MsgUnbanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsBlacklistAdmin(ctx, msg.Signer) {
		return nil, blacklist.ErrInvalidAdmin
	}

	k.DeleteAdversary(ctx, msg.Friend)

	return &blacklist.MsgUnbanResponse{}, ctx.EventManager().EmitTypedEvent(&blacklist.Unban{
		Friend: msg.Friend,
	})
}

//

func (k blacklistMsgServer) EnsureOwner(ctx sdk.Context, signer string) (string, error) {
	owner := k.GetBlacklistOwner(ctx)
	if owner == "" {
		return "", blacklist.ErrNoOwner
	}
	if signer != owner {
		return "", errors.Wrapf(blacklist.ErrInvalidOwner, "expected %s, got %s", owner, signer)
	}
	return owner, nil
}
