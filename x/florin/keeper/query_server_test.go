package keeper_test

import (
	"testing"

	"github.com/noble-assets/florin/x/florin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/stretchr/testify/require"
)

func TestOwnerQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owner with invalid request.
	_, err := server.Owner(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to query owner.
	res, err := server.Owner(goCtx, &types.QueryOwner{})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to query owner.
	res, err = server.Owner(goCtx, &types.QueryOwner{})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}

func TestSystemsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query systems with invalid request.
	_, err := server.Systems(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query systems with no state.
	res, err := server.Systems(goCtx, &types.QuerySystems{})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Systems)

	// ARRANGE: Set system accounts in state.
	system1, system2 := utils.TestAccount(), utils.TestAccount()
	k.SetSystem(ctx, system1.Address)
	k.SetSystem(ctx, system2.Address)

	// ACT: Attempt to query systems.
	res, err = server.Systems(goCtx, &types.QuerySystems{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Systems, 2)
	require.Contains(t, res.Systems, system1.Address)
	require.Contains(t, res.Systems, system2.Address)
}

func TestAdminsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query admins with invalid request.
	_, err := server.Admins(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query admins with no state.
	res, err := server.Admins(goCtx, &types.QueryAdmins{})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Admins)

	// ARRANGE: Set admin accounts in state.
	admin1, admin2 := utils.TestAccount(), utils.TestAccount()
	k.SetAdmin(ctx, admin1.Address)
	k.SetAdmin(ctx, admin2.Address)

	// ACT: Attempt to query admins.
	res, err = server.Admins(goCtx, &types.QueryAdmins{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Admins, 2)
	require.Contains(t, res.Admins, admin1.Address)
	require.Contains(t, res.Admins, admin2.Address)
}

func TestMaxMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query max mint allowance with invalid request.
	_, err := server.MaxMintAllowance(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query max mint allowance with no state.
	res, err := server.MaxMintAllowance(goCtx, &types.QueryMaxMintAllowance{})
	// ASSERT: The query should've succeeded, returning zero.
	require.NoError(t, err)
	require.True(t, res.MaxMintAllowance.IsZero())

	// ARRANGE: Set max mint allowance in state.
	k.SetMaxMintAllowance(ctx, MaxMintAllowance)

	// ACT: Attempt to query max mint allowance.
	res, err = server.MaxMintAllowance(goCtx, &types.QueryMaxMintAllowance{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, res.MaxMintAllowance)
}

func TestMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query mint allowance with invalid request.
	_, err := server.MintAllowance(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query mint allowance of random account.
	res, err := server.MintAllowance(goCtx, &types.QueryMintAllowance{
		Account: utils.TestAccount().Address,
	})
	// ASSERT: The query should've succeeded, returns zero.
	require.NoError(t, err)
	require.True(t, res.Allowance.IsZero())

	// ARRANGE: Set mint allowance in state.
	minter := utils.TestAccount()
	k.SetMintAllowance(ctx, minter.Address, One)

	// ACT: Attempt to query mint allowance.
	res, err = server.MintAllowance(goCtx, &types.QueryMintAllowance{
		Account: minter.Address,
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, res.Allowance)
}
