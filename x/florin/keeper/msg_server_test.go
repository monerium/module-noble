package keeper_test

import (
	"encoding/base64"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/stretchr/testify/require"
)

var (
	MaxMintAllowance, _ = sdk.NewIntFromString("50000000000000000000000000000000")
	One, _              = sdk.NewIntFromString("1000000000000000000")
)

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorIs(t, err, types.ErrNoPendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidPendingOwner)

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, k.GetOwner(ctx))
	require.Empty(t, k.GetPendingOwner(ctx))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.OwnershipTransferred", events[0].Type)
}

func TestAddAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add admin account with no owner set.
	_, err := server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add admin account with invalid signer.
	_, err = server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate an admin account.
	admin := utils.TestAccount()

	// ACT: Attempt to add admin account.
	_, err = server.AddAdminAccount(goCtx, &types.MsgAddAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsAdmin(ctx, admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.AdminAccountAdded", events[0].Type)
}

func TestAddSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add system account with no owner set.
	_, err := server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add system account with invalid signer.
	_, err = server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Generate a system account.
	system := utils.TestAccount()

	// ACT: Attempt to add system account.
	_, err = server.AddSystemAccount(goCtx, &types.MsgAddSystemAccount{
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, k.IsSystem(ctx, system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.SystemAccountAdded", events[0].Type)
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
		Accounts: make(map[string]authtypes.AccountI),
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
	k.SetSystem(ctx, system.Address)

	// ACT: Attempt to burn with invalid signer.
	_, err := server.Burn(goCtx, &types.MsgBurn{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to burn with no account in state.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: system.Address,
		From:   "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
	})
	// ASSERT: The action should've failed due to no account.
	require.ErrorIs(t, err, types.ErrNoPubKey)

	// ARRANGE: Set account in state, without pubkey.
	account.Accounts["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = &authtypes.BaseAccount{}

	// ACT: Attempt to burn with no pubkey in state.
	_, err = server.Burn(goCtx, &types.MsgBurn{
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
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Signature: signature,
	})
	// ASSERT: The action should've failed due to invalid signature.
	require.ErrorIs(t, err, types.ErrInvalidSignature)

	// ACT: Attempt to burn with insufficient balance.
	signature, _ = base64.StdEncoding.DecodeString("qe5dDxdOgY8B2LjMqnK5/5iRIFOCwdTu0G5ZQ66bHzVgP15V2Fb+fzOH0wPAUC5GUQ23M1cSvysulzKIbXY/4Q==")
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer:    system.Address,
		From:      "noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2",
		Amount:    One,
		Signature: signature,
	})
	// ASSERT: The action should've failed due to insufficient balance.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $EURe.
	bank.Balances["noble1rwvjzk28l38js7xx6mq23nrpghd8qqvxmj6ep2"] = sdk.NewCoins(sdk.NewCoin(k.Denom, One))

	// ACT: Attempt to burn.
	_, err = server.Burn(goCtx, &types.MsgBurn{
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
	k.SetSystem(ctx, system.Address)

	// ACT: Attempt to mint with invalid signer.
	_, err := server.Mint(goCtx, &types.MsgMint{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidSystem)

	// ACT: Attempt to mint with no allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: system.Address,
		Amount: One,
	})
	// ASSERT: The action should've failed due to no allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Set mint allowance in state.
	k.SetMintAllowance(ctx, system.Address, One)

	// ACT: Attempt to mint with insufficient allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: system.Address,
		Amount: One.MulRaw(2),
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorIs(t, err, types.ErrInsufficientAllowance)

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to mint.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: system.Address,
		To:     user.Address,
		Amount: One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, bank.Balances[user.Address].AmountOf(k.Denom))
	require.True(t, k.GetMintAllowance(ctx, system.Address).IsZero())
	events := ctx.EventManager().Events()
	require.Len(t, events, 2)
	require.Equal(t, "florin.v1.MintAllowance", events[1].Type)
}

func TestRemoveAdminAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove admin account with no owner set.
	_, err := server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove admin account with invalid signer.
	_, err = server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set admin in state.
	admin := utils.TestAccount()
	k.SetAdmin(ctx, admin.Address)
	require.True(t, k.IsAdmin(ctx, admin.Address))

	// ACT: Attempt to remove admin account.
	_, err = server.RemoveAdminAccount(goCtx, &types.MsgRemoveAdminAccount{
		Signer:  owner.Address,
		Account: admin.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsAdmin(ctx, admin.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.AdminAccountRemoved", events[0].Type)
}

func TestRemoveSystemAccount(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove system account with no owner set.
	_, err := server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove system account with invalid signer.
	_, err = server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ARRANGE: Set system in state.
	system := utils.TestAccount()
	k.SetSystem(ctx, system.Address)
	require.True(t, k.IsSystem(ctx, system.Address))

	// ACT: Attempt to remove system account.
	_, err = server.RemoveSystemAccount(goCtx, &types.MsgRemoveSystemAccount{
		Signer:  owner.Address,
		Account: system.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.False(t, k.IsSystem(ctx, system.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.SystemAccountRemoved", events[0].Type)
}

func TestSetMaxMintAllowance(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set max mint allowance with no owner set.
	_, err := server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to set max mint allowance with invalid signer.
	_, err = server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ACT: Attempt to set max mint allowance.
	_, err = server.SetMaxMintAllowance(goCtx, &types.MsgSetMaxMintAllowance{
		Signer: owner.Address,
		Amount: MaxMintAllowance,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, k.GetMaxMintAllowance(ctx))
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
	k.SetAdmin(ctx, admin.Address)
	// ARRANGE: Set max mint allowance in state.
	k.SetMaxMintAllowance(ctx, MaxMintAllowance)

	// ACT: Attempt to set mint allowance with invalid signer.
	_, err := server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidAdmin)

	// ACT: Attempt to set mint allowance with negative amount.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Signer: admin.Address,
		Amount: One.Neg(),
	})
	// ASSERT: The action should've failed due to negative amount.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ACT: Attempt to set mint allowance to more than max.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Signer: admin.Address,
		Amount: MaxMintAllowance.Add(One),
	})
	// ASSERT: The action should've failed due to more than max.
	require.ErrorIs(t, err, types.ErrInvalidAllowance)

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to set mint allowance.
	_, err = server.SetMintAllowance(goCtx, &types.MsgSetMintAllowance{
		Signer:  admin.Address,
		Account: minter.Address,
		Amount:  One,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, k.GetMintAllowance(ctx, minter.Address))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.MintAllowance", events[0].Type)
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(goCtx, &types.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorIs(t, err, types.ErrNoOwner)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorIs(t, err, types.ErrInvalidOwner)

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorIs(t, err, types.ErrSameOwner)

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, owner.Address, k.GetOwner(ctx))
	require.Equal(t, pendingOwner.Address, k.GetPendingOwner(ctx))
	events := ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.v1.OwnershipTransferStarted", events[0].Type)
}
