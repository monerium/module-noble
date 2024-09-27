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

package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/monerium/module-noble/v2/keeper"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/types/blacklist"
	"github.com/monerium/module-noble/v2/utils"
	"github.com/monerium/module-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestBlacklistAcceptOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(ctx, &blacklist.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorIs(t, err, blacklist.ErrNoPendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	_ = k.SetBlacklistPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &blacklist.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidPendingOwner)

	// ARRANGE: Set up a failing collection store for the attribute getter.
	tmpOwner := k.BlacklistOwner
	k.BlacklistOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.OwnerKey, "blacklistOwner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing BlacklistOwner collection store.
	_, err = server.AcceptOwnership(ctx, &blacklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlacklistOwner = tmpOwner

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmpPendingOwner := k.BlacklistPendingOwner
	k.BlacklistPendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.PendingOwnerKey, "blacklistPendingOwner", collections.StringValue,
	)

	// ACT: Attempt to accept ownership with failing BlacklistPendingOwner collection store.
	_, err = server.AcceptOwnership(ctx, &blacklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlacklistPendingOwner = tmpPendingOwner

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &blacklist.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlacklistOwner(ctx))
	require.Empty(t, k.GetBlacklistPendingOwner(ctx))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.OwnershipTransferred", events[0].Type)
}

func TestBlacklistAddAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ACT: Attempt to add admin account with no owner set.
	_, err := server.AddAdminAccount(ctx, &blacklist.MsgAddAdminAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, blacklist.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	_ = k.SetBlacklistOwner(ctx, owner.Address)

	// ACT: Attempt to add admin account with invalid signer.
	_, err = server.AddAdminAccount(ctx, &blacklist.MsgAddAdminAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidOwner)

	// ARRANGE: Generate an admin account.
	admin := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlacklistAdmins
	k.BlacklistAdmins = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.AdminPrefix, "blacklistAdmins", collections.StringKey,
	)

	// ACT: Attempt to add admin account with failing BlacklistAdmins collection store.
	_, err = server.AddAdminAccount(ctx, &blacklist.MsgAddAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlacklistAdmins = tmp

	// ACT: Attempt to add admin account.
	_, err = server.AddAdminAccount(ctx, &blacklist.MsgAddAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsBlacklistAdmin(ctx, admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.AdminAccountAdded", events[0].Type)
}

func TestBan(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	err := k.SetBlacklistAdmin(ctx, admin.Address)
	require.NoError(t, err)

	// ACT: Attempt to ban with invalid signer.
	_, err = server.Ban(ctx, &blacklist.MsgBan{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidAdmin)

	// ARRANGE: Generate an adversary account.
	adversary := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Adversaries
	k.Adversaries = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.AdversaryPrefix, "adversaries", collections.StringKey,
	)

	// ACT: Attempt to ban with failing Adversaries collection store.
	_, err = server.Ban(ctx, &blacklist.MsgBan{
		Signer:    admin.Address,
		Adversary: adversary.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Adversaries = tmp

	// ACT: Attempt to ban.
	_, err = server.Ban(ctx, &blacklist.MsgBan{
		Signer:    admin.Address,
		Adversary: adversary.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAdversary(ctx, adversary.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.Ban", events[0].Type)
}

func TestBlacklistRemoveAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ACT: Attempt to remove admin account with no owner set.
	_, err := server.RemoveAdminAccount(ctx, &blacklist.MsgRemoveAdminAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, blacklist.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	_ = k.SetBlacklistOwner(ctx, owner.Address)

	// ACT: Attempt to remove admin account with invalid signer.
	_, err = server.RemoveAdminAccount(ctx, &blacklist.MsgRemoveAdminAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidOwner)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	err = k.SetBlacklistAdmin(ctx, admin.Address)
	require.NoError(t, err)
	require.True(t, k.IsBlacklistAdmin(ctx, admin.Address))

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmp := k.BlacklistAdmins
	k.BlacklistAdmins = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.AdminPrefix, "blacklistAdmins", collections.StringKey,
	)

	// ACT: Attempt to remove admin account with failing BlacklistAdmins collection store.
	_, err = server.RemoveAdminAccount(ctx, &blacklist.MsgRemoveAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlacklistAdmins = tmp

	// ACT: Attempt to remove admin account.
	_, err = server.RemoveAdminAccount(ctx, &blacklist.MsgRemoveAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsBlacklistAdmin(ctx, admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.AdminAccountRemoved", events[0].Type)
}

func TestBlacklistTransferOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(ctx, &blacklist.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, blacklist.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetBlacklistOwner(ctx, owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &blacklist.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidOwner)

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(ctx, &blacklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorIs(t, err, blacklist.ErrSameOwner)

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.BlacklistPendingOwner
	k.BlacklistPendingOwner = collections.NewItem(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.PendingOwnerKey, "blacklistPendingOwner", collections.StringValue,
	)

	// ACT: Attempt to transfer ownership BlacklistPendingOwner collection store.
	_, err = server.TransferOwnership(ctx, &blacklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.BlacklistPendingOwner = tmp

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &blacklist.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, owner.Address, k.GetBlacklistOwner(ctx))
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetBlacklistPendingOwner(ctx))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.OwnershipTransferStarted", events[0].Type)
}

func TestUnban(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewBlacklistMsgServer(k)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	_ = k.SetBlacklistAdmin(ctx, admin.Address)

	// ACT: Attempt to unban with invalid signer.
	_, err := server.Unban(ctx, &blacklist.MsgUnban{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, blacklist.ErrInvalidAdmin)

	// ARRANGE: Set adversary in state.
	adversary := utils.TestAccount()
	_ = k.SetAdversary(ctx, adversary.Address)
	require.True(t, k.IsAdversary(ctx, adversary.Address))

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmp := k.Adversaries
	k.Adversaries = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		blacklist.AdversaryPrefix, "adversaries", collections.StringKey,
	)

	// ACT: Attempt to unban with failing Adversaries collection store.
	_, err = server.Unban(ctx, &blacklist.MsgUnban{
		Signer: admin.Address,
		Friend: adversary.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Adversaries = tmp

	// ACT: Attempt to unban.
	_, err = server.Unban(ctx, &blacklist.MsgUnban{
		Signer: admin.Address,
		Friend: adversary.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsAdversary(ctx, adversary.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.Unban", events[0].Type)
}
