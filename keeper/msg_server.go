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
	"bytes"
	"context"
	"fmt"

	"adr36.dev"
	"cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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

func (k msgServer) AcceptOwnership(ctx context.Context, msg *types.MsgAcceptOwnership) (*types.MsgAcceptOwnershipResponse, error) {
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
	if err := k.SetOwner(ctx, msg.Denom, msg.Signer); err != nil {
		return nil, errors.Wrapf(types.ErrInvalidOwner, "failed to set owner: %s", msg.Denom)
	}
	if err := k.DeletePendingOwner(ctx, msg.Denom); err != nil {
		return nil, errors.Wrapf(types.ErrInvalidPendingOwner, "failed to delete pending owner: %s", msg.Denom)
	}

	return &types.MsgAcceptOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.OwnershipTransferred{
		Denom:         msg.Denom,
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k msgServer) AddAdminAccount(ctx context.Context, msg *types.MsgAddAdminAccount) (*types.MsgAddAdminAccountResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.SetAdmin(ctx, msg.Denom, msg.Account); err != nil {
		return nil, err
	}

	return &types.MsgAddAdminAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.AdminAccountAdded{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) AddSystemAccount(ctx context.Context, msg *types.MsgAddSystemAccount) (*types.MsgAddSystemAccountResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.SetSystem(ctx, msg.Denom, msg.Account); err != nil {
		return nil, err
	}

	return &types.MsgAddSystemAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.SystemAccountAdded{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) AllowDenom(ctx context.Context, msg *types.MsgAllowDenom) (*types.MsgAllowDenomResponse, error) {
	if msg.Signer != k.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	supply := k.bankKeeper.GetSupply(ctx, msg.Denom)
	if !supply.IsZero() {
		return nil, types.ErrInvalidDenom
	}

	if err := k.SetAllowedDenom(ctx, msg.Denom); err != nil {
		return nil, err
	}
	if err := k.SetOwner(ctx, msg.Denom, msg.Owner); err != nil {
		return nil, err
	}

	return &types.MsgAllowDenomResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.DenomAllowed{
		Denom: msg.Denom,
		Owner: msg.Owner,
	})
}

func (k msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsSystem(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidSystem
	}

	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnpackAny(msg.PubKey, &pubKey); err != nil {
		return nil, errors.Wrap(err, "unable to unpack pubkey")
	}
	from, err := k.addressCodec.StringToBytes(msg.From)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode user address")
	}
	if !bytes.Equal(from, pubKey.Address()) {
		return nil, types.ErrInvalidPubKey
	}

	if !adr36.VerifySignature(
		pubKey,
		[]byte("I hereby declare that I am the address owner."),
		msg.Signature,
	) {
		return nil, types.ErrInvalidSignature
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from user to module")
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to burn from module")
	}

	return &types.MsgBurnResponse{}, nil
}

func (k msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
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
	if err := k.Keeper.SetMintAllowance(ctx, msg.Denom, msg.Signer, allowance); err != nil {
		return nil, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	to, err := k.addressCodec.StringToBytes(msg.To)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode user address")
	}
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to mint to module")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from module to user")
	}

	return &types.MsgMintResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.MintAllowance{
		Denom:   msg.Denom,
		Account: msg.Signer,
		Amount:  allowance,
	})
}

func (k msgServer) Recover(ctx context.Context, msg *types.MsgRecover) (*types.MsgRecoverResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsSystem(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidSystem
	}

	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnpackAny(msg.PubKey, &pubKey); err != nil {
		return nil, errors.Wrap(err, "unable to unpack pubkey")
	}
	from, err := k.addressCodec.StringToBytes(msg.From)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode user address %s", msg.From)
	}
	if !bytes.Equal(from, pubKey.Address()) {
		return nil, types.ErrInvalidPubKey
	}

	if !adr36.VerifySignature(
		pubKey,
		[]byte("I hereby declare that I am the address owner."),
		msg.Signature,
	) {
		return nil, types.ErrInvalidSignature
	}

	balance := k.bankKeeper.GetBalance(ctx, from, msg.Denom)
	if balance.IsZero() {
		return &types.MsgRecoverResponse{}, nil
	}

	to, err := k.addressCodec.StringToBytes(msg.To)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode user address %s", msg.To)
	}
	err = k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(balance))
	if err != nil {
		return nil, errors.Wrap(err, "unable to transfer from user to user")
	}

	return &types.MsgRecoverResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.Recovered{
		Denom:  msg.Denom,
		From:   msg.From,
		To:     msg.To,
		Amount: balance.Amount,
	})
}

func (k msgServer) RemoveAdminAccount(ctx context.Context, msg *types.MsgRemoveAdminAccount) (*types.MsgRemoveAdminAccountResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.DeleteAdmin(ctx, msg.Denom, msg.Account); err != nil {
		return nil, err
	}

	return &types.MsgRemoveAdminAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.AdminAccountRemoved{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) RemoveSystemAccount(ctx context.Context, msg *types.MsgRemoveSystemAccount) (*types.MsgRemoveSystemAccountResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.DeleteSystem(ctx, msg.Denom, msg.Account); err != nil {
		return nil, err
	}

	return &types.MsgRemoveSystemAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.SystemAccountRemoved{
		Denom:   msg.Denom,
		Account: msg.Account,
	})
}

func (k msgServer) SetMaxMintAllowance(ctx context.Context, msg *types.MsgSetMaxMintAllowance) (*types.MsgSetMaxMintAllowanceResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	_, err := k.EnsureOwner(ctx, msg.Denom, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.SetMaxMintAllowance(ctx, msg.Denom, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgSetMaxMintAllowanceResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.MaxMintAllowance{
		Denom:  msg.Denom,
		Amount: msg.Amount,
	})
}

func (k msgServer) SetMintAllowance(ctx context.Context, msg *types.MsgSetMintAllowance) (*types.MsgSetMintAllowanceResponse, error) {
	if !k.IsAllowedDenom(ctx, msg.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", msg.Denom)
	}
	if !k.IsAdmin(ctx, msg.Denom, msg.Signer) {
		return nil, types.ErrInvalidAdmin
	}

	if msg.Amount.IsNegative() || msg.Amount.GT(k.GetMaxMintAllowance(ctx, msg.Denom)) {
		return nil, types.ErrInvalidAllowance
	}

	if err := k.Keeper.SetMintAllowance(ctx, msg.Denom, msg.Account, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgSetMintAllowanceResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.MintAllowance{
		Denom:   msg.Denom,
		Account: msg.Account,
		Amount:  msg.Amount,
	})
}

func (k msgServer) TransferOwnership(ctx context.Context, msg *types.MsgTransferOwnership) (*types.MsgTransferOwnershipResponse, error) {
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

	if err := k.SetPendingOwner(ctx, msg.Denom, msg.NewOwner); err != nil {
		return nil, err
	}

	return &types.MsgTransferOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &types.OwnershipTransferStarted{
		Denom:         msg.Denom,
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

//

func (k msgServer) EnsureOwner(ctx context.Context, denom string, signer string) (string, error) {
	owner := k.GetOwner(ctx, denom)
	if owner == "" {
		return "", types.ErrNoOwner
	}
	if signer != owner {
		return "", errors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, signer)
	}
	return owner, nil
}
