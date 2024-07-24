package florin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
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

	k.SetAuthority(ctx, genesis.Authority)
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
		maxAllowance, _ := sdk.NewIntFromString(rawMaxAllowance)
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
		Authority:         k.GetAuthority(ctx),
		AllowedDenoms:     k.GetAllowedDenoms(ctx),
		Owners:            k.GetOwners(ctx),
		PendingOwners:     k.GetPendingOwners(ctx),
		Systems:           k.GetSystems(ctx),
		Admins:            k.GetAdmins(ctx),
		MintAllowances:    k.GetMintAllowances(ctx),
		MaxMintAllowances: k.GetMaxMintAllowances(ctx),
	}
}
