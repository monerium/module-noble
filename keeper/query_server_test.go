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

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/monerium/module-noble/v2/keeper"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/utils"
	"github.com/monerium/module-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestAuthorityQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query authority with invalid request.
	_, err := server.Authority(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query authority.
	res, err := server.Authority(ctx, &types.QueryAuthority{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, "authority", res.Authority)
}

func TestAllowedDenomsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query allowed denoms with invalid request.
	_, err := server.AllowedDenoms(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query allowed denoms.
	res, err := server.AllowedDenoms(ctx, &types.QueryAllowedDenoms{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.AllowedDenoms, 1)
	require.Contains(t, res.AllowedDenoms, "ueure")
}

func TestOwnersQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owners with invalid request.
	_, err := server.Owners(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set an owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to query owners.
	res, err := server.Owners(ctx, &types.QueryOwners{})
	// ASSERT: The query should've succeeded, with empty pending owners.
	require.NoError(t, err)
	require.Len(t, res.Owners, 1)
	require.Equal(t, owner.Address, res.Owners["ueure"])
	require.Empty(t, res.PendingOwners)

	// ARRANGE: Set a pending owner in state.
	pendingOwner := utils.TestAccount()
	err = k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)
	require.NoError(t, err)

	// ACT: Attempt to query owners.
	res, err = server.Owners(ctx, &types.QueryOwners{})
	// ASSERT: The query should've succeeded, with pending owners.
	require.NoError(t, err)
	require.Len(t, res.Owners, 1)
	require.Equal(t, owner.Address, res.Owners["ueure"])
	require.Len(t, res.PendingOwners, 1)
	require.Equal(t, pendingOwner.Address, res.PendingOwners["ueure"])
}

func TestOwnerQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query owner with invalid request.
	_, err := server.Owner(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query owner of not allowed denom.
	_, err = server.Owner(ctx, &types.QueryOwner{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ARRANGE: Set an owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, "ueure", owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to query owner.
	res, err := server.Owner(ctx, &types.QueryOwner{Denom: "ueure"})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set a pending owner in state.
	pendingOwner := utils.TestAccount()
	err = k.SetPendingOwner(ctx, "ueure", pendingOwner.Address)
	require.NoError(t, err)

	// ACT: Attempt to query owner.
	res, err = server.Owner(ctx, &types.QueryOwner{Denom: "ueure"})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}

func TestSystemsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query systems with invalid request.
	_, err := server.Systems(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query systems with no state.
	res, err := server.Systems(ctx, &types.QuerySystems{})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Systems)

	// ARRANGE: Set system accounts in state.
	system1, system2 := utils.TestAccount(), utils.TestAccount()
	err = k.SetSystem(ctx, "ueure", system1.Address)
	require.NoError(t, err)
	err = k.SetSystem(ctx, "ueure", system2.Address)
	require.NoError(t, err)

	// ACT: Attempt to query systems.
	res, err = server.Systems(ctx, &types.QuerySystems{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Systems, 2)
	require.Contains(t, res.Systems, types.Account{Denom: "ueure", Address: system1.Address})
	require.Contains(t, res.Systems, types.Account{Denom: "ueure", Address: system2.Address})
}

func TestSystemsByDenomQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query systems by denom with invalid request.
	_, err := server.SystemsByDenom(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query systems by denom with not allowed denom.
	_, err = server.SystemsByDenom(ctx, &types.QuerySystemsByDenom{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query systems by denom with no state.
	res, err := server.SystemsByDenom(ctx, &types.QuerySystemsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Systems)

	// ARRANGE: Set system accounts in state.
	system1, system2 := utils.TestAccount(), utils.TestAccount()
	err = k.SetSystem(ctx, "ueure", system1.Address)
	require.NoError(t, err)
	err = k.SetSystem(ctx, "ueure", system2.Address)
	require.NoError(t, err)

	// ACT: Attempt to query systems by denom.
	res, err = server.SystemsByDenom(ctx, &types.QuerySystemsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Systems, 2)
	require.Contains(t, res.Systems, system1.Address)
	require.Contains(t, res.Systems, system2.Address)
}

func TestAdminsQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query admins with invalid request.
	_, err := server.Admins(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query admins with no state.
	res, err := server.Admins(ctx, &types.QueryAdmins{})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Admins)

	// ARRANGE: Set admin accounts in state.
	admin1, admin2 := utils.TestAccount(), utils.TestAccount()
	err = k.SetAdmin(ctx, "ueure", admin1.Address)
	require.NoError(t, err)
	err = k.SetAdmin(ctx, "ueure", admin2.Address)
	require.NoError(t, err)

	// ACT: Attempt to query admins.
	res, err = server.Admins(ctx, &types.QueryAdmins{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Admins, 2)
	require.Contains(t, res.Admins, types.Account{Denom: "ueure", Address: admin1.Address})
	require.Contains(t, res.Admins, types.Account{Denom: "ueure", Address: admin2.Address})
}

func TestAdminsByDenomQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query admins by denom with invalid request.
	_, err := server.AdminsByDenom(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query admins by denom with not allowed denom.
	_, err = server.AdminsByDenom(ctx, &types.QueryAdminsByDenom{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query admins by denom with no state.
	res, err := server.AdminsByDenom(ctx, &types.QueryAdminsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded, returns empty.
	require.NoError(t, err)
	require.Empty(t, res.Admins)

	// ARRANGE: Set admin accounts in state.
	admin1, admin2 := utils.TestAccount(), utils.TestAccount()
	err = k.SetAdmin(ctx, "ueure", admin1.Address)
	require.NoError(t, err)
	err = k.SetAdmin(ctx, "ueure", admin2.Address)
	require.NoError(t, err)

	// ACT: Attempt to query admins by denom.
	res, err = server.AdminsByDenom(ctx, &types.QueryAdminsByDenom{Denom: "ueure"})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Admins, 2)
	require.Contains(t, res.Admins, admin1.Address)
	require.Contains(t, res.Admins, admin2.Address)
}

func TestMaxMintAllowancesQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query max mint allowances with invalid request.
	_, err := server.MaxMintAllowances(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query max mint allowances with no state.
	res, err := server.MaxMintAllowances(ctx, &types.QueryMaxMintAllowances{})
	// ASSERT: The query should've succeeded, returning zero.
	require.NoError(t, err)
	require.Empty(t, res.MaxMintAllowances)

	// ARRANGE: Set a max mint allowance in state.
	err = k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)
	require.NoError(t, err)

	// ACT: Attempt to query max mint allowances.
	res, err = server.MaxMintAllowances(ctx, &types.QueryMaxMintAllowances{})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.MaxMintAllowances, 1)
	require.Equal(t, MaxMintAllowance.String(), res.MaxMintAllowances["ueure"])
}

func TestMaxMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query max mint allowance with invalid request.
	_, err := server.MaxMintAllowance(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query max mint allowance of not allowed denom.
	_, err = server.MaxMintAllowance(ctx, &types.QueryMaxMintAllowance{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query max mint allowance with no state.
	res, err := server.MaxMintAllowance(ctx, &types.QueryMaxMintAllowance{Denom: "ueure"})
	// ASSERT: The query should return err as the key is missing
	require.NoError(t, err)
	require.True(t, res.MaxMintAllowance.IsZero())

	// ARRANGE: Set a max mint allowance in state.
	err = k.SetMaxMintAllowance(ctx, "ueure", MaxMintAllowance)
	require.NoError(t, err)

	// ACT: Attempt to query max mint allowance.
	res, err = server.MaxMintAllowance(ctx, &types.QueryMaxMintAllowance{Denom: "ueure"})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, MaxMintAllowance, res.MaxMintAllowance)
}

func TestMintAllowancesQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query mint allowances with invalid request.
	_, err := server.MintAllowances(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query mint allowances of not allowed denom.
	_, err = server.MintAllowances(ctx, &types.QueryMintAllowances{
		Denom: "uusde",
	})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ARRANGE: Set mint allowance in state.
	minter := utils.TestAccount()
	err = k.SetMintAllowance(ctx, "ueure", minter.Address, One)
	require.NoError(t, err)

	// ACT: Attempt to query mint allowances.
	res, err := server.MintAllowances(ctx, &types.QueryMintAllowances{
		Denom: "ueure",
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Len(t, res.Allowances, 1)
	require.Equal(t, One.String(), res.Allowances[minter.Address])
}

func TestMintAllowanceQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	server := keeper.NewQueryServer(k)

	// ACT: Attempt to query mint allowance with invalid request.
	_, err := server.MintAllowance(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query mint allowance with not allowed denom.
	_, err = server.MintAllowance(ctx, &types.QueryMintAllowance{Denom: "uusde"})
	// ASSERT: The query should've failed due to not allowed denom.
	require.ErrorContains(t, err, "uusde is not an allowed denom")

	// ACT: Attempt to query mint allowance of random account.
	res, err := server.MintAllowance(ctx, &types.QueryMintAllowance{
		Denom:   "ueure",
		Account: utils.TestAccount().Address,
	})
	// ASSERT: The query should've failed, returns err.
	require.NoError(t, err)
	require.True(t, res.Allowance.IsZero())

	// ARRANGE: Set mint allowance in state.
	minter := utils.TestAccount()
	err = k.SetMintAllowance(ctx, "ueure", minter.Address, One)
	require.NoError(t, err)

	// ACT: Attempt to query mint allowance.
	res, err = server.MintAllowance(ctx, &types.QueryMintAllowance{
		Denom:   "ueure",
		Account: minter.Address,
	})
	// ASSERT: The query should've succeeded.
	require.NoError(t, err)
	require.Equal(t, One, res.Allowance)
}
