// Copyright 2024 Monerium ehf.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keeper

import (
	"context"
	"fmt"

	"adr36.dev"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) AcceptOwnership(goCtx context.Context, msg *types.MsgAcceptOwnership) (*types.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	pendingOwner := k.GetPendingOwner(ctx, msg.Denom)
	if pendingOwner == "" {
		return nil, types.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(types.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetOwner(ctx, msg.Denom)
	k.SetOwner(ctx, msg.Denom, msg.Signer)
	k.DeletePendingOwner(ctx, msg.Denom)

	return &types.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferred{
		Denom:         msg.Denom,
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k msgServer) AddAdminAccount(goCtx context.Context, msg *types.MsgAddAdminAccount) (*types.MsgAddAdminAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.SetAdmin(ctx, msg.Denom, msg.Account)

	return &types.MsgAddAdminAccountResponse{}, ctx.EventManager().EmitTypedEvent(&types.AdminAccountAdded{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) AddSystemAccount(goCtx context.Context, msg *types.MsgAddSystemAccount) (*types.MsgAddSystemAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.SetSystem(ctx, msg.Denom, msg.Account)

	return &types.MsgAddSystemAccountResponse{}, ctx.EventManager().EmitTypedEvent(&types.SystemAccountAdded{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) AllowDenom(goCtx context.Context, msg *types.MsgAllowDenom) (*types.MsgAllowDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Signer != k.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	supply := k.bankKeeper.GetSupply(ctx, msg.Denom)
	if !supply.IsZero() {
		return nil, types.ErrInvalidDenom
	}

	k.SetAllowedDenom(ctx, msg.Denom)
	k.SetOwner(ctx, msg.Denom, msg.Owner)

	return &types.MsgAllowDenomResponse{}, ctx.EventManager().EmitTypedEvent(&types.DenomAllowed{
		Denom: msg.Denom,
		Owner: msg.Owner,
	})
}

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsSystem(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidSystem
	}

	address := sdk.MustAccAddressFromBech32(msg.From)
	account := k.accountKeeper.GetAccount(ctx, address)
	if account == nil || account.GetPubKey() == nil {
		return nil, types.ErrNoPubKey
	}

	if !adr36.VerifySignature(
		account.GetPubKey(),
		[]byte("I hereby declare that I am the address owner."),
		msg.Signature,
	) {
		return nil, types.ErrInvalidSignature
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	from := sdk.MustAccAddressFromBech32(msg.From)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from user to module")
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to burn from module")
	}

	return &types.MsgBurnResponse{}, nil
}

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsSystem(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidSystem
	}

	allowance := k.GetMintAllowance(ctx, msg.Denom, msg.Signer)
	if msg.Amount.GT(allowance) {
		return nil, types.ErrInsufficientAllowance
	}

	allowance = allowance.Sub(msg.Amount)
	k.Keeper.SetMintAllowance(ctx, msg.Denom, msg.Signer, allowance)

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	to := sdk.MustAccAddressFromBech32(msg.To)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to mint to module")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from module to user")
	}

	return &types.MsgMintResponse{}, ctx.EventManager().EmitTypedEvent(&types.MintAllowance{
		Denom:   msg.Denom,
		Account: msg.Signer,
		Amount:  allowance,
	})
}

func (k msgServer) Recover(goCtx context.Context, msg *types.MsgRecover) (*types.MsgRecoverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsSystem(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidSystem
	}

	from := sdk.MustAccAddressFromBech32(msg.From)
	account := k.accountKeeper.GetAccount(ctx, from)
	if account == nil || account.GetPubKey() == nil {
		return nil, types.ErrNoPubKey
	}

	if !adr36.VerifySignature(
		account.GetPubKey(),
		[]byte("I hereby declare that I am the address owner."),
		msg.Signature,
	) {
		return nil, types.ErrInvalidSignature
	}

	balance := k.bankKeeper.GetBalance(ctx, from, msg.Denom)
	if balance.IsZero() {
		return &types.MsgRecoverResponse{}, nil
	}

	to := sdk.MustAccAddressFromBech32(msg.To)
	err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(balance))
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from user to user")
	}

	return &types.MsgRecoverResponse{}, ctx.EventManager().EmitTypedEvent(&types.Recovered{
		Denom:  msg.Denom,
		From:   msg.From,
		To:     msg.To,
		Amount: balance.Amount,
	})
}

func (k msgServer) RemoveAdminAccount(goCtx context.Context, msg *types.MsgRemoveAdminAccount) (*types.MsgRemoveAdminAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.DeleteAdmin(ctx, msg.Denom, msg.Account)

	return &types.MsgRemoveAdminAccountResponse{}, ctx.EventManager().EmitTypedEvent(&types.AdminAccountRemoved{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) RemoveSystemAccount(goCtx context.Context, msg *types.MsgRemoveSystemAccount) (*types.MsgRemoveSystemAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.DeleteSystem(ctx, msg.Denom, msg.Account)

	return &types.MsgRemoveSystemAccountResponse{}, ctx.EventManager().EmitTypedEvent(&types.SystemAccountRemoved{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) SetMaxMintAllowance(goCtx context.Context, msg *types.MsgSetMaxMintAllowance) (*types.MsgSetMaxMintAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	k.Keeper.SetMaxMintAllowance(ctx, msg.Denom, msg.Amount)

	return &types.MsgSetMaxMintAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.MaxMintAllowance{
		Denom:  msg.Denom,
		Amount: msg.Amount,
	})
}

func (k msgServer) SetMintAllowance(goCtx context.Context, msg *types.MsgSetMintAllowance) (*types.MsgSetMintAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsAdmin(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidAdmin
	}

	maxMintAllowance := k.GetMaxMintAllowance(ctx, msg.Denom)
	if msg.Amount.IsNegative() || msg.Amount.GT(maxMintAllowance) {
		return nil, types.ErrInvalidAllowance
	}

	k.Keeper.SetMintAllowance(ctx, msg.Denom, msg.Account, msg.Amount)

	return &types.MsgSetMintAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.MintAllowance{
		Denom:   msg.Denom,
		Account: msg.Account,
		Amount:  msg.Amount,
	})
}

func (k msgServer) TransferOwnership(goCtx context.Context, msg *types.MsgTransferOwnership) (*types.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	owner, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if msg.NewOwner == owner {
		return nil, types.ErrSameOwner
	}

	k.SetPendingOwner(ctx, msg.Denom, msg.NewOwner)

	return &types.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferStarted{
		Denom:         msg.Denom,
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

//

func (k msgServer) EnsureOwner(ctx sdk.Context, denom string, signer string) (string, error) {
	owner := k.GetOwner(ctx, denom)
	if owner == "" {
		return "", types.ErrNoOwner
	}
	if signer != owner {
		return "", errors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, signer)
	}
	return owner, nil
}
