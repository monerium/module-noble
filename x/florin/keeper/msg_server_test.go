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

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/monerium/module-noble/v2/utils"
	"github.com/monerium/module-noble/v2/utils/mocks"
	"github.com/monerium/module-noble/v2/x/florin/keeper"
	"github.com/monerium/module-noble/v2/x/florin/types"
	"github.com/stretchr/testify/require"
)

var (
	MaxMintAllowance, _ = math.NewIntFromString("3000000000000000000000000")
	One, _              = math.NewIntFromString("1000000000000000000")
)

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with not allowed denom.
	_, err := server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorIs(t, err, types.ErrNoPendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidPendingOwner)

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Denom:  "ueure",
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetOwner(ctx, "ueure"))
	require.Empty(t, k.GetPendingOwner(ctx, "ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.OwnershipTransferred", events[0].Type)
}

func TestAddAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add admin account with not allowed denom.
	_, err := server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to add admin account with no owner set.
	_, err = server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to add admin account with invalid signer.
	_, err = server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate an admin account.
	admin := utils.TestAccount()

	// ACT: Attempt to add admin account.
	_, err = server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAdmin(ctx, "ueure", admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.AdminAccountAdded", events[0].Type)
}

func TestAddSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add system account with not allowed denom.
	_, err := server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to add system account with no owner set.
	_, err = server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to add system account with invalid signer.
	_, err = server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate a system account.
	system := utils.TestAccount()

	// ACT: Attempt to add system account.
	_, err = server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsSystem(ctx, "ueure", system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.SystemAccountAdded", events[0].Type)
}

