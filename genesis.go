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

package florin

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/keeper"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/types/blacklist"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genesis types.GenesisState) {
	if err := k.SetBlacklistOwner(ctx, genesis.BlacklistState.Owner); err != nil {
		panic(err)
	}
	if err := k.SetBlacklistPendingOwner(ctx, genesis.BlacklistState.PendingOwner); err != nil {
		panic(err)
	}
	for _, admin := range genesis.BlacklistState.Admins {
		if err := k.SetBlacklistAdmin(ctx, admin); err != nil {
			panic(err)
		}
	}
	for _, adversary := range genesis.BlacklistState.Adversaries {
		if err := k.SetAdversary(ctx, adversary); err != nil {
			panic(err)
		}
	}

	for _, denom := range genesis.AllowedDenoms {
		if err := k.SetAllowedDenom(ctx, denom); err != nil {
			panic(err)
		}
	}
	for denom, owner := range genesis.Owners {
		if err := k.SetOwner(ctx, denom, owner); err != nil {
			panic(err)
		}
	}
	for denom, pendingOwner := range genesis.PendingOwners {
		if err := k.SetPendingOwner(ctx, denom, pendingOwner); err != nil {
			panic(err)
		}
	}
	for _, system := range genesis.Systems {
		if err := k.SetSystem(ctx, system.Denom, system.Address); err != nil {
			panic(err)
		}
	}
	for _, admin := range genesis.Admins {
		if err := k.SetAdmin(ctx, admin.Denom, admin.Address); err != nil {
			panic(err)
		}
	}
	for _, item := range genesis.MintAllowances {
		if err := k.SetMintAllowance(ctx, item.Denom, item.Address, item.Allowance); err != nil {
			panic(err)
		}
	}
	for denom, rawMaxAllowance := range genesis.MaxMintAllowances {
		maxAllowance, _ := math.NewIntFromString(rawMaxAllowance)
		if err := k.SetMaxMintAllowance(ctx, denom, maxAllowance); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlacklistState: blacklist.GenesisState{
			Owner:        k.GetBlacklistOwner(ctx),
			PendingOwner: k.GetBlacklistPendingOwner(ctx),
			Admins:       k.GetBlacklistAdmins(ctx),
			Adversaries:  k.GetAdversaries(ctx),
		},
		AllowedDenoms:     k.GetAllowedDenoms(ctx),
		Owners:            k.GetOwners(ctx),
		PendingOwners:     k.GetPendingOwners(ctx),
		Systems:           k.GetSystems(ctx),
		Admins:            k.GetAdmins(ctx),
		MintAllowances:    k.GetMintAllowances(ctx),
		MaxMintAllowances: k.GetMaxMintAllowances(ctx),
	}
}
