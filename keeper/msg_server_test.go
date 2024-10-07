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
	"encoding/base64"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/keeper"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/utils"
	"github.com/monerium/module-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

var (
	MaxMintAllowance, _ = math.NewIntFromString("3000000000000000000000000")
	One, _              = math.NewIntFromString("1000000000000000000")
)

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with not allowed denom.
	_, err := server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorIs(t, err, types.ErrNoPendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	err = k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)
	require.NoError(t, err)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidPendingOwner)

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmpOwner := k.Owner
	k.Owner = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.OwnerPrefix, "owner", collections.StringKey, collections.StringValue,
	)

	// ACT: Attempt to transfer ownership with failing Owner collection store.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Owner = tmpOwner

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmpPendingOwner := k.PendingOwner
	k.PendingOwner = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.PendingOwnerPrefix, "pendingOwner", collections.StringKey, collections.StringValue,
	)

	// ACT: Attempt to transfer ownership with failing PendingOwner collection store.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.PendingOwner = tmpPendingOwner

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(ctx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetOwner(ctx, "ueure"))
	require.Empty(t, k.GetPendingOwner(ctx, "ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.OwnershipTransferred", events[0].Type)
}

func TestAddAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add admin account with not allowed denom.
	_, err := server.AddAdminAccount(ctx, &types.MsgAddAdminAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to add admin account with no owner set.
	_, err = server.AddAdminAccount(ctx, &types.MsgAddAdminAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to add admin account with invalid signer.
	_, err = server.AddAdminAccount(ctx, &types.MsgAddAdminAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate an admin account.
	admin := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Admins
	k.Admins = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.AdminPrefix, "admins", collections.PairKeyCodec(collections.StringKey, collections.StringKey),
	)

	// ACT: Attempt to add admin account with failing Admins collection store.
	_, err = server.AddAdminAccount(ctx, &types.MsgAddAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Admins = tmp

	// ACT: Attempt to add admin account.
	_, err = server.AddAdminAccount(ctx, &types.MsgAddAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAdmin(ctx, "ueure", admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.AdminAccountAdded", events[0].Type)
}

func TestAddSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add system account with not allowed denom.
	_, err := server.AddSystemAccount(ctx, &types.MsgAddSystemAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to add system account with no owner set.
	_, err = server.AddSystemAccount(ctx, &types.MsgAddSystemAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to add system account with invalid signer.
	_, err = server.AddSystemAccount(ctx, &types.MsgAddSystemAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate a system account.
	system := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.Systems
	k.Systems = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.SystemPrefix, "systems", collections.PairKeyCodec(collections.StringKey, collections.StringKey),
	)

	// ACT: Attempt to add system account with failing Systems collection store.
	_, err = server.AddSystemAccount(ctx, &types.MsgAddSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Systems = tmp

	// ACT: Attempt to add system account.
	_, err = server.AddSystemAccount(ctx, &types.MsgAddSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsSystem(ctx, "ueure", system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.SystemAccountAdded", events[0].Type)
}

func TestAllowDenom(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to allow denom with invalid signer.
	_, err := server.AllowDenom(ctx, &types.MsgAllowDenom{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidAuthority)

	// ACT: Attempt to allow denom already in use.
	_, err = server.AllowDenom(ctx, &types.MsgAllowDenom{
		Signer: "authority",
		Denom:  "uusdc",
	})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorIs(t, err, types.ErrInvalidDenom)

	// ARRANGE: Generate an owner account.
	owner := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmpAllowedDenoms := k.AllowedDenoms
	k.AllowedDenoms = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.AllowedDenomPrefix, "allowedDenoms", collections.StringKey,
	)

	// ACT: Attempt to allow denom with failing AllowedDenoms collection store.
	_, err = server.AllowDenom(ctx, &types.MsgAllowDenom{
		Signer: "authority",
		Denom:  "uusde",
		Owner:  owner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.AllowedDenoms = tmpAllowedDenoms

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmpOwner := k.Owner
	k.Owner = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.OwnerPrefix, "owner", collections.StringKey, collections.StringValue,
	)

	// ACT: Attempt to allow denom with failing Owner collection store.
	_, err = server.AllowDenom(ctx, &types.MsgAllowDenom{
		Signer: "authority",
		Denom:  "uusde",
		Owner:  owner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Owner = tmpOwner

	// ACT: Attempt to allow denom.
	_, err = server.AllowDenom(ctx, &types.MsgAllowDenom{
		Signer: "authority",
		Denom:  "uusde",
		Owner:  owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAllowedDenom(ctx, "uusde"))
	require.Equal(t, owner.Address, k.GetOwner(ctx, "uusde"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.DenomAllowed", events[0].Type)
}

func TestBurn(t *testing.T) {
	// The below signature was generated using Keplr's signArbitrary function.
	//
	// share bubble good swarm sustain leaf burst build spirit inflict undo shadow antique warm soft praise foam slab laptop hint giggle also book treat
	//
	// {
	//     "pub_key": {
	//         "type": "tendermint/PubKeySecp256k1",
	//         "value": "AlE8CxHR19ID5lxrVtTxSgJFlK3T+eYtyDM/vBA3Fowr"
	//     },
	//     "signature": "qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q=="
	// }
	bz, _ := base64.StdEncoding.DecodeString("AlE8CxHR19ID5lxrVtTxSgJFlK3T+eYtyDM/vBA3Fowr")
	pubKey, _ := codectypes.NewAnyWithValue(&secp256k1.PubKey{Key: bz})

	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.FlorinWithKeepers(bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	err := k.SetSystem(ctx, "ueure", system.Address)
	require.NoError(t, err)

	// ACT: Attempt to burn invalid denom.
	_, err = server.Burn(ctx, &types.MsgBurn{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to burn with invalid signer.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to burn with invalid any.
	invalidPubKey, _ := codectypes.NewAnyWithValue(&types.MsgBurn{})
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		PubKey: invalidPubKey,
	})
	// ASSERT: The action should've failed due to invalid any.
	require.ErrorContains(t, err, "unable to unpack pubkey")

	// ACT: Attempt to burn from invalid user address.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		From:   utils.TestAccount().Invalid,
		PubKey: pubKey,
	})
	// ASSERT: The action should've failed due to invalid user address.
	require.ErrorContains(t, err, "unable to decode user address")

	// ACT: Attempt to burn with missing public key.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		From:   utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid public key.
	require.ErrorIs(t, err, types.ErrInvalidPubKey)

	// ACT: Attempt to burn with invalid public key.
	invalidPubKey, _ = codectypes.NewAnyWithValue(secp256k1.GenPrivKey().PubKey())
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		PubKey: invalidPubKey,
	})
	// ASSERT: The action should've failed due to invalid public key.
	require.ErrorIs(t, err, types.ErrInvalidPubKey)

	// ACT: Attempt to burn with invalid signature.
	signature, _ := base64.StdEncoding.DecodeString("QBrRfIqjdBvXx9zaBcuiE9P5SVesxFO/He3deyx2OE0NoSNqwmSb7b5iP2UhZRI1duiOeho3+NETUkCBv14zjQ==")
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've failed due to invalid signature.
	require.ErrorIs(t, err, types.ErrInvalidSignature)

	// ACT: Attempt to burn with insufficient balance.
	signature, _ = base64.StdEncoding.DecodeString("qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q==")
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Amount:    One,
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've failed due to insufficient balance.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $EURe.
	bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = sdk.NewCoins(sdk.NewCoin("ueure", One))

	// ACT: Attempt to burn.
	_, err = server.Burn(ctx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Amount:    One,
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"].IsZero())
}

func TestMint(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.FlorinWithKeepers(bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	err := k.SetSystem(ctx, "ueure", system.Address)
	require.NoError(t, err)

	// ACT: Attempt to mint invalid denom.
	_, err = server.Mint(ctx, &types.MsgMint{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to mint with invalid signer.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to mint with no allowance.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		Amount: One,
	})
	// ASSERT: The action should've failed due to no allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Set mint allowance in state.
	err = k.SetMintAllowance(ctx, "ueure", system.Address, One)
	require.NoError(t, err)

	// ACT: Attempt to mint with insufficient allowance.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		Amount: One.MulRaw(2),
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.MintAllowance
	k.MintAllowance = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MintAllowancePrefix, "mintAllowance", collections.PairKeyCodec(collections.StringKey, collections.StringKey), collections.BytesValue,
	)

	// ACT: Attempt to mint with failing MintAllowance collection store.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		To:     user.Address,
		Amount: One,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.MintAllowance = tmp

	// ACT: Attempt to mint to invalid user address.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		To:     user.Invalid,
		Amount: One,
	})
	// ASSERT: The action should've failed due to invalid user address.
	require.ErrorContains(t, err, "unable to decode user address")

	// ARRANGE: Reset mint allowance in state.
	err = k.SetMintAllowance(ctx, "ueure", system.Address, One)
	require.NoError(t, err)

	// ACT: Attempt to mint.
	_, err = server.Mint(ctx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		To:     user.Address,
		Amount: One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, bank.Balances[user.Address].AmountOf("ueure"))
	require.NoError(t, err)
	require.True(t, k.GetMintAllowance(ctx, "ueure", system.Address).IsZero())
	events := ctx.EventManager().Events()
	require.Len(t, events, 2)
	require.Equal(t, "florin.v2.MintAllowance", events[1].Type)
}

func TestRecover(t *testing.T) {
	// The below signature was generated using Keplr's signArbitrary function.
	//
	// share bubble good swarm sustain leaf burst build spirit inflict undo shadow antique warm soft praise foam slab laptop hint giggle also book treat
	//
	// {
	//     "pub_key": {
	//         "type": "tendermint/PubKeySecp256k1",
	//         "value": "AlE8CxHR19ID5lxrVtTxSgJFlK3T+eYtyDM/vBA3Fowr"
	//     },
	//     "signature": "qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q=="
	// }
	bz, _ := base64.StdEncoding.DecodeString("AlE8CxHR19ID5lxrVtTxSgJFlK3T+eYtyDM/vBA3Fowr")
	pubKey, _ := codectypes.NewAnyWithValue(&secp256k1.PubKey{Key: bz})

	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.FlorinWithKeepers(bank)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	err := k.SetSystem(ctx, "ueure", system.Address)
	require.NoError(t, err)

	// ACT: Attempt to recover invalid denom.
	_, err = server.Recover(ctx, &types.MsgRecover{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to recover with invalid signer.
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to recover with invalid any.
	invalidPubKey, _ := codectypes.NewAnyWithValue(&types.MsgRecover{})
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		PubKey: invalidPubKey,
	})
	// ASSERT: The action should've failed due to invalid any.
	require.ErrorContains(t, err, "unable to unpack pubkey")

	// ACT: Attempt to recover from invalid user address.
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		From:   utils.TestAccount().Invalid,
		PubKey: pubKey,
	})
	// ASSERT: The action should've failed due to invalid user address.
	require.ErrorContains(t, err, "unable to decode user address")

	// ACT: Attempt to recover with invalid public key.
	invalidPubKey, _ = codectypes.NewAnyWithValue(secp256k1.GenPrivKey().PubKey())
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		PubKey: invalidPubKey,
	})
	// ASSERT: The action should've failed due to invalid public key.
	require.ErrorIs(t, err, types.ErrInvalidPubKey)

	// ACT: Attempt to recover with missing public key.
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to invalid public key.
	require.ErrorIs(t, err, types.ErrInvalidPubKey)

	// ACT: Attempt to recover with invalid signature.
	signature, _ := base64.StdEncoding.DecodeString("QBrRfIqjdBvXx9zaBcuiE9P5SVesxFO/He3deyx2OE0NoSNqwmSb7b5iP2UhZRI1duiOeho3+NETUkCBv14zjQ==")
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've failed due to invalid signature.
	require.ErrorIs(t, err, types.ErrInvalidSignature)

	// ARRANGE: Generate a recipient address.
	recipient := utils.TestAccount()

	// ACT: Attempt to recover with no balance.
	signature, _ = base64.StdEncoding.DecodeString("qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q==")
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		To:        recipient.Address,
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Give user 1 $EURe.
	bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = sdk.NewCoins(sdk.NewCoin("ueure", One))

	// ACT: Attempt to recover to invalid user address.
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		To:        recipient.Invalid,
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've failed due to invalid user address.
	require.ErrorContains(t, err, "unable to decode user address")

	// ACT: Attempt to recover.
	_, err = server.Recover(ctx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		To:        recipient.Address,
		Signature: signature,
		PubKey:    pubKey,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"].IsZero())
	require.Equal(t, One, bank.Balances[recipient.Address].AmountOf("ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 2)
	require.Equal(t, "florin.v2.Recovered", events[1].Type)
}

func TestRemoveAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove admin account with not allowed denom.
	_, err := server.RemoveAdminAccount(ctx, &types.MsgRemoveAdminAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to remove admin account with no owner set.
	_, err = server.RemoveAdminAccount(ctx, &types.MsgRemoveAdminAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to remove admin account with invalid signer.
	_, err = server.RemoveAdminAccount(ctx, &types.MsgRemoveAdminAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	err = k.SetAdmin(ctx, "ueure", admin.Address)
	require.NoError(t, err)
	require.True(t, k.IsAdmin(ctx, "ueure", admin.Address))

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmp := k.Admins
	k.Admins = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.AdminPrefix, "admins", collections.PairKeyCodec(collections.StringKey, collections.StringKey),
	)

	// ACT: Attempt to remove admin account with failing Admins collection store.
	_, err = server.RemoveAdminAccount(ctx, &types.MsgRemoveAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Admins = tmp

	// ACT: Attempt to remove admin account.
	_, err = server.RemoveAdminAccount(ctx, &types.MsgRemoveAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsAdmin(ctx, "ueure", admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.AdminAccountRemoved", events[0].Type)
}

func TestRemoveSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove system account with not allowed denom.
	_, err := server.RemoveSystemAccount(ctx, &types.MsgRemoveSystemAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to remove system account with no owner set.
	_, err = server.RemoveSystemAccount(ctx, &types.MsgRemoveSystemAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to remove system account with invalid signer.
	_, err = server.RemoveSystemAccount(ctx, &types.MsgRemoveSystemAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	err = k.SetSystem(ctx, "ueure", system.Address)
	require.NoError(t, err)
	require.True(t, k.IsSystem(ctx, "ueure", system.Address))

	// ARRANGE: Set up a failing collection store for the attribute deleter.
	tmp := k.Systems
	k.Systems = collections.NewKeySet(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Delete, utils.GetKVStore(ctx, types.ModuleName))),
		types.SystemPrefix, "systems", collections.PairKeyCodec(collections.StringKey, collections.StringKey),
	)

	// ACT: Attempt to remove system account with failing Systems collection store.
	_, err = server.RemoveSystemAccount(ctx, &types.MsgRemoveSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've failed due to collection store deleter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.Systems = tmp

	// ACT: Attempt to remove system account.
	_, err = server.RemoveSystemAccount(ctx, &types.MsgRemoveSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsSystem(ctx, "ueure", system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.SystemAccountRemoved", events[0].Type)
}

func TestSetMaxMintAllowance(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set max mint allowance with not allowed denom.
	_, err := server.SetMaxMintAllowance(ctx, &types.MsgSetMaxMintAllowance{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to set max mint allowance with no owner set.
	_, err = server.SetMaxMintAllowance(ctx, &types.MsgSetMaxMintAllowance{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to set max mint allowance with invalid signer.
	_, err = server.SetMaxMintAllowance(ctx, &types.MsgSetMaxMintAllowance{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.MaxMintAllowance
	k.MaxMintAllowance = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MaxMintAllowancePrefix, "maxMintAllowance", collections.StringKey, collections.BytesValue,
	)

	// ACT: Attempt to set max mint allowance with failing MaxMintAllowance collection store.
	_, err = server.SetMaxMintAllowance(ctx, &types.MsgSetMaxMintAllowance{
		Denom:  "ueure",
		Signer: owner.Address,
		Amount: MaxMintAllowance,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.MaxMintAllowance = tmp

	// ACT: Attempt to set max mint allowance.
	_, err = server.SetMaxMintAllowance(ctx, &types.MsgSetMaxMintAllowance{
		Denom:  "ueure",
		Signer: owner.Address,
		Amount: MaxMintAllowance,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, k.GetMaxMintAllowance(ctx, "ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.MaxMintAllowance", events[0].Type)
}

func TestSetMintAllowance(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	err := k.SetAdmin(ctx, "ueure", admin.Address)
	require.NoError(t, err)
	// ARRANGE: Set max mint allowance in state.
	err = k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)
	require.NoError(t, err)

	// ACT: Attempt to set mint allowance with invalid denom.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to set mint allowance with invalid signer.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidAdmin)

	// ACT: Attempt to set mint allowance with negative amount.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: admin.Address,
		Amount: One.Neg(),
	})
	// ASSERT: The action should've failed due to negative amount.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ACT: Attempt to set mint allowance to more than max.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: admin.Address,
		Amount: MaxMintAllowance.Add(One),
	})
	// ASSERT: The action should've failed due to more than max.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.MintAllowance
	k.MintAllowance = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.MintAllowancePrefix, "mintAllowance", collections.PairKeyCodec(collections.StringKey, collections.StringKey), collections.BytesValue,
	)

	// ACT: Attempt to set mint allowance with failing MintAllowance collection store.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom:   "ueure",
		Signer:  admin.Address,
		Account: minter.Address,
		Amount:  One,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.MintAllowance = tmp

	// ACT: Attempt to set mint allowance.
	_, err = server.SetMintAllowance(ctx, &types.MsgSetMintAllowance{
		Denom:   "ueure",
		Signer:  admin.Address,
		Account: minter.Address,
		Amount:  One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, One, k.GetMintAllowance(ctx, "ueure", minter.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.MintAllowance", events[0].Type)
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with not allowed denom.
	_, err := server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to transfer ownership with no owner set.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom:    "ueure",
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorIs(t, err, types.ErrSameOwner)

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ARRANGE: Set up a failing collection store for the attribute setter.
	tmp := k.PendingOwner
	k.PendingOwner = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PendingOwnerPrefix, "pendingOwner", collections.StringKey, collections.StringValue,
	)

	// ACT: Attempt to transfer ownership with failing PendingOwner collection store.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom:    "ueure",
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've failed due to collection store setter error.
	require.Error(t, err, mocks.ErrorStoreAccess)
	k.PendingOwner = tmp

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(ctx, &types.MsgTransferOwnership{
		Denom:    "ueure",
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, owner.Address, k.GetOwner(ctx, "ueure"))
	require.Equal(t, pendingOwner.Address, k.GetPendingOwner(ctx, "ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v2.OwnershipTransferStarted", events[0].Type)
}