func TestAllowDenom(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to allow denom with no authority set.
	_, err := server.AllowDenom(goCtx, &types.MsgAllowDenom{})
	// ASSERT: The action should've failed due to no authority set.
	require.ErrorIs(t, err, types.ErrNoAuthority)

	// ARRANGE: Set authority in state.
	authority := utils.TestAccount()
	k.SetAuthority(ctx, authority.Address)

	// ACT: Attempt to allow denom with invalid signer.
	_, err = server.AllowDenom(goCtx, &types.MsgAllowDenom{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidAuthority)

	// ACT: Attempt to allow denom already in use.
	_, err = server.AllowDenom(goCtx, &types.MsgAllowDenom{
		Signer: authority.Address,
		Denom:  "uusdc",
	})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorIs(t, err, types.ErrInvalidDenom)

	// ARRANGE: Generate an owner account.
	owner := utils.TestAccount()

	// ACT: Attempt to allow denom.
	_, err = server.AllowDenom(goCtx, &types.MsgAllowDenom{
		Signer: authority.Address,
		Denom:  "uusde",
		Owner:  owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAllowedDenom(ctx, "uusde"))
	require.Equal(t, owner.Address, k.GetOwner(ctx, "uusde"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.DenomAllowed", events[0].Type)
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
	pubkey, _ := codectypes.NewAnyWithValue(&secp256k1.PubKey{Key: bz})

	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.FlorinWithKeepers(account, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	k.SetSystem(ctx, "ueure", system.Address)

	// ACT: Attempt to burn invalid denom.
	_, err := server.Burn(goCtx, &types.MsgBurn{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to burn with invalid signer.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to burn with no account in state.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to no account.
	require.ErrorIs(t, err, types.ErrNoPubKey)

	// ARRANGE: Set account in state, without pubkey.
	account.Accounts["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = &authtypes.BaseAccount{}

	// ACT: Attempt to burn with no pubkey in state.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to no pubkey.
	require.ErrorIs(t, err, types.ErrNoPubKey)

	// ARRANGE: Set pubkey in state.
	account.Accounts["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = &authtypes.BaseAccount{
		Address:       "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		PubKey:        pubkey,
		AccountNumber: 0,
		Sequence:      0,
	}

	// ACT: Attempt to burn with invalid signature.
	signature, _ := base64.StdEncoding.DecodeString("QBrRfIqjdBvXx9zaBcuiE9P5SVesxFO/He3deyx2OE0NoSNqwmSb7b5iP2UhZRI1duiOeho3+NETUkCBv14zjQ==")
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Signature: signature,
	})
	// ASSERT: The action should've failed due to invalid signature.
	require.ErrorIs(t, err, types.ErrInvalidSignature)

	// ACT: Attempt to burn with insufficient balance.
	signature, _ = base64.StdEncoding.DecodeString("qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q==")
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Amount:    One,
		Signature: signature,
	})
	// ASSERT: The action should've failed due to insufficient balance.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $EURe.
	bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = sdk.NewCoins(sdk.NewCoin("ueure", One))

	// ACT: Attempt to burn.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Amount:    One,
		Signature: signature,
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
	k, ctx := mocks.FlorinWithKeepers(mocks.AccountKeeper{}, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	k.SetSystem(ctx, "ueure", system.Address)

	// ACT: Attempt to mint invalid denom.
	_, err := server.Mint(goCtx, &types.MsgMint{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to mint with invalid signer.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to mint with no allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		Amount: One,
	})
	// ASSERT: The action should've failed due to no allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Set mint allowance in state.
	k.SetMintAllowance(ctx, "ueure", system.Address, One)

	// ACT: Attempt to mint with insufficient allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		Amount: One.MulRaw(2),
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to mint.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Denom:  "ueure",
		Signer: system.Address,
		To:     user.Address,
		Amount: One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, bank.Balances[user.Address].AmountOf("ueure"))
	require.True(t, k.GetMintAllowance(ctx, "ueure", system.Address).IsZero())
	events := ctx.EventManager().Events()
	require.Len(t, events, 2)
	require.Equal(t, "florin.v1.MintAllowance", events[1].Type)
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
	pubkey, _ := codectypes.NewAnyWithValue(&secp256k1.PubKey{Key: bz})

	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.FlorinWithKeepers(account, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	k.SetSystem(ctx, "ueure", system.Address)

	// ACT: Attempt to recover invalid denom.
	_, err := server.Recover(goCtx, &types.MsgRecover{Denom: "uusde"})
	// ASSERT: The action should've failed due to invalid denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to recover with invalid signer.
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to recover with no account in state.
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to no account.
	require.ErrorIs(t, err, types.ErrNoPubKey)

	// ARRANGE: Set account in state, without pubkey.
	account.Accounts["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = &authtypes.BaseAccount{}

	// ACT: Attempt to recover with no pubkey in state.
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:  "ueure",
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to no pubkey.
	require.ErrorIs(t, err, types.ErrNoPubKey)

	// ARRANGE: Set pubkey in state.
	account.Accounts["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = &authtypes.BaseAccount{
		Address:       "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		PubKey:        pubkey,
		AccountNumber: 0,
		Sequence:      0,
	}

	// ACT: Attempt to recover with invalid signature.
	signature, _ := base64.StdEncoding.DecodeString("QBrRfIqjdBvXx9zaBcuiE9P5SVesxFO/He3deyx2OE0NoSNqwmSb7b5iP2UhZRI1duiOeho3+NETUkCBv14zjQ==")
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Signature: signature,
	})
	// ASSERT: The action should've failed due to invalid signature.
	require.ErrorIs(t, err, types.ErrInvalidSignature)

	// ARRANGE: Generate a recipient address.
	recipient := utils.TestAccount()

	// ACT: Attempt to recover with no balance.
	signature, _ = base64.StdEncoding.DecodeString("qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q==")
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		To:        recipient.Address,
		Signature: signature,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)

	// ARRANGE: Give user 1 $EURe.
	bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = sdk.NewCoins(sdk.NewCoin("ueure", One))

	// ACT: Attempt to recover.
	_, err = server.Recover(goCtx, &types.MsgRecover{
		Denom:     "ueure",
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		To:        recipient.Address,
		Signature: signature,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"].IsZero())
	require.Equal(t, One, bank.Balances[recipient.Address].AmountOf("ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 2)
	require.Equal(t, "florin.v1.Recovered", events[1].Type)
}

func TestRemoveAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove admin account with not allowed denom.
	_, err := server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to remove admin account with no owner set.
	_, err = server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to remove admin account with invalid signer.
	_, err = server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	k.SetAdmin(ctx, "ueure", admin.Address)
	require.True(t, k.IsAdmin(ctx, "ueure", admin.Address))

	// ACT: Attempt to remove admin account.
	_, err = server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsAdmin(ctx, "ueure", admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.AdminAccountRemoved", events[0].Type)
}

func TestRemoveSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove system account with not allowed denom.
	_, err := server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to remove system account with no owner set.
	_, err = server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to remove system account with invalid signer.
	_, err = server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	k.SetSystem(ctx, "ueure", system.Address)
	require.True(t, k.IsSystem(ctx, "ueure", system.Address))

	// ACT: Attempt to remove system account.
	_, err = server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Denom:   "ueure",
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsSystem(ctx, "ueure", system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.SystemAccountRemoved", events[0].Type)
}

func TestSetMaxMintAllowance(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set max mint allowance with not allowed denom.
	_, err := server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to set max mint allowance with no owner set.
	_, err = server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to set max mint allowance with invalid signer.
	_, err = server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ACT: Attempt to set max mint allowance.
	_, err = server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Denom:  "ueure",
		Signer: owner.Address,
		Amount: MaxMintAllowance,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, k.GetMaxMintAllowance(ctx, "ueure"))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.MaxMintAllowance", events[0].Type)
}

func TestSetMintAllowance(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	k.SetAdmin(ctx, "ueure", admin.Address)
	// ARRANGE: Set max mint allowance in state.
	k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)

	// ACT: Attempt to set mint allowance with invalid denom.
	_, err := server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to set mint allowance with invalid signer.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidAdmin)

	// ACT: Attempt to set mint allowance with negative amount.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: admin.Address,
		Amount: One.Neg(),
	})
	// ASSERT: The action should've failed due to negative amount.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ACT: Attempt to set mint allowance to more than max.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Denom:  "ueure",
		Signer: admin.Address,
		Amount: MaxMintAllowance.Add(One),
	})
	// ASSERT: The action should've failed due to more than max.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to set mint allowance.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Denom:   "ueure",
		Signer:  admin.Address,
		Account: minter.Address,
		Amount:  One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, k.GetMintAllowance(ctx, "ueure", minter.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.MintAllowance", events[0].Type)
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with not allowed denom.
	_, err := server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Denom: "uusde",
	})
	// ASSERT: The action should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to transfer ownership with no owner set.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Denom: "ueure",
	})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Denom:  "ueure",
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Denom:    "ueure",
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorIs(t, err, types.ErrSameOwner)

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
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
	require.Equal(t, "florin.v1.OwnershipTransferStarted", events[0].Type)
}
