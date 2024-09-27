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

	"cosmossdk.io/errors"
	"github.com/monerium/module-noble/v2/types/blacklist"
)

var _ blacklist.MsgServer = &blacklistMsgServer{}

type blacklistMsgServer struct {
	*Keeper
}

func NewBlacklistMsgServer(keeper *Keeper) blacklist.MsgServer {
	return &blacklistMsgServer{Keeper: keeper}
}

func (k blacklistMsgServer) AcceptOwnership(ctx context.Context, msg *blacklist.MsgAcceptOwnership) (*blacklist.MsgAcceptOwnershipResponse, error) {
	pendingOwner := k.GetBlacklistPendingOwner(ctx)
	if pendingOwner == "" {
		return nil, blacklist.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(blacklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	owner := k.GetBlacklistOwner(ctx)
	if err := k.SetBlacklistOwner(ctx, msg.Signer); err != nil {
		return nil, errors.Wrap(err, "failed to set blacklist owner")
	}
	if err := k.DeleteBlacklistPendingOwner(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to delete blacklist pending owner")
	}

	return &blacklist.MsgAcceptOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k blacklistMsgServer) AddAdminAccount(ctx context.Context, msg *blacklist.MsgAddAdminAccount) (*blacklist.MsgAddAdminAccountResponse, error) {
	_, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.SetBlacklistAdmin(ctx, msg.Account); err != nil {
		return nil, errors.Wrapf(err, "failed to set blacklist admin: %s", msg.Account)
	}

	return &blacklist.MsgAddAdminAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.AdminAccountAdded{
		Account: msg.Account,
	})
}

func (k blacklistMsgServer) Ban(ctx context.Context, msg *blacklist.MsgBan) (*blacklist.MsgBanResponse, error) {
	if !k.IsBlacklistAdmin(ctx, msg.Signer) {
		return nil, blacklist.ErrInvalidAdmin
	}

	if err := k.SetAdversary(ctx, msg.Adversary); err != nil {
		return nil, errors.Wrapf(err, "failed to set blacklist adversary: %s", msg.Adversary)
	}

	return &blacklist.MsgBanResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.Ban{
		Adversary: msg.Adversary,
	})
}

func (k blacklistMsgServer) RemoveAdminAccount(ctx context.Context, msg *blacklist.MsgRemoveAdminAccount) (*blacklist.MsgRemoveAdminAccountResponse, error) {
	_, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	if err := k.DeleteBlacklistAdmin(ctx, msg.Account); err != nil {
		return nil, errors.Wrapf(err, "failed to delete blacklist admin: %s", msg.Account)
	}

	return &blacklist.MsgRemoveAdminAccountResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.AdminAccountRemoved{
		Account: msg.Account,
	})
}

func (k blacklistMsgServer) TransferOwnership(ctx context.Context, msg *blacklist.MsgTransferOwnership) (*blacklist.MsgTransferOwnershipResponse, error) {
	owner, err := k.EnsureOwner(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	if msg.NewOwner == owner {
		return nil, blacklist.ErrSameOwner
	}

	if err := k.SetBlacklistPendingOwner(ctx, msg.NewOwner); err != nil {
		return nil, errors.Wrapf(err, "failed to set blacklist pending owner: %s", msg.NewOwner)
	}

	return &blacklist.MsgTransferOwnershipResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k blacklistMsgServer) Unban(ctx context.Context, msg *blacklist.MsgUnban) (*blacklist.MsgUnbanResponse, error) {
	if !k.IsBlacklistAdmin(ctx, msg.Signer) {
		return nil, blacklist.ErrInvalidAdmin
	}

	if err := k.DeleteAdversary(ctx, msg.Friend); err != nil {
		return nil, errors.Wrapf(err, "failed to delete blacklist adversary: %s", msg.Friend)
	}

	return &blacklist.MsgUnbanResponse{}, k.eventService.EventManager(ctx).Emit(ctx, &blacklist.Unban{
		Friend: msg.Friend,
	})
}

//

func (k blacklistMsgServer) EnsureOwner(ctx context.Context, signer string) (string, error) {
	owner := k.GetBlacklistOwner(ctx)
	if owner == "" {
		return "", blacklist.ErrNoOwner
	}
	if signer != owner {
		return "", errors.Wrapf(blacklist.ErrInvalidOwner, "expected %s, got %s", owner, signer)
	}
	return owner, nil
}
