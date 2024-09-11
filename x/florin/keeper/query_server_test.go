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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/monerium/module-noble/utils"
	"github.com/monerium/module-noble/utils/mocks"
	"github.com/monerium/module-noble/x/florin/keeper"
	"github.com/monerium/module-noble/x/florin/types"
	"github.com/stretchr/testify/require"
)

func TestAllowedDenomsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query allowed denoms with invalid request.
	_, err := server.AllowedDenoms(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query allowed denoms.
	res, err := server.AllowedDenoms(goCtx, &types.QueryAllowedDenoms{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.AllowedDenoms, 1)
	require.Contains(t, res.AllowedDenoms, "ueure")
}

func TestOwnersQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owners with invalid request.
	_, err := server.Owners(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set an owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to query owners.
	res, err := server.Owners(goCtx, &types.QueryOwners{})
	// ASSERT: The query should've succeeded, with empty pending owners.
	require.NoError(t, err)
	require.Len(t, res.Owners, 1)
	require.Equal(t, owner.Address, res.Owners["ueure"])
	require.Empty(t, res.PendingOwners)

	// ARRANGE: Set a pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)

	// ACT: Attempt to query owners.
	res, err = server.Owners(goCtx, &types.QueryOwners{})
	// ASSERT: The query should've succeeded, with pending owners.
	require.NoError(t, err)
	require.Len(t, res.Owners, 1)
	require.Equal(t, owner.Address, res.Owners["ueure"])
	require.Len(t, res.PendingOwners, 1)
	require.Equal(t, pendingOwner.Address, res.PendingOwners["ueure"])
}

func TestOwnerQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owner with invalid request.
	_, err := server.Owner(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query owner of not allowed denom.
	_, err = server.Owner(goCtx, &types.QueryOwner{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ARRANGE: Set an owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, "ueure", owner.Address)

	// ACT: Attempt to query owner.
	res, err := server.Owner(goCtx, &types.QueryOwner{Denom: "ueure"})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set a pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)

	// ACT: Attempt to query owner.
	res, err = server.Owner(goCtx, &types.QueryOwner{Denom: "ueure"})
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
	k.SetSystem(ctx, "ueure", system1.Address)
	k.SetSystem(ctx, "ueure", system2.Address)

	// ACT: Attempt to query systems.
	res, err = server.Systems(goCtx, &types.QuerySystems{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Systems, 2)
	require.Contains(t, res.Systems, types.Account{Denom: "ueure", Address: system1.Address})
	require.Contains(t, res.Systems, types.Account{Denom: "ueure", Address: system2.Address})
}

func TestSystemsByDenomQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query systems by denom with invalid request.
	_, err := server.SystemsByDenom(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query systems by denom with not allowed denom.
	_, err = server.SystemsByDenom(goCtx, &types.QuerySystemsByDenom{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query systems by denom with no state.
	res, err := server.SystemsByDenom(goCtx, &types.QuerySystemsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Systems)

	// ARRANGE: Set system accounts in state.
	system1, system2 := utils.TestAccount(), utils.TestAccount()
	k.SetSystem(ctx, "ueure", system1.Address)
	k.SetSystem(ctx, "ueure", system2.Address)

	// ACT: Attempt to query systems by denom.
	res, err = server.SystemsByDenom(goCtx, &types.QuerySystemsByDenom{Denom: "ueure"})
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
	k.SetAdmin(ctx, "ueure", admin1.Address)
	k.SetAdmin(ctx, "ueure", admin2.Address)

	// ACT: Attempt to query admins.
	res, err = server.Admins(goCtx, &types.QueryAdmins{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Admins, 2)
	require.Contains(t, res.Admins, types.Account{Denom: "ueure", Address: admin1.Address})
	require.Contains(t, res.Admins, types.Account{Denom: "ueure", Address: admin2.Address})
}

func TestAdminsByDenomQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query admins by denom with invalid request.
	_, err := server.AdminsByDenom(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query admins by denom with not allowed denom.
	_, err = server.AdminsByDenom(goCtx, &types.QueryAdminsByDenom{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query admins by denom with no state.
	res, err := server.AdminsByDenom(goCtx, &types.QueryAdminsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Admins)

	// ARRANGE: Set admin accounts in state.
	admin1, admin2 := utils.TestAccount(), utils.TestAccount()
	k.SetAdmin(ctx, "ueure", admin1.Address)
	k.SetAdmin(ctx, "ueure", admin2.Address)

	// ACT: Attempt to query admins by denom.
	res, err = server.AdminsByDenom(goCtx, &types.QueryAdminsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Admins, 2)
	require.Contains(t, res.Admins, admin1.Address)
	require.Contains(t, res.Admins, admin2.Address)
}

func TestMaxMintAllowancesQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query max mint allowances with invalid request.
	_, err := server.MaxMintAllowances(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query max mint allowances with no state.
	res, err := server.MaxMintAllowances(goCtx, &types.QueryMaxMintAllowances{})
	// ASSERT: The query should've succeeded, returning zero.
	require.NoError(t, err)
	require.Empty(t, res.MaxMintAllowances)

	// ARRANGE: Set a max mint allowance in state.
	k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)

	// ACT: Attempt to query max mint allowances.
	res, err = server.MaxMintAllowances(goCtx, &types.QueryMaxMintAllowances{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.MaxMintAllowances, 1)
	require.Equal(t, MaxMintAllowance.String(), res.MaxMintAllowances["ueure"])
}

func TestMaxMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query max mint allowance with invalid request.
	_, err := server.MaxMintAllowance(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query max mint allowance of not allowed denom.
	_, err = server.MaxMintAllowance(goCtx, &types.QueryMaxMintAllowance{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query max mint allowance with no state.
	res, err := server.MaxMintAllowance(goCtx, &types.QueryMaxMintAllowance{Denom: "ueure"})
	// ASSERT: The query should've succeeded, returning zero.
	require.NoError(t, err)
	require.True(t, res.MaxMintAllowance.IsZero())

	// ARRANGE: Set a max mint allowance in state.
	k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)

	// ACT: Attempt to query max mint allowance.
	res, err = server.MaxMintAllowance(goCtx, &types.QueryMaxMintAllowance{Denom: "ueure"})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, res.MaxMintAllowance)
}

func TestMintAllowancesQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query mint allowances with invalid request.
	_, err := server.MintAllowances(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query mint allowances of not allowed denom.
	_, err = server.MintAllowances(goCtx, &types.QueryMintAllowances{
		Denom: "uusde",
	})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ARRANGE: Set mint allowance in state.
	minter := utils.TestAccount()
	k.SetMintAllowance(ctx, "ueure", minter.Address, One)

	// ACT: Attempt to query mint allowances.
	res, err := server.MintAllowances(goCtx, &types.QueryMintAllowances{
		Denom: "ueure",
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Allowances, 1)
	require.Equal(t, One.String(), res.Allowances[minter.Address])
}

func TestMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query mint allowance with invalid request.
	_, err := server.MintAllowance(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query mint allowance with not allowed denom.
	_, err = server.MintAllowance(goCtx, &types.QueryMintAllowance{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query mint allowance of random account.
	res, err := server.MintAllowance(goCtx, &types.QueryMintAllowance{
		Denom:   "ueure",
		Account: utils.TestAccount().Address,
	})
	// ASSERT: The query should've succeeded, returns zero.
	require.NoError(t, err)
	require.True(t, res.Allowance.IsZero())

	// ARRANGE: Set mint allowance in state.
	minter := utils.TestAccount()
	k.SetMintAllowance(ctx, "ueure", minter.Address, One)

	// ACT: Attempt to query mint allowance.
	res, err = server.MintAllowance(goCtx, &types.QueryMintAllowance{
		Denom:   "ueure",
		Account: minter.Address,
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, res.Allowance)
}
