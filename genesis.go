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
	k.SetBlacklistOwner(ctx, genesis.BlacklistState.Owner)
	k.SetBlacklistPendingOwner(ctx, genesis.BlacklistState.PendingOwner)
	for _, admin := range genesis.BlacklistState.Admins {
		k.SetBlacklistAdmin(ctx, admin)
	}
	for _, adversary := range genesis.BlacklistState.Adversaries {
		k.SetAdversary(ctx, adversary)
	}

	for _, denom := range genesis.AllowedDenoms {
		k.SetAllowedDenom(ctx, denom)
	}
	for denom, owner := range genesis.Owners {
		k.SetOwner(ctx, denom, owner)
	}
	for denom, pendingOwner := range genesis.PendingOwners {
		k.SetPendingOwner(ctx, denom, pendingOwner)
	}
	for _, system := range genesis.Systems {
		k.SetSystem(ctx, system.Denom, system.Address)
	}
	for _, admin := range genesis.Admins {
		k.SetAdmin(ctx, admin.Denom, admin.Address)
	}
	for _, item := range genesis.MintAllowances {
		k.SetMintAllowance(ctx, item.Denom, item.Address, item.Allowance)
	}
	for denom, rawMaxAllowance := range genesis.MaxMintAllowances {
		maxAllowance, _ := math.NewIntFromString(rawMaxAllowance)
		k.SetMaxMintAllowance(ctx, denom, maxAllowance)
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
